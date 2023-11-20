package main

import (
	"context"
	"dynamoSimplified/config"
	"dynamoSimplified/hash"
	pb "dynamoSimplified/pb"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	defaultTimeout  = time.Second
	replicaError    = "error with replica operation: %v"
	connectionError = "failed to connect to node: %v"
)

// GRPCOperation represents a function type for gRPC operations.
type GRPCOperation func(ctx context.Context, client pb.KeyValueStoreClient, kv *pb.KeyValue, result chan<- *pb.KeyValue) error

// Make a gRPC read call.
func performRead(
	ctx context.Context,
	client pb.KeyValueStoreClient,
	kv *pb.KeyValue,
	result chan<- *pb.KeyValue,
) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in closed channel when performing Read", r)
		}
	}()

	r, err := client.Read(ctx, &pb.ReadRequest{Key: kv.Key, IsReplica: true})
	if err != nil {
		return fmt.Errorf(replicaError, err)
	}
	if r.Success && len(r.GetKeyValue()) == 1 {
		log.Print("node read succes received at perform read")
		result <- r.KeyValue[0]
	} else {
		log.Print(r.Message)
		result <- &pb.KeyValue{Key: "Read", Value: "Read Failed!"}
		return fmt.Errorf(replicaError, r.Message)
	}
	return nil
}

// Make a gRPC write call.
func performWrite(
	ctx context.Context,
	client pb.KeyValueStoreClient,
	kv *pb.KeyValue,
	result chan<- *pb.KeyValue,
) error {
	// Placeholder: Add your write call and its specifics here.
	// Example: `w, err := client.Write(ctx, &pb.WriteRequest{Key: key, Value: value})`

	// Here's a hypothetical write success message. Adjust it to match your actual API.
	// result <- pb.KeyValue{Key: key, Value: "Write Successful!"}
	r, err := client.Write(ctx, &pb.WriteRequest{KeyValue: kv, IsReplica: true})
	if err != nil {
		return fmt.Errorf(replicaError, err)
	}
	if r.Success {
		result <- r.KeyValue[0]
	} else {
		return fmt.Errorf(replicaError, r.Message)
	}
	return nil
}

// Make a gRPC write call.
func performOwnerNodeWrite(
	ctx context.Context,
	client pb.KeyValueStoreClient,
	node *pb.Node,
	kv *pb.KeyValue,
	succesfulSend chan<- *pb.Node,
) error {

	// Placeholder: Add your write call and its specifics here.
	// Example: `w, err := client.Write(ctx, &pb.WriteRequest{Key: key, Value: value})`

	// Here's a hypothetical write success message. Adjust it to match your actual API.
	// result <- pb.KeyValue{Key: key, Value: "Write Successful!"}
	r, err := client.Write(ctx, &pb.WriteRequest{KeyValue: kv, IsReplica: true})
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
	kv *pb.KeyValue,
) error {
	//TODO: call the write and handle error and return
	r, err := client.HintedHandoff(ctx, &pb.HintedHandoffWriteRequest{KeyValue: kv, Node: node})
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
	kv *pb.KeyValue,
	timeout time.Duration,
	succesfulSend chan<- *pb.Node,
) error {

	callCtx, callCancel := context.WithTimeout(ctx, timeout)
	defer callCancel()

	conn, err := CreateGRPCConnection(node.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)
	return performOwnerNodeWrite(callCtx, client, node, kv, succesfulSend)
}

// grpcCall performs the given gRPC operation on the specified node.
func grpcCall(
	ctx context.Context,
	cancel context.CancelFunc,
	node *pb.Node,
	kv *pb.KeyValue,
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

	log.Printf("coordinator connecting to node %v", node.Address)
	conn, err := CreateGRPCConnection(node.Address)
	if err != nil {
		//TODO: update membership list.
		//check op its a write req or read req
		// if write
		//TODO need to check what the error is and if needed perform hinted handoff
		// create connection
		// set op to performHintedHandoff

		switch op {
		case config.READ:
			return err
		case config.WRITE:

			successor := hash.GetSuccessiveNode(node, nodes, op)
			conn, err = CreateGRPCConnection(successor.Address)

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
	nodes []*pb.HintedHandoffWriteRequest,
) []*pb.Node {
	succesfulSend := make(chan *pb.Node)
	defer close(succesfulSend)

	var wg sync.WaitGroup

	for _, node := range nodes {
		wg.Add(1)
		go func(n *pb.HintedHandoffWriteRequest) {
			defer wg.Done()
			log.Printf("GRCP call to hintedhandoff %v", n.Node)

			if err := hintedHandoffGrpcCall(ctx, n.Node, n.KeyValue, defaultTimeout, succesfulSend); err != nil {
				log.Println("Error in gRPC call:", err)
			}
		}(node)
	}

	go func() {
		wg.Wait()
		close(succesfulSend)
	}()

	var collectedSentNodes []*pb.Node

	go func() {
		for res := range succesfulSend {
			collectedSentNodes = append(collectedSentNodes, res)
		}
	}()

	return collectedSentNodes
}

// Sends requests to the appropriate replicas.
func SendRequestToReplica(
	kv *pb.KeyValue,
	nodes hash.NodeSlice,
	op config.Operation,
	currAddr string,
	coordsuccess bool,
	respChan chan<- []*pb.KeyValue,
) {
	targetNodes, err := hash.GetNodesFromKey(hash.GenHash(kv.Key), nodes)
	if err != nil {
		log.Println("Error obtaining nodes for key:", err)
		respChan <- nil
	}

	// var operation GRPCOperation
	var requiredResponses int32
	var wg sync.WaitGroup

	switch op {
	case config.READ: //TODO: timeout on required responses
		// operation = performRead
		if coordsuccess {
			requiredResponses = int32(config.R - 1)
		} else {
			requiredResponses = int32(config.R)
		}
	case config.WRITE:
		// operation = performWrite
		if coordsuccess {
			requiredResponses = int32(config.W - 1)
		} else {
			requiredResponses = int32(config.W)
		}
	}

	result := make(chan *pb.KeyValue)
	defer close(result)

	// var responseCounter int32
	// done := make(chan bool)
	// defer close(done)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure all resources are cleaned up

	// Monitoring goroutine
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for range result {
	// 		if atomic.AddInt32(&responseCounter, 1) >= requiredResponses {
	// 			log.Print("required responses received")
	// 			done <- true
	// 			break
	// 		}
	// 		log.Print(responseCounter)
	// 	}
	// }()
	log.Print(len(targetNodes))
	for _, node := range targetNodes {
		if node.Address == currAddr {
			continue
		}
		wg.Add(1)
		go func(n *pb.Node) {
			defer wg.Done()
			if err := grpcCall(ctx, cancel, n, kv, defaultTimeout, op, nodes, result); err != nil {
				log.Println("Error in gRPC call:", err)
			}
		}(node)
	}

	// Collect results until the desired number of responses is reached
	var collectedResults []*pb.KeyValue
	for {
		res := <-result
		collectedResults = append(collectedResults, res)

		if len(collectedResults) >= int(requiredResponses) {
			respChan <- collectedResults
			break
		}
	}
	wg.Wait()

}
