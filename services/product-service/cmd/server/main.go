package main

import (
	"context"
	"log"

	"github.com/ploezy/ecommerce-platform/product-service/config"
	"github.com/ploezy/ecommerce-platform/product-service/internal/model"
	"github.com/ploezy/ecommerce-platform/product-service/internal/repository"
	"github.com/ploezy/ecommerce-platform/product-service/pkg/database"
	"github.com/ploezy/ecommerce-platform/product-service/pkg/redis"
	"gorm.io/gorm"
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

	// Run Auto Migration
	if err := database.AutoMigrate(db); err != nil{
		log.Fatalf("Failed to migrate database: %v",err)
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
	log.Println("Redis connection:", redisClient != nil)
	testRepository(db)
	log.Println("All connections successful!")
}
func testRepository(db *gorm.DB) {
	ctx := context.Background()
	repo := repository.NewProductRepository(db)

	// Test 1: Create Product
	log.Println("Test 1: Create Product")
	product := &model.Product{
		Name:        "iPhone 15 Pro",
		Description: "Latest Apple smartphone with A17 Pro chip",
		Price:       39900.00,
		Stock:       50,
		Category:    "Electronics",
		Images:      []string{"iphone15-1.jpg", "iphone15-2.jpg"},
	}

	if err := repo.Create(ctx, product); err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}
	log.Printf("Product created successfully with ID: %d", product.ID)

	// Test 2: Find By ID
	log.Println("Test 2: Find Product by ID")
	foundProduct, err := repo.FindByID(ctx, product.ID)
	if err != nil {
		log.Fatalf("Failed to find product: %v", err)
	}
	log.Printf("Product found: %s (Price: %.2f THB)", foundProduct.Name, foundProduct.Price)
	
	// Test 3: Create more products
	log.Println("Test 3: Create more products")
	products := []model.Product{
		{
			Name:        "MacBook Pro 16",
			Description: "Powerful laptop for professionals",
			Price:       89900.00,
			Stock:       30,
			Category:    "Electronics",
			Images:      []string{"macbook-1.jpg"},
		},
		{
			Name:        "AirPods Pro",
			Description: "Wireless earbuds with noise cancellation",
			Price:       8900.00,
			Stock:       100,
			Category:    "Electronics",
			Images:      []string{"airpods-1.jpg"},
		},
	}

	for _, p := range products {
		prod := p
		if err := repo.Create(ctx, &prod); err != nil {
			log.Printf("Failed to create product: %v", err)
		} else {
			log.Printf("Created: %s", prod.Name)
		}
	}

	// Test 4: Find All with Pagination
	log.Println("Test 4: Find All Products (Pagination)")
	allProducts, total, err := repo.FindAll(ctx, 0, 10)
	if err != nil {
		log.Fatalf("Failed to find all products: %v", err)
	}
	log.Printf("Found %d products (Total: %d)", len(allProducts), total)
	for _, p := range allProducts {
		log.Printf("   - %s: %.2f THB (Stock: %d)", p.Name, p.Price, p.Stock)
	}

	// Test 5: Search Products
	log.Println("Test 5: Search Products")
	searchResults, searchTotal, err := repo.Search(ctx, "Pro", 0, 10)
	if err != nil {
		log.Fatalf("Failed to search products: %v", err)
	}
	log.Printf("Search results for 'Pro': %d products found", searchTotal)
	for _, p := range searchResults {
		log.Printf("   - %s: %.2f THB", p.Name, p.Price)
	}

	// Test 6: Update Product
	log.Println("Test 6: Update Product")
	foundProduct.Price = 37900.00
	foundProduct.Stock = 45
	if err := repo.Update(ctx, foundProduct); err != nil {
		log.Fatalf("Failed to update product: %v", err)
	}
	log.Printf("✅ Product updated: New price = %.2f THB, New stock = %d", foundProduct.Price, foundProduct.Stock)

	// Test 7: Delete Product
	log.Println("Test 7: Delete Product")
	if err := repo.Delete(ctx, foundProduct.ID); err != nil {
		log.Fatalf("Failed to delete product: %v", err)
	}
	log.Printf("Product deleted (soft delete) with ID: %d", foundProduct.ID)

	// Verify deletion
	_, err = repo.FindByID(ctx, foundProduct.ID)
	if err != nil {
		log.Printf("Confirmed: Product not found after deletion (soft delete working)")
	}
}