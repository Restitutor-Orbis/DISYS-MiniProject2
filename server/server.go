package main

import (
	"context"
	"fmt"
	"log"
	"net"
	t "time"

	chittychat "github.com/Restitutor-Orbis/DISYS-MiniProject2/chittychat"

	"google.golang.org/grpc"
)

type Server struct {
	chittychat.UnimplementedGetCurrentTimeServer
}

func (s *Server) GetTime(ctx context.Context, in *chittychat.GetTimeRequest) (*chittychat.GetTimeReply, error) {
	fmt.Printf("Received get time request")
	return &chittychat.GetTimeReply{Reply: t.Now().String()}, nil
}

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}

	grpcServer := grpc.NewServer()
	chittychat.RegisterGetCurrentTimeServer(grpcServer, &Server{})

	fmt.Println("Server is set up on port 9080")

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
