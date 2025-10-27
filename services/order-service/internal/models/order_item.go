package models

import (
	"time"

	"gorm.io/gorm"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderID   uint           `gorm:"not null;index" json:"order_id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Subtotal  float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for OrderItem model
func (OrderItem) TableName() string {
	return "order_items"
}

// CalculateSubtotal calculates the subtotal (price * quantity)
func (oi *OrderItem) CalculateSubtotal() {
	oi.Subtotal = oi.Price * float64(oi.Quantity)
}