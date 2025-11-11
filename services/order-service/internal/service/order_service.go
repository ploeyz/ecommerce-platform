package service

import (
    "context"
    "github.com/ploezy/ecommerce-platform/order-service/internal/models"
    "github.com/ploezy/ecommerce-platform/order-service/internal/repository"
    grpcclient "github.com/ploezy/ecommerce-platform/order-service/internal/grpc/client"
    "github.com/ploezy/ecommerce-platform/order-service/pkg/kafka"
    "errors"
    "fmt"
    "gorm.io/gorm"
)

type OrderService interface {
    CreateOrder(ctx context.Context, userID uint, req *models.CreateOrderRequest) (*models.Order, error)
    GetOrderByID(ctx context.Context, orderID, userID uint) (*models.Order, error)
    GetUserOrders(ctx context.Context, userID uint, page, limit int) ([]models.Order, int64, error)
    UpdateOrderStatus(ctx context.Context, orderID uint, status string) error
    CancelOrder(ctx context.Context, orderID, userID uint) error
}

type orderService struct {
    repo           repository.OrderRepository
    db             *gorm.DB
    userClient     *grpcclient.UserClient
    productClient  *grpcclient.ProductClient
    kafkaProducer  *kafka.Producer
}

func NewOrderService(
    repo repository.OrderRepository,
    db *gorm.DB,
    userClient *grpcclient.UserClient,
    productClient *grpcclient.ProductClient,
    kafkaProducer *kafka.Producer,
) OrderService {
    return &orderService{
        repo:          repo,
        db:            db,
        userClient:    userClient,
        productClient: productClient,
        kafkaProducer: kafkaProducer,
    }
}

// GetOrderByID retrieves an order by ID with authorization check
func (s *orderService) GetOrderByID(ctx context.Context, orderID, userID uint) (*models.Order, error) {
    order, err := s.repo.FindByID(ctx, orderID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("order not found")
        }
        return nil, fmt.Errorf("failed to get order: %w", err)
    }
    
    if order.UserID != userID {
        return nil, errors.New("unauthorized: order does not belong to this user")
    }
    
    return order, nil
}

// GetUserOrders retrieves all orders for a user with pagination
func (s *orderService) GetUserOrders(ctx context.Context, userID uint, page, limit int) ([]models.Order, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }
    
    offset := (page - 1) * limit
    
    orders, total, err := s.repo.FindByUserID(ctx, userID, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get user orders: %w", err)
    }
    
    return orders, total, nil
}

// UpdateOrderStatus updates the status of an order with validation
func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID uint, status string) error {
    validStatuses := map[string]bool{
        models.OrderStatusPending:    true,
        models.OrderStatusProcessing: true,
        models.OrderStatusShipped:    true,
        models.OrderStatusDelivered:  true,
        models.OrderStatusCancelled:  true,
    }
    
    if !validStatuses[status] {
        return errors.New("invalid order status")
    }
    
    order, err := s.repo.FindByID(ctx, orderID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("order not found")
        }
        return fmt.Errorf("failed to get order: %w", err)
    }
    
    if !s.isValidStatusTransition(order.Status, status) {
        return fmt.Errorf("cannot change status from %s to %s", order.Status, status)
    }
    
    if err := s.repo.UpdateStatus(ctx, orderID, status); err != nil {
        return fmt.Errorf("failed to update order status: %w", err)
    }
    
    event := kafka.OrderStatusChangedEvent{
        OrderID:   orderID,
        OldStatus: order.Status,
        NewStatus: status,
    }
    
    if err := s.kafkaProducer.SendOrderStatusChanged(event); err != nil {
        fmt.Printf("Warning: failed to send kafka event: %v\n", err)
    }
    
    return nil
}

// CancelOrder cancels an order and restores product stock
func (s *orderService) CancelOrder(ctx context.Context, orderID, userID uint) error {
    order, err := s.repo.FindByID(ctx, orderID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("order not found")
        }
        return fmt.Errorf("failed to get order: %w", err)
    }
    
    if order.UserID != userID {
        return errors.New("unauthorized: order does not belong to this user")
    }
    
    if order.Status != models.OrderStatusPending {
        return fmt.Errorf("cannot cancel order with status: %s (only pending orders can be cancelled)", order.Status)
    }
    
    if err := s.repo.UpdateStatus(ctx, orderID, models.OrderStatusCancelled); err != nil {
        return fmt.Errorf("failed to cancel order: %w", err)
    }
    
    for _, item := range order.Items {
        _, err := s.productClient.UpdateStock(ctx, uint32(item.ProductID), int32(item.Quantity))
        if err != nil {
            fmt.Printf("Warning: failed to restore stock for product %d: %v\n", item.ProductID, err)
        }
    }
    
    event := kafka.OrderCancelledEvent{
        OrderID: orderID,
        UserID:  userID,
    }
    
    if err := s.kafkaProducer.SendOrderCancelled(event); err != nil {
        fmt.Printf("Warning: failed to send kafka event: %v\n", err)
    }
    
    return nil
}

// CreateOrder creates a new order with validation and stock management
func (s *orderService) CreateOrder(ctx context.Context, userID uint, req *models.CreateOrderRequest) (*models.Order, error) {
    if len(req.Items) == 0 {
        return nil, errors.New("order must have at least one item")
    }
    
    tx := s.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    var totalAmount float64
    orderItems := make([]models.OrderItem, 0, len(req.Items))
    
    for _, item := range req.Items {
        if item.Quantity <= 0 {
            tx.Rollback()
            return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
        }
        
        productResp, err := s.productClient.GetProduct(ctx, uint32(item.ProductID))
        if err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to get product %d: %w", item.ProductID, err)
        }
        
        stockResp, err := s.productClient.CheckStock(ctx, uint32(item.ProductID), int32(item.Quantity))
        if err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to check stock for product %d: %w", item.ProductID, err)
        }
        
        // ใช้ CurrentStock แทน Stock
        if !stockResp.Available {
            tx.Rollback()
            return nil, fmt.Errorf("product %d has insufficient stock (available: %d, requested: %d)", 
                item.ProductID, stockResp.CurrentStock, item.Quantity)
        }
        
        // ใช้ Price ได้เลย
        price := productResp.Product.Price
        subtotal := price * float64(item.Quantity)
        totalAmount += subtotal
        
        orderItem := models.OrderItem{
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
            Price:     price,
            Subtotal:  subtotal,
        }
        orderItems = append(orderItems, orderItem)
    }
    
    order := &models.Order{
        UserID:      userID,
        TotalAmount: totalAmount,
        Status:      models.OrderStatusPending,
        Items:       orderItems,
    }
    
    if err := tx.Create(order).Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to create order: %w", err)
    }
    
    for _, item := range orderItems {
        _, err := s.productClient.UpdateStock(ctx, uint32(item.ProductID), -int32(item.Quantity))
        if err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to update stock for product %d: %w", item.ProductID, err)
        }
    }
    
    if err := tx.Commit().Error; err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    event := kafka.OrderCreatedEvent{
        OrderID:     order.ID,
        UserID:      order.UserID,
        TotalAmount: order.TotalAmount,
    }
    
    if err := s.kafkaProducer.SendOrderCreated(event); err != nil {
        fmt.Printf("Warning: failed to send kafka event: %v\n", err)
    }
    
    return order, nil
}

// isValidStatusTransition checks if status transition is valid
func (s *orderService) isValidStatusTransition(oldStatus, newStatus string) bool {
    if oldStatus == newStatus {
        return false
    }
    
    validTransitions := map[string][]string{
        models.OrderStatusPending:    {models.OrderStatusProcessing, models.OrderStatusCancelled},
        models.OrderStatusProcessing: {models.OrderStatusShipped},
        models.OrderStatusShipped:    {models.OrderStatusDelivered},
        models.OrderStatusDelivered:  {},
        models.OrderStatusCancelled:  {},
    }
    
    allowedStatuses := validTransitions[oldStatus]
    for _, allowed := range allowedStatuses {
        if newStatus == allowed {
            return true
        }
    }
    
    return false
}