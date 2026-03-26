package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/hepstore/tribecart/proto/tribecart/v1"
)

type server struct {
	pb.UnimplementedPaymentServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	log.Printf("Processing payment for order %s with amount %f", req.OrderId, req.Amount)
	// A real implementation would integrate with a payment provider
	return &pb.ProcessPaymentResponse{
		TransactionId: "some-transaction-id",
		Success:       true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, NewServer())
	reflection.Register(s)
	log.Println("Server listening at", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
