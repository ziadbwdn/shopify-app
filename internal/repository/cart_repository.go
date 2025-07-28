package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"shopify-app/internal/contract"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
)

// cartRepository implements the contract.CartRepository interface
type cartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new instance of the cart repository
func NewCartRepository(db *gorm.DB) contract.CartRepository {
	return &cartRepository{db: db}
}

// GetOrCreateCart retrieves or creates a cart for a user
func (r *cartRepository) GetOrCreateCart(ctx context.Context, userID utils.BinaryUUID) (*entities.Cart, *exception.AppError) {
	var cart entities.Cart
	if err := r.db.WithContext(ctx).Where(entities.Cart{UserID: userID}).FirstOrCreate(&cart).Error; err != nil {
		return nil, exception.NewAppError(err, "failed to get or create cart")
	}
	return &cart, nil
}

// GetCartWithItems retrieves a cart with all its items for a user
func (r *cartRepository) GetCartWithItems(ctx context.Context, userID utils.BinaryUUID) (*entities.Cart, *exception.AppError) {
	var cart entities.Cart
	err := r.db.WithContext(ctx).
		Preload("CartItems").
		Preload("CartItems.Menu").
		Where("user_id = ?", userID).
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no cart, create one and return it (it will be empty)
			return r.GetOrCreateCart(ctx, userID)
		}
		return nil, exception.NewAppError(err, "failed to get cart with items")
	}
	return &cart, nil
}

// AddItemToCart adds an item to the cart or updates quantity if it exists
func (r *cartRepository) AddItemToCart(ctx context.Context, cartID, menuID utils.BinaryUUID, quantity int, price *utils.GormDecimal) *exception.AppError {
	// Check if the item already exists in the cart
	var existingItem entities.CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND menu_id = ?", cartID, menuID).
		First(&existingItem).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Item does not exist, create a new one
			newItem := entities.CartItem{
				CartID:   cartID,
				MenuID:   menuID,
				Quantity: quantity,
				Price:    price,
			}
			if createErr := r.db.WithContext(ctx).Create(&newItem).Error; createErr != nil {
				return exception.NewAppError(createErr, "failed to add new item to cart")
			}
			return nil
		}
		// Any other error
		return exception.NewAppError(err, "failed to check for existing cart item")
	}

	// Item exists, update its quantity
	newQuantity := existingItem.Quantity + quantity
	if updateErr := r.db.WithContext(ctx).Model(&existingItem).Update("quantity", newQuantity).Error; updateErr != nil {
		return exception.NewAppError(updateErr, "failed to update item quantity in cart")
	}

	return nil
}

// UpdateCartItemQuantity updates the quantity of a specific cart item
func (r *cartRepository) UpdateCartItemQuantity(ctx context.Context, cartItemID utils.BinaryUUID, quantity int) *exception.AppError {
	if err := r.db.WithContext(ctx).Model(&entities.CartItem{}).Where("id = ?", cartItemID).Update("quantity", quantity).Error; err != nil {
		return exception.NewAppError(err, "failed to update cart item quantity")
	}
	return nil
}

// RemoveCartItem removes a specific item from the cart
func (r *cartRepository) RemoveCartItem(ctx context.Context, cartItemID utils.BinaryUUID) *exception.AppError {
	if err := r.db.WithContext(ctx).Delete(&entities.CartItem{}, "id = ?", cartItemID).Error; err != nil {
		return exception.NewAppError(err, "failed to remove cart item")
	}
	return nil
}

// ClearCart removes all items from a user's cart
func (r *cartRepository) ClearCart(ctx context.Context, userID utils.BinaryUUID) *exception.AppError {
	// First, get the cart ID for the user
	cart, err := r.GetOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}

	// Delete all items associated with that cart ID
	if deleteErr := r.db.WithContext(ctx).Where("cart_id = ?", cart.ID).Delete(&entities.CartItem{}).Error; deleteErr != nil {
		return exception.NewAppError(deleteErr, "failed to clear cart")
	}
	return nil
}

// GetCartItem retrieves a specific cart item
func (r *cartRepository) GetCartItem(ctx context.Context, cartItemID utils.BinaryUUID) (*entities.CartItem, *exception.AppError) {
	var item entities.CartItem
	if err := r.db.WithContext(ctx).First(&item, "id = ?", cartItemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(err, "cart item not found")
		}
		return nil, exception.NewAppError(err, "failed to get cart item")
	}
	return &item, nil
}

// GetCartItemByMenuID retrieves a cart item by cart ID and menu ID
func (r *cartRepository) GetCartItemByMenuID(ctx context.Context, cartID, menuID utils.BinaryUUID) (*entities.CartItem, *exception.AppError) {
	var item entities.CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND menu_id = ?", cartID, menuID).
		First(&item).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if not found, as it's a check, not an error
		}
		return nil, exception.NewAppError(err, "failed to get cart item by menu id")
	}
	return &item, nil
}

// GetCartTotal calculates the total amount for a cart
func (r *cartRepository) GetCartTotal(ctx context.Context, cartID utils.BinaryUUID) (*utils.GormDecimal, *exception.AppError) {
	var items []entities.CartItem
	if err := r.db.WithContext(ctx).Where("cart_id = ?", cartID).Find(&items).Error; err != nil {
		return nil, exception.NewAppError(err, "failed to get cart items for total calculation")
	}

	total := utils.NewGormDecimal("0")
	for _, item := range items {
		subtotal := item.GetSubtotal()
		total = utils.MustGormDecimalAdd(*total, *subtotal)
	}

	return total, nil
}

// ValidateCartOwnership validates that a cart item belongs to a specific user
func (r *cartRepository) ValidateCartOwnership(ctx context.Context, cartItemID, userID utils.BinaryUUID) (bool, *exception.AppError) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.CartItem{}).
		Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("cart_items.id = ? AND carts.user_id = ?", cartItemID, userID).
		Count(&count).Error

	if err != nil {
		return false, exception.NewAppError(err, "failed to validate cart ownership")
	}

	return count > 0, nil
}
