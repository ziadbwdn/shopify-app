// internal/contract/user_contract.go
package contract

import (
	"context"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
)

// UserRepository defines the contract for user data access operations
type UserRepository interface {
	// CreateUser creates a new user in the database || ready to ignore
	CreateUser(ctx context.Context, user *entities.User) *exception.AppError
	
	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, id utils.BinaryUUID) (*entities.User, *exception.AppError)
	
	// GetUserByEmail retrieves a user by their email address
	GetUserByEmail(ctx context.Context, email string) (*entities.User, *exception.AppError)
	
	// UpdateUser updates an existing user's information
	UpdateUser(ctx context.Context, user *entities.User) *exception.AppError
	
	// DeleteUser soft deletes a user (if needed for admin operations)
	DeleteUser(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	
	// EmailExists checks if an email already exists in the database
	EmailExists(ctx context.Context, email string) (bool, *exception.AppError)
	
	// GetUsersByRole retrieves users by their role (for admin operations)
	GetUsersByRole(ctx context.Context, role entities.UserRole, offset, limit int) ([]entities.User, int64, *exception.AppError)
}

// UserService defines the contract for user business logic operations
type UserService interface {
	// Register handles user registration with validation and password hashing
	Register(ctx context.Context, email, password string, role entities.UserRole) (*entities.User, string, *exception.AppError)
	
	// Login handles user authentication and JWT token generation
	Login(ctx context.Context, email, password string) (*entities.User, string, *exception.AppError)
	
	// GetUserProfile retrieves user profile information
	GetUserProfile(ctx context.Context, userID utils.BinaryUUID) (*entities.User, *exception.AppError)
	
	// UpdateUserProfile updates user profile information
	UpdateUserProfile(ctx context.Context, userID utils.BinaryUUID, email string) (*entities.User, *exception.AppError)
	
	// ChangePassword handles password change with validation
	ChangePassword(ctx context.Context, userID utils.BinaryUUID, currentPassword, newPassword string) *exception.AppError
	
	// ValidateUserCredentials validates user credentials without generating token
	ValidateUserCredentials(ctx context.Context, email, password string) (*entities.User, *exception.AppError)
	
	// GetUsersByRole retrieves users by role (admin operation)
	GetUsersByRole(ctx context.Context, role entities.UserRole, offset, limit int) ([]entities.User, int64, *exception.AppError)
}