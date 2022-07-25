package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/rikatz/grpctest/api"
)

var id string

func main() {
	flag.StringVar(&id, "id", "node1", "")
	flag.Parse()
	rand.Seed(time.Now().UnixMicro())
	conn, err := grpc.Dial("localhost:12345", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewEventClient(conn)
	runPublishEvent(client)

}

func runPublishEvent(client pb.EventClient) {
	ctx := context.TODO()
	stream, err := client.PublishEvent(ctx)
	if err != nil {
		log.Fatalf("client.RecordRoute failed: %v", err)
	}
	for {
		number := rand.Intn(100-1) + 1
		msg := fmt.Sprintf("%s - %d", id, number)
		reason := "blo123"
		if number%5 == 0 {
			reason = "bla"
		}
		message := &pb.EventMessage{
			Backend: &pb.BackendName{
				Name:      "test1",
				Namespace: "bla123",
			},
			Reason:    reason,
			Eventtype: "xpto123",
			Message:   msg,
		}
		if err := stream.Send(message); err != nil {
			log.Fatalf("client.PublishEvent: stream.Send(%v) failed: %v", message, err)
		}
		time.Sleep(1 * time.Second)
	}
}
