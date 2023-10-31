package node

import (
	"context"
	pb "dynamoSimplified/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	// "time"
)

func ServePTP(ctx context.Context, node *PTPNode, seed_node *PTPNode, group *sync.WaitGroup, join *sync.WaitGroup) {
	defer group.Done()

	server := NewPTPServer(node.Address)
	lis, err := net.Listen("tcp", node.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening at %v", node.Address)
	// setup gRPC
	s := grpc.NewServer()
	pb.RegisterPTPServer(s, server)
	// pb.RegisterKeyValueStoreServer(s, server)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ch := make(chan bool)

	// begin background tasks

	// handle RPCs
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		ch <- true
	}()

	// seed node
	if seed_node != nil {
		// create grpc client
		conn, err := grpc.Dial(seed_node.Address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()

		client := pb.NewPTPClient(conn)
		// join the seed node
		resp, err := client.AddNewNode(ctx, convertPTPNodetoPBNode(node))
		if err != nil {
			log.Fatalf("%d failed to join %d at %v, retrying...", node.Id, seed_node.Id, seed_node.Address)
		} else {
			log.Printf("%d joined successfully", node.Id)
		}
		server.mux.Lock()
		server.MembershipList = convertPBNodetoNode(resp.Nodes)
		server.mux.Unlock()
	}

	go server.Serve(ctx)

	// node has joined
	if join != nil {
		join.Done()
	}

	// stop listening if context is done
	select {
	case <-ch:
		return
	case <-ctx.Done():
		s.Stop()
		return
	}
}

func convertPTPNodetoPBNode(node *PTPNode) *pb.Node {
	return &pb.Node{
		Id:      node.Id,
		Address: node.Address,
	}
}

func convertPBNodetoNode(node []*pb.Node) []*PTPNode {
	pbNode := make([]*PTPNode, len(node))
	for i, n := range node {
		pbNode[i] = &PTPNode{n.Id, n.Address}
	}
	return pbNode
}
