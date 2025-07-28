package handler

import (
	"shopify-app/internal/api/dto"
	"shopify-app/internal/contract"
	"shopify-app/internal/utils"
	"shopify-app/pkg/gin_helper"
	"shopify-app/pkg/web_response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService contract.OrderService
}

func NewOrderHandler(orderService contract.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	userID, _ := c.Get("userID")
	order, err := h.orderService.CheckoutCart(c.Request.Context(), userID.(utils.BinaryUUID))
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	web_response.Success(c, order)
}

func (h *OrderHandler) GetOrderHistory(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userID, _ := c.Get("userID")
	orders, count, err := h.orderService.GetOrderHistory(c.Request.Context(), userID.(utils.BinaryUUID), offset, limit)
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	web_response.Success(c, gin.H{"orders": orders, "count": count})
}

func (h *OrderHandler) GetOrderDetails(c *gin.Context) {
	id, err := utils.UUIDFromParam(c, "id")
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	userID, _ := c.Get("userID")
	order, appErr := h.orderService.GetOrderDetails(c.Request.Context(), userID.(utils.BinaryUUID), id)
	if appErr != nil {
		web_response.HandleError(c, appErr)
		return
	}
	web_response.Success(c, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := utils.UUIDFromParam(c, "id")
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	var req dto.UpdateOrderStatusRequest
	if err := gin_helper.BindAndValidate(c, &req); err != nil {
		web_response.HandleError(c, err)
		return
	}
	order, appErr := h.orderService.UpdateOrderStatus(c.Request.Context(), id, req.Status)
	if appErr != nil {
		web_response.HandleError(c, appErr)
		return
	}
	web_response.Success(c, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := utils.UUIDFromParam(c, "id")
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	userID, _ := c.Get("userID")
	order, appErr := h.orderService.CancelOrder(c.Request.Context(), userID.(utils.BinaryUUID), id)
	if appErr != nil {
		web_response.HandleError(c, appErr)
		return
	}
	web_response.Success(c, order)
}
