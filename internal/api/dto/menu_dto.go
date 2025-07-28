package dto

// CreateMenuRequest defines the request body for creating a new menu item
type CreateMenuRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Category    string  `json:"category" validate:"required,min=2,max=100"`
	Stock       int     `json:"stock" validate:"gte=0"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
}

// UpdateMenuRequest defines the request body for updating a menu item
type UpdateMenuRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Category    string  `json:"category" validate:"required,min=2,max=100"`
	Stock       int     `json:"stock" validate:"gte=0"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
}
