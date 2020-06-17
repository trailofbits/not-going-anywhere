package main

import (
    "context"
    "log"
    "time"
    "fmt"

    // external dependencies
    "google.golang.org/grpc"

    // our ProtoBuffers internal libraries
    pb "github.com/trailofbits/not-going-anywhere/internal/ping"
)

const (
    address = "localhost:5001"
)

func main() {
    conn, err := grpc.Dial(address, grpc.WithInsecure())

    if err != nil {
        log.Fatalf("could not connect to server; make sure you have 'friends_server_net' running: %v", err)
    }

    defer conn.Close()

    client := pb.NewPingNGAClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    for idx := 0; idx < 10; idx++  {
        data := fmt.Sprintf("secret%d", idx)
        resp, _ := client.PingSensitive(ctx, &pb.Message{Data: data})
        log.Printf("received token: %v", resp.GetData())
    }
}
