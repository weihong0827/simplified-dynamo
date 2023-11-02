package main

import (
	"fmt"
	"google.golang.org/grpc"
)

func CreateGRPCConnection(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf(connectionError, err)
	}
	return conn, nil
}
