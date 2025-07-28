package handler

import (
	"shopify-app/internal/api/dto"
	"shopify-app/internal/contract"
	"shopify-app/internal/utils"
	"shopify-app/pkg/gin_helper"
	"shopify-app/pkg/web_response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService contract.UserService
}

func NewUserHandler(userService contract.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.userService.GetUserProfile(c.Request.Context(), userID.(utils.BinaryUUID))
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	web_response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := gin_helper.BindAndValidate(c, &req); err != nil {
		web_response.HandleError(c, err)
		return
	}
	userID, _ := c.Get("userID")
	user, err := h.userService.UpdateUserProfile(c.Request.Context(), userID.(utils.BinaryUUID), req.Email)
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	web_response.Success(c, user)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := gin_helper.BindAndValidate(c, &req); err != nil {
		web_response.HandleError(c, err)
		return
	}
	userID, _ := c.Get("userID")
	err := h.userService.ChangePassword(c.Request.Context(), userID.(utils.BinaryUUID), req.CurrentPassword, req.NewPassword)
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	web_response.Success(c, "password changed successfully")
}
