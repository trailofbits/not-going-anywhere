package main

import (
    "context"
    "log"
    "time"
    "fmt"

    // external dependencies
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    // our ProtoBuffers internal libraries
    pb "github.com/trailofbits/not-going-anywhere/internal/ping"
)

const (
    address = "localhost:5002"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	log.Print("loaded cert...")

    conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds), grpc.WithTimeout(1000 * time.Millisecond))

	log.Print("gRPC dialed...")

    if err != nil {
        log.Fatalf("could not connect to server; make sure you have 'ping_server_secure' running: %v", err)
    }

    defer conn.Close()

    client := pb.NewPingNGAClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

	log.Print("client connected...")

    for idx := 0; idx < 10; idx++  {
        data := fmt.Sprintf("secret%d", idx)
        resp, _ := client.PingSensitive(ctx, &pb.Message{Data: data})
        log.Printf("received token: %v", resp.GetData())
    }
}
