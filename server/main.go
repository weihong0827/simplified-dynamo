package main

import (
	"context"
	"dynamoSimplified/config"
	hash "dynamoSimplified/hash"
	pb "dynamoSimplified/pb"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// server is used to implement dynamo.KeyValueStoreServer.
// TODO: store data in memory first
type Server struct {
	pb.UnimplementedKeyValueStoreServer
	id             uint32
	addr           string
	mu             *sync.RWMutex // protects the following
	store          map[uint32]pb.KeyValue
	membershipList *pb.MembershipList
	vectorClocks   map[string]pb.VectorClock
}

type NodeConnection struct {
	Address *pb.Node
	Conn    *grpc.ClientConn
}

type GetResponse struct {
	Message string `json:"message"`
}

func NewServer(addr string) *Server {
	return &Server{
		id:    hash.GenHash(addr),
		addr:  addr,
		mu:    &sync.RWMutex{},
		store: make(map[uint32]pb.KeyValue),
		membershipList: &pb.MembershipList{Nodes: []*pb.Node{
			&pb.Node{
				Id:        hash.GenHash(addr),
				Address:   addr,
				Timestamp: timestamppb.Now(),
				IsAlive:   true,
			},
		}},
		vectorClocks: make(map[string]pb.VectorClock),
	}
}
func (s *Server) Forward(ctx context.Context, in *pb.WriteRequest) (*pb.WriteResponse, error) {
	s.mu.Lock()
	// defer s.mu.Unlock()
	nodeID := s.id
	key := in.KeyValue.Key
	targetNodes, _ := hash.GetNodesFromKey(hash.GenHash(key), s.membershipList.Nodes)
	for _, node := range targetNodes {
		if node.Id == nodeID {
			s.mu.Unlock()
			return s.Write(ctx, in)
		}
	}
	forwardNode, _ := s.getFastestRespondingServer(targetNodes)
	coordNode := pb.NewKeyValueStoreClient(forwardNode.Conn)
	s.mu.Unlock()
	return coordNode.Write(ctx, in)
}

func (s *Server) BulkWrite(ctx context.Context, in *pb.BulkWriteRequest) (*pb.Empty, error) {
	log.Printf("Receive Bulk write request at %d", s.id)

	for _, kv := range in.KeyValue {
		idx := hash.GenHash(kv.Key)
		s.store[idx] = *kv
	}
	return &pb.Empty{}, nil
}

func (s *Server) InitiateKeyRangeChange(
	newNode *pb.Node,
) {
	log.Printf("Initiating key range change")
	// Get n nodes
	// We are gettinging 1 node before the coordinatorNode
	// and n nodes after the coordinatorNode including the coordinatorNode
	var offsets []int
	offsets = append(offsets, -1)
	for i := 0; i < config.N; i++ {
		offsets = append(offsets, i)
	}

	nodes, err := hash.GetNodeFromKeyWithOffSet(offsets, s.id, s.membershipList.Nodes)
	if err != nil {
		log.Println("Error When assigning key range change:", err)
		return
	}
	//modify coordinatorNode
	log.Printf(
		"Transfering data from node %d to new node %d for key range %d to %d",
		s.id,
		newNode.Id,
		nodes[0].Id,
		nodes[config.N-1].Id,
	)
	Transfer(s.store, nodes[0].Id, nodes[config.N-1].Id, newNode)

	log.Printf(
		"Deleting data from node %d  for key range %d to %d",
		s.id,
		newNode.Id,
		nodes[config.N-1].Id,
	)

	s.Delete(context.Background(), &pb.ReplicaDeleteRequest{
		Start: newNode.Id,
		End:   nodes[config.N-1].Id,
	})
	log.Println("Coordinator node modification completed")

	// Delete from other nodes
	// Only delete if there are more than config.N + 1 nodes in the network

	if len(s.membershipList.Nodes) <= config.N+1 {
		log.Println("There are not enough nodes in the network to delete from")
		log.Println("Key range change completed")
		return
	}
	for i := 2; i <= config.N; i++ {
		node := nodes[i]
		startIdx := i - config.N + 1
		endIdx := startIdx + 1
		if i == config.N {
			log.Printf(
				"Deleting data from node %d  for key range %d to %d",
				node.Id,
				nodes[startIdx].Id,
				newNode.Id,
			)
			DeleteReplicaFromTarget(node, nodes[startIdx].Id, newNode.Id)
		} else {
			log.Printf(
				"Deleting data from node %d  for key range %d to %d",
				node.Id,
				nodes[startIdx].Id,
				nodes[endIdx].Id,
			)
			DeleteReplicaFromTarget(node, nodes[startIdx].Id, nodes[endIdx].Id)
		}
	}
	log.Println("Key range change completed")

}

// Write implements dynamo.KeyValueStoreServer
func (s *Server) Write(ctx context.Context, in *pb.WriteRequest) (*pb.WriteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Printf("write request received for %s", in.KeyValue.Key)
	// this should be the ID of the current node
	nodeID := s.id
	key := in.KeyValue.Key

	targetNodes, _ := hash.GetNodesFromKey(hash.GenHash(key), s.membershipList.Nodes)
	for _, node := range targetNodes { //check that you are indeed responsible for this key
		if node.Id == nodeID {
			if !in.IsReplica {
				var currentClock *pb.VectorClock
				// Update vector clock
				kv, found := s.store[hash.GenHash(key)]
				if !found {
					currentClock = &pb.VectorClock{Timestamps: make(map[uint32]*pb.ClockStruct)}
					currentClock.Timestamps[nodeID] = &pb.ClockStruct{
						ClokcVal:  1,
						Timestamp: timestamppb.Now(),
					}
				} else {
					currentClock = kv.VectorClock
					nodeClock, cfound := currentClock.Timestamps[nodeID]
					if cfound {
						nodeClock.ClokcVal += 1
						nodeClock.Timestamp = timestamppb.Now()
					} else {
						currentClock.Timestamps[nodeID] = &pb.ClockStruct{ClokcVal: 1, Timestamp: timestamppb.Now()}
					}
				}
				in.KeyValue.VectorClock = currentClock
				s.store[hash.GenHash(key)] = *in.KeyValue
				value, ok := s.store[hash.GenHash(key)]
				respChan := make(chan []*pb.KeyValue)
				go SendRequestToReplica(
					&value,
					s.membershipList.Nodes,
					config.WRITE,
					s.addr,
					ok,
					respChan,
				) //how to detect when write fails?
				replicaResult := <-respChan
				close(respChan)
				result := append(replicaResult, &value)
				log.Print("coordinator, required number of nodes hv written")
				return &pb.WriteResponse{KeyValue: result, Success: true}, nil
				//TODO: implement timeout when waited to long to get write success. or detect write failure
			}
			s.store[hash.GenHash(key)] = *in.KeyValue
			value, _ := s.store[hash.GenHash(key)]
			// for _, val := range s.store {
			// 	log.Print(s.addr, " node store: ", val.Key, val.Value, val.VectorClock)
			// }
			return &pb.WriteResponse{KeyValue: []*pb.KeyValue{&value}, Success: true}, nil
		}
	}
	return &pb.WriteResponse{Success: false, Message: "not responsible for this key"}, nil

}

func (s *Server) HintedHandoffWriteRequest(
	ctx context.Context,
	in *pb.HintedHandoffWriteRequest,
) (*pb.WriteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// TODO
	return &pb.WriteResponse{Success: true}, nil
}

func (s *Server) Delete(ctx context.Context, in *pb.ReplicaDeleteRequest) (*pb.Empty, error) {
	//TODO: Compare and update the membershiplist
	log.Printf("Delete request received at node %d from range %d to %d", s.id, in.Start, in.End)

	for key := range s.store {
		if key >= in.Start && key <= in.End {
			delete(s.store, key)
		}
	}
	return &pb.Empty{}, nil
}

// Read implements dynamo.KeyValueStoreServer
func (s *Server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Printf("node %v read request received for %v", s.id, in.Key)

	key := in.Key
	kv := pb.KeyValue{Key: key}

	value, ok := s.store[hash.GenHash(key)]
	if !in.IsReplica { //coordinator might not be responsible but try find anyways lmao
		respChan := make(chan []*pb.KeyValue)
		go SendRequestToReplica(&kv, s.membershipList.Nodes, config.READ, s.addr, ok, respChan)
		replicaResult := <-respChan
		close(respChan)
		if ok {
			replicaResult = append(replicaResult, &value) //contains the addresses of all stores
		}
		//compare vector clocks
		result := CompareVectorClocks(replicaResult)
		log.Printf("coordinator result of read %v", result)
		if len(result) == 0 {
			return &pb.ReadResponse{
				KeyValue: result,
				Success:  false,
				Message:  "Key Value store does not exist in the database",
			}, nil // TODO: raise error
		}
		return &pb.ReadResponse{KeyValue: result, Success: true}, nil
	}

	if !ok {
		return &pb.ReadResponse{Success: false, Message: "Key not found"}, nil
	}

	return &pb.ReadResponse{KeyValue: []*pb.KeyValue{&value}, Success: true}, nil
}

// Join implements dynamo.NodeServServer
func (s *Server) Join(ctx context.Context, in *pb.Node) (*pb.JoinResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("Join request received from %v", in.Address)

	// Update membership list
	s.membershipList.Nodes = append(s.membershipList.Nodes, in)

	// update key range
	s.InitiateKeyRangeChange(in)

	// Send membership list to joining node
	return &pb.JoinResponse{
		MembershipList: &pb.MembershipList{
			Nodes: s.membershipList.Nodes,
		},
		Data: nil,
	}, nil

}

// Gossip implements dynamo.KeyValueStoreServer
func (s *Server) Gossip(ctx context.Context, in *pb.GossipMessage) (*pb.GossipAck, error) {
	// Update nodes based on the received gossip message
	s.mu.Lock()
	s.membershipList = ReconcileMembershipList(s.membershipList, in.MembershipList)
	log.Println("Membership list:")
	for _, node := range s.membershipList.Nodes {
		if node.IsAlive {
			log.Printf("Node %v is alive", node.Address)
		} else {
			log.Printf("Node %v is dead", node.Address)
		}
	}
	s.mu.Unlock()

	return &pb.GossipAck{Success: true}, nil
}

func ReconcileMembershipList(
	list1 *pb.MembershipList,
	list2 *pb.MembershipList,
) *pb.MembershipList {
	var mp = make(map[uint32]*pb.Node)

	for _, node := range list1.Nodes {
		mp[node.Id] = node
	}

	for _, node := range list2.Nodes {
		// log.Printf("node %v", node)
		if _, ok := mp[node.Id]; !ok {
			// log.Printf("node %v", node)
			mp[node.Id] = node
		} else {
			// log.Printf("node in map %v", mp[node.Id])
			if mp[node.Id].Timestamp.Seconds < node.Timestamp.Seconds {
				mp[node.Id] = node
			}
		}
	}

	var newList = pb.MembershipList{Nodes: []*pb.Node{}}
	for _, node := range mp {
		newList.Nodes = append(newList.Nodes, node)
	}

	return &pb.MembershipList{Nodes: newList.Nodes}
}

// create a method to periodically send gossip message to other nodes
func (s *Server) SendGossip(ctx context.Context) {
	for {
		// randomly pick one other node from membership list
		// send gossip to that node
		s.mu.RLock()
		targetNode := s.membershipList.Nodes[rand.Intn(len(s.membershipList.Nodes))]
		s.mu.RUnlock()
		if targetNode.Address == s.addr {
			continue
		}

		// create grpc client
		conn, err := grpc.Dial(targetNode.Address, grpc.WithInsecure())
		if err != nil {
			log.Printf("fail to dial: %v", err)
			// update membership list to change isAlive to false
			s.mu.Lock()
			for _, node := range s.membershipList.Nodes {
				if node.Address == targetNode.Address {
					node.IsAlive = false
					node.Timestamp = timestamppb.Now()
					s.mu.Unlock()
					break
				}
			}
			continue
		}
		defer conn.Close()

		client := pb.NewKeyValueStoreClient(conn)

		// send gossip message
		s.mu.RLock()
		membershipList := s.membershipList
		s.mu.RUnlock()
		resp, err := client.Gossip(ctx, &pb.GossipMessage{MembershipList: membershipList})
		if err != nil {
			log.Printf("fail to send gossip to %v", targetNode.Address)
			// update membership list to change isAlive to false
			s.mu.Lock()
			for _, node := range s.membershipList.Nodes {
				if node.Address == targetNode.Address {
					node.IsAlive = false
					node.Timestamp = timestamppb.Now()
					s.mu.Unlock()
					break
				}
			}
			continue
		}

		if resp.Success {
			s.mu.Lock()
			for _, node := range s.membershipList.Nodes {
				if node.Address == targetNode.Address {
					node.IsAlive = true
					node.Timestamp = timestamppb.Now()
					s.mu.Unlock()
					break
				}
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func (s *Server) getFastestRespondingServer(servers []*pb.Node) (*NodeConnection, error) {
	// Create a channel to receive the first responding server
	ch := make(chan *NodeConnection, len(servers))

	// Ping all servers concurrently
	for _, server := range servers {
		go func(pNode *pb.Node) {
			conn, err1 := CreateGRPCConnection(pNode.Address)
			log.Print("Get fastest responding node: Pnode address", pNode.Address)
			if err1 != nil {
				s.updateMembershipList(true, pNode)
			}
			client := pb.NewKeyValueStoreClient(conn)
			connectedNode := &NodeConnection{Address: pNode, Conn: conn}
			// Call your gRPC ping method here (replace with your actual method)
			_, err2 := client.Ping(context.Background(), &pb.PingRequest{})
			if err2 == nil {
				ch <- connectedNode
			}
		}(server)
	}

	select {
	case server := <-ch:
		return server, nil
	case <-time.After(3 * time.Second):
		return nil, status.Error(codes.DeadlineExceeded, "No server responded in time")
	}
}

func (s *Server) updateMembershipList(true bool, targetNode *pb.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, node := range s.membershipList.Nodes {
		if node.Address == targetNode.Address {
			node.IsAlive = false
			node.Timestamp = timestamppb.Now()
			break
		}
	}
}

func (s *Server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{}, nil
}

func getSeedNodeAddr(webclient string) string {
	// call get rest api to webclient address
	// get the seed node address
	// Make a GET request to the API. {message: "addr"}
	resp, err := http.Get(webclient)
	if err != nil {
		log.Fatalf("Failed to make the request: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read the response body: %v", err)
	}

	log.Printf("Response body: %v", string(body))

	var data GetResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("Failed to unmarshal the response body: %v", err)
	}

	return data.Message
}

var (
	addr      = flag.String("addr", "127.0.0.1:50051", "the addr to serve on")
	webclient = flag.String("webclient", "", "the addr of the seed node")
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
	pb.RegisterKeyValueStoreServer(grpcServer, server)

	log.Printf("Server listening at %v", lis.Addr())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	addrToJoin := getSeedNodeAddr(*webclient)
	// join the seed node if not empty
	if addrToJoin != "" {
		// create grpc client
		log.Printf("Joining %v", addrToJoin)
		conn, err := grpc.Dial(addrToJoin, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()

		client := pb.NewKeyValueStoreClient(conn)

		// join the seed node
		resp, err := client.Join(
			context.Background(),
			&pb.Node{
				Id:        hash.GenHash(*addr),
				Address:   *addr,
				Timestamp: timestamppb.Now(),
				IsAlive:   true,
			},
		)
		if err != nil {
			log.Fatalf("%d failed to join %d at %v, retrying...", server.id, addrToJoin, addrToJoin)
		} else {
			log.Printf("%d joined successfully", server.id)
		}

		server.mu.Lock()
		server.membershipList = resp.MembershipList
		server.mu.Unlock()
	}

	log.Printf("Starting gossip...")

	// start gossiping
	go server.SendGossip(context.Background())
	for {
	}
}
