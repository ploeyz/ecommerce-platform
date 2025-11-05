package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ploezy/ecommerce-platform/order-service/config"
	"github.com/ploezy/ecommerce-platform/order-service/internal/grpc/client"
	"github.com/ploezy/ecommerce-platform/order-service/pkg/database"
	"github.com/ploezy/ecommerce-platform/order-service/pkg/kafka"
	"github.com/ploezy/ecommerce-platform/order-service/pkg/redis"
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
	fmt.Println("\n Connecting to database...")
	err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf(" Database connection failed: %v", err)
	}

	// Run migrations
	fmt.Println()
	err = database.AutoMigrate()
	if err != nil{
		log.Fatalf("Database migration failed: %v",err)
	}
	
	// Connect to Redis
	fmt.Println("\n Connecting to Redis...")
	err = redis.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf(" Redis connection failed: %v", err)
	}
	
	// Initialize Kafka Producer
	fmt.Println("\n Initializing Kafka producer...")
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf(" Kafka producer initialization failed: %v", err)
	}
	defer producer.Close()

	fmt.Println("\n All connections successful!")

	fmt.Println("\n⏳ Connecting to User Service gRPC...")
	userClient, err := client.NewUserClient(cfg.UserServiceGRPCURL)
	if err != nil {
		log.Printf("⚠️ User Service gRPC connection failed: %v", err)
		log.Println("⚠️ User Service gRPC server may not be running yet")
	} else {
		defer userClient.Close()
	}
	



	// Test Kafka by sending a test event
	fmt.Println("\n Testing Kafka producer...")
	testEvent := kafka.OrderCreatedEvent{
		OrderID:     1,
		UserID:      1,
		TotalAmount: 1500.00,
		Status:      "pending",
		Items: []kafka.OrderItemEvent{
			{
				ProductID: 1,
				Quantity:  2,
				Price:     500.00,
				Subtotal:  1000.00,
			},
			{
				ProductID: 2,
				Quantity:  1,
				Price:     500.00,
				Subtotal:  500.00,
			},
		},
		CreatedAt: time.Now(),
	}

	err = producer.PublishEvent(cfg.KafkaTopicOrderCreated, "order-1", testEvent)
	if err != nil {
		log.Printf(" Failed to publish test event: %v", err)
	} else {
		fmt.Println(" Test event published successfully!")
	}

	fmt.Println("\n Setup completed!")
}