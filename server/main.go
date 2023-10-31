package main

import (
	"context"
	"dynamoSimplified/config"
	pb "dynamoSimplified/pb"
	utils "dynamoSimplified/utils"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

// server is used to implement dynamo.KeyValueStoreServer.
// TODO: store data in memory first
type Server struct {
	pb.UnimplementedNodeServServer
	id             uint32
	addr           string
	mu             *sync.RWMutex // protects the following
	store          map[string]pb.KeyValue
	membershipList *pb.MembershipList
	vectorClocks   map[string]pb.VectorClock
}

func NewServer(addr string) *Server {
	return &Server{
		id:    utils.GenHash(addr),
		addr:  addr,
		mu:    &sync.RWMutex{},
		store: make(map[string]pb.KeyValue),
		membershipList: &pb.MembershipList{Nodes: []*pb.Node{
			&pb.Node{Id: utils.GenHash(addr), Address: addr},
		}, Timestamp: timestamppb.Now()},
		vectorClocks: make(map[string]pb.VectorClock),
	}
}

// Write implements dynamo.KeyValueStoreServer
func (s *Server) Write(ctx context.Context, in *pb.WriteRequest) (*pb.WriteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := in.KeyValue.Key

	// Update vector clock
	currentClock, found := s.vectorClocks[key]
	if !found {
		currentClock = pb.VectorClock{Timestamps: make(map[string]*pb.ClockStruct)}
	}
	nodeID := "nodeID" // this should be the ID of the current node
	currentClock.Timestamps[nodeID] = time.Now().UnixNano()
	// Use appropriate time for your use case
	s.vectorClocks[key] = currentClock

	// Store the new value
	in.KeyValue.VectorClock = &currentClock
	s.store[key] = *in.KeyValue

	// Replicate write to W-1 other nodes (assuming the first write is the current node)

	// Make a gRPC call to Write method of the other node
	// ...

	return &pb.WriteResponse{Success: true}, nil
}

// Read implements dynamo.KeyValueStoreServer
func (s *Server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Printf("read request received for %v", in.Key)

	key := in.Key
	value, ok := s.store[key]
	if !ok {
		return &pb.ReadResponse{Success: false, Message: "Key not found"}, nil
	}
	if in.IsReplica {
		replicaResult := SendRequestToReplica(key, s.membershipList.Nodes, config.READ, s.addr)
		result := append(replicaResult, &value)
		return &pb.ReadResponse{KeyValue: result, Success: true}, nil
	}

	return &pb.ReadResponse{KeyValue: *pb.KeyValue{&value}, Success: true}, nil
}

// Gossip implements dynamo.KeyValueStoreServer
func (s *Server) Gossip(ctx context.Context, in *pb.GossipMessage) (*pb.GossipAck, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update nodes based on the received gossip message
	// ...

	return &pb.GossipAck{Success: true}, nil
}

// Join implements dynamo.NodeServServer
func (s *Server) Join(ctx context.Context, in *pb.Node) (*pb.MembershipList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("Join request received from %v", in.Address)

	// Update membership list
	s.membershipList.Nodes = append(s.membershipList.Nodes, in)
	s.membershipList.Timestamp = timestamppb.Now()

	// Send membership list to joining node
	return &pb.MembershipList{Nodes: s.membershipList.Nodes, Timestamp: s.membershipList.Timestamp}, nil
}

// Gossip implements dynamo.NodeServServer
func (s *Server) Gossip(ctx context.Context, in *pb.GossipMessage) (*pb.GossipAck, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if in.MembershipList.Timestamp.Seconds < s.membershipList.Timestamp.Seconds {
		return &pb.GossipAck{Success: false}, nil
	}

	// Update nodes based on the received gossip message
	s.membershipList = in.MembershipList

	// log membershipList
	fmt.Printf("membershipList at node %v\n", s.addr)
	for _, node := range s.membershipList.Nodes {
		fmt.Println(node.Address)
	}

	return &pb.GossipAck{Success: true}, nil
}

// create a method to periodically send gossip message to other nodes
func (s *Server) SendGossip(ctx context.Context) {
	for {
		s.mu.RLock()

		// randomly pick one other node from membership list
		// send gossip to that node
		targetNode := s.membershipList.Nodes[rand.Intn(len(s.membershipList.Nodes))]
		if targetNode.Address == s.addr {
			s.mu.RUnlock()
			time.Sleep(time.Second * 5)
			continue
		}

		// create grpc client
		conn, err := grpc.Dial(targetNode.Address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()

		client := pb.NewNodeServClient(conn)

		// send gossip message
		client.Gossip(ctx, &pb.GossipMessage{MembershipList: s.membershipList})

		s.mu.RUnlock()

		time.Sleep(time.Second * 5)
	}
}

var (
	addr      = flag.String("addr", "127.0.0.1:50051", "the addr to serve on")
	seed_addr = flag.String("seed_addr", "", "the addr of the seed node")
	sleep     = flag.Duration("sleep", time.Second*5, "duration between changes in health")

	system = "" // empty string represents the health of the system
)

func main() {
	flag.Parse()

	// create a server
	server := NewServer(*addr)

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Register the server with the gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterNodeServServer(grpcServer, server)

	// Register the health check server with the gRPC server
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(grpcServer, healthcheck)

	log.Printf("Server listening at %v", lis.Addr())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// join the seed node if not empty
	if *seed_addr != "" {

		// create grpc client
		conn, err := grpc.Dial(*seed_addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()

		client := pb.NewNodeServClient(conn)

		// join the seed node
		resp, err := client.Join(context.Background(), &pb.Node{Id: utils.GenHash(*addr), Address: *addr})
		if err != nil {
			log.Fatalf("%d failed to join %d at %v, retrying...", server.id, seed_addr, seed_addr)
		} else {
			log.Printf("%d joined successfully", server.id)
		}

		server.mu.Lock()
		server.membershipList = resp
		server.mu.Unlock()
	}

	log.Printf("Starting gossip...")

	// start gossiping
	go server.SendGossip(context.Background())
	for {
	}
}
