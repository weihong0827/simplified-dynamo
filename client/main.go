package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	pb "grpcSimplified/pb"
)

var serviceConfig = `{
	"loadBalancingPolicy": "round_robin",
	"healthCheckConfig": {
		"serviceName": ""
	}
}`

func main() {
	flag.Parse()

	r := manual.NewBuilderWithScheme("whatever")
	r.InitialState(resolver.State{
		Addresses: []resolver.Address{
			{Addr: "localhost:50051"},
			{Addr: "localhost:50052"},
		},
	})
	// TODO: Custom loadBalancingPolicy to implement consistent hashing

	address := fmt.Sprintf("%s:///unused", r.Scheme())

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithResolvers(r),
		grpc.WithDefaultServiceConfig(serviceConfig),
	}

	conn, err := grpc.Dial(address, options...)
	if err != nil {
		log.Fatalf("grpc.Dial(%q): %v", address, err)
	}
	defer conn.Close()
	client := pb.NewKeyValueStoreClient(conn)
	client.Write(context.Background(), &pb.WriteRequest{Key: "foo", Value: "bar"})
	client.Read(context.Background(), &pb.ReadRequest{Key: "foo"})

}
