package model

import "github.com/lib/pq"

// CreateProductRequest is the request for creating a product
type CreateProductRequest struct {
	Name        string         `json:"name" binding:"required" example:"iPhone 15 Pro Max"`
	Description string         `json:"description" example:"Latest Apple flagship smartphone"`
	Price       float64        `json:"price" binding:"required,gt=0" example:"45900"`
	Stock       int            `json:"stock" binding:"required,gte=0" example:"50"`
	Category    string         `json:"category" binding:"required" example:"Electronics"`
	Images      pq.StringArray `json:"images" swaggertype:"array,string" example:"image1.jpg,image2.jpg"`
}

// UpdateProductRequest is the request for updating a product
type UpdateProductRequest struct {
	Name        string         `json:"name" example:"iPhone 15 Pro Max"`
	Description string         `json:"description" example:"Updated description"`
	Price       float64        `json:"price" binding:"omitempty,gt=0" example:"43900"`
	Stock       int            `json:"stock" binding:"omitempty,gte=0" example:"45"`
	Category    string         `json:"category" example:"Electronics"`
	Images      pq.StringArray `json:"images" swaggertype:"array,string" example:"image1.jpg,image2.jpg"`
}

// ProductResponse is the response for product
type ProductResponse struct {
	ID          uint           `json:"id" example:"1"`
	Name        string         `json:"name" example:"iPhone 15 Pro Max"`
	Description string         `json:"description" example:"Latest Apple flagship smartphone"`
	Price       float64        `json:"price" example:"45900"`
	Stock       int            `json:"stock" example:"50"`
	Category    string         `json:"category" example:"Electronics"`
	Images      pq.StringArray `json:"images" swaggertype:"array,string" example:"image1.jpg,image2.jpg"`
	CreatedAt   string         `json:"created_at" example:"2025-11-07 15:30:00"`
	UpdatedAt   string         `json:"updated_at" example:"2025-11-07 15:30:00"`
}

// PaginationResponse is the response for paginated data
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total" example:"100"`
	Page       int         `json:"page" example:"1"`
	Limit      int         `json:"limit" example:"10"`
	TotalPages int         `json:"total_pages" example:"10"`
}

// SearchRequest is the request for searching products
type SearchRequest struct {
	Keyword string `json:"keyword" form:"keyword" example:"iPhone"`
	Page    int    `json:"page" form:"page" example:"1"`
	Limit   int    `json:"limit" form:"limit" example:"10"`
}