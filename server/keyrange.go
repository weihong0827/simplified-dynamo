package main

import (
	"context"
	"errors"
	"log"

	pb "dynamoSimplified/pb"
)

func Transfer(
	store map[uint32]pb.KeyValue,
	start uint32,
	end uint32,
	targetNode *pb.Node,
) (*pb.Empty, error) {
	log.Print("bulk write function 1")
	dataToTransfer := []*pb.KeyValue{}
	for key, value := range store {
		if IsKeyInRange(key, start, end) {
			dataToTransfer = append(dataToTransfer, &value)
		}
	}
	err := BulkWriteToTarget(dataToTransfer, targetNode)
	if err != nil {
		log.Println("Error when transferring data:", err)
		return &pb.Empty{}, errors.New("error when transferring data")
	}
	return &pb.Empty{}, nil
}

func BulkWriteToTarget(kvToTransfer []*pb.KeyValue, targetNode *pb.Node) error {
	log.Printf("bulk write function 2 to Address %s", targetNode.Address)
	conn, err := CreateGRPCConnection(targetNode.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)
	_, err = client.BulkWrite(context.Background(), &pb.BulkWriteRequest{
		KeyValue: kvToTransfer,
	})
	return err
}

func DeleteReplicaFromTarget(target *pb.Node, start uint32, end uint32) error {
	conn, err := CreateGRPCConnection(target.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)
	_, err = client.Delete(context.Background(), &pb.ReplicaDeleteRequest{
		Start: start,
		End:   end,
	})
	return err
}

// func retryTransfer(kv *pb.KeyValue, targetNode *pb.Node) {
// 	for i := 0; i < maxRetries; i++ {
// 		err := writeToTarget(kv, targetNode)
// 		if err != nil {
// 			log.Println("Retry failed:", err)
// 			time.Sleep(1 * time.Second)
// 		} else {
// 			return
// 		}
// 	}
// 	log.Printf("Failed to transfer after %d retries\n", maxRetries)
// }
