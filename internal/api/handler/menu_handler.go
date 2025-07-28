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

type MenuHandler struct {
	menuService contract.MenuService
}

func NewMenuHandler(menuService contract.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var req dto.CreateMenuRequest
	if err := gin_helper.BindAndValidate(c, &req); err != nil {
		web_response.HandleError(c, err)
		return
	}
	price, _ := utils.Float64ToGormDecimal(req.Price)
	menu, err := h.menuService.AddMenu(c.Request.Context(), req.Name, req.Description, price, req.Category, req.Stock, req.ImageURL)
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	web_response.Success(c, menu)
}

func (h *MenuHandler) GetMenus(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	category := c.Query("category")
	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active_only", "true"))

	menus, count, err := h.menuService.GetMenus(c.Request.Context(), offset, limit, search, category, activeOnly)
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	web_response.Success(c, gin.H{"menus": menus, "count": count})
}

func (h *MenuHandler) GetMenuByID(c *gin.Context) {
	id, err := utils.UUIDFromParam(c, "id")
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	menu, appErr := h.menuService.GetMenuByID(c.Request.Context(), id)
	if appErr != nil {
		web_response.HandleError(c, appErr)
		return
	}
	web_response.Success(c, menu)
}

func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id, err := utils.UUIDFromParam(c, "id")
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	var req dto.UpdateMenuRequest
	if err := gin_helper.BindAndValidate(c, &req); err != nil {
		web_response.HandleError(c, err)
		return
	}
	price, _ := utils.Float64ToGormDecimal(req.Price)
	menu, appErr := h.menuService.UpdateMenu(c.Request.Context(), id, req.Name, req.Description, price, req.Category, req.Stock, req.ImageURL)
	if appErr != nil {
		web_response.HandleError(c, appErr)
		return
	}
	web_response.Success(c, menu)
}

func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id, err := utils.UUIDFromParam(c, "id")
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	appErr := h.menuService.DeleteMenu(c.Request.Context(), id)
	if appErr != nil {
		web_response.HandleError(c, appErr)
		return
	}
	web_response.Success(c, "menu deleted successfully")
}
