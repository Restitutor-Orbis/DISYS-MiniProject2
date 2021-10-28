package main

import (
	"context"
	"fmt"
	"log"
	"net"

	//"sync"

	t "time"

	chittychat "github.com/Restitutor-Orbis/DISYS-MiniProject2/chittychat"

	"google.golang.org/grpc"
)

type Server struct {
	chittychat.UnimplementedChittyChatServer
}

var sliceOfStreams []chittychat.ChittyChat_BroadcastServer

//init new map to track users
//var userIDtoNameMap map[int]string = make(map[int]string)

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9081")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}

	grpcServer := grpc.NewServer()
	chittychat.RegisterChittyChatServer(grpcServer, &Server{})

	fmt.Println("Server is set up on port 9080")

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

	BroadcastToAllClients(in)

	return &chittychat.PublishReply{}, nil
}

func (s *Server) Broadcast(in *chittychat.BroadcastRequest, broadcastServer chittychat.ChittyChat_BroadcastServer) error {
	fmt.Println("Initializing", in.UserId)

	message := chittychat.BroadcastReply{
		User:    "Bo",
		Message: "has joined",
		Time:    t.Now().Format("15:04:05"),
	}

	sliceOfStreams = append(sliceOfStreams, broadcastServer)

	fmt.Println("Added stream to server")

	broadcastServer.Send(&message)

	return nil
}

func BroadcastToAllClients(message *chittychat.PublishRequest) {

	//var mutex = &sync.Mutex{}

	fmt.Println("Broadcasting to", len(sliceOfStreams), "people")

	broadcastReply := chittychat.BroadcastReply{
		User:    message.User,
		Message: message.Message,
		Time:    message.Time,
	}

	//send message to every known stream
	for i := range sliceOfStreams {

		//mutex.Lock()

		sliceOfStreams[i].Send(&broadcastReply)
		//fmt.Println(sliceOfStreams[i].Context())

		//mutex.Unlock()

		fmt.Println("Sending message to user", i)
	}
}

/* func addClientToMap(name string) {

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
*/
