package utils

import (
	"crypto/md5"
	pb "dynamoSimplified/pb"
	"encoding/binary"
	"errors"
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

func GetAddressFromNode(number uint32, nodes NodeSlice) (*pb.Node, error) {
	//BUG: Cannot pass the test check logic for search
	if len(nodes) == 0 {
		return nil, errors.New("no nodes available")
	}
	sort.Sort(nodes)

	// Binary search: find the range containing the number.
	index := sort.Search(len(nodes), func(i int) bool { return nodes[i].Id > number }) - 1
	if index == 0 {
		return nodes[len(nodes)-1], nil
	} else {
		return nodes[index], nil
	}
}

func GenHash(key string) uint32 {
	h := md5.New()
	h.Write([]byte(key))
	hashBytes := h.Sum(nil)
	return binary.BigEndian.Uint32(hashBytes[:4])
}
