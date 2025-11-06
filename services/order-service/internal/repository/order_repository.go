package repository

import (
	"context"

	"github.com/ploezy/ecommerce-platform/order-service/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
    Create(ctx context.Context, order *models.Order) error
    FindByID(ctx context.Context, id uint) (*models.Order, error)
    FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]models.Order, int64, error)
    Update(ctx context.Context, order *models.Order) error
    UpdateStatus(ctx context.Context, orderID uint, status string) error
}

type orderRepository struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
    return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
    return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) FindByID(ctx context.Context, id uint) (*models.Order, error){
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&order, id ).Error

	if err != nil{
		return nil, err
	}
	return &order,nil
}

func (r *orderRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]models.Order , int64 , error ){
	var orders []models.Order
    var total int64
    
    // นับจำนวน orders ทั้งหมดของ user
    if err := r.db.WithContext(ctx).
        Model(&models.Order{}).
        Where("user_id = ?", userID).
        Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // ดึง orders พร้อม pagination
    err := r.db.WithContext(ctx).
        Preload("Items").
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&orders).Error
    
    if err != nil {
        return nil, 0, err
    }
    
    return orders, total, nil
}

func (r *orderRepository) Update(ctx context.Context, order *models.Order) error {
    return r.db.WithContext(ctx).Save(order).Error
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID uint, status string) error {
    return r.db.WithContext(ctx).
        Model(&models.Order{}).
        Where("id = ?", orderID).
        Update("status", status).Error
}