package handler

import "github.com/gin-gonic/gin"

func SetupRouter(productHandler *ProductHandler) *gin.Engine{
	router := gin.Default()

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
			// Public routes
			products.GET("", productHandler.GetAllProducts)           // GET /api/v1/products
			products.GET("/:id", productHandler.GetProductByID)       // GET /api/v1/products/:id
			products.GET("/search", productHandler.SearchProducts)    // GET /api/v1/products/search

			// Protected routes (will add auth middleware later)
			products.POST("", productHandler.CreateProduct)           // POST /api/v1/products
			products.PUT("/:id", productHandler.UpdateProduct)        // PUT /api/v1/products/:id
			products.DELETE("/:id", productHandler.DeleteProduct)     // DELETE /api/v1/products/:id
		}
	}
	return router
}