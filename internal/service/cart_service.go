package service

import (
	"context"
	"shopify-app/internal/contract"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
)

type cartService struct {
	cartRepo contract.CartRepository
	menuRepo contract.MenuRepository
}

func NewCartService(cartRepo contract.CartRepository, menuRepo contract.MenuRepository) contract.CartService {
	return &cartService{cartRepo: cartRepo, menuRepo: menuRepo}
}

func (s *cartService) AddItemToCart(ctx context.Context, userID, menuID utils.BinaryUUID, quantity int) *exception.AppError {
	if quantity <= 0 {
		return exception.NewAppError(nil, "quantity must be positive")
	}

	menu, err := s.menuRepo.GetMenuByID(ctx, menuID)
	if err != nil {
		return err
	}

	if !menu.IsInStock(quantity) {
		return exception.NewAppError(nil, "item is out of stock")
	}

	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}

	return s.cartRepo.AddItemToCart(ctx, cart.ID, menuID, quantity, menu.Price)
}

func (s *cartService) GetUserCart(ctx context.Context, userID utils.BinaryUUID) (*entities.Cart, *utils.GormDecimal, *exception.AppError) {
	cart, err := s.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	total, err := s.cartRepo.GetCartTotal(ctx, cart.ID)
	if err != nil {
		return nil, nil, err
	}

	return cart, total, nil
}

func (s *cartService) UpdateCartItem(ctx context.Context, userID, cartItemID utils.BinaryUUID, quantity int) *exception.AppError {
	if quantity <= 0 {
		return s.RemoveCartItem(ctx, userID, cartItemID)
	}

	owned, err := s.cartRepo.ValidateCartOwnership(ctx, cartItemID, userID)
	if err != nil {
		return err
	}
	if !owned {
		return exception.NewAppError(nil, "cart item not found or not owned by user")
	}

	item, err := s.cartRepo.GetCartItem(ctx, cartItemID)
	if err != nil {
		return err
	}

	menu, err := s.menuRepo.GetMenuByID(ctx, item.MenuID)
	if err != nil {
		return err
	}

	if !menu.IsInStock(quantity) {
		return exception.NewAppError(nil, "not enough stock")
	}

	return s.cartRepo.UpdateCartItemQuantity(ctx, cartItemID, quantity)
}

func (s *cartService) RemoveCartItem(ctx context.Context, userID, cartItemID utils.BinaryUUID) *exception.AppError {
	owned, err := s.cartRepo.ValidateCartOwnership(ctx, cartItemID, userID)
	if err != nil {
		return err
	}
	if !owned {
		return exception.NewAppError(nil, "cart item not found or not owned by user")
	}

	return s.cartRepo.RemoveCartItem(ctx, cartItemID)
}

func (s *cartService) ClearCart(ctx context.Context, userID utils.BinaryUUID) *exception.AppError {
	return s.cartRepo.ClearCart(ctx, userID)
}

func (s *cartService) ValidateCartForCheckout(ctx context.Context, userID utils.BinaryUUID) (*entities.Cart, *utils.GormDecimal, *exception.AppError) {
	cart, total, err := s.GetUserCart(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	if len(cart.CartItems) == 0 {
		return nil, nil, exception.NewAppError(nil, "cart is empty")
	}

	for _, item := range cart.CartItems {
		if !item.Menu.IsInStock(item.Quantity) {
			return nil, nil, exception.NewAppError(nil, "one or more items are out of stock")
		}
	}

	return cart, total, nil
}

func (s *cartService) GetCartItemCount(ctx context.Context, userID utils.BinaryUUID) (int, *exception.AppError) {
	cart, err := s.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		return 0, err
	}
	return len(cart.CartItems), nil
}

func (s *cartService) SyncCartItemPrices(ctx context.Context, userID utils.BinaryUUID) *exception.AppError {
	cart, err := s.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		return err
	}

	for _, item := range cart.CartItems {
		if !utils.GormDecimalEquals(*item.Price, *item.Menu.Price) {
			item.Price = item.Menu.Price
			if err := s.cartRepo.UpdateCartItemQuantity(ctx, item.ID, item.Quantity); err != nil {
				// In a real app, you might want to handle this more gracefully
				return exception.NewAppError(err, "failed to sync cart item price")
			}
		}
	}
	return nil
}
