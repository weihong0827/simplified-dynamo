package hash

import (
	"crypto/md5"
	"dynamoSimplified/config"
	"dynamoSimplified/hash"
	pb "dynamoSimplified/pb"
	"encoding/binary"
	"sort"
)

type NodeSlice []*pb.Node

type NodeKvSlice []*pb.HintedHandoffWriteRequest

func (a NodeSlice) Len() int { return len(a) }

func (a NodeSlice) Less(
	i, j int,
) bool {
	return a[i].Id < a[j].Id
}                                 // Assuming 'Start' is a field in 'pb.Node'
func (a NodeSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func GetNodesFromKey(key uint32, nodes NodeSlice, op config.Operation) ([]*pb.Node, error) {
	if len(nodes) == 0 {
		return nil, ErrNoNodesAvailable
	}

	var n int
	var result []*pb.Node
	n = config.N

	sort.Sort(nodes)

	// Binary search: find the range containing the number.
	index := sort.Search(len(nodes), func(i int) bool { return nodes[i].Id > key }) - 1

	for i := 0; i < n; i++ {
		indexToAdd := index + i
		if indexToAdd == len(nodes) {
			indexToAdd = 0
		}
		result = append(result, nodes[indexToAdd])

	}
	return result, nil
}

// TODO: write function to get successive k nodes for hinted handoff
func GetSuccessiveNode(
	key uint32,
	nodes hash.NodeSlice,
	op config.Operation,
) *pb.Node {

	var successor *pb.Node
	for _, node := range nodes {
		// successor after n
		// check if active
		// else go next

		n := config.N

		sort.Sort(nodes)

		index := sort.Search(len(nodes), func(i int) bool { return nodes[i].Id > (key + uint32(n)) }) - 1

		for i := 0; i < len(nodes); i++ {
			indexToUse := index + i
			if indexToUse == len(nodes) {
				// indexToUse = 0
			}
			if nodes[indexToUse].active == true {
				successor = nodes[indexToUse]
				break
			}
		}
	}
	return successor
}

func GenHash(key string) uint32 {
	h := md5.New()
	h.Write([]byte(key))
	hashBytes := h.Sum(nil)
	return binary.BigEndian.Uint32(hashBytes[:4])
}
