package main

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
    "log"

	t "time"
    "github.com/Restitutor-Orbis/DISYS-MiniProject2/chittychat"
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
    c := chittychat.NewGetCurrentTimeClient(conn)
 
    for {
		SendRequest(c)
		t.Sleep(5*t.Second)
	}
}

func SendRequest(c chittychat.GetCurrentTimeClient) {
    // Between the curly brackets are nothing, because the .proto file expects no input.
    message := chittychat.GetTimeRequest{}

    response, err := c.GetTime(context.Background(), &message)
    if err != nil {
        log.Fatalf("Error when calling XXX: %s", err)
    }

    fmt.Printf("Response from the Server: %s \n", response.Reply)
}