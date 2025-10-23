package repository

import (
	"context"

	"github.com/ploezy/ecommerce-platform/product-service/internal/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	FindByID(ctx context.Context, id uint) (*model.Product, error)
	FindAll(ctx context.Context, offset, limit int) ([]model.Product, int64, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, keyword string, offset, limit int) ([]model.Product, int64, error)
}