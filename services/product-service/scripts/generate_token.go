package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func main() {
	secretKey := "your-super-secret-jwt-key-change-this-in-production"

	// Admin token
	adminClaims := JWTClaims{
		UserID: 1,
		Email:  "admin@example.com",
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	adminToken := jwt.NewWithClaims(jwt.SigningMethodHS256, adminClaims)
	adminTokenString, _ := adminToken.SignedString([]byte(secretKey))

	fmt.Println("Admin Token (valid for 24 hours):")
	fmt.Println(adminTokenString)
	fmt.Println()

	// User token (not admin)
	userClaims := JWTClaims{
		UserID: 2,
		Email:  "user@example.com",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	userTokenString, _ := userToken.SignedString([]byte(secretKey))

	fmt.Println("User Token (valid for 24 hours):")
	fmt.Println(userTokenString)
}