package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name" binding:"required"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null" json:"price" binding:"required,gt=0"`
	Stock       int            `gorm:"not null;default:0" json:"stock" binding:"required,gte=0"`
	Category    string         `gorm:"size:100" json:"category" binding:"required"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Product model

func (Product) TableName() string {
	return "products"
}
