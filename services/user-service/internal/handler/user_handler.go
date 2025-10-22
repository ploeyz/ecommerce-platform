package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ploezy/ecommerce-platform/user-service/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler{
	return &UserHandler{service: service}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

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