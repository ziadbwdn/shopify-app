// internal/contract/order_contract.go
package contract

import (
	"context"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
	"time"
)

// OrderRepository defines the contract for order data access operations
type OrderRepository interface {
	// CreateOrder creates a new order in the database
	CreateOrder(ctx context.Context, order *entities.Order) *exception.AppError
	
	// CreateOrderItems creates order items in batch
	CreateOrderItems(ctx context.Context, orderItems []entities.OrderItem) *exception.AppError
	
	// GetOrderByID retrieves an order by its ID
	GetOrderByID(ctx context.Context, id utils.BinaryUUID) (*entities.Order, *exception.AppError)
	
	// GetOrderWithItems retrieves an order with all its items
	GetOrderWithItems(ctx context.Context, id utils.BinaryUUID) (*entities.Order, *exception.AppError)
	
	// GetOrdersByUserID retrieves orders for a specific user with pagination
	GetOrdersByUserID(ctx context.Context, userID utils.BinaryUUID, offset, limit int) ([]entities.Order, int64, *exception.AppError)
	
	// GetAllOrders retrieves all orders with pagination (admin operation)
	GetAllOrders(ctx context.Context, offset, limit int, status entities.OrderStatus) ([]entities.Order, int64, *exception.AppError)
	
	// UpdateOrderStatus updates the status of an order
	UpdateOrderStatus(ctx context.Context, id utils.BinaryUUID, status entities.OrderStatus) *exception.AppError
	
	// GetOrdersByDateRange retrieves orders within a date range
	GetOrdersByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entities.Order, *exception.AppError)
	
	// GetOrdersByStatus retrieves orders by status
	GetOrdersByStatus(ctx context.Context, status entities.OrderStatus, offset, limit int) ([]entities.Order, int64, *exception.AppError)
	
	// ValidateOrderOwnership validates that an order belongs to a specific user
	ValidateOrderOwnership(ctx context.Context, orderID, userID utils.BinaryUUID) (bool, *exception.AppError)
}

// OrderService defines the contract for order business logic operations
type OrderService interface {
	// CheckoutCart processes cart checkout and creates an order
	CheckoutCart(ctx context.Context, userID utils.BinaryUUID) (*entities.Order, *exception.AppError)
	
	// GetOrderHistory retrieves user's order history with pagination
	GetOrderHistory(ctx context.Context, userID utils.BinaryUUID, offset, limit int) ([]entities.Order, int64, *exception.AppError)
	
	// GetOrderDetails retrieves detailed information about a specific order
	GetOrderDetails(ctx context.Context, userID, orderID utils.BinaryUUID) (*entities.Order, *exception.AppError)
	
	// UpdateOrderStatus updates order status with validation (admin operation)
	UpdateOrderStatus(ctx context.Context, orderID utils.BinaryUUID, status entities.OrderStatus) (*entities.Order, *exception.AppError)
	
	// CancelOrder cancels an order if possible
	CancelOrder(ctx context.Context, userID, orderID utils.BinaryUUID) (*entities.Order, *exception.AppError)
	
	// GetAllOrders retrieves all orders for admin with filtering
	GetAllOrders(ctx context.Context, offset, limit int, status entities.OrderStatus) ([]entities.Order, int64, *exception.AppError)
	
	// GetOrdersByDateRange retrieves orders within date range (admin operation)
	GetOrdersByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entities.Order, *exception.AppError)
	
	// ValidateOrderAccess validates user access to order
	ValidateOrderAccess(ctx context.Context, userID, orderID utils.BinaryUUID, userRole entities.UserRole) (bool, *exception.AppError)
}