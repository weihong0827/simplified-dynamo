package main

import (
	"context"
	"dynamoSimplified/config"
	pb "dynamoSimplified/pb"
	"dynamoSimplified/utils"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"sync/atomic"
	"time"
)

const (
	defaultTimeout  = time.Second
	replicaError    = "error with replica operation: %v"
	connectionError = "failed to connect to node: %v"
)

// GRPCOperation represents a function type for gRPC operations.
type GRPCOperation func(ctx context.Context, client pb.KeyValueStoreClient, key string, result chan<- *pb.KeyValue) error

// Make a gRPC read call.
func performRead(
	ctx context.Context,
	client pb.KeyValueStoreClient,
	key string,
	result chan<- *pb.KeyValue,
) error {

	r, err := client.Read(ctx, &pb.ReadRequest{Key: key, IsReplica: false})
	if err != nil {
		return fmt.Errorf(replicaError, err)
	}
	if r.Success && len(r.GetKeyValue()) == 1 {
		result <- r.KeyValue[0]
	} else {
		return fmt.Errorf(replicaError, "unexpected response format")
	}
	return nil
}

// Make a gRPC write call.
func performWrite(
	ctx context.Context,
	client pb.KeyValueStoreClient,
	key string,
	result chan<- *pb.KeyValue,
) error {

	// Placeholder: Add your write call and its specifics here.
	// Example: `w, err := client.Write(ctx, &pb.WriteRequest{Key: key, Value: value})`

	// Here's a hypothetical write success message. Adjust it to match your actual API.
	// result <- pb.KeyValue{Key: key, Value: "Write Successful!"}
	return nil
}

// grpcCall performs the given gRPC operation on the specified node.
func grpcCall(
	ctx context.Context,
	node *pb.Node,
	key string,
	op GRPCOperation,
	timeout time.Duration,
	result chan<- *pb.KeyValue,
) error {
	callCtx, callCancel := context.WithTimeout(ctx, timeout)
	defer callCancel()

	conn, err := createGRPCConnection(node.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)
	return op(callCtx, client, key, result)
}

// Sends requests to the appropriate replicas.
func SendRequestToReplica(
	key string,
	nodes utils.NodeSlice,
	op config.Operation,
	currentMachinePort uint32,
) []*pb.KeyValue {
	targetNodes, err := utils.GetNodesFromKey(utils.GenHash(key), nodes, config.READ)
	if err != nil {
		log.Println("Error obtaining nodes for key:", err)
		return nil
	}

	var operation GRPCOperation
	var requiredResponses int32
	switch op {
	case config.READ:
		operation = performRead
		requiredResponses = int32(config.R)
	case config.WRITE:
		operation = performWrite
		requiredResponses = int32(config.W)
	}

	result := make(chan *pb.KeyValue)
	defer close(result)

	var responseCounter int32
	done := make(chan bool)
	defer close(done)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure all resources are cleaned up

	// Monitoring goroutine
	go func() {
		for range result {
			if atomic.AddInt32(&responseCounter, 1) >= requiredResponses {
				done <- true
				break
			}
		}
	}()

	for _, node := range targetNodes {
		if node.Port == currentMachinePort {
			continue
		}
		go func(n *pb.Node) {
			if err := grpcCall(ctx, n, key, operation, defaultTimeout, result); err != nil {
				log.Println("Error in gRPC call:", err)
			}
		}(node)
	}

	// Collect results until the desired number of responses is reached
	var collectedResults []*pb.KeyValue
collect:
	for {
		select {
		case res := <-result:
			collectedResults = append(collectedResults, res)
		case <-done:
			break collect
		}
	}

	return collectedResults
}
