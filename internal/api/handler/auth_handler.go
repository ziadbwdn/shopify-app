package handler

import (
	"shopify-app/internal/api/dto"
	"shopify-app/internal/contract"
	"shopify-app/internal/entities"
	"shopify-app/pkg/gin_helper"
	"shopify-app/pkg/web_response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService contract.UserService
}

func NewAuthHandler(userService contract.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := gin_helper.BindAndValidate(c, &req); err != nil {
		web_response.HandleError(c, err)
		return
	}

	role := entities.RoleCustomer
	if req.Role == string(entities.RoleAdmin) {
		role = entities.RoleAdmin
	}

	user, token, err := h.userService.Register(c.Request.Context(), req.Email, req.Password, role)
	if err != nil {
		web_response.HandleError(c, err)
		return
	}

	web_response.Success(c, gin.H{"user": user, "token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := gin_helper.BindAndValidate(c, &req); err != nil {
		web_response.HandleError(c, err)
		return
	}

	user, token, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		web_response.HandleError(c, err)
		return
	}

	web_response.Success(c, gin.H{"user": user, "token": token})
}
