package main

import (
	"fmt"
	"gRPC/greet/greetpb"
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

	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
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
