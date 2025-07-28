// internal/entities/user.go
package entities

import (
	"shopify-app/internal/utils"
	"time"
	"gorm.io/gorm"
)

// UserRole defines the role types for users
type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
)

// User represents the user entity in the database
type User struct {
	ID        utils.BinaryUUID `gorm:"type:binary(16);primaryKey" json:"id"`
	Email     string           `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string           `gorm:"type:varchar(255);not null" json:"-"` // Never expose password in JSON
	Role      UserRole         `gorm:"type:enum('admin','customer');not null;default:'customer'" json:"role"`
	CreatedAt time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relationships
	Cart   *Cart    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"cart,omitempty"`
	Orders []Order  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"orders,omitempty"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to generate UUID before creating user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == (utils.BinaryUUID{}) {
		u.ID = utils.NewBinaryUUID()
	}
	return nil
}