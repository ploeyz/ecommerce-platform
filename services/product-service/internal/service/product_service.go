package service

import (
	"context"

	"github.com/ploezy/ecommerce-platform/product-service/internal/model"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.ProductResponse, error)
	GetProductByID(ctx context.Context, id uint) (*model.ProductResponse, error)
	GetAllProducts(ctx context.Context, page, limit int) (*model.PaginationResponse, error)
	UpdateProduct(ctx context.Context, id uint, req *model.UpdateProductRequest) (*model.ProductResponse, error)
	DeleteProduct(ctx context.Context, id uint) error
	SearchProducts(ctx context.Context, keyword string, page, limit int) (*model.PaginationResponse, error)
}