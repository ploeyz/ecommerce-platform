package main

import (
	"log"

	"github.com/ploezy/ecommerce-platform/product-service/config"
	"github.com/ploezy/ecommerce-platform/product-service/pkg/database"
	"github.com/ploezy/ecommerce-platform/product-service/pkg/redis"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("✅ Configuration loaded successfully")
	log.Printf("Server will run on port: %s", cfg.Server.Port)

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
	log.Println("Database connection:", db != nil)

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
	log.Println("Redis connection:", redisClient != nil)

	log.Println("All connections successful!")
}