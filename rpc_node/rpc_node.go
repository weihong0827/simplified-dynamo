package main

import (
	"context"
	"dynamoSimplified/node"
	"sync"
)

const (
	address  = "127.0.0.1:50051"
	address2 = "127.0.0.1:50052"
	address3 = "127.0.0.1:50053"
)

func main() {
	group := &sync.WaitGroup{}
	group.Add(1)
	node0 := node.NewPTPNode(address)
	node1 := node.NewPTPNode(address2)
	node2 := node.NewPTPNode(address3)
	go node.ServePTP(context.Background(), node0, nil, group, nil)
	go node.ServePTP(context.Background(), node1, node0, group, nil)
	go node.ServePTP(context.Background(), node2, node0, group, nil)
	group.Wait()
}
