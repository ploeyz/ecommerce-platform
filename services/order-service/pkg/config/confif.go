package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string
	GRPCPort   string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Kafka
	KafkaBrokers                 string
	KafkaTopicOrderCreated       string
	KafkaTopicOrderStatusChanged string
	KafkaTopicOrderCancelled     string

	// gRPC Services
	UserServiceGRPCURL    string
	ProductServiceGRPCURL string

	// JWT
	JWTSecret string
}

func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		// Server
		ServerPort: getEnv("SERVER_PORT", "8083"),
		GRPCPort:   getEnv("GRPC_PORT", "50051"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "ecom_user"),
		DBPassword: getEnv("DB_PASSWORD", "ecom_pass"),
		DBName:     getEnv("DB_NAME", "ecom_db"),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       2, // Default DB 0

		// Kafka
		KafkaBrokers:                 getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopicOrderCreated:       getEnv("KAFKA_TOPIC_ORDER_CREATED", "order.created"),
		KafkaTopicOrderStatusChanged: getEnv("KAFKA_TOPIC_ORDER_STATUS_CHANGED", "order.status_changed"),
		KafkaTopicOrderCancelled:     getEnv("KAFKA_TOPIC_ORDER_CANCELLED", "order.cancelled"),

		// gRPC Services
		UserServiceGRPCURL:    getEnv("USER_SERVICE_GRPC_URL", "localhost:50052"),
		ProductServiceGRPCURL: getEnv("PRODUCT_SERVICE_GRPC_URL", "localhost:50053"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "your-super-secret-key"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
