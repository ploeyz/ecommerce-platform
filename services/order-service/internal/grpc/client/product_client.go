package client

import (
	"context"
	"fmt"
	"time"
	"log"
	
	pb "github.com/ploezy/ecommerce-platform/order-service/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	client pb.ProductServiceClient
	conn   *grpc.ClientConn
}

var productClient *ProductClient

// NewProductClient creates a new gRPC client for Product Service
func NewProductClient(address string) (*ProductClient, error) {
	// Create connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dial to Product Service gRPC server
	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	// Create gRPC client
	client := pb.NewProductServiceClient(conn)

	productClient = &ProductClient{
		client: client,
		conn:   conn,
	}

	log.Printf("✅ Product Service gRPC client connected to %s", address)
	return productClient, nil
}

// GetProductClient returns the singleton product client instance
func GetProductClient() *ProductClient {
	return productClient
}

// GetProduct retrieves product information by product ID
func (c *ProductClient) GetProduct(ctx context.Context, productID uint32) (*pb.GetProductResponse, error) {
	req := &pb.GetProductRequest{
		ProductId: productID,
	}

	resp, err := c.client.GetProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return resp, nil
}

// CheckStock checks if product has enough stock
func (c *ProductClient) CheckStock(ctx context.Context, productID uint32, quantity int32) (*pb.CheckStockResponse, error) {
	req := &pb.CheckStockRequest{
		ProductId: productID,
		Quantity:  quantity,
	}

	resp, err := c.client.CheckStock(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check stock: %w", err)
	}

	return resp, nil
}

// UpdateStock updates product stock
func (c *ProductClient) UpdateStock(ctx context.Context, productID uint32, quantity int32) (*pb.UpdateStockResponse, error) {
	req := &pb.UpdateStockRequest{
		ProductId: productID,
		Quantity:  quantity,
	}

	resp, err := c.client.UpdateStock(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	return resp, nil
}

// Close closes the gRPC connection
func (c *ProductClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}