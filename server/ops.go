package main

import (
	"context"
	"dynamoSimplified/config"
	"dynamoSimplified/hash"
	pb "dynamoSimplified/pb"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

const (
	defaultTimeout  = time.Second
	replicaError    = "error with replica operation: %v"
	connectionError = "failed to connect to node: %v"
)

// GRPCOperation represents a function type for gRPC operations.
type GRPCOperation func(ctx context.Context, client pb.KeyValueStoreClient, kv pb.KeyValue, result chan<- *pb.KeyValue) error

// Make a gRPC read call.
func performRead(
	ctx context.Context,
	client pb.KeyValueStoreClient,
	kv pb.KeyValue,
	result chan<- *pb.KeyValue,
) error {

	r, err := client.Read(ctx, &pb.ReadRequest{Key: kv.Key, IsReplica: true})
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
	kv pb.KeyValue,
	result chan<- *pb.KeyValue,
) error {

	// Placeholder: Add your write call and its specifics here.
	// Example: `w, err := client.Write(ctx, &pb.WriteRequest{Key: key, Value: value})`

	// Here's a hypothetical write success message. Adjust it to match your actual API.
	// result <- pb.KeyValue{Key: key, Value: "Write Successful!"}
	r, err := client.Write(ctx, &pb.WriteRequest{KeyValue: &kv, IsReplica: true})
	if err != nil {
		return fmt.Errorf(replicaError, err)
	}
	if r.Success {
		result <- r.KeyValue[0]
	} else {
		return fmt.Errorf(replicaError, "unexpected response format")
	}
	return nil
}

func createGRPCConnection(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf(connectionError, err)
	}
	return conn, nil
}

// Make a gRPC write call.
func performOriNodeWrite(
	ctx context.Context,
	client pb.KeyValueStoreClient,
	node *pb.Node,
	kv pb.KeyValue,
	succesfulSend chan<- *pb.Node,
) error {

	// Placeholder: Add your write call and its specifics here.
	// Example: `w, err := client.Write(ctx, &pb.WriteRequest{Key: key, Value: value})`

	// Here's a hypothetical write success message. Adjust it to match your actual API.
	// result <- pb.KeyValue{Key: key, Value: "Write Successful!"}
	r, err := client.Write(ctx, &pb.WriteRequest{KeyValue: &kv, IsReplica: true})
	if err != nil {
		return fmt.Errorf(replicaError, err)
	}
	if r.Success {
		succesfulSend <- node
	} else {
		return fmt.Errorf(replicaError, "unexpected response format")
	}
	return nil
}

func performHintedHandoffWrite(
	ctx context.Context,
	client pb.KeyValueStoreClient,
	node *pb.Node,
	kv pb.KeyValue,
) error {
	//TODO: call the write and handle error and return
	r, err := client.HintedHandoffWrite(ctx, &pb.HintedHandoffWriteRequest{KeyValue: &kv, Node: node})
	if err != nil {
		return fmt.Errorf(replicaError, err)
	}
	if r.Success {
		// return fmt.Errorf("Hinted Handoff Write successful")
		return nil
	} else {
		return fmt.Errorf(replicaError, "unexpected response format")
	}
}

func hintedHandoffGrpcCall(ctx context.Context,
	node *pb.Node,
	kv pb.KeyValue,
	timeout time.Duration,
	succesfulSend chan<- *pb.Node,
) error {
	callCtx, callCancel := context.WithTimeout(ctx, timeout)
	defer callCancel()

	conn, err := createGRPCConnection(node.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)
	return performOriNodeWrite(callCtx, client, node, kv, succesfulSend)
}

// grpcCall performs the given gRPC operation on the specified node.
func grpcCall(
	ctx context.Context,
	node *pb.Node,
	kv pb.KeyValue,
	timeout time.Duration,
	op config.Operation,
	nodes hash.NodeSlice,
	result chan<- *pb.KeyValue,
) error {

	var operation GRPCOperation
	switch op {
	case config.READ: //TODO: timeout on required responses
		operation = performRead
	case config.WRITE:
		operation = performWrite
	}

	callCtx, callCancel := context.WithTimeout(ctx, timeout)
	defer callCancel()

	conn, err := createGRPCConnection(node.Address)
	if err != nil {
		// --> TODO: update membership list.

		//check op its a write req or read req
		// if write
		//TODO need to check what the error is and if needed perform hinted handoff
		// create connection
		// set op to performHintedHandoff

		switch op {
		case config.READ:
			return err
		case config.WRITE:

			successor := hash.GetSuccessiveNode(node.Id, nodes, op)
			conn, err = createGRPCConnection(successor.Address)

			defer conn.Close()

			client := pb.NewKeyValueStoreClient(conn)
			return performHintedHandoffWrite(callCtx, client, node, kv)
		}
	}
	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)
	return operation(callCtx, client, kv, result)
}

func SendToNode(
	ctx context.Context,
	nodes hash.NodeKvSlice,
) []*pb.Node {

	done := make(chan bool)
	defer close(done)

	succesfulSend := make(chan *pb.Node)
	defer close(succesfulSend)

	for _, node := range nodes {

		// TODO close hinted handoff after all make grpc call
		// defer done <- true
		go func(n *pb.HintedHandoffWriteRequest) {
			if err := hintedHandoffGrpcCall(ctx, n.Node, *n.KeyValue, defaultTimeout, succesfulSend); err != nil {
				log.Println("Error in gRPC call:", err)
			}

		}(node)
	}

	var collectedSentNodes []*pb.Node
collect:
	for {
		select {
		case res := <-succesfulSend:
			collectedSentNodes = append(collectedSentNodes, res)
		case <-done:
			break collect
		}

	}
	return collectedSentNodes
}

// Sends requests to the appropriate replicas.
func SendRequestToReplica(
	kv pb.KeyValue,
	nodes hash.NodeSlice,
	op config.Operation,
	currAddr string,
) []*pb.KeyValue {
	targetNodes, err := hash.GetNodesFromKey(hash.GenHash(kv.Key), nodes, config.READ)
	if err != nil {
		log.Println("Error obtaining nodes for key:", err)
		return nil
	}

	var requiredResponses int32

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
		if node.Address == currAddr {
			continue
		}
		go func(n *pb.Node) {
			if err := grpcCall(ctx, n, kv, defaultTimeout, op, nodes, result); err != nil {
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
