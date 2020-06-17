package main

import (
    "context"
    "net"
    "log"

    // external dependencies
    "google.golang.org/grpc"
    "github.com/google/uuid"

    // our ProtoBuffers internal libraries
    pb "github.com/trailofbits/not-going-anywhere/internal/ping"
)

const (
    port = ":5001"
)

type server struct {
    pb.UnimplementedPingNGAServer
}

func (s *server) PingSensitive(ctx context.Context, in *pb.Message) (*pb.Response, error) {
    log.Printf("received: %v", in.GetData())
    uid, _ := uuid.NewRandom()
    ret := pb.Response{Data: uid.String()}
    return &ret, nil
}

func main() {
    netserv, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to bind to port %d: %v", port, err)
    }

    grpcserv := grpc.NewServer()
    pb.RegisterPingNGAServer(grpcserv, &server{})

    if err := grpcserv.Serve(netserv); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
