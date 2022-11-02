package main

import (
	"log"
	"net"
	"os"

	"github.com/examples/bookstore/rpc"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on port %s", port)
	grpcServer := grpc.NewServer()
	rpc.RegisterBookstoreServer(grpcServer, NewBookstoreServer())
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
