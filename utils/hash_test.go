package utils

import (
	pb "dynamoSimplified/pb"
	"errors"
	"reflect"
	"testing"
)

// TODO: update tests
// Test_GetAddressFromNode_EmptyNodeSlice checks the function's response with an empty slice of nodes.
func Test_GetNodesFromKey_EmptyNodeSlice(t *testing.T) {
	nodes := NodeSlice{}
	var key uint32 = 1

	node, err := GetNodesFromKey(key, nodes)

	if node != nil || !errors.Is(err, ErrNoNodesAvailable) {
		t.Errorf("Expected no nodes available error, got %v and node %v", err, node)
	}
}

// Test_GetAddressFromNode_ValidSingleKey checks behavior with a single node in the slice.
func Test_GetAddressFromNode_ValidSingleKey(t *testing.T) {
	node := &pb.Node{Start: 1}
	nodes := NodeSlice{node}
	var key uint32 = 3

	result, err := GetAddressFromNode(key, nodes)

	if err != nil || !reflect.DeepEqual(result, node) {
		t.Errorf("Expected node %v, got node %v with error %v", node, result, err)
	}
}

// Test_GetAddressFromNode_ValidMultipleKeys tests function with multiple nodes and multiple keys.
func Test_GetAddressFromNode_ValidMultipleKeys(t *testing.T) {
	nodes := NodeSlice{
		&pb.Node{Start: 1},
		&pb.Node{Start: 3},
		&pb.Node{Start: 5},
	}

	keys := []uint32{2, 4, 6}
	expectedNodes := []*pb.Node{nodes[0], nodes[1], nodes[2]}

	for i, key := range keys {
		result, err := GetAddressFromNode(key, nodes)
		if err != nil || !reflect.DeepEqual(result, expectedNodes[i]) {
			t.Errorf(
				"For key %d, expected node %v, got %v with error %v",
				key,
				expectedNodes[i],
				result,
				err,
			)
		}
	}
}

// Test_GetAddressFromNode_OutOfRangeKey tests how function handles a key that results in a hash out of nodes' range.
func Test_GetAddressFromNode_OutOfRangeKey(t *testing.T) {
	nodes := NodeSlice{
		&pb.Node{Start: 1},
		&pb.Node{Start: 100},
	}

	var key uint32 = 200

	// Assuming that the 'GenHash' function generates a hash out of the range of node starts
	result, err := GetAddressFromNode(key, nodes)

	if err != nil || result != nodes[len(nodes)-1] {
		t.Errorf("Expected last node in slice, got node %v with error %v", result, err)
	}
}

// Test_GetAddressFromNode_UnsortedNodes tests function with unsorted nodes.
func Test_GetAddressFromNode_UnsortedNodes(t *testing.T) {
	nodes := NodeSlice{
		&pb.Node{Start: 5},
		&pb.Node{Start: 3},
		&pb.Node{Start: 1},
	}

	var key uint32 = 2

	// The function is expected to sort nodes, so it should not return an error or incorrect node
	result, err := GetAddressFromNode(key, nodes)
	expectedNode := nodes[0]

	if err != nil || !reflect.DeepEqual(result, expectedNode) {
		t.Errorf("Expected correct node based on sorting, got node %v with error %v", result, err)
	}
}
