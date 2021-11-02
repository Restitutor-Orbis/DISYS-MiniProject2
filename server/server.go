package main

import (
	"context"
	"log"
	"net"
	"strconv"

	chittychat "github.com/Restitutor-Orbis/DISYS-MiniProject2/chittychat"

	"google.golang.org/grpc"
)

type Server struct {
	chittychat.UnimplementedChittyChatServer
}

//slice of all known connections
var sliceOfStreams []chittychat.ChittyChat_SubscribeServer

//server's lamport time
//init to 0
var lamportTime = chittychat.LamportTime{
	Time: 0,
}

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9081")
	if err != nil {
		log.Fatalf(GetTimeAsString(), "Failed to listen on port 9081: %v", err)
	}

	grpcServer := grpc.NewServer()
	chittychat.RegisterChittyChatServer(grpcServer, &Server{})

	log.Println(GetTimeAsString(), "Server is set up on port 9081")

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf(GetTimeAsString(), "failed to server %v", err)
	}

	//grpc listen and serve
	err = grpcServer.Serve(list)
	if err != nil {
		log.Fatalf(GetTimeAsString(), "Failed to start gRPC Server :: %v", err)
	}
}

func (s *Server) Publish(ctx context.Context, in *chittychat.PublishRequest) (*chittychat.PublishReply, error) {
	//increment lamport time
	lamportTime.UpdateTime(in.Time)
	lamportTime.Time++

	log.Println(GetTimeAsString(), "Received PublishRequest from", in.User)

	broadcastReply := chittychat.SubscribeReply{
		User:    in.User,
		Message: in.Message,
		Time:    lamportTime.Time,
	}

	BroadcastToAllClients(&broadcastReply)

	return &chittychat.PublishReply{}, nil
}

func (s *Server) Subscribe(in *chittychat.SubscribeRequest, subscriptionServer chittychat.ChittyChat_SubscribeServer) error {

	//increment lamport time
	lamportTime.UpdateTime(in.Time)
	lamportTime.Time++

	message := chittychat.SubscribeReply{
		User:    in.Username,
		Message: "has joined the chat",
		Time:    lamportTime.Time,
	}

	//save this stream
	//this should maybe be handled in a separate go routine, to prevent the server from being killed off?
	sliceOfStreams = append(sliceOfStreams, subscriptionServer)

	log.Println(GetTimeAsString(), "Added new client stream to server")

	BroadcastToAllClients(&message)

	//prevent function from terminating
	//keeps the stream connection alive
	for {
		select {
		case <-subscriptionServer.Context().Done():
			broadcastReply := chittychat.SubscribeReply{
				User:    message.User,
				Message: "has left the chat",
				Time:    message.Time,
			}

			log.Println(GetTimeAsString(), message.User+" has disconnected")

			BroadcastToAllClients(&broadcastReply)

			/* for _, element := range sliceOfStreams {
				if element == subscriptionServer {
					element = nil
				}
			} */

			return nil
		}
	}
}

func BroadcastToAllClients(message *chittychat.SubscribeReply) {

	//increment lamport time
	lamportTime.Time++
	message.Time = lamportTime.Time

	//send message to every known stream
	for _, element := range sliceOfStreams {
		if element != nil {
			log.Println(GetTimeAsString(), "Broadcasting to client")
			element.Send(message)
		}
	}
}

//doesn't require lamport increment
func GetTimeAsString() string {
	timeToString := strconv.FormatInt(int64(lamportTime.Time), 10)
	return "[" + timeToString + "]"
}
