package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/ploezy/ecommerce-platform/order-service/config"
)

var Client *redis.Client
var Ctx = context.Background()

func ConnectRedis(cfg *config.Config) error {
	// Create Redis client
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test connection
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return nil
}

func GetClient() *redis.Client {
	return Client
}

func GetContext() context.Context {
	return Ctx
}
