package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	t "time"

	"github.com/Restitutor-Orbis/DISYS-MiniProject2/chittychat"
)

func main() {
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	client := chittychat.NewChittyChatClient(conn)

	for {
		SendRequest(client)
		t.Sleep(5 * t.Second)
	}
}

func SendRequest(client chittychat.ChittyChatClient) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	message := chittychat.PublishRequest{
		User:    "Bo",
		Message: "Hej med dig",
		Time:    t.Now().GoString(),
	}

	client.Publish(context.Background(), &message)
}

func PrintBroadcasts(client chittychat.ChittyChatClient, request chittychat.BroadcastRequest, id int32) {
	clientMessage := chittychat.BroadcastRequest{
		UserId: id,
	}

	stream, err := client.Broadcast(context.Background(), &clientMessage)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	for {
		requestToPrint, err := stream.Recv()

		if err != nil {
			log.Fatalf("cannot receive %v", err)
		}

		fmt.Println(requestToPrint.User, requestToPrint.Message, requestToPrint.Time)
	}
}
