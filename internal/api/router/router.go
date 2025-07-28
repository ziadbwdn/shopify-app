package router

import (
	"shopify-app/internal/api/handler"
	"shopify-app/internal/config"
	"shopify-app/internal/contract"
	"shopify-app/internal/entities"
	"shopify-app/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(
	cfg *config.Config,
	userService contract.UserService,
	menuService contract.MenuService,
	cartService contract.CartService,
	orderService contract.OrderService,
	reportService contract.ReportService,
) *gin.Engine {
	r := gin.Default()

	authHandler := handler.NewAuthHandler(userService)
	userHandler := handler.NewUserHandler(userService)
	menuHandler := handler.NewMenuHandler(menuService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)
	reportHandler := handler.NewReportHandler(reportService)

	// Public routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// Authenticated routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		// User routes
		userRoutes := api.Group("/user")
		{
			userRoutes.GET("/profile", userHandler.GetProfile)
			userRoutes.PUT("/profile", userHandler.UpdateProfile)
			userRoutes.POST("/change-password", userHandler.ChangePassword)
		}

		// Menu routes (publicly readable)
		menuRoutes := api.Group("/menus")
		{
			menuRoutes.GET("/", menuHandler.GetMenus)
			menuRoutes.GET("/:id", menuHandler.GetMenuByID)
		}

		// Admin-only menu routes
		adminMenuRoutes := api.Group("/admin/menus")
		adminMenuRoutes.Use(middleware.RoleMiddleware(entities.RoleAdmin))
		{
			adminMenuRoutes.POST("/", menuHandler.CreateMenu)
			adminMenuRoutes.PUT("/:id", menuHandler.UpdateMenu)
			adminMenuRoutes.DELETE("/:id", menuHandler.DeleteMenu)
		}

		// Cart routes
		cartRoutes := api.Group("/cart")
		{
			cartRoutes.GET("/", cartHandler.GetCart)
			cartRoutes.POST("/items", cartHandler.AddToCart)
			cartRoutes.PUT("/items/:id", cartHandler.UpdateCartItem)
			cartRoutes.DELETE("/items/:id", cartHandler.RemoveCartItem)
			cartRoutes.DELETE("/", cartHandler.ClearCart)
		}

		// Order routes
		orderRoutes := api.Group("/orders")
		{
			orderRoutes.POST("/checkout", orderHandler.Checkout)
			orderRoutes.GET("/", orderHandler.GetOrderHistory)
			orderRoutes.GET("/:id", orderHandler.GetOrderDetails)
			orderRoutes.POST("/:id/cancel", orderHandler.CancelOrder)
		}

		// Admin-only order routes
		adminOrderRoutes := api.Group("/admin/orders")
		adminOrderRoutes.Use(middleware.RoleMiddleware(entities.RoleAdmin))
		{
			adminOrderRoutes.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		}

		// Report routes (admin only)
		reportRoutes := api.Group("/reports")
		reportRoutes.Use(middleware.RoleMiddleware(entities.RoleAdmin))
		{
			reportRoutes.GET("/sales", reportHandler.GetSalesReport)
			reportRoutes.GET("/bestsellers", reportHandler.GetBestSellingItems)
		}
	}

	return r
}
