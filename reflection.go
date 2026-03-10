package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
)

const address = "localhost:50052"

type Stream = grpc_reflection_v1.ServerReflection_ServerReflectionInfoClient

func createReflectionClient() (*grpc.ClientConn, Stream, error) {
	// Connect to your port-forwarded address
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
		return nil, nil, err
	}

	// Create a reflection client
	client := grpc_reflection_v1.NewServerReflectionClient(conn)
	stream, err := client.ServerReflectionInfo(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
		return nil, nil, err
	}

	return conn, stream, nil
}

func getServices(stream Stream) {
	// Send a request to list services
	err := stream.Send(&grpc_reflection_v1.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1.ServerReflectionRequest_ListServices{
			ListServices: "*", // Get everything
		},
	})
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	// Receive and print the response
	resp, err := stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive response: %v", err)
	}

	services := resp.GetListServicesResponse().Service
	fmt.Printf("Found %d services on %s:\n", len(services), address)

	for _, s := range services {
		fmt.Printf(" - Service: %s\n", s.Name)
	}
}

func getMethods() {
}
