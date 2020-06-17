package main

import (
    "context"
    "net"
    "log"

    // external dependencies
    "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
    "github.com/google/uuid"

    // our ProtoBuffers internal libraries
    pb "github.com/trailofbits/not-going-anywhere/internal/ping"
)

const (
    port = ":5002"
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

	creds, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

    grpcserv := grpc.NewServer(grpc.Creds(creds))
    pb.RegisterPingNGAServer(grpcserv, &server{})

    if err := grpcserv.Serve(netserv); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
