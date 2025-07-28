package dto

import "shopify-app/internal/entities"

// UpdateOrderStatusRequest defines the request body for updating an order's status
type UpdateOrderStatusRequest struct {
	Status entities.OrderStatus `json:"status" validate:"required,oneof=pending confirmed preparing ready delivered cancelled"`
}
