package main

import (
	"log"

	"github.com/ploezy/ecommerce-platform/user-service/config"
	"github.com/ploezy/ecommerce-platform/user-service/internal/model"
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
	log.Println("User Service Started Successfully")
}