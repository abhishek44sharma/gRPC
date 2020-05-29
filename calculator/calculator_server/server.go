package main

import (
	"context"
	"fmt"
	"gRPC/calculator/calculatorpb"
	"io"
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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Compute Average function invoked to compute average\n")
	count, sum := 0, int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//We've reached end of the input
			result := &calculatorpb.AverageResponse{
				Average: float32(sum) / float32(count),
			}
			return stream.SendAndClose(result)
		} else if err != nil {
			log.Fatalf("Error while receiving request from stream: %v\n", err)
		}
		count += 1
		sum += req.GetNumber()
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	maxNumber := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			log.Fatalf("Error while receiving request from stream: %v\n", err)
			return err
		}
		number := req.GetNumber()
		if number > maxNumber {
			maxNumber = number
			err := stream.Send(&calculatorpb.FindMaxResponse{
				Number: maxNumber,
			})
			if err != nil {
				log.Fatalf("Error while sending data to client: %v\n", err)
				return err
			}
		}
	}
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
