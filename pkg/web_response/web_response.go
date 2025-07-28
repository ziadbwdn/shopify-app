package web_response

import (
	"net/http"
	"shopify-app/internal/exception"

	"github.com/gin-gonic/gin"
)

// WebResponse is the standardized structure for all API responses.
type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

// ErrorDetail is the standardized structure for the nested error object.
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success sends a standardized success response.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, WebResponse{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   data,
	})
}

// HandleError sends a standardized error response from an AppError.
func HandleError(c *gin.Context, err error) {
	appErr, ok := err.(*exception.AppError)
	if !ok {
		appErr = exception.NewInternalServerError(err.Error())
	}

	status := appErr.HTTPStatus()

	errorDetail := ErrorDetail{
		Code:    string(appErr.Code),
		Message: appErr.Message,
		Details: appErr.Details,
	}

	c.AbortWithStatusJSON(status, WebResponse{
		Code:   status,
		Status: http.StatusText(status),
		Error:  errorDetail,
	})
}

// NewUnauthorizedError creates a new unauthorized error.
func NewUnauthorizedError(message string) *exception.AppError {
	return exception.NewAppError(nil, message, exception.CodeUnauthorized)
}

// NewForbiddenError creates a new forbidden error.
func NewForbiddenError(message string) *exception.AppError {
	return exception.NewAppError(nil, message, exception.CodeForbidden)
}
