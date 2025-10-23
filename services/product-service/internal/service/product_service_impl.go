package service

import (
	"context"
	"errors"
	"math"

	"github.com/ploezy/ecommerce-platform/product-service/internal/model"
	"github.com/ploezy/ecommerce-platform/product-service/internal/repository"
	"gorm.io/gorm"
)

type productService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// CreateProduct creates a new product
func (s *productService) CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.ProductResponse, error) {
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Images:      req.Images,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// GetProductByID gets a product by ID
func (s *productService) GetProductByID(ctx context.Context,id uint) (*model.ProductResponse, error){
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return s.toProductResponse(product),nil
}

// GetAllProducts gets all products with pagination
func (s *productService) GetAllProducts(ctx context.Context, page, limit int) (*model.PaginationResponse, error){
	//Set default value
	if page < 1{
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, total, err := s.repo.FindAll(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response
	var productResponses []model.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, *s.toProductResponse(&p))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &model.PaginationResponse{
		Data:       productResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateProduct updates a product
func (s *productService) UpdateProduct(ctx context.Context, id uint, req *model.UpdateProductRequest) (*model.ProductResponse, error) {
	// Find existing product
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if len(req.Images) > 0 {
		product.Images = req.Images
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// DeleteProduct deletes a product
func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	// Check if product exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

// SearchProducts searches products by keyword
func (s *productService) SearchProducts(ctx context.Context, keyword string, page, limit int) (*model.PaginationResponse, error) {
	// Set default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, total, err := s.repo.Search(ctx, keyword, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response
	var productResponses []model.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, *s.toProductResponse(&p))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &model.PaginationResponse{
		Data:       productResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// Helper function to convert Product to ProductResponse
func (s *productService) toProductResponse(product *model.Product) *model.ProductResponse {
	return &model.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		Images:      product.Images,
		CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}