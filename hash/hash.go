package hash

import (
	"crypto/md5"
	"dynamoSimplified/config"
	pb "dynamoSimplified/pb"
	"encoding/binary"
	"log"
	"sort"
)

type NodeSlice []*pb.Node

func (a NodeSlice) Len() int { return len(a) }

func (a NodeSlice) Less(
	i, j int,
) bool {
	return a[i].Id < a[j].Id
}                                 // Assuming 'Start' is a field in 'pb.Node'
func (a NodeSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func GetNodesFromKey(key uint32, nodes NodeSlice) ([]*pb.Node, error) {
	if len(nodes) == 0 {
		return nil, ErrNoNodesAvailable
	}

	var offsets []int

	for i := 0; i < config.N; i++ {
		offsets = append(offsets, i)
	}
	return GetNodeFromKeyWithOffSet(offsets, key, nodes)
}

func GetNodeFromKeyWithOffSet(
	offsets []int,
	key uint32,
	nodes NodeSlice,
) ([]*pb.Node, error) {
	if len(nodes) == 0 {
		return nil, ErrNoNodesAvailable
	}

	var result []*pb.Node

	sort.Sort(nodes)

	// Binary search: find the range containing the number.
	index := sort.Search(len(nodes), func(i int) bool { return nodes[i].Id > key }) - 1
	if index == -1 {
		index = len(nodes) - 1
	}
	for offset, _ := range offsets {
		log.Printf("offset: %d", offset)
		indexToAdd := (index + offset) % len(nodes)
		log.Printf("indexToAdd: %d", indexToAdd)
		result = append(result, nodes[indexToAdd])
	}
	log.Printf("result: %v", result)
	return result, nil

}

func GetResponsibleNode(key uint32, nodes NodeSlice) (string, error) {
	if len(nodes) == 0 {
		return "", nil
	}

	sort.Sort(nodes)
	index := sort.Search(len(nodes), func(i int) bool { return nodes[i].Id > key }) - 1

	if index == -1 {
		index = len(nodes) - 1
	}

	return nodes[index].Address, nil
}

//TODO: write function to get successive k nodes for hinted handoff

func GenHash(key string) uint32 {
	h := md5.New()
	h.Write([]byte(key))
	hashBytes := h.Sum(nil)
	return binary.BigEndian.Uint32(hashBytes[:4])
}
