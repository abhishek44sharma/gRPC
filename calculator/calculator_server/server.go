package main

import (
	"context"
	"fmt"
	"gRPC/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

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

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Prime number decomposition request received: %v\n", req)
	number := req.GetNumber()
	k := int32(2)

	for number > 1 {
		if number%k == 0 {
			resp := &calculatorpb.DecompositionResponse{
				DecomposedNumber: k,
			}
			fmt.Printf("Sent number %d to stream\n", k)
			stream.Send(resp)
			number = number / k
		} else {
			k = k + 1
		}
	}
	return nil
}

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
