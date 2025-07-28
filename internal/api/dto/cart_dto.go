package dto

import "shopify-app/internal/utils"

// AddToCartRequest defines the request body for adding an item to the cart
type AddToCartRequest struct {
	MenuID   utils.BinaryUUID `json:"menu_id" validate:"required"`
	Quantity int              `json:"quantity" validate:"required,gt=0"`
}

// UpdateCartItemRequest defines the request body for updating a cart item's quantity
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,gte=0"`
}
