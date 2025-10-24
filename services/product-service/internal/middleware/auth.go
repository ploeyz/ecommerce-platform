package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ploezy/ecommerce-platform/product-service/pkg/auth"
)

type AuthMiddleware struct {
	jwtHelper *auth.JWTHelper
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtHelper *auth.JWTHelper) *AuthMiddleware{
	return &AuthMiddleware{jwtHelper: jwtHelper}
}

// Authenticate middleware to verify JWT token
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid authorization format. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := m.jwtHelper.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RequireAdmin middleware to check if user is admin
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "User role not found",
			})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Access denied. Admin role required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}