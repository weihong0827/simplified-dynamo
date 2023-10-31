package utils

import (
	"crypto/md5"
	"dynamoSimplified/config"
	pb "dynamoSimplified/pb"
	"encoding/binary"
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

func GetNodesFromKey(key uint32, nodes NodeSlice, op config.Operation) ([]*pb.Node, error) {
	if len(nodes) == 0 {
		return nil, ErrNoNodesAvailable
	}

	var n int
	var result []*pb.Node

	switch op {
	case config.READ:
		n = config.R
	case config.WRITE:
		n = config.W
	}

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

func GenHash(key string) uint32 {
	h := md5.New()
	h.Write([]byte(key))
	hashBytes := h.Sum(nil)
	return binary.BigEndian.Uint32(hashBytes[:4])
}
