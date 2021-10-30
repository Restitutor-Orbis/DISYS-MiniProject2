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
		User:    "Bo",
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
		t.Sleep(1000 * t.Hour)
	}
}

func StoreClientStream(broadcastServer chittychat.ChittyChat_BroadcastServer) {
	sliceOfStreams = append(sliceOfStreams, broadcastServer)

	fmt.Println("Added", broadcastServer, "to slice")

	broadcastReply := chittychat.BroadcastReply{
		User:    "Jens",
		Message: "Sendt fra StoreClientStream",
		Time:    t.Now().Format("15:04:05"),
	}

	//broadcastServer.Send(&broadcastReply)

	for i, v := range sliceOfStreams {

		//mutex.Lock()

		sliceOfStreams[i].Send(&broadcastReply)
		//fmt.Println(sliceOfStreams[i].Context())

		//mutex.Unlock()

		fmt.Println(i, v)
	}
}

func BroadcastToAllClients(message *chittychat.BroadcastReply) {

	//var mutex = &sync.Mutex{}

	fmt.Println("Broadcasting to", len(sliceOfStreams), "people")

	//send message to every known stream
	for _, element := range sliceOfStreams {
		element.Send(message)
	}
}
