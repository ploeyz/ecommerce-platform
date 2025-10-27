package models

import (
	"time"

	"gorm.io/gorm"
)

// Order status constants
const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusShipped    = "shipped"
	OrderStatusDelivered  = "delivered"
	OrderStatusCancelled  = "cancelled"
)

// Order represents an order in the system
type Order struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	TotalAmount float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status      string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Items       []OrderItem    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
// TableName specifies the table name for Order model
func (Order) TableName() string {
	return "orders"
}

// IsValidStatus checks if the status is valid
func IsValidStatus(status string) bool {
	validStatuses := []string{
		OrderStatusPending,
		OrderStatusProcessing,
		OrderStatusShipped,
		OrderStatusDelivered,
		OrderStatusCancelled,
	}

	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// CanTransitionTo checks if status can transition to newStatus
func (o *Order) CanTransitionTo(newStatus string) bool {
	// Define valid transitions
	validTransitions := map[string][]string{
		OrderStatusPending: {
			OrderStatusProcessing,
			OrderStatusCancelled,
		},
		OrderStatusProcessing: {
			OrderStatusShipped,
			OrderStatusCancelled,
		},
		OrderStatusShipped: {
			OrderStatusDelivered,
		},
		OrderStatusDelivered: {},
		OrderStatusCancelled: {},
	}

	allowedStatuses, exists := validTransitions[o.Status]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}
	return false
}

// CanBeCancelled checks if order can be cancelled
func (o *Order) CanBeCancelled() bool {
	return o.Status == OrderStatusPending
}