package main

import (
	"context"
	"fmt"
	"log"
	t "time"

	cc "github.com/Restitutor-Orbis/DISYS-MiniProject2/ChittyChat"

	"google.golang.org/grpc"
)

func main() {
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := cc.NewChittyChatClient(conn)

	fmt.Println("")
	fmt.Println("---")

	for {
		SendPublishRequest(c)
		t.Sleep(5 * t.Second)
	}
}

func SendPublishRequest(c cc.ChittyChatClient) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	message := cc.PublishMessage{Message: "hej"}

	response, err := c.Publish(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Publish: %s", err)
	}

	fmt.Printf("Message recieved: %s\n", response.Reply)
}
