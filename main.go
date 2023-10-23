package main

import (
	"context"
	pb "dynamoSimplified/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

// server is used to implement dynamo.KeyValueStoreServer.
type server struct {
	pb.UnimplementedKeyValueStoreServer
	mu           sync.RWMutex // protects the following
	store        map[string]pb.KeyValue
	nodes        []pb.Node // known nodes in the ring
	vectorClocks map[string]pb.VectorClock
}

// Write implements dynamo.KeyValueStoreServer
func (s *server) Write(ctx context.Context, in *pb.WriteRequest) (*pb.WriteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := in.KeyValue.Key

	// Update vector clock
	currentClock, found := s.vectorClocks[key]
	if !found {
		currentClock = pb.VectorClock{Timestamps: make(map[string]int64)}
	}
	nodeID := "nodeID" // this should be the ID of the current node
	currentClock.Timestamps[nodeID] = time.Now().
		UnixNano()
		// Use appropriate time for your use case
	s.vectorClocks[key] = currentClock

	// Store the new value
	in.KeyValue.VectorClock = &currentClock
	s.store[key] = *in.KeyValue

	// Replicate write to W-1 other nodes (assuming the first write is the current node)
	for i, node := range s.nodes {
		if i >= W-1 {
			break
		}
		// Make a gRPC call to Write method of the other node
		// ...
	}

	return &pb.WriteResponse{Success: true}, nil
}

// Read implements dynamo.KeyValueStoreServer
func (s *server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := in.Key
	value, ok := s.store[key]
	if !ok {
		return &pb.ReadResponse{Success: false, Message: "Key not found"}, nil
	}

	// If the current node is not the coordinator, forward the read request to the coordinator
	// ...

	// Otherwise, read from R nodes
	// ...

	return &pb.ReadResponse{KeyValue: &value, Success: true}, nil
}

// Gossip implements dynamo.KeyValueStoreServer
func (s *server) Gossip(ctx context.Context, in *pb.GossipMessage) (*pb.GossipAck, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update nodes based on the received gossip message
	// ...

	return &pb.GossipAck{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterKeyValueStoreServer(grpcServer, &server{
		store:        make(map[string]pb.KeyValue),
		vectorClocks: make(map[string]pb.VectorClock),
		// initialize other necessary fields
	})
	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
