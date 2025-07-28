package service

import (
	"context"
	"shopify-app/internal/config"
	"shopify-app/internal/contract"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
	"shopify-app/pkg/jwt"
)

type userService struct {
	userRepo contract.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo contract.UserRepository, cfg *config.Config) contract.UserService {
	return &userService{userRepo: userRepo, cfg: cfg}
}

func (s *userService) Register(ctx context.Context, email, password string, role entities.UserRole) (*entities.User, string, *exception.AppError) {
	if err := utils.ValidatePasswordWithRegex(password); err != nil {
		return nil, "", exception.NewAppError(err, "invalid password")
	}

	exists, err := s.userRepo.EmailExists(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", exception.NewAppError(nil, "email already exists")
	}

	hashedPassword, hashErr := utils.HashPassword(password)
	if hashErr != nil {
		return nil, "", exception.NewAppError(hashErr, "failed to hash password")
	}

	user := &entities.User{
		Email:    email,
		Password: hashedPassword,
		Role:     role,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, "", err
	}

	token, tokenErr := jwt.GenerateToken(user.ID, user.Email, string(user.Role), s.cfg.JWTSecret)
	if tokenErr != nil {
		return nil, "", exception.NewAppError(tokenErr, "failed to generate token")
	}

	return user, token, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (*entities.User, string, *exception.AppError) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", exception.NewAppError(nil, "invalid credentials")
	}

	token, tokenErr := jwt.GenerateToken(user.ID, user.Email, string(user.Role), s.cfg.JWTSecret)
	if tokenErr != nil {
		return nil, "", exception.NewAppError(tokenErr, "failed to generate token")
	}

	return user, token, nil
}

func (s *userService) GetUserProfile(ctx context.Context, userID utils.BinaryUUID) (*entities.User, *exception.AppError) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func (s *userService) UpdateUserProfile(ctx context.Context, userID utils.BinaryUUID, email string) (*entities.User, *exception.AppError) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if email is changing and if the new one is already taken
	if email != user.Email {
		exists, err := s.userRepo.EmailExists(ctx, email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, exception.NewAppError(nil, "email already taken")
		}
		user.Email = email
	}

	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID utils.BinaryUUID, currentPassword, newPassword string) *exception.AppError {
	if err := utils.ValidatePasswordWithRegex(newPassword); err != nil {
		return exception.NewAppError(err, "invalid new password")
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(currentPassword, user.Password) {
		return exception.NewAppError(nil, "invalid current password")
	}

	hashedPassword, hashErr := utils.HashPassword(newPassword)
	if hashErr != nil {
		return exception.NewAppError(hashErr, "failed to hash new password")
	}

	user.Password = hashedPassword
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *userService) ValidateUserCredentials(ctx context.Context, email, password string) (*entities.User, *exception.AppError) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, exception.NewAppError(nil, "invalid credentials")
	}
	return user, nil
}

func (s *userService) GetUsersByRole(ctx context.Context, role entities.UserRole, offset, limit int) ([]entities.User, int64, *exception.AppError) {
	return s.userRepo.GetUsersByRole(ctx, role, offset, limit)
}
