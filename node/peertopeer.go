package node

import (
	"context"
	pb "dynamoSimplified/pb"
	utils "dynamoSimplified/utils"
	"log"
	"sync"
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

		}
	}
}

// node 2 join node 1
func (s *PTPServer) AddNewNode(ctx context.Context, in *pb.Node) (*pb.MembershipList, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.MembershipList = append(s.MembershipList, &PTPNode{in.Id, in.Address})
	log.Printf("Node %d joined", in.Id)
	// print new membershiplist
	for _, node := range s.MembershipList {
		log.Printf("Node %d: %s", node.Id, node.Address)
	}
	return &pb.MembershipList{Nodes: convertToPBNode(s.MembershipList)}, nil
}

func (s *PTPServer) Gossip(ctx context.Context, in *pb.GossipMessage) (*pb.GossipAck, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	// Compare membership list

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

// node 0 create server

// node 1 join node 0
// node1.membershiplist -> node1, node0 t=x

// node 2 join node 0
// node2.membershiplist -> node2, node0 t=x+10

// gossip node1 -> node0
// node0.membershiplist -> node1, node0 t=x+20

// gossip node2 -> node0
// node0.membershiplist -> node2, node1, node0

// node2 mati, node0 tau
// node0.membershiplist -> node1, node0

// node0 gossip -> node1
