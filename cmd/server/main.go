package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/rikatz/grpctest/api"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedEventServer
	pb.UnimplementedConfigurationServer
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 12345))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterEventServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}

func (s *server) PublishEvent(stream pb.Event_PublishEventServer) error {
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.EventReturn{})
		}
		if err != nil {
			return err
		}
		fmt.Printf("\n%+v // %s // %s\n", event.Backend, event.Eventtype, event.Message)
	}
}
