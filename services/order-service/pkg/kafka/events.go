package kafka

import "time"

// OrderCreatedEvent represents an order creation event
type OrderCreatedEvent struct {
	OrderID     uint      `json:"order_id"`
	UserID      uint      `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	Items       []OrderItemEvent `json:"items"`
	CreatedAt   time.Time `json:"created_at"`
}

// OrderItemEvent represents an order item in the event
type OrderItemEvent struct {
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

// OrderStatusChangedEvent represents an order status change event
type OrderStatusChangedEvent struct {
	OrderID   uint      `json:"order_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderCancelledEvent represents an order cancellation event
type OrderCancelledEvent struct {
	OrderID     uint      `json:"order_id"`
	UserID      uint      `json:"user_id"`
	Reason      string    `json:"reason"`
	CancelledAt time.Time `json:"cancelled_at"`
}