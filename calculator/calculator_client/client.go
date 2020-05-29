package main

import (
	"context"
	"fmt"
	"gRPC/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50023", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to setup connection: %v\n", err)
	}

	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBiDiStreaming(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.Request{
		FirstNumber:  20,
		SecondNumber: 5,
	}
	resp, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while invoking Sum service: %v\n", err)
	}
	fmt.Println(resp)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberRequest{
		Number: 120,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition service: %v\n", err)
	}
	log.Printf("Decomposed number: ")
	for {
		mesg, err := resStream.Recv()
		if err == io.EOF {
			// We've reached end of response
			break
		} else if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}
		log.Print(mesg.GetDecomposedNumber())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting to do client streaming RPC...\n")
	requests := []*calculatorpb.AverageRequest{
		&calculatorpb.AverageRequest{
			Number: 1,
		},
		&calculatorpb.AverageRequest{
			Number: 2,
		},
		&calculatorpb.AverageRequest{
			Number: 3,
		},
		&calculatorpb.AverageRequest{
			Number: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling compute average: %v\n", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request to calculate Average: %v\n", req)
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending request to stream: %v\n", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from Compute Average: %v\n", err)
	}
	fmt.Printf("Response received: %v\n", res)
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a bi-di streaming RPC...")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream\n")
		return
	}

	requests := []*calculatorpb.FindMaxRequest{
		&calculatorpb.FindMaxRequest{
			Number: 1,
		},
		&calculatorpb.FindMaxRequest{
			Number: 5,
		},
		&calculatorpb.FindMaxRequest{
			Number: 3,
		},
		&calculatorpb.FindMaxRequest{
			Number: 6,
		},
		&calculatorpb.FindMaxRequest{
			Number: 2,
		},
		&calculatorpb.FindMaxRequest{
			Number: 20,
		},
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending request to stream: %v\n", req)
			if err := stream.Send(req); err != nil {
				log.Fatalf("Error while sending data to stream: %v\n", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Error while receving data from stream: %v\n", err)
			}
			fmt.Printf("Response received: %v\n", resp.GetNumber())
		}
		close(waitc)
	}()

	<-waitc
}
