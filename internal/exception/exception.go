package exception

import (
	"fmt"
	"net/http"
)

// ErrorCode defines the type for standardized error codes.
type ErrorCode string

const (
	// CodeValidation indicates an error in input validation.
	CodeValidation ErrorCode = "VALIDATION_ERROR"
	// CodeNotFound indicates that a requested resource was not found.
	CodeNotFound ErrorCode = "NOT_FOUND"
	// CodeUnauthorized indicates a failure in authentication (e.g., bad token, missing credentials).
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	// CodeForbidden indicates a failure in authorization (e.g., insufficient permissions).
	CodeForbidden ErrorCode = "FORBIDDEN"
	// CodeDatabaseError indicates a problem with the database.
	CodeDatabaseError ErrorCode = "DATABASE_ERROR"
	// CodeInternalServerError indicates an unexpected server-side error.
	CodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
)

// AppError is the standard error structure for the application.
type AppError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Err     error       `json:"-"`
}

// Error returns the string representation of the AppError.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %s, message: %s, cause: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

// HTTPStatus returns the corresponding HTTP status code for an AppError.
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case CodeValidation:
		return http.StatusBadRequest
	case CodeNotFound:
		return http.StatusNotFound
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeDatabaseError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// NewAppError creates a new AppError.
func NewAppError(err error, message string, code ...ErrorCode) *AppError {
	appErr := &AppError{
		Message: message,
		Err:     err,
	}
	if len(code) > 0 {
		appErr.Code = code[0]
	} else {
		appErr.Code = CodeInternalServerError
	}
	return appErr
}

// NewInternalServerError creates a new internal server error.
func NewInternalServerError(message string) *AppError {
	return NewAppError(nil, message, CodeInternalServerError)
}

// NewValidationError creates a new validation error.
func NewValidationError(message string, details ...interface{}) *AppError {
	err := NewAppError(nil, message, CodeValidation)
	if len(details) > 0 {
		err.Details = details
	}
	return err
}
