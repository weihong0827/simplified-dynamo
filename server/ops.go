package main

import (
	"context"
	"dynamoSimplified/config"
	"dynamoSimplified/hash"
	pb "dynamoSimplified/pb"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc/status"
)

const (
	defaultTimeout  = time.Second
	replicaError    = "error with replica operation: %v"
	connectionError = "failed to connect to node: %v"
	// simulateFailure = true
)

// GRPCOperation represents a function type for gRPC operations.
type GRPCOperation func(ctx context.Context, client pb.KeyValueStoreClient, kv *pb.KeyValue, result chan<- *pb.KeyValue, nodeError chan<- uint32, nodeId uint32) error

// Make a gRPC read call.
func performRead(
	ctx context.Context,
	client pb.KeyValueStoreClient,
	kv *pb.KeyValue,
	result chan<- *pb.KeyValue,
	nodeErr chan<- uint32,
	nodeId uint32,
) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in closed channel when performing Read", r)
		}
	}()

	r, err := client.Read(ctx, &pb.ReadRequest{Key: kv.Key, IsReplica: true})
	if err != nil {
		if errors.Is(err, status.Error(505, nodeFailure)) {
			nodeErr <- nodeId
		}
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
	nodeErr chan<- uint32,
	nodeId uint32,
) error {
	// Placeholder: Add your write call and its specifics here.
	// Example: `w, err := client.Write(ctx, &pb.WriteRequest{Key: key, Value: value})`

	// Here's a hypothetical write success message. Adjust it to match your actual API.
	// result <- pb.KeyValue{Key: key, Value: "Write Successful!"}
	r, err := client.Write(ctx, &pb.WriteRequest{KeyValue: kv, IsReplica: true})
	if err != nil {
		if errors.Is(err, status.Error(505, nodeFailure)) {
			nodeErr <- nodeId
		}
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
// func performOwnerNodeWrite(
// 	ctx context.Context,
// 	client pb.KeyValueStoreClient,
// 	node *pb.Node,
// 	kv *pb.KeyValue,
// 	succesfulSend chan<- *pb.Node,
// ) error {

// 	// Placeholder: Add your write call and its specifics here.
// 	// Example: `w, err := client.Write(ctx, &pb.WriteRequest{Key: key, Value: value})`

// 	// Here's a hypothetical write success message. Adjust it to match your actual API.
// 	// result <- pb.KeyValue{Key: key, Value: "Write Successful!"}
// 	r, err := client.Write(ctx, &pb.WriteRequest{KeyValue: kv, IsReplica: true})
// 	if err != nil {
// 		return fmt.Errorf(replicaError, err)
// 	}
// 	if r.Success {
// 		succesfulSend <- node
// 	} else {
// 		return fmt.Errorf(replicaError, "unexpected response format")
// 	}
// 	return nil
// }

func performHintedHandoffWrite(
	ctx context.Context,
	client pb.KeyValueStoreClient,
	node *pb.Node,
	kv *pb.KeyValue,
	result chan<- *pb.KeyValue,
	nodeErr chan<- uint32,
	nodeId uint32,
) error {
	//TODO: call the write and handle error and return
	r, err := client.HintedHandoffWrite(ctx, &pb.HintedHandoffWriteRequest{KeyValue: kv, Nodeid: node.Id})
	if err != nil {
		return fmt.Errorf(replicaError, err)
	}
	if r.Success {
		// return fmt.Errorf("Hinted Handoff Write successful")
		result <- r.KeyValue
		return nil
	} else {
		return fmt.Errorf(replicaError, "unexpected response format")
	}
}

func hintedHandoffGrpcCall(ctx context.Context,
	node *pb.Node,
	kv *pb.KeyValue,
) error {

	_, callCancel := context.WithTimeout(ctx, defaultTimeout)
	defer callCancel()

	conn, err := CreateGRPCConnection(node.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)
	r, err := client.Write(ctx, &pb.WriteRequest{KeyValue: kv, IsReplica: true})
	if err != nil {
		return fmt.Errorf(replicaError, err)
	}
	if !r.Success {
		return fmt.Errorf(replicaError, "unexpected response format")
	}
	return nil
}

// grpcCall performs the given gRPC operation on the specified node.
func grpcCall(
	ctx context.Context,
	cancel context.CancelFunc,
	node *pb.Node,
	kv *pb.KeyValue,
	op config.Operation,
	nodes hash.NodeSlice,
	result chan<- *pb.KeyValue,
	nodeErrors chan<- uint32,
) error {

	var operation GRPCOperation
	switch op {
	case config.READ: //TODO: timeout on required responses
		operation = performRead
	case config.WRITE:
		operation = performWrite
	}

	callCtx, callCancel := context.WithTimeout(ctx, defaultTimeout)
	defer callCancel()

	log.Printf("coordinator connecting to node %v", node.Address)
	conn, err := CreateGRPCConnection(node.Address)
	client := pb.NewKeyValueStoreClient(conn)
	err = operation(callCtx, client, kv, result, nodeErrors, node.Id)
	conn.Close()
	if err != nil && op == config.WRITE {
		log.Print("Error in creating grpc connection, finding hinted handoff node:", err)
		offsetid := config.N - 1
		var successors []*pb.Node
		for err != nil && offsetid != 0 { // if fail to establish contact with node, try to contact the next node
			offsetid = (offsetid + 1) % len(nodes)
			log.Print(offsetid)
			successors, _ = hash.GetNodeFromKeyWithOffSet([]int{offsetid}, hash.GenHash(kv.Key), nodes)
			log.Print("contacting successor node: ", successors[0].Address)
			conn, err = CreateGRPCConnection(successors[0].Address)
			client := pb.NewKeyValueStoreClient(conn)
			err = performHintedHandoffWrite(callCtx, client, node, kv, result, nodeErrors, node.Id)
			// log.Print("Hinted handoff error ", err)
			conn.Close()
		}
		log.Print("Hinted handoff node found and performed: ", successors[0].Address)
		return err

	}
	return err
}

// func SendToNode(
// 	ctx context.Context,
// 	nodes []*pb.HintedHandoffWriteRequest,
// ) []*pb.Node {
// 	succesfulSend := make(chan *pb.Node)
// 	defer close(succesfulSend)

// 	var wg sync.WaitGroup

// 	for _, node := range nodes {
// 		wg.Add(1)
// 		go func(n *pb.HintedHandoffWriteRequest) {
// 			defer wg.Done()
// 			log.Printf("GRCP call to hintedhandoff %v", n.Node)

// 			if err := hintedHandoffGrpcCall(ctx, n.Node, n.KeyValue); err != nil {
// 				log.Println("Error in gRPC call:", err)
// 			}
// 		}(node)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(succesfulSend)
// 	}()

// 	var collectedSentNodes []*pb.Node

// 	go func() {
// 		for res := range succesfulSend {
// 			collectedSentNodes = append(collectedSentNodes, res)
// 		}
// 	}()

// 	return collectedSentNodes
// }

// Sends requests to the appropriate replicas.
func SendRequestToReplica(
	kv *pb.KeyValue,
	nodes hash.NodeSlice,
	op config.Operation,
	currAddr string,
	coordsuccess bool,
	respChan chan<- []*pb.KeyValue,
	errorChan chan<- []uint32,
) {
	targetNodes, err := hash.GetNodesFromKey(hash.GenHash(kv.Key), nodes)
	if err != nil {
		log.Println("Error obtaining nodes for key:", err)
		respChan <- nil
	}

	// simulate failure
	// if simulateFailure {

	// 	target := targetNodes[0].Address
	// 	log.Printf("Selected crash node %v, %s", targetNodes[0], target)

	// 	for _, node := range nodes {
	// 		if node.Address == target {
	// 			for {
	// 				if node.IsAlive == false {
	// 					log.Printf("Target Node has been crashed forcefully. Node ID: %d, Node Address: %s", node.Id, node.Address)
	// 					time.Sleep(time.Second * 5)
	// 					log.Printf("========\n Programme continued after forced crash simulation \n========\n")

	// 					break
	// 				}
	// 			}
	// 		}
	// 	}

	// 	// for {
	// 	// 	time.Sleep(time.Second * 5)

	// 	// 	log.Printf("GRPC connection crash node ")

	// 	// 	conn, err := CreateGRPCConnection(target)
	// 	// 	if err != nil {
	// 	// 		log.Printf("GRPC call failed. Target Node has been crashed forcefully. Node Address: %s", target)
	// 	// 		time.Sleep(time.Second * 5)
	// 	// 		conn.Close()

	// 	// 		log.Printf("========\n Programme continued after forced crash simulation \n========\n")

	// 	// 		break
	// 	// 	} else {
	// 	// 		conn.Close()
	// 	// 	}
	// 	// }

	// }

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

	nodeErrors := make(chan uint32)
	defer close(nodeErrors)
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
			if err := grpcCall(ctx, cancel, n, kv, op, nodes, result, nodeErrors); err != nil {
				log.Println("Error in gRPC call:", err)
			}
		}(node)
	}

	// Collect results until the desired number of responses is reached
	var collectedResults []*pb.KeyValue
	var errorResults []uint32
	for {
		select {
		case res := <-result:
			collectedResults = append(collectedResults, res)
		case grpCallError := <-nodeErrors:
			errorResults = append(errorResults, grpCallError)
		}
		if len(collectedResults) >= int(requiredResponses) {
			respChan <- collectedResults
			errorChan <- errorResults
			break
		}
	}
	wg.Wait()

}
