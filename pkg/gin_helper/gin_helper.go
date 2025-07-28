package gin_helper

import (
	"shopify-app/internal/exception"
	//"shopify-app/pkg/web_response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindAndValidate binds the request body to the given struct and validates it.
// If binding or validation fails, it handles the error response and returns true.
func BindAndValidate(c *gin.Context, obj interface{}) error {
	// Bind the request body
	if err := c.ShouldBindJSON(obj); err != nil {
		return exception.NewAppError(err, "Invalid request body", exception.CodeValidation)
	}

	// Validate the struct
	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		// In a real app, you might want to format the validation errors nicely
		return exception.NewAppError(err, "Validation failed", exception.CodeValidation)
	}

	return nil
}
