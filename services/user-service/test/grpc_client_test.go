package test

import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/ploezy/ecommerce-platform/user-service/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestUserServiceGRPC(t *testing.T) {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Test 1: Register
	t.Run("Register User", func(t *testing.T) {
		req := &pb.RegisterRequest{
			Email:     "grpc@example.com",
			Password:  "password123",
			FirstName: "gRPC",
			LastName:  "User",
		}

		resp, err := client.Register(ctx, req)
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}

		log.Printf("✅ Register Response: %+v", resp)
		
		if resp.Email != req.Email {
			t.Errorf("Expected email %s, got %s", req.Email, resp.Email)
		}
	})

	// Test 2: Login
	t.Run("Login User", func(t *testing.T) {
		req := &pb.LoginRequest{
			Email:    "grpc@example.com",
			Password: "password123",
		}

		resp, err := client.Login(ctx, req)
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}

		log.Printf("✅ Login Response: Token=%s, UserID=%d", resp.Token, resp.UserId)
		
		if resp.Token == "" {
			t.Error("Expected token, got empty string")
		}

		// Test 3: Validate Token
		t.Run("Validate Token", func(t *testing.T) {
			validateReq := &pb.ValidateTokenRequest{
				Token: resp.Token,
			}

			validateResp, err := client.ValidateToken(ctx, validateReq)
			if err != nil {
				t.Fatalf("ValidateToken failed: %v", err)
			}

			log.Printf("✅ ValidateToken Response: Valid=%v, UserID=%d", validateResp.Valid, validateResp.UserId)

			if !validateResp.Valid {
				t.Error("Expected token to be valid")
			}
		})

		// Test 4: Get User By ID
		t.Run("Get User By ID", func(t *testing.T) {
			getUserReq := &pb.GetUserByIDRequest{
				Id: resp.UserId,
			}

			getUserResp, err := client.GetUserByID(ctx, getUserReq)
			if err != nil {
				t.Fatalf("GetUserByID failed: %v", err)
			}

			log.Printf("✅ GetUserByID Response: %+v", getUserResp)

			if getUserResp.Email != "grpc@example.com" {
				t.Errorf("Expected email grpc@example.com, got %s", getUserResp.Email)
			}
		})

		// Test 5: Get User By Email
		t.Run("Get User By Email", func(t *testing.T) {
			getUserReq := &pb.GetUserByEmailRequest{
				Email: "grpc@example.com",
			}

			getUserResp, err := client.GetUserByEmail(ctx, getUserReq)
			if err != nil {
				t.Fatalf("GetUserByEmail failed: %v", err)
			}

			log.Printf("✅ GetUserByEmail Response: %+v", getUserResp)

			if getUserResp.Id != resp.UserId {
				t.Errorf("Expected user ID %d, got %d", resp.UserId, getUserResp.Id)
			}
		})
	})
}