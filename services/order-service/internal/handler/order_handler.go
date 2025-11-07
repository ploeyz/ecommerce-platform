package handler

import (
    "net/http"
    "strconv"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/ploezy/ecommerce-platform/order-service/internal/models"
    "github.com/ploezy/ecommerce-platform/order-service/internal/service"
)

type OrderHandler struct {
    service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
    return &OrderHandler{
        service: service,
    }
}

// getUserIDFromContext extracts user ID from JWT context
func (h *OrderHandler) getUserIDFromContext(c *gin.Context) (uint, error) {
    userID, exists := c.Get("user_id")
    if !exists {
        return 0, gin.Error{
            Err:  http.ErrAbortHandler,
            Meta: "user_id not found in context",
        }
    }
    
    uid, ok := userID.(uint)
    if !ok {
        return 0, gin.Error{
            Err:  http.ErrAbortHandler,
            Meta: "invalid user_id type",
        }
    }
    
    return uid, nil
}

// contains checks if a string contains a substring
func contains(str, substr string) bool {
    return strings.Contains(str, substr)
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with items
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Order data"
// @Success 201 {object} map[string]interface{} "Order created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    // ดึง user ID จาก JWT context
    userID, err := h.getUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "unauthorized",
        })
        return
    }
    
    // Parse request body
    var req models.CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid request body",
            "details": err.Error(),
        })
        return
    }
    
    // Validate request
    if len(req.Items) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "order must have at least one item",
        })
        return
    }
    

    order, err := h.service.CreateOrder(c.Request.Context(), userID, &req)
    if err != nil {
        // Check error type for appropriate status code
        statusCode := http.StatusInternalServerError
        errorMessage := err.Error()
        
        // Handle specific errors
        if contains(errorMessage, "insufficient stock") {
            statusCode = http.StatusBadRequest
        } else if contains(errorMessage, "not found") {
            statusCode = http.StatusNotFound
        }
        
        c.JSON(statusCode, gin.H{
            "error": errorMessage,
        })
        return
    }
    
    // Return success response
    c.JSON(http.StatusCreated, gin.H{
        "message": "order created successfully",
        "data": order,
    })
}

// GetOrders godoc
// @Summary Get user orders
// @Description Get all orders for the authenticated user with pagination
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "Orders retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
    // ดึง user ID จาก JWT context
    userID, err := h.getUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "unauthorized",
        })
        return
    }
    
    // Parse pagination parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    // Call service to get orders
    orders, total, err := h.service.GetUserOrders(c.Request.Context(), userID, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    // Calculate pagination metadata
    totalPages := (int(total) + limit - 1) / limit
    
    // Return success response
    c.JSON(http.StatusOK, gin.H{
        "message": "orders retrieved successfully",
        "data": orders,
        "pagination": gin.H{
            "current_page": page,
            "per_page":     limit,
            "total":        total,
            "total_pages":  totalPages,
        },
    })
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Description Get a specific order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]interface{} "Order retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid order ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
    // ดึง user ID จาก JWT context
    userID, err := h.getUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "unauthorized",
        })
        return
    }
    
    // Parse order ID from URL parameter
    orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid order id",
        })
        return
    }
    
    // Call service to get order
    order, err := h.service.GetOrderByID(c.Request.Context(), uint(orderID), userID)
    if err != nil {
        statusCode := http.StatusInternalServerError
        errorMessage := err.Error()
        
        // Handle specific errors
        if contains(errorMessage, "not found") {
            statusCode = http.StatusNotFound
        } else if contains(errorMessage, "unauthorized") {
            statusCode = http.StatusForbidden
        }
        
        c.JSON(statusCode, gin.H{
            "error": errorMessage,
        })
        return
    }
    
    // Return success response
    c.JSON(http.StatusOK, gin.H{
        "message": "order retrieved successfully",
        "data": order,
    })
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an order (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body map[string]string true "New status"
// @Success 200 {object} map[string]interface{} "Order status updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /admin/orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
    // Parse order ID from URL parameter
    orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid order id",
        })
        return
    }
    
    // Parse request body
    var req struct {
        Status string `json:"status" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid request body",
            "details": err.Error(),
        })
        return
    }
    
    // Call service to update status
    err = h.service.UpdateOrderStatus(c.Request.Context(), uint(orderID), req.Status)
    if err != nil {
        statusCode := http.StatusInternalServerError
        errorMessage := err.Error()
        
        // Handle specific errors
        if contains(errorMessage, "not found") {
            statusCode = http.StatusNotFound
        } else if contains(errorMessage, "invalid") || contains(errorMessage, "cannot change") {
            statusCode = http.StatusBadRequest
        }
        
        c.JSON(statusCode, gin.H{
            "error": errorMessage,
        })
        return
    }
    
    // Return success response
    c.JSON(http.StatusOK, gin.H{
        "message": "order status updated successfully",
    })
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel a pending order and restore product stock
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]interface{} "Order cancelled successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
    // ดึง user ID จาก JWT context
    userID, err := h.getUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "unauthorized",
        })
        return
    }
    
    // Parse order ID from URL parameter
    orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid order id",
        })
        return
    }
    
    // Call service to cancel order
    err = h.service.CancelOrder(c.Request.Context(), uint(orderID), userID)
    if err != nil {
        statusCode := http.StatusInternalServerError
        errorMessage := err.Error()
        
        // Handle specific errors
        if contains(errorMessage, "not found") {
            statusCode = http.StatusNotFound
        } else if contains(errorMessage, "unauthorized") {
            statusCode = http.StatusForbidden
        } else if contains(errorMessage, "cannot cancel") {
            statusCode = http.StatusBadRequest
        }
        
        c.JSON(statusCode, gin.H{
            "error": errorMessage,
        })
        return
    }
    
    // Return success response
    c.JSON(http.StatusOK, gin.H{
        "message": "order cancelled successfully",
    })
}