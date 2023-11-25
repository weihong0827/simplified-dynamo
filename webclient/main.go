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

type ReadWriteResponse struct {
	Key             string                     `JSON:"key"`
	ReplicaResponse []ReadWriteReplicaResponse `JSON:"replica_response"`
	Hashvalue       uint32                     `JSON:"hashvalue"`
}

type ReadWriteReplicaResponse struct {
	Value       string           `JSON:"value"`
	VectorClock map[uint32]int64 `JSON:"vector_clock"`
	NodeId      uint32           `JSON:"node_id"`
}

var (
	servers []Server
	mutex   = &sync.Mutex{}
)

func main() {
	// Initialize the list of backend servers NEED TO RUN AT LEAST 4 NODES
	servers = []Server{
		// {&pb.Node{Id: hash.GenHash("127.0.0.1:50051"), Address: "127.0.0.1:50051"}, nil},
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
		log.Printf("client initiated %s", fastestServer.Address.Address)

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

		result := ConvertPbReadResponse(c.Query("key"), resp.KeyValue)

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
		log.Printf("Fastest Node %v:", fastestServer.Address.Address)
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

		log.Printf("%v", resp)

		c.JSON(http.StatusOK, gin.H{"message": "Write Successful!"})
	})

	router.POST("/kill", func(c *gin.Context) {
		var killNodeAddress KillNode
		if err := c.ShouldBindJSON(&killNodeAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		conn, err := grpc.Dial(killNodeAddress.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		client := pb.NewKeyValueStoreClient(conn)
		_, err = client.KillNode(context.Background(), &pb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Killed node at address: " + killNodeAddress.Address})
	})

	router.POST("/revive", func(c *gin.Context) {
		var reviveNodeAddress ReviveNode
		if err := c.ShouldBindJSON(&reviveNodeAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		conn, err := grpc.Dial(reviveNodeAddress.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		client := pb.NewKeyValueStoreClient(conn)
		_, err = client.ReviveNode(context.Background(), &pb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Revive node at address: " + reviveNodeAddress.Address})
	})
	router.Run(":8080")
}

func getFastestRespondingServer() (*Server, error) {
	// Create a channel to receive the first responding server
	ch := make(chan Server, 1)

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
				log.Print("Server responded: ", s.Address.Address)
				ch <- s
			} else {
				log.Print("Error in ping: ", err)
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

func ConvertPbReadResponse(key string, keyValue map[uint32]*pb.KeyValue) ReadWriteResponse {
	log.Print("Keyvalue", keyValue)

	result := ReadWriteResponse{
		Key:             key,
		ReplicaResponse: make([]ReadWriteReplicaResponse, 0),
		Hashvalue:       hash.GenHash(key),
	}
	for sourceNodeId, kv := range keyValue {
		m := make(map[uint32]int64)
		for k, v := range kv.VectorClock.Timestamps {
			m[k] = v.ClokcVal
		}
		replicaResponse := ReadWriteReplicaResponse{
			Value:       kv.Value,
			VectorClock: m,
			NodeId:      sourceNodeId,
		}
		result.ReplicaResponse = append(result.ReplicaResponse, replicaResponse)
	}
	log.Print("Replica Response", result.ReplicaResponse)

	return result
}

func ConvertPbWriteResponse(keyValue []*pb.KeyValue) ReadWriteResponse {
	var result ReadWriteResponse
	result = ReadWriteResponse{
		Key:             keyValue[0].Key,
		ReplicaResponse: make([]ReadWriteReplicaResponse, 0),
		Hashvalue:       hash.GenHash(keyValue[0].Key),
	}
	for _, kv := range keyValue {
		m := make(map[uint32]int64)
		for k, v := range kv.VectorClock.Timestamps {
			m[k] = v.ClokcVal
		}
		replicaResponse := ReadWriteReplicaResponse{
			Value:       kv.Value,
			VectorClock: m,
		}
		result.ReplicaResponse = append(result.ReplicaResponse, replicaResponse)
	}
	log.Printf("Replica Response %v", result.ReplicaResponse)
	return result
}
