package service

import (
	"context"
	"shopify-app/internal/contract"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
	"time"
)

type orderService struct {
	orderRepo contract.OrderRepository
	cartSvc   contract.CartService
	menuRepo  contract.MenuRepository
}

func NewOrderService(orderRepo contract.OrderRepository, cartSvc contract.CartService, menuRepo contract.MenuRepository) contract.OrderService {
	return &orderService{orderRepo: orderRepo, cartSvc: cartSvc, menuRepo: menuRepo}
}

func (s *orderService) CheckoutCart(ctx context.Context, userID utils.BinaryUUID) (*entities.Order, *exception.AppError) {
	cart, total, err := s.cartSvc.ValidateCartForCheckout(ctx, userID)
	if err != nil {
		return nil, err
	}

	// In a real app, this whole block should be a single database transaction
	order := &entities.Order{
		UserID:      userID,
		TotalAmount: total,
		Status:      entities.StatusPending,
	}

	if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	var orderItems []entities.OrderItem
	var stockReduction = make(map[utils.BinaryUUID]int)
	for _, cartItem := range cart.CartItems {
		orderItems = append(orderItems, entities.OrderItem{
			OrderID:  order.ID,
			MenuID:   cartItem.MenuID,
			Quantity: cartItem.Quantity,
			Price:    cartItem.Price,
			MenuName: cartItem.Menu.Name,
		})
		stockReduction[cartItem.MenuID] = cartItem.Quantity
	}

	if err := s.orderRepo.CreateOrderItems(ctx, orderItems); err != nil {
		// Here you would ideally roll back the order creation
		return nil, err
	}

	// Reduce stock
	for menuID, quantity := range stockReduction {
		if err := s.menuRepo.ReduceMenuStock(ctx, menuID, quantity); err != nil {
			// And here you would roll back everything
			return nil, exception.NewAppError(err, "failed to reduce stock during checkout")
		}
	}

	// Clear the user's cart
	if err := s.cartSvc.ClearCart(ctx, userID); err != nil {
		// This is less critical, maybe just log it
	}

	return s.orderRepo.GetOrderWithItems(ctx, order.ID)
}

func (s *orderService) GetOrderHistory(ctx context.Context, userID utils.BinaryUUID, offset, limit int) ([]entities.Order, int64, *exception.AppError) {
	return s.orderRepo.GetOrdersByUserID(ctx, userID, offset, limit)
}

func (s *orderService) GetOrderDetails(ctx context.Context, userID, orderID utils.BinaryUUID) (*entities.Order, *exception.AppError) {
	owned, err := s.orderRepo.ValidateOrderOwnership(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, exception.NewAppError(nil, "order not found or not owned by user")
	}
	return s.orderRepo.GetOrderWithItems(ctx, orderID)
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID utils.BinaryUUID, status entities.OrderStatus) (*entities.Order, *exception.AppError) {
	if _, err := s.orderRepo.GetOrderByID(ctx, orderID); err != nil {
		return nil, err
	}

	if err := s.orderRepo.UpdateOrderStatus(ctx, orderID, status); err != nil {
		return nil, err
	}
	return s.orderRepo.GetOrderWithItems(ctx, orderID)
}

func (s *orderService) CancelOrder(ctx context.Context, userID, orderID utils.BinaryUUID) (*entities.Order, *exception.AppError) {
	owned, err := s.orderRepo.ValidateOrderOwnership(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, exception.NewAppError(nil, "order not found or not owned by user")
	}

	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if !order.CanBeCancelled() {
		return nil, exception.NewAppError(nil, "order cannot be cancelled in its current state")
	}

	return s.UpdateOrderStatus(ctx, orderID, entities.StatusCancelled)
}

func (s *orderService) GetAllOrders(ctx context.Context, offset, limit int, status entities.OrderStatus) ([]entities.Order, int64, *exception.AppError) {
	return s.orderRepo.GetAllOrders(ctx, offset, limit, status)
}

func (s *orderService) GetOrdersByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entities.Order, *exception.AppError) {
	return s.orderRepo.GetOrdersByDateRange(ctx, startDate, endDate)
}

func (s *orderService) ValidateOrderAccess(ctx context.Context, userID, orderID utils.BinaryUUID, userRole entities.UserRole) (bool, *exception.AppError) {
	if userRole == entities.RoleAdmin {
		return true, nil
	}
	return s.orderRepo.ValidateOrderOwnership(ctx, orderID, userID)
}
