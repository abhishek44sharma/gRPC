package main

import (
	"context"
	"fmt"
	"gRPC/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked: %v\n", req)
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

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("Long Greet function invoked with streaming request\n")
	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// We've reached end of stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		} else if err != nil {
			log.Fatalf("Error while reading from client stream: %v\n", err)
		}
		result += req.GetGreeting().GetFirstName() + " "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("Greet Everyone function invoked with streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + " !"
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error while sending data to client: %v\n", err)
			return err
		}
	}
}

func main() {
	fmt.Println("Starting Greet listener")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Println("Started Greet listener")

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}
}
