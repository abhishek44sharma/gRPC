package main

import (
	"context"
	"fmt"
	"gRPC/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50023", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to setup connection: %v\n", err)
	}

	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)

	// Unary
	// doUnary(c)

	// Server Streaming
	doServerStreaming(c)
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
