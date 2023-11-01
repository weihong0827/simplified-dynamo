package webclient

import (
	"context"
	"dynamoSimplified/config"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"sync"
	"time"

	utils "dynamoSimplified/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "dynamoSimplified/pb" // replace with your actual gRPC protos package
)

type Server struct {
	Address *pb.Node
	Conn    *grpc.ClientConn
}

var servers []Server
var mutex = &sync.Mutex{}

func main() {
	// Initialize the list of backend servers
	servers = []Server{
		{&pb.Node{Id: utils.GenHash("localhost:50051"), Address: "localhost:50051"}, nil},
		{&pb.Node{Id: utils.GenHash("localhost:50052"), Address: "localhost:50052"}, nil},
		{&pb.Node{Id: utils.GenHash("localhost:50053"), Address: "localhost:50053"}, nil},
		{&pb.Node{Id: utils.GenHash("localhost:50054"), Address: "localhost:50054"}, nil},
		{&pb.Node{Id: utils.GenHash("localhost:50055"), Address: "localhost:50055"}, nil},
	}

	// Establish gRPC connections to all servers
	for i, server := range servers {
		conn, err := grpc.Dial(server.Address.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to server at %s: %v", server.Address, err)
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
		conn, err := grpc.Dial(fastestServer.Address.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

		client := pb.NewKeyValueStoreClient(conn)

		// TODO: Call your gRPC method here (replace with your actual method), Also make protobuf message
		resp, err := client.Read(context.Background(), &pb.ReadRequest{Key: c.Query("key")})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Fatalf("Failed to read key: %v with error %v", c.Query("key"), err)
			return
		}

		// Forward the response from the backend server to the client
		c.JSON(http.StatusOK, gin.H{"message": resp.Message})
	})

	router.GET("/addNode", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()

		port := c.Query("port")
		hash := utils.GenHash("127.0.0.1:" + port)
		nodes, err := utils.GetNodesFromKey(hash, getServersAddresses(servers), config.READ)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		nodeToContact := nodes[0]

		c.JSON(http.StatusOK, gin.H{"message": nodeToContact.Address})
	})

	router.PUT("/put", func(c *gin.Context) {
		fastestServer, err := getFastestRespondingServer()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Establish a gRPC connection to the fastest server
		conn, err := grpc.Dial(fastestServer.Address.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

		client := pb.NewKeyValueStoreClient(conn)

		resp, err := client.Write(context.Background(), &pb.WriteRequest{KeyValue: &pb.KeyValue{Key: c.Query("key"), Value: c.Query("value")}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Fatalf("Failed to write key: %v with error %v", c.Query("key"), err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": resp.Message})
	})

	router.Run(":8080")
}

func getFastestRespondingServer() (*Server, error) {
	// Create a channel to receive the first responding server
	ch := make(chan *Server, len(servers))

	// Ping all servers concurrently
	for _, server := range servers {
		go func(s *Server) {
			client := pb.NewKeyValueStoreClient(s.Conn)

			// Call your gRPC ping method here (replace with your actual method)
			_, err := client.Ping(context.Background(), &pb.PingRequest{})
			if err == nil {
				ch <- s
			}
		}(&server)
	}

	select {
	case server := <-ch:
		return server, nil
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
