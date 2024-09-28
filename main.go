package main

import (
	"log"
	"net"

	grpcServers "github.com/alpardfm/go-grpc-rest/api/grpc"
	"github.com/alpardfm/go-grpc-rest/api/rest"
	"github.com/alpardfm/go-grpc-rest/db"
	"github.com/alpardfm/go-grpc-rest/pb"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Mengimpor driver PostgreSQL
	"google.golang.org/grpc"
)

func main() {
	db.InitDB()

	// Run gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterProductServiceServer(grpcServer, &grpcServers.ProductServer{})
		log.Println("gRPC server listening on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Run REST API server
	router := gin.Default()
	rest.RegisterRoutes(router)
	log.Println("REST server listening on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to serve REST server: %v", err)
	}
}
