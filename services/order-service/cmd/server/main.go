package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ploezy/ecommerce-platform/order-service/config"
	"github.com/ploezy/ecommerce-platform/order-service/internal/grpc/client"
	"github.com/ploezy/ecommerce-platform/order-service/internal/handler"
	"github.com/ploezy/ecommerce-platform/order-service/internal/repository"
	"github.com/ploezy/ecommerce-platform/order-service/internal/service"
	"github.com/ploezy/ecommerce-platform/order-service/pkg/database"
	"github.com/ploezy/ecommerce-platform/order-service/pkg/kafka"
	"github.com/ploezy/ecommerce-platform/order-service/pkg/redis"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	log.Println("===== Order Service Configuration =====")
	log.Printf("Server Port: %s\n", cfg.ServerPort)
	log.Printf("gRPC Port: %s\n", cfg.GRPCPort)
	log.Printf("Database: %s@%s:%s/%s\n", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
	log.Printf("Redis: %s:%s (DB: %d)\n", cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	log.Printf("Kafka Brokers: %s\n", cfg.KafkaBrokers)
	log.Printf("User Service gRPC: %s\n", cfg.UserServiceGRPCURL)
	log.Printf("Product Service gRPC: %s\n", cfg.ProductServiceGRPCURL)
	log.Println("========================================")

	// Connect to database
	log.Println("\nConnecting to database...")
	err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Run migrations
	log.Println("Running database migrations...")
	err = database.AutoMigrate()
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	// Connect to Redis
	log.Println("\nConnecting to Redis...")
	err = redis.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	// Initialize Kafka Producer
	log.Println("\nInitializing Kafka producer...")
	kafkaProducer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf("Kafka producer initialization failed: %v", err)
	}
	defer kafkaProducer.Close()

	// Initialize User Service gRPC Client
	log.Println("\nConnecting to User Service gRPC...")
	userClient, err := client.NewUserClient(cfg.UserServiceGRPCURL)
	if err != nil {
		log.Fatalf("User Service gRPC connection failed: %v", err)
	}
	defer userClient.Close()

	// Initialize Product Service gRPC Client
	log.Println("Connecting to Product Service gRPC...")
	productClient, err := client.NewProductClient(cfg.ProductServiceGRPCURL)
	if err != nil {
		log.Fatalf("Product Service gRPC connection failed: %v", err)
	}
	defer productClient.Close()

	log.Println("\nAll connections successful!")

	// Get database instance
	db := database.GetDB()

	// Initialize layers: Repository -> Service -> Handler
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, db, userClient, productClient, kafkaProducer)
	orderHandler := handler.NewOrderHandler(orderService)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "order-service",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Order routes (Protected - require JWT)
		orders := v1.Group("/orders")
		orders.Use(AuthMiddleware(userClient)) // JWT Middleware
		{
			orders.POST("", orderHandler.CreateOrder)            // Create order
			orders.GET("", orderHandler.GetOrders)               // Get user orders (pagination)
			orders.GET("/:id", orderHandler.GetOrderByID)        // Get order by ID
			orders.POST("/:id/cancel", orderHandler.CancelOrder) // Cancel order
		}

		// Admin routes (Protected - require JWT)
		admin := v1.Group("/admin/orders")
		admin.Use(AuthMiddleware(userClient)) // JWT Middleware
		{
			admin.PUT("/:id/status", orderHandler.UpdateOrderStatus) // Update order status
		}
	}

	// Start server in goroutine
	serverAddr := ":" + cfg.ServerPort
	log.Printf("\nOrder Service starting on %s\n", serverAddr)

	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\nShutting down Order Service...")
	log.Println("Order Service stopped")
}

// AuthMiddleware validates JWT token via User Service
func AuthMiddleware(userClient *client.UserClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Extract token (format: "Bearer <token>")
		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			c.JSON(401, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		// Validate token via User Service gRPC
		resp, err := userClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		if !resp.Valid {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user_id in context for handlers
		c.Set("user_id", uint(resp.UserId))
		c.Next()
	}
}