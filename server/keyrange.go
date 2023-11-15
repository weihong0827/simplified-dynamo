package main

/*
import (
	"context"
	"dynamoSimplified/config"
	"dynamoSimplified/hash"
	pb "dynamoSimplified/pb"
	"log"
	"time"
)

const maxRetries = 10

// 1. Define the strategy interface
type KeyRangeModifier interface {
	Modify(key uint32, value pb.KeyValue)
}

type DeleteOperation struct {
	TargetNode *pb.Node
}

func (d *DeleteOperation) Modify(key uint32, value pb.KeyValue) {

}

type TransferOperation struct {
	TargetNode *pb.Node
}

func (t *TransferOperation) Modify(store *map[uint32]pb.KeyValue, key uint32, value pb.KeyValue) {
	err := writeToTarget(&value, t.TargetNode)
	if err != nil {
		log.Println("Failed to transfer, starting retry:", err)
		go retryTransfer(&value, t.TargetNode)
	}
}

func ModifyData(
	start uint32,
	end uint32,
	operation KeyRangeModifier,
) {
	for key, value := range *store {
		if key >= start && key <= end {
			operation.Modify(key, value)
		}
	}
}

func InitiateKeyRangeChange(
	membershipList *pb.MembershipList,
	coordinatorNode *pb.Node,
	newNode *pb.Node,
	store *map[uint32]pb.KeyValue,
) {
	// Get n nodes
	// We are gettinging 1 node before the coordinatorNode
	// and n nodes after the coordinatorNode including the coordinatorNode
	var offsets []int
	offsets = append(offsets, -1)
	for i := 0; i < config.N; i++ {
		offsets = append(offsets, i)
	}

	nodes, err := hash.GetNodeFromKeyWithOffSet(offsets, coordinatorNode.Id, membershipList.Nodes)
	if err != nil {
		log.Println("Error When assigning key range change:", err)
		return
	}
	//modify coordinatorNode

}
func DeleteFromTarget(kv *pb.KeyValue, targetNode *pb.Node) error {
	conn, err := CreateGRPCConnection(targetNode.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)
	_, err = client.Delete(context.Background(), &pb.ReadRequest{
		KeyValue: kv,
	})
	return err
}

func WriteToTarget(kv *pb.KeyValue, targetNode *pb.Node) error {
	conn, err := CreateGRPCConnection(targetNode.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)
	_, err = client.Write(context.Background(), &pb.WriteRequest{
		IsReplica: false,
		KeyValue:  kv,
	})
	return err
}

func retryTransfer(kv *pb.KeyValue, targetNode *pb.Node) {
	for i := 0; i < maxRetries; i++ {
		err := writeToTarget(kv, targetNode)
		if err != nil {
			log.Println("Retry failed:", err)
			time.Sleep(1 * time.Second)
		} else {
			return
		}
	}
	log.Printf("Failed to transfer after %d retries\n", maxRetries)
}

*/
