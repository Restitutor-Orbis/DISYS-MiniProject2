package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"google.golang.org/grpc"

	t "time"

	"github.com/Restitutor-Orbis/DISYS-MiniProject2/chittychat"
)

func main() {
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9081", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	client := chittychat.NewChittyChatClient(conn)

	//
	clientMessage := chittychat.BroadcastRequest{
		UserId: rand.Int31n(10000),
	}

	stream, err := client.Broadcast(context.Background(), &clientMessage)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	//read from server
	go PrintBroadcasts(client, stream)

	for {
		t.Sleep(3 * t.Second)
		PublishToServer(client)
	}
}

func PublishToServer(client chittychat.ChittyChatClient) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	message := chittychat.PublishRequest{
		User:    "Bo",
		Message: "Hej med dig",
		Time:    t.Now().GoString(),
	}

	client.Publish(context.Background(), &message)
}

func PrintBroadcasts(client chittychat.ChittyChatClient, stream chittychat.ChittyChat_BroadcastClient) {

	for {
		requestToPrint, err := stream.Recv()

		fmt.Println(requestToPrint)

		if err != nil {

		}

		if requestToPrint != nil {
			fmt.Println(requestToPrint.User, requestToPrint.Message, requestToPrint.Time)
		}

		t.Sleep(1 * t.Second)
	}
}
