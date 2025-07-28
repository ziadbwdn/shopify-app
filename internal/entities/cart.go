// internal/entities/cart.go
package entities

import (
	"shopify-app/internal/utils"
	"time"
	"strconv"
	"gorm.io/gorm"
	pbdecimal "google.golang.org/genproto/googleapis/type/decimal"
)

// Cart represents the shopping cart entity in the database
type Cart struct {
	ID        utils.BinaryUUID `gorm:"type:binary(16);primaryKey" json:"id"`
	UserID    utils.BinaryUUID `gorm:"type:binary(16);not null;uniqueIndex" json:"user_id"`
	CreatedAt time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relationships
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	CartItems []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"cart_items,omitempty"`
}

// TableName returns the table name for the Cart entity
func (Cart) TableName() string {
	return "carts"
}

// BeforeCreate hook to generate UUID before creating cart
func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == (utils.BinaryUUID{}) {
		c.ID = utils.NewBinaryUUID()
	}
	return nil
}

// CartItem represents individual items in a shopping cart
type CartItem struct {
	ID        utils.BinaryUUID   `gorm:"type:binary(16);primaryKey" json:"id"`
	CartID    utils.BinaryUUID   `gorm:"type:binary(16);not null;index" json:"cart_id"`
	MenuID    utils.BinaryUUID   `gorm:"type:binary(16);not null;index" json:"menu_id"`
	Quantity  int                `gorm:"type:int;not null;default:1" json:"quantity" validate:"required,gt=0"`
	Price     *utils.GormDecimal `gorm:"type:decimal(10,2);not null" json:"price"` // Price snapshot at time of adding to cart
	CreatedAt time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relationships
	Cart Cart `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"cart,omitempty"`
	Menu Menu `gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE" json:"menu,omitempty"`
}

// TableName returns the table name for the CartItem entity
func (CartItem) TableName() string {
	return "cart_items"
}

// BeforeCreate hook to generate UUID before creating cart item
func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == (utils.BinaryUUID{}) {
		ci.ID = utils.NewBinaryUUID()
	}
	return nil
}

// GetSubtotal calculates the subtotal for this cart item (price * quantity)
func (ci *CartItem) GetSubtotal() *utils.GormDecimal {
	if ci.Price == nil {
		return &utils.GormDecimal{Internal: pbdecimal.Decimal{Value: "0"}}
	}
	// Simple string-based calculation for decimal
	priceFloat, _ := strconv.ParseFloat(ci.Price.Internal.Value, 64)
	subtotal := priceFloat * float64(ci.Quantity)
	return &utils.GormDecimal{Internal: pbdecimal.Decimal{Value: strconv.FormatFloat(subtotal, 'f', 2, 64)}}
}