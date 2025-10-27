package models

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
}

// CreateOrderItemRequest represents an item in the create order request
type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required,min=1"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// UpdateOrderStatusRequest represents the request to update order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// CancelOrderRequest represents the request to cancel an order
type CancelOrderRequest struct {
	Reason string `json:"reason"`
}

// OrderResponse represents the response for an order
type OrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	TotalAmount float64             `json:"total_amount"`
	Status      string              `json:"status"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}

// OrderItemResponse represents an item in the order response
type OrderItemResponse struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

// PaginationQuery represents pagination parameters
type PaginationQuery struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

// PaginatedOrdersResponse represents paginated orders response
type PaginatedOrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}