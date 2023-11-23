package main

import (
	"context"

	// "dynamoSimplified/config"
	"log"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	hash "dynamoSimplified/hash"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "dynamoSimplified/pb" // replace with your actual gRPC protos package
)

type KillNode struct {
	Address string
}

type ReviveNode struct {
	Address string
}

type Server struct {
	Address *pb.Node
	Conn    *grpc.ClientConn
}

type Response struct {
	Key             string            `JSON:"key"`
	ReplicaResponse []ReplicaResponse `JSON:"replica_response"`
	Hashvalue       uint32            `JSON:"hashvalue"`
}

type ReplicaResponse struct {
	Value       string           `JSON:"value"`
	VectorClock map[uint32]int64 `JSON:"vector_clock"`
}

var servers []Server
var mutex = &sync.Mutex{}

func main() {
	// Initialize the list of backend servers NEED TO RUN AT LEAST 4 NODES
	servers = []Server{
		// {&pb.Node{Id: hash.GenHash("127.0.0.1:50051"), Address: "127.0.0.1:50051"}, nil},
		// {&pb.Node{Id: hash.GenHash("127.0.0.1:50052"), Address: "127.0.0.1:50052"}, nil},
		// {&pb.Node{Id: hash.GenHash("127.0.0.1:50053"), Address: "127.0.0.1:50053"}, nil},
	}

	// Establish gRPC connections to all servers
	for i, server := range servers {
		conn, err := grpc.Dial(
			server.Address.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Printf("Failed to connect to server at %s: %v", server.Address, err)
		}
		servers[i].Conn = conn
	}

	router := gin.Default()

	router.GET("/get", func(c *gin.Context) {
		fastestServer, err := getFastestRespondingServer()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Establish a gRPC connection to the fastest server
		conn, err := grpc.Dial(
			fastestServer.Address.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		client := pb.NewKeyValueStoreClient(conn)

		// TODO: Call your gRPC method here (replace with your actual method), Also make protobuf message
		resp, err := client.Read(
			context.Background(),
			&pb.ReadRequest{Key: c.Query("key"), IsReplica: false},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("Failed to read key: %v with error %v", c.Query("key"), err)
			return
		}

		result := ConvertPbReadResponse(resp.KeyValue)

		log.Printf("%v", result)

		// Forward the response from the backend server to the client
		c.JSON(http.StatusOK, gin.H{"message": result})
	})

	router.GET("/addNode", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()

		address := c.Query("address")
		hashVal := hash.GenHash(address)
		node, nodeRespErr := hash.GetResponsibleNode(hashVal, getServersAddresses(servers))
		if nodeRespErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": nodeRespErr.Error()})
			return
		}
		conn, connErr := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if connErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": connErr.Error()})
			return
		}
		servers = append(servers, Server{&pb.Node{Id: hashVal, Address: address}, conn})

		c.JSON(http.StatusOK, gin.H{"message": node})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	router.PUT("/put", func(c *gin.Context) {
		fastestServer, err := getFastestRespondingServer()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Establish a gRPC connection to the fastest server
		log.Print(fastestServer.Address.Address)
		conn, err := grpc.Dial(
			fastestServer.Address.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		client := pb.NewKeyValueStoreClient(conn)
		resp, err := client.Forward(
			context.Background(),
			&pb.WriteRequest{
				KeyValue:  &pb.KeyValue{Key: c.Query("key"), Value: c.Query("value")},
				IsReplica: false,
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("Failed to write key: %v with error %v", c.Query("key"), err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": resp.Message})
	})

	// router.POST("/kill", func(c *gin.Context) {
	// 	var killNodeAddress KillNode
	// 	if err := c.ShouldBindJSON(&killNodeAddress); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	// connect to the node address that we want to kill
	// 	for _, server := range servers {
	// 		if server.Address.Address == killNodeAddress.Address {
	// 			// Create a gRPC connection to the server
	// 			client := pb.NewKeyValueStoreClient(server.Conn)
	// 			log.Print("Client created")
	// 			// TODO: Check how KillNode is implemented.
	// 			_, err := client.KillNode(context.Background())

	// 			if err != nil {
	// 				log.Printf("Failed to kill node: %v with error %v", killNodeAddress.Address, err)
	// 				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to kill node"})
	// 				return
	// 			}
	// 		}
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{"message": "Node successfully killed!"})

	// router.POST("/revive", func(c *gin.Context) {
	// 	var reviveNodeAddress ReviveNode
	// 	if err := c.ShouldBindJSON(&reviveNodeAddress); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	// connect to the node address that we want to kill
	// 	for _, server := range servers {
	// 		if server.Address.Address == reviveNodeAddress.Address {
	// 			// Create a gRPC connection to the server
	// 			client := pb.NewKeyValueStoreClient(server.Conn)
	// 			log.Print("Client created")
	// 			// TODO: Check how ReviveNode is implemented.
	// 			_, err := client.ReviveNode(context.Background())

	// 			if err != nil {
	// 				log.Printf("Failed to revive node: %v with error %v", killNodeAddress.Address, err)
	// 				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to revive node"})
	// 				return
	// 			}
	// 		}
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{"message": "Node successfully revived!"})
	// })})

	router.Run(":8080")

}

func getFastestRespondingServer() (*Server, error) {
	// Create a channel to receive the first responding server
	ch := make(chan Server, len(servers))

	log.Print("Finding fastest responding server!")
	// Ping all servers concurrently
	for _, server := range servers {
		log.Print("Pinging server", server.Address.Address)
		go func(s Server) {
			// Create a gRPC connection to the server
			client := pb.NewKeyValueStoreClient(s.Conn)
			log.Print("Client created")
			// Call your gRPC ping method here (replace with your actual method)
			_, err := client.Ping(context.Background(), &pb.PingRequest{})
			if err == nil {
				ch <- s
			} else {
				log.Print("Error in ping", err)
			}
		}(server)
	}

	select {
	case server := <-ch:
		return &server, nil
	case <-time.After(3 * time.Second):
		return nil, status.Error(codes.DeadlineExceeded, "No server responded in time")
	}
}

func getServersAddresses(servers []Server) []*pb.Node {
	addresses := make([]*pb.Node, len(servers))
	for i, server := range servers {
		addresses[i] = server.Address
	}
	return addresses
}

func ConvertPbReadResponse(keyValue []*pb.KeyValue) Response {
	log.Print("Keyvalue", keyValue)
	var result Response
	result = Response{
		Key:             keyValue[0].Key,
		ReplicaResponse: make([]ReplicaResponse, len(keyValue)),
		Hashvalue:       hash.GenHash(keyValue[0].Key),
	}
	log.Print("Key", keyValue[0].Key)
	log.Print("Hash", hash.GenHash(keyValue[0].Key))
	for _, kv := range keyValue {
		m := make(map[uint32]int64)
		for k, v := range kv.VectorClock.Timestamps {
			m[k] = v.ClokcVal
		}
		replicaResponse := ReplicaResponse{
			Value:       kv.Value,
			VectorClock: m,
		}
		result.ReplicaResponse = append(result.ReplicaResponse, replicaResponse)
	}
	log.Print("Replica Response", result.ReplicaResponse)

	return result
}
