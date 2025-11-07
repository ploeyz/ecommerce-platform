package handler

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
	"github.com/ploezy/ecommerce-platform/product-service/internal/middleware"
)

func SetupRouter(productHandler *ProductHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine{
	router := gin.Default()

	// Swagger documentation with custom config
	url := ginSwagger.URL("http://localhost:8082/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		url,
		ginSwagger.DefaultModelsExpandDepth(-1), // ซ่อน Models section
	))
	
	// Health check
	router.GET("/health", func(c *gin.Context){
		c.JSON(200, gin.H{
			"status":"ok",
			"service" : "prodouct-service",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			// Public routes (no authentication required)
			products.GET("", productHandler.GetAllProducts)           // GET /api/v1/products
			products.GET("/search", productHandler.SearchProducts)    // GET /api/v1/products/search
			products.GET("/:id", productHandler.GetProductByID)       // GET /api/v1/products/:id

			// Protected routes (authentication + admin role required)
			protected := products.Group("")
			protected.Use(authMiddleware.Authenticate())
			protected.Use(authMiddleware.RequireAdmin())
			{
				protected.POST("", productHandler.CreateProduct)      // POST /api/v1/products
				protected.PUT("/:id", productHandler.UpdateProduct)   // PUT /api/v1/products/:id
				protected.DELETE("/:id", productHandler.DeleteProduct) // DELETE /api/v1/products/:id
			}
		}
	}
	return router
}