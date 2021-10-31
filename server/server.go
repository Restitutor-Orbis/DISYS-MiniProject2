package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	//"sync"

	t "time"

	chittychat "github.com/Restitutor-Orbis/DISYS-MiniProject2/chittychat"

	"google.golang.org/grpc"
)

type Server struct {
	chittychat.UnimplementedChittyChatServer
}

var sliceOfStreams []chittychat.ChittyChat_BroadcastServer
var UserIDtoUsername = make(map[int32]string)

//init new map to track users
//var userIDtoNameMap map[int]string = make(map[int]string)

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9081")
	if err != nil {
		log.Fatalf("Failed to listen on port 9081: %v", err)
	}

	grpcServer := grpc.NewServer()
	chittychat.RegisterChittyChatServer(grpcServer, &Server{})

	fmt.Println("Server is set up on port 9081")

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}

	//grpc listen and serve
	err = grpcServer.Serve(list)
	if err != nil {
		log.Fatalf("Failed to start gRPC Server :: %v", err)
	}
}

func (s *Server) Publish(ctx context.Context, in *chittychat.PublishRequest) (*chittychat.PublishReply, error) {

	fmt.Println("Received PublishRequest from", in.User)

	//addClientToMap(in.User)

	broadcastReply := chittychat.BroadcastReply{
		User:    in.User,
		Message: in.Message,
		Time:    in.Time,
	}

	BroadcastToAllClients(&broadcastReply)

	return &chittychat.PublishReply{}, nil
}

func (s *Server) Broadcast(in *chittychat.BroadcastRequest, broadcastServer chittychat.ChittyChat_BroadcastServer) error {
	fmt.Println("Initializing", in.UserId)

	message := chittychat.BroadcastReply{
		User:    strconv.FormatInt(int64(in.UserId), 10),
		Message: "has joined",
		Time:    t.Now().Format("15:04:05"),
	}

	//save this stream
	//this should maybe be handled in a separate go routine, to prevent the server from being killed off?
	sliceOfStreams = append(sliceOfStreams, broadcastServer)

	fmt.Println("Added stream to server")

	BroadcastToAllClients(&message)

	//prevent function from terminating
	//keeps the stream connection alive
	for {
		select {
		case <-broadcastServer.Context().Done():
			broadcastReply := chittychat.BroadcastReply{
				User:    message.User,
				Message: "has left the chat",
				Time:    message.Time,
			}

			BroadcastToAllClients(&broadcastReply)

			/* for _, element := range sliceOfStreams {
				if element == broadcastServer {
					element = nil
				}
			} */

			return nil
			//element = nil
		}
	}
}

func BroadcastToAllClients(message *chittychat.BroadcastReply) {
	//send message to every known stream
	for _, element := range sliceOfStreams {
		if element != nil {
			element.Send(message)
		}
	}
}
