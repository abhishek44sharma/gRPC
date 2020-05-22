package main

import (
	"context"
	"fmt"
	"gRPC/greet/greetpb"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with: %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	resp := &greetpb.GreetResponse{
		Result: result,
	}
	return resp, nil
}

func (s *server) GreetManyTimes(req *greetpb.ManyTimesGreetRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet Many Times request received: %v\n", req)
	name := req.GetGreeting().GetFirstName()
	for i := 1; i < 10; i++ {
		result := name + " " + strconv.Itoa(i)
		resp := &greetpb.ManyTimesGreetResponse{
			Greeting: result,
		}
		if err := stream.Send(resp); err != nil {
			log.Fatalf("Error while seding response to stream: %v\n", err)
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}
}
