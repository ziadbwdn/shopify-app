// internal/entities/menu.go
package entities

import (
	"shopify-app/internal/utils"
	"time"
	"gorm.io/gorm"
)

// Menu represents the food menu entity in the database
type Menu struct {
	ID          utils.BinaryUUID  `gorm:"type:binary(16);primaryKey" json:"id"`
	Name        string            `gorm:"type:varchar(255);not null;index" json:"name" validate:"required,min=2,max=255"`
	Description string            `gorm:"type:text" json:"description"`
	Price       *utils.GormDecimal `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gt=0"`
	Category    string            `gorm:"type:varchar(100);not null;index" json:"category" validate:"required,min=2,max=100"`
	Stock       int               `gorm:"type:int;not null;default:0" json:"stock" validate:"gte=0"`
	ImageURL    string            `gorm:"type:varchar(500)" json:"image_url"`
	IsActive    bool              `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relationships
	CartItems  []CartItem  `gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE" json:"cart_items,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:MenuID;constraint:OnDelete:RESTRICT" json:"order_items,omitempty"`
}

// TableName returns the table name for the Menu entity
func (Menu) TableName() string {
	return "menus"
}

// BeforeCreate hook to generate UUID before creating menu
func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	if m.ID == (utils.BinaryUUID{}) {
		m.ID = utils.NewBinaryUUID()
	}
	return nil
}

// IsInStock checks if the menu item has sufficient stock
func (m *Menu) IsInStock(requestedQuantity int) bool {
	return m.IsActive && m.Stock >= requestedQuantity
}

// ReduceStock reduces the stock by the given quantity
func (m *Menu) ReduceStock(quantity int) error {
	if !m.IsInStock(quantity) {
		return gorm.ErrInvalidValue
	}
	m.Stock -= quantity
	return nil
}