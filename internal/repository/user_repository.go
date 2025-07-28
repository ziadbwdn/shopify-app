package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"shopify-app/internal/contract"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
)

// userRepository implements the contract.UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of the user repository
func NewUserRepository(db *gorm.DB) contract.UserRepository {
	return &userRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *userRepository) CreateUser(ctx context.Context, user *entities.User) *exception.AppError {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return exception.NewAppError(err, "failed to create user")
	}
	return nil
}

// GetUserByID retrieves a user by their ID
func (r *userRepository) GetUserByID(ctx context.Context, id utils.BinaryUUID) (*entities.User, *exception.AppError) {
	var user entities.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(err, "user not found")
		}
		return nil, exception.NewAppError(err, "failed to get user by id")
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email address
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, *exception.AppError) {
	var user entities.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(err, "user not found")
		}
		return nil, exception.NewAppError(err, "failed to get user by email")
	}
	return &user, nil
}

// UpdateUser updates an existing user's information
func (r *userRepository) UpdateUser(ctx context.Context, user *entities.User) *exception.AppError {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return exception.NewAppError(err, "failed to update user")
	}
	return nil
}

// DeleteUser soft deletes a user (if the entity has gorm.DeletedAt field)
// Note: The current User entity does not have soft delete. This will be a hard delete.
func (r *userRepository) DeleteUser(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	if err := r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error; err != nil {
		return exception.NewAppError(err, "failed to delete user")
	}
	return nil
}

// EmailExists checks if an email already exists in the database
func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, *exception.AppError) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, exception.NewAppError(err, "failed to check if email exists")
	}
	return count > 0, nil
}

// GetUsersByRole retrieves users by their role with pagination
func (r *userRepository) GetUsersByRole(ctx context.Context, role entities.UserRole, offset, limit int) ([]entities.User, int64, *exception.AppError) {
	var users []entities.User
	var count int64

	query := r.db.WithContext(ctx).Model(&entities.User{}).Where("role = ?", role)

	// Get total count
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to count users by role")
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to get users by role")
	}

	return users, count, nil
}
