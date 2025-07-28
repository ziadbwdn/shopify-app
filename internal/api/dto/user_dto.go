package dto

// ChangePasswordRequest defines the request body for changing a password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}

// UpdateProfileRequest defines the request body for updating a user's profile
type UpdateProfileRequest struct {
	Email string `json:"email" validate:"required,email"`
}
