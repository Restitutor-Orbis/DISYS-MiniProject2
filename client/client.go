package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

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

	//Read user input in terminal
	go ReadFromTerminal(client)

	//read from server
	go PrintBroadcastsFromServer(client)

	/* for {
		t.Sleep(3 * t.Second)
		PublishToServer(client)
	} */

	//make sure client doesn't close
	for {
		t.Sleep(1000 * t.Hour)
	}
}

func ReadFromTerminal(client chittychat.ChittyChatClient) {
	for {
		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalf("Failed to read from console")
		}

		clientMessage = strings.Trim(clientMessage, "\r\n")

		publishRequest := chittychat.PublishRequest{
			User:    strconv.FormatInt(rand.Int63n(10000), 10),
			Message: clientMessage,
			Time:    t.Now().Format("2006-01-02 15:04:05"),
		}

		PublishToServer(client, publishRequest)
	}
}

//call grpc method
func PublishToServer(client chittychat.ChittyChatClient, message chittychat.PublishRequest) {
	client.Publish(context.Background(), &message)
}

func PrintBroadcastsFromServer(client chittychat.ChittyChatClient) {

	clientMessage := chittychat.BroadcastRequest{
		UserId: rand.Int31n(10000),
	}

	stream, err := client.Broadcast(context.Background(), &clientMessage)
	if err != nil {
		log.Fatalf("Error while opening stream %v", err)
	}

	for {
		messageToPrint, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		fmt.Println("["+messageToPrint.Time+"]", messageToPrint.User+":", messageToPrint.Message)
	}
}
