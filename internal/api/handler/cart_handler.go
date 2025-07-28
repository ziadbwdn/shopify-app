package handler

import (
	"shopify-app/internal/api/dto"
	"shopify-app/internal/contract"
	"shopify-app/internal/utils"
	"shopify-app/pkg/gin_helper"
	"shopify-app/pkg/web_response"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService contract.CartService
}

func NewCartHandler(cartService contract.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	var req dto.AddToCartRequest
	if err := gin_helper.BindAndValidate(c, &req); err != nil {
		web_response.HandleError(c, err)
		return
	}
	userID, _ := c.Get("userID")
	err := h.cartService.AddItemToCart(c.Request.Context(), userID.(utils.BinaryUUID), req.MenuID, req.Quantity)
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	web_response.Success(c, "item added to cart")
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, _ := c.Get("userID")
	cart, total, err := h.cartService.GetUserCart(c.Request.Context(), userID.(utils.BinaryUUID))
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	web_response.Success(c, gin.H{"cart": cart, "total": total})
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	id, err := utils.UUIDFromParam(c, "id")
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	var req dto.UpdateCartItemRequest
	if err := gin_helper.BindAndValidate(c, &req); err != nil {
		web_response.HandleError(c, err)
		return
	}
	userID, _ := c.Get("userID")
	appErr := h.cartService.UpdateCartItem(c.Request.Context(), userID.(utils.BinaryUUID), id, req.Quantity)
	if appErr != nil {
		web_response.HandleError(c, appErr)
		return
	}
	web_response.Success(c, "cart item updated")
}

func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	id, err := utils.UUIDFromParam(c, "id")
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	userID, _ := c.Get("userID")
	appErr := h.cartService.RemoveCartItem(c.Request.Context(), userID.(utils.BinaryUUID), id)
	if appErr != nil {
		web_response.HandleError(c, appErr)
		return
	}
	web_response.Success(c, "cart item removed")
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, _ := c.Get("userID")
	err := h.cartService.ClearCart(c.Request.Context(), userID.(utils.BinaryUUID))
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	web_response.Success(c, "cart cleared")
}
