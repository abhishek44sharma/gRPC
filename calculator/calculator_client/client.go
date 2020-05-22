package main

import (
	"context"
	"fmt"
	"gRPC/calculator/calculatorpb"
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
	doUnary(c)
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
