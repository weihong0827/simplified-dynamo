package node

import (
	"context"
	pb "dynamoSimplified/pb"
	utils "dynamoSimplified/utils"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"sync"
	"time"
)

type PTPNode struct {
	Id      uint32 // index
	Address string // localhost:port
}

type PTPServer struct {
	pb.UnimplementedPTPServer
	self           *PTPNode          // self node
	MembershipList []*PTPNode        // membership list
	mux            *sync.Mutex       // mutex for membership list
	storage        map[string]string // local storage
}

func NewPTPNode(addr string) *PTPNode {
	return &PTPNode{
		Id:      utils.GenHash(addr),
		Address: addr,
	}
}

func NewPTPServer(addr string) *PTPServer {
	self := &PTPNode{
		utils.GenHash(addr),
		addr,
	}
	membershipList := []*PTPNode{self}
	return &PTPServer{
		self:           self,
		MembershipList: membershipList,
		mux:            &sync.Mutex{},
		storage:        make(map[string]string),
	}
}

func (s *PTPServer) Serve(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			{
				func() {
					time.Sleep(3 * time.Second)
					log.Printf("Sending gossip...")
					s.SendGossip(ctx)
					time.Sleep(1 * time.Second)
				}()
			}
			{
				func() {
					// print my membership list
					s.mux.Lock()
					defer s.mux.Unlock()
					for _, node := range s.MembershipList {
						log.Printf("Node %d: %s", node.Id, node.Address)
					}
				}()
			}
			// TODO: gossip

		}
	}
}

func (s *PTPServer) AddNewNode(ctx context.Context, in *pb.Node) (*pb.MembershipList, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	// add new node to membership list
	s.MembershipList = append(s.MembershipList, &PTPNode{in.Id, in.Address})
	log.Printf("Node %s joined", in.Address)
	// print new membershiplist
	for _, node := range s.MembershipList {
		log.Printf("Node %d: %s", node.Id, node.Address)
	}
	return &pb.MembershipList{Nodes: convertToPBNode(s.MembershipList)}, nil
}

func (s *PTPServer) Gossip(ctx context.Context, in *pb.GossipMessage) (*pb.GossipAck, error) {
	s.mux.Lock()
	// Compare membership list
	s.MembershipList = convertPBNodetoNode(in.Nodes)
	s.mux.Unlock()
	// Update membership list
	return &pb.GossipAck{}, nil
}

func convertToPBNode(node []*PTPNode) []*pb.Node {
	pbNode := make([]*pb.Node, len(node))
	for i, n := range node {
		pbNode[i] = &pb.Node{Id: n.Id, Address: n.Address}
	}
	return pbNode
}

// func (s *PTPServer) Join(ctx context.Context, node *PTPNode) error {
// 	s.mux.Lock()
// 	s.membershipList = append(s.membershipList, node)
// 	s.mux.Unlock()
// 	return nil
// 	s.AddNewNode(ctx, node)
// }

// function to send gossip to other nodes
func (s *PTPServer) SendGossip(ctx context.Context) {
	// randomly pick one other node from membership list
	// send gossip to that node

	log.Printf("Inside gossip...")
	s.mux.Lock()
	// pick random node
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(s.MembershipList))
	randNode := s.MembershipList[randIndex]
	// create grpc client
	log.Printf("Gossiping to %d at %v", randNode.Id, randNode.Address)
	conn, err := grpc.Dial(randNode.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPTPClient(conn)
	// send gossip
	memList := convertToPBNode(s.MembershipList)
	s.mux.Unlock()
	resp, err := client.Gossip(ctx, &pb.GossipMessage{Nodes: memList})
	if err != nil {
		log.Fatalf("%d failed to gossip to %d at %v, retrying...", s.self.Id, randNode.Id, randNode.Address)
	} else {
		log.Printf("%d gossiped successfully", s.self.Id)
	}
	// log response
	log.Printf("Gossip response: %v", resp)
	// update membership list
}
