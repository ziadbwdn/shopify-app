package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"shopify-app/internal/contract"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
	"time"
)

// orderRepository implements the contract.OrderRepository interface
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new instance of the order repository
func NewOrderRepository(db *gorm.DB) contract.OrderRepository {
	return &orderRepository{db: db}
}

// CreateOrder creates a new order in the database
func (r *orderRepository) CreateOrder(ctx context.Context, order *entities.Order) *exception.AppError {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return exception.NewAppError(err, "failed to create order")
	}
	return nil
}

// CreateOrderItems creates order items in batch
func (r *orderRepository) CreateOrderItems(ctx context.Context, orderItems []entities.OrderItem) *exception.AppError {
	if err := r.db.WithContext(ctx).Create(&orderItems).Error; err != nil {
		return exception.NewAppError(err, "failed to create order items")
	}
	return nil
}

// GetOrderByID retrieves an order by its ID
func (r *orderRepository) GetOrderByID(ctx context.Context, id utils.BinaryUUID) (*entities.Order, *exception.AppError) {
	var order entities.Order
	if err := r.db.WithContext(ctx).First(&order, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(err, "order not found")
		}
		return nil, exception.NewAppError(err, "failed to get order by id")
	}
	return &order, nil
}

// GetOrderWithItems retrieves an order with all its items
func (r *orderRepository) GetOrderWithItems(ctx context.Context, id utils.BinaryUUID) (*entities.Order, *exception.AppError) {
	var order entities.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems").
		Preload("OrderItems.Menu").
		First(&order, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(err, "order not found")
		}
		return nil, exception.NewAppError(err, "failed to get order with items")
	}
	return &order, nil
}

// GetOrdersByUserID retrieves orders for a specific user with pagination
func (r *orderRepository) GetOrdersByUserID(ctx context.Context, userID utils.BinaryUUID, offset, limit int) ([]entities.Order, int64, *exception.AppError) {
	var orders []entities.Order
	var count int64

	query := r.db.WithContext(ctx).Model(&entities.Order{}).Where("user_id = ?", userID)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to count orders by user id")
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to get orders by user id")
	}

	return orders, count, nil
}

// GetAllOrders retrieves all orders with pagination (admin operation)
func (r *orderRepository) GetAllOrders(ctx context.Context, offset, limit int, status entities.OrderStatus) ([]entities.Order, int64, *exception.AppError) {
	var orders []entities.Order
	var count int64

	query := r.db.WithContext(ctx).Model(&entities.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to count all orders")
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to get all orders")
	}

	return orders, count, nil
}

// UpdateOrderStatus updates the status of an order
func (r *orderRepository) UpdateOrderStatus(ctx context.Context, id utils.BinaryUUID, status entities.OrderStatus) *exception.AppError {
	if err := r.db.WithContext(ctx).Model(&entities.Order{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return exception.NewAppError(err, "failed to update order status")
	}
	return nil
}

// GetOrdersByDateRange retrieves orders within a date range
func (r *orderRepository) GetOrdersByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entities.Order, *exception.AppError) {
	var orders []entities.Order
	err := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, exception.NewAppError(err, "failed to get orders by date range")
	}
	return orders, nil
}

// GetOrdersByStatus retrieves orders by status
func (r *orderRepository) GetOrdersByStatus(ctx context.Context, status entities.OrderStatus, offset, limit int) ([]entities.Order, int64, *exception.AppError) {
	var orders []entities.Order
	var count int64

	query := r.db.WithContext(ctx).Model(&entities.Order{}).Where("status = ?", status)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to count orders by status")
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to get orders by status")
	}

	return orders, count, nil
}

// ValidateOrderOwnership validates that an order belongs to a specific user
func (r *orderRepository) ValidateOrderOwnership(ctx context.Context, orderID, userID utils.BinaryUUID) (bool, *exception.AppError) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Order{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		Count(&count).Error

	if err != nil {
		return false, exception.NewAppError(err, "failed to validate order ownership")
	}

	return count > 0, nil
}
