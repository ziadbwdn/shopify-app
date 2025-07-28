package dto

import (
	"shopify-app/internal/entities"
	"shopify-app/internal/utils"
)

// User DTOs
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"oneof=admin customer"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}

// Menu DTOs
type CreateMenuRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Category    string  `json:"category" validate:"required"`
	Stock       int     `json:"stock" validate:"gte=0"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
}

type UpdateMenuRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Category    string  `json:"category" validate:"required"`
	Stock       int     `json:"stock" validate:"gte=0"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
}

// Cart DTOs
type AddToCartRequest struct {
	MenuID   utils.BinaryUUID `json:"menu_id" validate:"required"`
	Quantity int              `json:"quantity" validate:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,gte=0"`
}

// Order DTOs
type UpdateOrderStatusRequest struct {
	Status entities.OrderStatus `json:"status" validate:"required,oneof=pending confirmed preparing ready delivered cancelled"`
}
