package server

import (
	"fmt"
	"log"
	"net"
	grpcHandler "github.com/ploezy/ecommerce-platform/product-service/internal/grpc/handler"
	pb "github.com/ploezy/ecommerce-platform/product-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server  *grpc.Server
	handler *grpcHandler.ProductGRPCHandler
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(handler *grpcHandler.ProductGRPCHandler) *GRPCServer {
	server := grpc.NewServer()
	
	// Register Product Service
	pb.RegisterProductServiceServer(server, handler)
	
	// Register reflection service (for tools like grpcurl)
	reflection.Register(server)

	return &GRPCServer{
		server:  server,
		handler: handler,
	}
}

// Start starts the gRPC server
func (s *GRPCServer) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	log.Printf("gRPC Server is running on port %s\n", port)
	
	if err := s.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

// Stop stops the gRPC server gracefully
func (s *GRPCServer) Stop() {
	log.Println("Stopping gRPC server...")
	s.server.GracefulStop()
}