package main

import (
	"fmt"
	"gRPC/greet/greetpb"
	"io"
	"log"
	"time"

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
	fmt.Printf("Connection created\n")

	//doUnary(c)
	//doServerStreaming(c)
	// doClientStreaming(c)
	doBiDiStreaming(c)
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
	req := &greetpb.ManyTimesGreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Abhishek",
			LastName:  "Sharma",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet Many Times Service: %v\n", err)
	}
	for {
		mesg, err := resStream.Recv()
		if err == io.EOF {
			// We've reached end of the streaming
			break
		} else if err != nil {
			log.Fatalf("Error while reading response from stream: %v\n", err)
		}
		fmt.Println(mesg.GetGreeting())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do client streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Abhishek",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Preeti",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rahul",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ankit",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Happy",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling Long Greet: %v\n", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		//Sending request to LongGreet
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from Long Greet: %v\n", err)
	}
	fmt.Printf("Response from Long Greet: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do bi di streaming RPC...")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v\n", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Abhishek",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Preeti",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rahul",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ankit",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Happy",
			},
		},
	}

	waitc := make(chan struct{})

	//go routine to send request to server
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending request to stream: %v\n", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("Error while sending data to stream: %v\n", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// go routine to received response from server
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Error while receiving data from stream: %v\n", err)
				break
			}
			fmt.Printf("Received response: %v\n", resp.GetResult())
		}
		close(waitc)
	}()

	//block until everything is done
	<-waitc
}
