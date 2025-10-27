package main

import (
	"fmt"
	"log"

	"github.com/ploezy/ecommerce-platform/order-service/pkg/config"
	"github.com/ploezy/ecommerce-platform/order-service/pkg/database"
)

// Print config to verify

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Print config to verify
	fmt.Println("===== Order Service Configuration =====")
	fmt.Printf("Server Port: %s\n", cfg.ServerPort)
	fmt.Printf("gRPC Port: %s\n", cfg.GRPCPort)
	fmt.Printf("Database: %s@%s:%s/%s\n", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
	fmt.Printf("Redis: %s:%s (DB: %d)\n", cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	fmt.Printf("Kafka Brokers: %s\n", cfg.KafkaBrokers)
	fmt.Printf("User Service gRPC: %s\n", cfg.UserServiceGRPCURL)
	fmt.Printf("Product Service gRPC: %s\n", cfg.ProductServiceGRPCURL)
	fmt.Println("========================================")

	// Connect to database
	fmt.Println("\nConnecting to database...")
	err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	fmt.Println("\nAll connections successful!")
}