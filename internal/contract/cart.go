// internal/contract/cart_contract.go
package contract

import (
	"context"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
)

// CartRepository defines the contract for cart data access operations
type CartRepository interface {
	// GetOrCreateCart retrieves or creates a cart for a user
	GetOrCreateCart(ctx context.Context, userID utils.BinaryUUID) (*entities.Cart, *exception.AppError)
	
	// GetCartWithItems retrieves a cart with all its items for a user
	GetCartWithItems(ctx context.Context, userID utils.BinaryUUID) (*entities.Cart, *exception.AppError)
	
	// AddItemToCart adds an item to the cart or updates quantity if exists
	AddItemToCart(ctx context.Context, cartID, menuID utils.BinaryUUID, quantity int, price *utils.GormDecimal) *exception.AppError
	
	// UpdateCartItemQuantity updates the quantity of a specific cart item
	UpdateCartItemQuantity(ctx context.Context, cartItemID utils.BinaryUUID, quantity int) *exception.AppError
	
	// RemoveCartItem removes a specific item from the cart
	RemoveCartItem(ctx context.Context, cartItemID utils.BinaryUUID) *exception.AppError
	
	// ClearCart removes all items from a user's cart
	ClearCart(ctx context.Context, userID utils.BinaryUUID) *exception.AppError
	
	// GetCartItem retrieves a specific cart item
	GetCartItem(ctx context.Context, cartItemID utils.BinaryUUID) (*entities.CartItem, *exception.AppError)
	
	// GetCartItemByMenuID retrieves a cart item by cart ID and menu ID
	GetCartItemByMenuID(ctx context.Context, cartID, menuID utils.BinaryUUID) (*entities.CartItem, *exception.AppError)
	
	// GetCartTotal calculates the total amount for a cart
	GetCartTotal(ctx context.Context, cartID utils.BinaryUUID) (*utils.GormDecimal, *exception.AppError)
	
	// ValidateCartOwnership validates that a cart item belongs to a specific user
	ValidateCartOwnership(ctx context.Context, cartItemID, userID utils.BinaryUUID) (bool, *exception.AppError)
}

// CartService defines the contract for cart business logic operations
type CartService interface {
	// AddItemToCart handles adding an item to cart with validation
	AddItemToCart(ctx context.Context, userID, menuID utils.BinaryUUID, quantity int) *exception.AppError
	
	// GetUserCart retrieves user's cart with all items and calculations
	GetUserCart(ctx context.Context, userID utils.BinaryUUID) (*entities.Cart, *utils.GormDecimal, *exception.AppError)
	
	// UpdateCartItem handles updating cart item quantity with validation
	UpdateCartItem(ctx context.Context, userID, cartItemID utils.BinaryUUID, quantity int) *exception.AppError
	
	// RemoveCartItem handles removing an item from cart with validation
	RemoveCartItem(ctx context.Context, userID, cartItemID utils.BinaryUUID) *exception.AppError
	
	// ClearCart clears all items from user's cart
	ClearCart(ctx context.Context, userID utils.BinaryUUID) *exception.AppError
	
	// ValidateCartForCheckout validates cart items before checkout
	ValidateCartForCheckout(ctx context.Context, userID utils.BinaryUUID) (*entities.Cart, *utils.GormDecimal, *exception.AppError)
	
	// GetCartItemCount returns the total number of items in user's cart
	GetCartItemCount(ctx context.Context, userID utils.BinaryUUID) (int, *exception.AppError)
	
	// SyncCartItemPrices updates cart item prices with current menu prices
	SyncCartItemPrices(ctx context.Context, userID utils.BinaryUUID) *exception.AppError
}