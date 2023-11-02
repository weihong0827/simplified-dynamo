package main

import (
	"context"
	pb "dynamoSimplified/pb"
	"errors"
	"log"
)

func Transfer(
	store map[uint32]pb.KeyValue,
	start uint32,
	end uint32,
	targetNode *pb.Node,
) (*pb.Empty, error) {
	dataToTransfer := []*pb.KeyValue{}
	for key, value := range store {
		if key >= start && key < end {
			dataToTransfer = append(dataToTransfer, &value)
		}
	}
	err := BulkWriteToTarget(dataToTransfer, targetNode)
	if err != nil {
		log.Println("Error when transferring data:", err)
		return &pb.Empty{}, errors.New("Error when transferring data")
	}
	return &pb.Empty{}, nil
}

func BulkWriteToTarget(kvToTransfer []*pb.KeyValue, targetNode *pb.Node) error {
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
