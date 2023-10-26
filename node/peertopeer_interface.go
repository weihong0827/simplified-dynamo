package node

import (
	"context"
	pb "dynamoSimplified/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

func ServePTP(ctx context.Context, node *PTPNode, bootstrap_node *PTPNode, group *sync.WaitGroup, join *sync.WaitGroup) {
	// setup logger
	defer group.Done()
	// setup Chord instances
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
	go server.Serve(ctx)

	// handle RPCs
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		ch <- true
	}()

	// bootstrap node
	if bootstrap_node != nil {
		for {
			memlist, err := server.AddNewNode(ctx, convertPTPNodetoPBNode(bootstrap_node))
			server.MembershipList = convertPBNodetoNode(memlist.Nodes)
			if err != nil {
				log.Fatalf("%d failed to join %d at %v, retrying...", node.Id, bootstrap_node.Id, bootstrap_node.Address)
			} else {
				log.Printf("%d joined successfully", node.Id)
				break
			}
			time.Sleep(time.Second)
		}
	}

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
