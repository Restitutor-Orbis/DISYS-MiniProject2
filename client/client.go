package main

import (
	"context"
	"fmt"
	"io"
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

	//read from server
	go PrintBroadcasts(client)

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
		Time:    t.Now().Format("2006-01-02 15:04:05"),
	}

	client.Publish(context.Background(), &message)
}

func PrintBroadcasts(client chittychat.ChittyChatClient) {

	var count int

	clientMessage := chittychat.BroadcastRequest{
		UserId: rand.Int31n(10000),
	}

	stream, err := client.Broadcast(context.Background(), &clientMessage)
	if err != nil {
		log.Fatalf("Error while opening stream %v", err)
	}

	for {

		messageToPrint, err := stream.Recv()

		//fmt.Println(requestToPrint)

		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		count++
		fmt.Println("["+messageToPrint.Time+"]", messageToPrint.User+":", messageToPrint.Message)

		t.Sleep(1 * t.Second)
		fmt.Println(count)
	}
}
