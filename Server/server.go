package main

import (
	"context"
	"fmt"
	"log"
	"net"

	cc "github.com/Restitutor-Orbis/DISYS-MiniProject2/ChittyChat"

	"google.golang.org/grpc"
)

type Server struct {
	//cc.ChittyChatServer
	cc.UnimplementedChittyChatServer
}

func (s *Server) Publish(ctx context.Context, in *cc.PublishMessage) (*cc.PublishReply, error) {
	fmt.Printf("Received PublishMessage request\n")
	var message string = in.Message
	return &cc.PublishReply{Reply: message}, nil
}

func main() {
	fmt.Println("Starting server")

	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()
	cc.RegisterChittyChatServer(grpcServer, &Server{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
