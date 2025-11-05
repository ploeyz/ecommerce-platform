package client

import (
	"context"
	"fmt"
	"log"
	"time"
	pb "github.com/ploezy/ecommerce-platform/order-service/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

var userClient *UserClient

// NewUserClient creates a new gRPC client for User Service
func NewUserClient(address string) (*UserClient, error) {
	// Create connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dial to User Service gRPC server
	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	// Create gRPC client
	client := pb.NewUserServiceClient(conn)

	userClient = &UserClient{
		client: client,
		conn: 	conn,
	}
	log.Printf("User Service gRPC client connected to %s",address)
	return userClient,nil
}

// GetUserClient returns the singleton user client instance
func GetUserClient() *UserClient{
	return userClient
}

// GetUSer retrieves user information by user ID
func (c *UserClient) GetUser(ctx context.Context,userID uint32) (*pb.GetUserResponse, error) {
	req := &pb.GetUserRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp, nil
}

// ValidateToken validates a JWT token
func (c *UserClient) ValidateToken(ctx context.Context, token string) (*pb.ValidateTokenResponse, error){
	req := &pb.ValidateTokenRequest{
		Token: token,
	}
	resp, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate toekn: %w", err)
	}
	return resp, nil
}

// Close close the gRPC connection
func(c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}