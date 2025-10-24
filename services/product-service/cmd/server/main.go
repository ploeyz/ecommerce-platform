package main

import (
	"context"
	"log"

	"github.com/lib/pq"
	"github.com/ploezy/ecommerce-platform/product-service/config"
	"github.com/ploezy/ecommerce-platform/product-service/internal/handler"
	"github.com/ploezy/ecommerce-platform/product-service/internal/middleware"
	"github.com/ploezy/ecommerce-platform/product-service/internal/model"
	"github.com/ploezy/ecommerce-platform/product-service/internal/repository"
	"github.com/ploezy/ecommerce-platform/product-service/internal/service"
	"github.com/ploezy/ecommerce-platform/product-service/pkg/auth"
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
	//testRepository(db)
	//testService(db)

	// Initialize JWT helper
	jwtHelper := auth.NewJWTHelper(cfg.JWT.Secret)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtHelper)
	// Initialize cache service
	cacheService := redis.NewCacheService(redisClient)

	// Initialize router
	productRepo 	:= repository.NewProductRepository(db)
	productService 	:= service.NewProductService(productRepo, cacheService)
	productHandler 	:= handler.NewProductHandler(productService)

	// Setup router
	router := handler.SetupRouter(productHandler,authMiddleware)

	// Start server
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Product Service is running on http://localhost%s\n", serverAddr)
	log.Println("API Documentation:")
	log.Println("   GET    /health")
	log.Println("   GET    /api/v1/products")
	log.Println("   GET    /api/v1/products/:id")
	log.Println("   GET    /api/v1/products/search?keyword=xxx")
	log.Println("   POST   /api/v1/products")
	log.Println("   PUT    /api/v1/products/:id")
	log.Println("   DELETE /api/v1/products/:id")

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Println("All connections successful!")
}

func testService(db *gorm.DB,  cache *redis.CacheService) {
	ctx := context.Background()

	// Create repository and service
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo, cache)

	// Test 1: Create Product
	log.Println("Test 1: Create Product via Service")
	createReq := &model.CreateProductRequest{
		Name:        "Samsung Galaxy S24 Ultra",
		Description: "Latest flagship smartphone from Samsung",
		Price:       44900.00,
		Stock:       75,
		Category:    "Electronics",
		Images:      pq.StringArray{"s24-1.jpg", "s24-2.jpg"},
	}

	product, err := svc.CreateProduct(ctx, createReq)
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}
	log.Printf("Product created: ID=%d, Name=%s, Price=%.2f THB", product.ID, product.Name, product.Price)

	// Test 2: Get Product by ID
	log.Println("Test 2: Get Product by ID")
	foundProduct, err := svc.GetProductByID(ctx, product.ID)
	if err != nil {
		log.Fatalf("Failed to get product: %v", err)
	}
	log.Printf("Product found: %s (Stock: %d)", foundProduct.Name, foundProduct.Stock)

	// Test 3: Create more products
	log.Println("Test 3: Create more products")
	products := []model.CreateProductRequest{
		{
			Name:        "Sony WH-1000XM5",
			Description: "Premium noise-canceling headphones",
			Price:       13900.00,
			Stock:       40,
			Category:    "Electronics",
			Images:      pq.StringArray{"sony-1.jpg"},
		},
		{
			Name:        "iPad Pro 12.9",
			Description: "Powerful tablet with M2 chip",
			Price:       45900.00,
			Stock:       25,
			Category:    "Electronics",
			Images:      pq.StringArray{"ipad-1.jpg", "ipad-2.jpg"},
		},
	}

	for _, req := range products {
		r := req
		p, err := svc.CreateProduct(ctx, &r)
		if err != nil {
			log.Printf("Failed to create product: %v", err)
		} else {
			log.Printf("Created: %s (%.2f THB)", p.Name, p.Price)
		}
	}

	// Test 4: Get All Products with Pagination
	log.Println("Test 4: Get All Products (Page 1, Limit 2)")
	allProducts, err := svc.GetAllProducts(ctx, 1, 2)
	if err != nil {
		log.Fatalf("Failed to get all products: %v", err)
	}
	log.Printf("Total: %d, Page: %d/%d, Showing: %d products",
		allProducts.Total, allProducts.Page, allProducts.TotalPages, allProducts.Limit)

	if products, ok := allProducts.Data.([]model.ProductResponse); ok {
		for _, p := range products {
			log.Printf("   - %s: %.2f THB", p.Name, p.Price)
		}
	}

	// Test 5: Search Products
	log.Println("Test 5: Search Products (keyword: 'Pro')")
	searchResults, err := svc.SearchProducts(ctx, "Pro", 1, 10)
	if err != nil {
		log.Fatalf("Failed to search products: %v", err)
	}
	log.Printf("Found %d products", searchResults.Total)

	if products, ok := searchResults.Data.([]model.ProductResponse); ok {
		for _, p := range products {
			log.Printf("   - %s: %.2f THB", p.Name, p.Price)
		}
	}

	// Test 6: Update Product
	log.Println("Test 6: Update Product")
	updateReq := &model.UpdateProductRequest{
		Price: 42900.00,
		Stock: 80,
	}

	updatedProduct, err := svc.UpdateProduct(ctx, product.ID, updateReq)
	if err != nil {
		log.Fatalf("Failed to update product: %v", err)
	}
	log.Printf("Product updated: %s - New Price: %.2f THB, New Stock: %d",
		updatedProduct.Name, updatedProduct.Price, updatedProduct.Stock)

	// Test 7: Delete Product
	log.Println("Test 7: Delete Product")
	if err := svc.DeleteProduct(ctx, product.ID); err != nil {
		log.Fatalf("Failed to delete product: %v", err)
	}
	log.Printf("Product deleted successfully")

	// Verify deletion
	_, err = svc.GetProductByID(ctx, product.ID)
	if err != nil {
		log.Printf("Confirmed: %v", err)
	}
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