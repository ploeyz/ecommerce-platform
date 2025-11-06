package main

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	pb "github.com/ploezy/ecommerce-platform/user-service/proto/user"
	usergrpc "github.com/ploezy/ecommerce-platform/user-service/internal/grpc"
	"github.com/ploezy/ecommerce-platform/user-service/config"
	"github.com/ploezy/ecommerce-platform/user-service/internal/handler"
	"github.com/ploezy/ecommerce-platform/user-service/internal/middleware"
	"github.com/ploezy/ecommerce-platform/user-service/internal/model"
	"github.com/ploezy/ecommerce-platform/user-service/internal/repository"
	"github.com/ploezy/ecommerce-platform/user-service/internal/service"
	"github.com/ploezy/ecommerce-platform/user-service/pkg/database"
	"google.golang.org/grpc"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	//Connect to database
	db,err := database.NewPostgresDB(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	//Auto migrate
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, cfg.JWTSecret)
	
	// Start gRPC Server in goroutine
	go startGRPCServer(userService, cfg.JWTSecret, cfg.GRPCPort)

	// Start REST API Server
	startRESTServer(userHandler, cfg)
}
func startGRPCServer(userService service.UserService, jwtSecret string, grpcPort string) {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, usergrpc.NewUserGRPCServer(userService, jwtSecret))

	log.Printf("gRPC Server running on port %s", grpcPort) 
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
func startRESTServer(userHandler *handler.UserHandler, cfg *config.Config) {
	r := gin.Default()

	// Public Routes
	api := r.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
	}

	// Protected Routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/profile", userHandler.GetProfile)
	}

	log.Printf("REST API Server running on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start REST server:", err)
	}
}