package repository

import (
	"context"

	"github.com/ploezy/ecommerce-platform/product-service/internal/model"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create creates a new product
func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// FindByID finds a product by ID
func (r *productRepository) FindByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindAll finds all products with pagination
func (r *productRepository) FindAll(ctx context.Context, offset, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Update updates a product
func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete soft deletes a product
func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, id).Error
}

// Search searches products by name or category
func (r *productRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("name ILIKE ? OR category ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}