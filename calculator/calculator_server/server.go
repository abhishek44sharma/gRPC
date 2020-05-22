package main

import (
	"context"
	"fmt"
	"gRPC/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func (s *server) Sum(ctx context.Context, req *calculatorpb.Request) (*calculatorpb.Response, error) {
	fmt.Printf("Sum request received: %v\n", req)
	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()

	sum := firstNumber + secondNumber
	resp := &calculatorpb.Response{
		Sum: sum,
	}
	return resp, nil
}

type server struct{}

func main() {
	fmt.Println("Starting calculator listener")
	lis, err := net.Listen("tcp", "0.0.0.0:50023")
	if err != nil {
		log.Fatalf("Failed to start listner: %v\n", err)
	}
	fmt.Println("Started calculator server")

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server request: %v\n", err)
	}
}
