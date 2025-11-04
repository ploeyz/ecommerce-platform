package main

import (
	"log"
	"os"
	"os/signal"
	"github.com/ploezy/ecommerce-platform/product-service/config"
	grpcHandler "github.com/ploezy/ecommerce-platform/product-service/internal/grpc/handler"
	grpcServer "github.com/ploezy/ecommerce-platform/product-service/internal/grpc/server"
	"github.com/ploezy/ecommerce-platform/product-service/internal/handler"
	"github.com/ploezy/ecommerce-platform/product-service/internal/middleware"
	"github.com/ploezy/ecommerce-platform/product-service/internal/repository"
	"github.com/ploezy/ecommerce-platform/product-service/internal/service"
	"github.com/ploezy/ecommerce-platform/product-service/pkg/auth"
	"github.com/ploezy/ecommerce-platform/product-service/pkg/database"
	"github.com/ploezy/ecommerce-platform/product-service/pkg/redis"
	"syscall"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Configuration loaded successfully")

	// Connect to PostgreSQL
	db, err := database.ConnectPostgres(database.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run Auto Migration
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Connect to Redis
	redisClient, err := redis.ConnectRedis(redis.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Printf("Redis connected: %v\n", redisClient != nil)

	// Initialize JWT helper
	jwtHelper := auth.NewJWTHelper(cfg.JWT.Secret)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtHelper)

	// Initialize cache service
	cacheService := redis.NewCacheService(redisClient)

	// Initialize layers
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo, cacheService)

	// HTTP Handler
	httpHandler := handler.NewProductHandler(productService)

	// gRPC Handler
	grpcProductHandler := grpcHandler.NewProductGRPCHandler(productService)

	// Start gRPC Server in goroutine
	grpcSrv := grpcServer.NewGRPCServer(grpcProductHandler)
	go func() {
		if err := grpcSrv.Start(cfg.Server.GRPCPort); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Setup HTTP router
	router := handler.SetupRouter(httpHandler, authMiddleware)

	// Start HTTP Server in goroutine
	go func() {
		serverAddr := ":" + cfg.Server.Port
		log.Printf("HTTP Server is running on http://localhost%s\n", serverAddr)
		log.Println("REST API Documentation:")
		log.Println("   PUBLIC ROUTES:")
		log.Println("   GET    /health")
		log.Println("   GET    /api/v1/products")
		log.Println("   GET    /api/v1/products/:id")
		log.Println("   GET    /api/v1/products/search?keyword=xxx")
		log.Println("   PROTECTED ROUTES (Admin Only):")
		log.Println("   POST   /api/v1/products")
		log.Println("   PUT    /api/v1/products/:id")
		log.Println("   DELETE /api/v1/products/:id")

		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
	grpcSrv.Stop()
	log.Println("Servers stopped gracefully")
}