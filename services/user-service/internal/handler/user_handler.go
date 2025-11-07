package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ploezy/ecommerce-platform/user-service/internal/service"
)

type UserHandler struct {
	service service.UserService
	jwtSecret string
}

func NewUserHandler(service service.UserService,jwtSecret string) *UserHandler{
	return &UserHandler{
		service: service,
		jwtSecret: jwtSecret,
	}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email" example:"user@example.com"`
	Password  string `json:"password" binding:"required,min=6" example:"password123"`
	FirstName string `json:"first_name" binding:"required" example:"John"`
	LastName  string `json:"last_name" binding:"required" example:"Doe"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	user, err := h.service.Register(req.Email,req.Password,req.FirstName,req.LastName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusCreated,gin.H{
		"message" :"User registered sunccessful",
		"user" : user,
	})
}

// Login godoc
// @Summary Login user
// @Description Login with email and password to get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{} "Login successful with token"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /login [post]
func (h *UserHandler) Login (c *gin.Context){
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}
	token,user, err := h.service.Login(req.Email,req.Password,h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest ,gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusOK , gin.H{
		"message":"Login successful",
		"token":   token,
		"user":    user,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile (requires authentication)
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Profile retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /profile [get]
func (h *UserHandler) GetProfile(c *gin.Context){
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"user_id": userID,
		"email":   c.GetString("email"),
		"role":    c.GetString("role"),
	})
}