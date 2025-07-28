package middleware

import (
	"shopify-app/internal/config"
	"shopify-app/internal/entities"
	"shopify-app/pkg/jwt"
	"shopify-app/pkg/web_response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			web_response.HandleError(c, web_response.NewUnauthorizedError("missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			web_response.HandleError(c, web_response.NewUnauthorizedError("invalid authorization header format"))
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(parts[1], cfg.JWTSecret)
		if err != nil {
			web_response.HandleError(c, web_response.NewUnauthorizedError(err.Error()))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(requiredRole entities.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			web_response.HandleError(c, web_response.NewForbiddenError("user role not found in context"))
			c.Abort()
			return
		}

		userRole := entities.UserRole(role.(string))
		if userRole != requiredRole && userRole != entities.RoleAdmin { // Admins can do anything
			web_response.HandleError(c, web_response.NewForbiddenError("insufficient permissions"))
			c.Abort()
			return
		}
		c.Next()
	}
}