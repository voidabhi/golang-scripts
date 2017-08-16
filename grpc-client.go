package main

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/jzelinskie/grpc/simple"
)

func main() {
	conn, err := grpc.Dial("localhost:6000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	defer conn.Close()

	client := pb.NewSimpleServiceClient(conn)
	stream, err := client.SimpleRPC(context.Background())
	waitc := make(chan struct{})

	msg := &pb.SimpleData{"sup"}
	go func() {
		for {
			log.Println("Sleeping...")
			time.Sleep(2 * time.Second)
			log.Println("Sending msg...")
			stream.Send(msg)
		}
	}()
	<-waitc
	stream.CloseSend()
}
