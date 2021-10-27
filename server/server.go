package main

import (
	"context"
	"fmt"
	"log"
	"net"

	chittychat "github.com/Restitutor-Orbis/DISYS-MiniProject2/chittychat"

	"google.golang.org/grpc"
)

type Server struct {
	chittychat.UnimplementedChittyChatServer
}

//init new map to track users
var userIDtoNameMap map[int]string = make(map[int]string)

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}

	grpcServer := grpc.NewServer()
	chittychat.RegisterChittyChatServer(grpcServer, &Server{})

	fmt.Println("Server is set up on port 9080")

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func (s *Server) Publish(ctx context.Context, in *chittychat.PublishRequest) (*chittychat.PublishReply, error) {

	fmt.Println("Received PublishRequest from", in.User)

	addClientToMap(in.User)

	return &chittychat.PublishReply{}, nil
}

/* func (s *Server) Broadcast(ctx context.Context, in *chittychat.BroadcastRequest) (*chittychat.BroadcastReply, error) {
	fmt.Println("Sending message from", in.User, "to")

	return &chittychat.BroadcastReply{}, nil
} */

func addClientToMap(name string) {

	var id int

	//check if user already exists
	for key := range userIDtoNameMap {
		if userIDtoNameMap[key] == name {
			//stop function if username already exists
			return
		}
	}

	//set id index to length
	id = len(userIDtoNameMap)
	fmt.Println("User", name, "added to client")

	//add to map
	userIDtoNameMap[id] = name
}
