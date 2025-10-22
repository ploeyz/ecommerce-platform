package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ploezy/ecommerce-platform/user-service/pkg/auth"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}
		
		token := parts[1]
		// Validate token
		claims, err := auth.ValidateToken(token,jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id",claims.UserID)
		c.Set("email",claims.Email)
		c.Set("role",claims.Role)

		c.Next()
	}
}