package main

import (
	"context"
	"dynamoSimplified/config"
	hash "dynamoSimplified/hash"
	pb "dynamoSimplified/pb"
	"flag"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// server is used to implement dynamo.KeyValueStoreServer.
// TODO: store data in memory first
type Server struct {
	pb.UnimplementedNodeServServer
	id                   uint32
	addr                 string
	mu                   *sync.RWMutex // protects the following
	store                map[string]pb.KeyValue
	membershipList       *pb.MembershipList
	hintedMembershipList *pb.HintedHandoffMembershipList
	vectorClocks         map[string]pb.VectorClock
}

func NewServer(addr string) *Server {
	return &Server{
		id:    hash.GenHash(addr),
		addr:  addr,
		mu:    &sync.RWMutex{},
		store: make(map[string]pb.KeyValue),
		membershipList: &pb.MembershipList{Nodes: []*pb.Node{
			&pb.Node{Id: hash.GenHash(addr), Address: addr},
		}, Timestamp: timestamppb.Now()},
		hintedMembershipList: &pb.HintedHandoffMembershipList{
			Hinted: []*pb.HintedHandoffWriteRequest{
				&pb.KeyValue{Key: hash.GenHash(addr), Value: nil},
				&pb.Node{Id: hash.GenHash(addr), Address: addr},
			},
		},
		vectorClocks: make(map[string]pb.VectorClock),
	}
}

// Write implements dynamo.KeyValueStoreServer
func (s *Server) Write(ctx context.Context, in *pb.WriteRequest) (*pb.WriteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// this should be the ID of the current node
	nodeID := s.id
	key := in.KeyValue.Key

	// Use appropriate time for your use case
	// s.vectorClocks[key] = currentClock

	// Store the new value

	// Replicate write to W-1 other nodes (assuming the first write is the current node)

	// Make a gRPC call to Write method of the other node
	// ...
	if !in.IsReplica {
		var currentClock *pb.VectorClock
		// Update vector clock
		kv, found := s.store[key]
		if !found {
			currentClock = &pb.VectorClock{Timestamps: make(map[uint32]*pb.ClockStruct)}
			currentClock.Timestamps[nodeID] = &pb.ClockStruct{ClokcVal: 1, Timestamp: timestamppb.Now()}
		} else {
			currentClock = kv.VectorClock
			currentClock.Timestamps[nodeID].ClokcVal += 1
			currentClock.Timestamps[nodeID].Timestamp = timestamppb.Now()
		}
		in.KeyValue.VectorClock = currentClock
		s.store[key] = *in.KeyValue
		value, _ := s.store[key]
		replicaResult := SendRequestToReplica(kv, s.membershipList.Nodes, config.WRITE, s.addr) //replica result is an array
		result := append(replicaResult, &value)
		return &pb.WriteResponse{KeyValue: result, Success: true}, nil
		//TODO: implement timeout when waited to long to get write success. or detect write failure
	} else {
		value, _ := s.store[key]
		return &pb.WriteResponse{KeyValue: []*pb.KeyValue{&value}, Success: true}, nil
	}

	return &pb.WriteResponse{Success: true}, nil
}

func contains(slice []*pb.Node, node *pb.Node) bool {
	for _, e := range slice {
		if e == node {
			return true
		}
	}
	return false
}

func (s *Server) checkHintedMembershipList(ctx context.Context) {
	// sepearate function periodically check this struct has anything
	// if yes send back to the actual node that it is meant for
	// on every server receive

	// periodically

	ticker := time.NewTicker(30 * time.Second) // Adjust the interval as needed

	for {
		select {
		case <-ticker.C:
			if s.hintedMembershipList != nil {
				result := SendToNode(ctx, s.hintedMembershipList.Hinted)

				// update list
				var updatedList *pb.HintedHandoffMembershipList
				for _, n := range s.hintedMembershipList {
					if !contains(result, n.Node) {
						updatedList = append(updatedList, n.Node)
					}
				}
				s.hintedMembershipList = updatedList
			}
		}
	}
}

func (s *Server) HintedHandoffWrite(ctx context.Context, in *pb.HintedHandoffWriteRequest) (*pb.HintedHandoffWriteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// // store this hinted information in a seperated struct
	// // whcih contains the actual kv and the node

	inactiveNode := in.Node

	for _, node := range s.membershipList.Nodes {
		if node == inactiveNode {
			node.active = false
			break
		}
	}

	s.hintedMembershipList = append(s.hintedMembershipList, in)

	return &pb.HintedHandoffWriteResponse{KeyValue: []*pb.KeyValue{in.KeyValue}, Node: in.Node, Success: true}, nil
	// return &pb.WriteResponse{KeyValue: []*pb.KeyValue{&value}, Success: true}, nil
}

// Read implements dynamo.KeyValueStoreServer
func (s *Server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Printf("read request received for %v", in.Key)

	key := in.Key
	kv := pb.KeyValue{Key: key}

	value, ok := s.store[key]
	if !ok {
		return &pb.ReadResponse{Success: false, Message: "Key not found"}, nil
	}
	if !in.IsReplica {
		replicaResult := SendRequestToReplica(kv, s.membershipList.Nodes, config.READ, s.addr)
		result := append(replicaResult, &value) //contains the addresses of all stores
		//compare vector clocks
		result = CompareVectorClocks(result)
		return &pb.ReadResponse{KeyValue: result, Success: true}, nil
	}

	return &pb.ReadResponse{KeyValue: []*pb.KeyValue{&value}, Success: true}, nil
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
	log.Printf("membershipList at node %v\n", s.addr)
	for _, node := range s.membershipList.Nodes {
		log.Println(node.Address)
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

	log.Printf("Server listening at %v", lis.Addr())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go server.checkHintedMembershipList(context.Background())

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
		resp, err := client.Join(context.Background(), &pb.Node{Id: hash.GenHash(*addr), Address: *addr})
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
