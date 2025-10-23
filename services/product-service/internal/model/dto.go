package model

import "github.com/lib/pq"

// CreateProductRequest is the request for creating a product
type CreateProductRequest struct {
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	Price       float64        `json:"price" binding:"required,gt=0"`
	Stock       int            `json:"stock" binding:"required,gte=0"`
	Category    string         `json:"category" binding:"required"`
	Images      pq.StringArray `json:"images"`
	
}

// UpdateProductRequest is the request for updating a product
type UpdateProductRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price" binding:"omitempty,gt=0"`
	Stock       int            `json:"stock" binding:"omitempty,gte=0"`
	Category    string         `json:"category"`
	Images      pq.StringArray `json:"images"`
}

// ProductResponse is the response for product
type ProductResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Stock       int            `json:"stock"`
	Category    string         `json:"category"`
	Images      pq.StringArray `json:"images"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

// PaginationResponse is the response for paginated data
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// SearchRequest is the request for searching products
type SearchRequest struct {
	Keyword string `json:"keyword" form:"keyword"`
	Page    int    `json:"page" form:"page"`
	Limit   int    `json:"limit" form:"limit"`
}