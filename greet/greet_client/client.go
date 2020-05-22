package main

import (
	"fmt"
	"gRPC/greet/greetpb"
	"io"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to setup connection")
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Connection created: %f", c)

	//doUnary(c)

	doServerStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do Unary gRPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "ABHISHEK",
			LastName:  "SHARMA",
		},
	}
	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling gRPC: %v", err)
	}
	log.Printf("Response from Greet: %v", resp.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do ServerStreaming gRPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Abhishek",
			LastName:  "Sharma",
		},
	}

	respStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes Service: %v\n", err)
	}

	for {
		mesg, err := respStream.Recv()
		if err == io.EOF {
			// We've reached end of the file
			break
		} else if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}
		log.Printf("Response from GreetManyTimes: %s\n", mesg.GetResult())
	}
}
