package main

import (
	"context"
	"dynamoSimplified/config"
	hash "dynamoSimplified/hash"
	pb "dynamoSimplified/pb"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// server is used to implement dynamo.KeyValueStoreServer.
// TODO: store data in memory first
type server struct {
	pb.UnimplementedKeyValueStoreServer
	mu             sync.RWMutex // protects the following
	port           uint32
	store          map[string]pb.KeyValue
	membershipList pb.MembershipList
	// vectorClocks   map[string]pb.VectorClock
}

func NewServer() *server {
	return &server{
		store: make(map[string]pb.KeyValue),
		// vectorClocks: make(map[string]pb.VectorClock),
	}
}

// Write implements dynamo.KeyValueStoreServer
func (s *server) Write(ctx context.Context, in *pb.WriteRequest) (*pb.WriteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// this should be the ID of the current node
	nodeID := hash.GenHash(strconv.FormatUint(uint64(s.port), 10))
	key := in.KeyValue.Key

	var currentClock *pb.VectorClock
	// Update vector clock
	kv, found := s.store[key]
	if !found {
		currentClock = &pb.VectorClock{Timestamps: make(map[uint32]*pb.ClockStruct)}
		currentClock.Timestamps[nodeID] = &pb.ClockStruct{ClokcVal: 1, Timestamp: timestamppb.Now()}
		in.KeyValue.VectorClock = currentClock
	} else {
		currentClock = kv.VectorClock
		currentClock.Timestamps[nodeID].ClokcVal += 1
		currentClock.Timestamps[nodeID].Timestamp = timestamppb.Now()
		in.KeyValue.VectorClock = currentClock
	}

	// Use appropriate time for your use case
	// s.vectorClocks[key] = currentClock
	in.KeyValue.VectorClock = currentClock
	// Store the new value

	s.store[key] = *in.KeyValue
	value, _ := s.store[key]

	// Replicate write to W-1 other nodes (assuming the first write is the current node)

	// Make a gRPC call to Write method of the other node
	// ...
	if !in.IsReplica {
		replicaResult := SendRequestToReplica(kv, s.membershipList.Nodes, config.WRITE, s.port) //replica result is an array
		result := append(replicaResult, &value)
		return &pb.WriteResponse{KeyValue: result, Success: true}, nil
		//TODO: implement timeout when waited to long to get write success. or detect write failure
	}

	return &pb.WriteResponse{Success: true}, nil
}

// Read implements dynamo.KeyValueStoreServer
func (s *server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := in.Key
	kv := pb.KeyValue{Key: key}

	value, ok := s.store[key]
	if !ok {
		return &pb.ReadResponse{Success: false, Message: "Key not found"}, nil
	}
	if !in.IsReplica {
		replicaResult := SendRequestToReplica(kv, s.membershipList.Nodes, config.READ, s.port)
		result := append(replicaResult, &value) //contains the addresses of all stores
		//compare vector clocks
		result = CompareVectorClocks(result)
		return &pb.ReadResponse{KeyValue: result, Success: true}, nil
	}

	return &pb.ReadResponse{KeyValue: []*pb.KeyValue{&value}, Success: true}, nil
}

// Gossip implements dynamo.KeyValueStoreServer
func (s *server) Gossip(ctx context.Context, in *pb.GossipMessage) (*pb.GossipAck, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update nodes based on the received gossip message
	// ... //TODO: implement gossip response protocol
	return &pb.GossipAck{Success: true}, nil
}

var (
	port  = flag.Int("port", 50051, "the port to serve on")
	sleep = flag.Duration("sleep", time.Second*5, "duration between changes in health")

	system = "" // empty string represents the health of the system
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register the server with the gRPC server
	pb.RegisterKeyValueStoreServer(grpcServer, NewServer())

	// Register the health check server with the gRPC server
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(grpcServer, healthcheck)

	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
