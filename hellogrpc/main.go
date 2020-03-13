package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "github.com/fransoaardi/hellogrpc/proto"
)

type helloServer struct {
	pb.UnimplementedHelloServer
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
	}

	server := grpc.NewServer()
	srv := new(helloServer)
	pb.RegisterHelloServer(server, srv)

	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
	}
}

func (*helloServer) SayHello(ctx context.Context, in *pb.Greeting) (*pb.Introduction, error) {
	version := "v1-hellogrpc"
	hostname, _ := os.Hostname()

	name := in.Name
	out := pb.Introduction{
		Name:     name,
		Version:  version,
		Hostname: hostname,
	}

	return &out, nil
}
