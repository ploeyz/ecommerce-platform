package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ploezy/ecommerce-platform/user-service/config"
	"github.com/ploezy/ecommerce-platform/user-service/internal/handler"
	"github.com/ploezy/ecommerce-platform/user-service/internal/model"
	"github.com/ploezy/ecommerce-platform/user-service/internal/repository"
	"github.com/ploezy/ecommerce-platform/user-service/internal/service"
	"github.com/ploezy/ecommerce-platform/user-service/pkg/database"
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
	userhandler := handler.NewUserHandler(userService)
	
	// Setup Gin

	r := gin.Default()

	// Routes

	api := r.Group("/api/v1")
	{
		api.POST("/register", userhandler.Register)
	}
	// Start server
	log.Printf("✅ User Service running on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}