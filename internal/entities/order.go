// internal/entities/order.go
package entities

import (
	"shopify-app/internal/utils"
	"time"
	"strconv"
	"gorm.io/gorm"
	pbdecimal "google.golang.org/genproto/googleapis/type/decimal"
)

// OrderStatus defines the status types for orders
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusPreparing OrderStatus = "preparing"
	StatusReady     OrderStatus = "ready"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

// Order represents the order entity in the database
type Order struct {
	ID          utils.BinaryUUID   `gorm:"type:binary(16);primaryKey" json:"id"`
	UserID      utils.BinaryUUID   `gorm:"type:binary(16);not null;index" json:"user_id"`
	TotalAmount *utils.GormDecimal `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status      OrderStatus        `gorm:"type:enum('pending','confirmed','preparing','ready','delivered','cancelled');not null;default:'pending'" json:"status"`
	CreatedAt   time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relationships
	User       User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items,omitempty"`
}

// TableName returns the table name for the Order entity
func (Order) TableName() string {
	return "orders"
}

// BeforeCreate hook to generate UUID before creating order
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == (utils.BinaryUUID{}) {
		o.ID = utils.NewBinaryUUID()
	}
	return nil
}

// OrderItem represents individual items in an order
type OrderItem struct {
	ID        utils.BinaryUUID   `gorm:"type:binary(16);primaryKey" json:"id"`
	OrderID   utils.BinaryUUID   `gorm:"type:binary(16);not null;index" json:"order_id"`
	MenuID    utils.BinaryUUID   `gorm:"type:binary(16);not null;index" json:"menu_id"`
	Quantity  int                `gorm:"type:int;not null" json:"quantity"`
	Price     *utils.GormDecimal `gorm:"type:decimal(10,2);not null" json:"price"` // Price snapshot at time of order
	MenuName  string             `gorm:"type:varchar(255);not null" json:"menu_name"` // Menu name snapshot
	CreatedAt time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relationships
	Order Order `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
	Menu  Menu  `gorm:"foreignKey:MenuID;constraint:OnDelete:RESTRICT" json:"menu,omitempty"`
}

// TableName returns the table name for the OrderItem entity
func (OrderItem) TableName() string {
	return "order_items"
}

// BeforeCreate hook to generate UUID before creating order item
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == (utils.BinaryUUID{}) {
		oi.ID = utils.NewBinaryUUID()
	}
	return nil
}

// GetSubtotal calculates the subtotal for this order item (price * quantity)
func (oi *OrderItem) GetSubtotal() *utils.GormDecimal {
	if oi.Price == nil {
		return &utils.GormDecimal{Internal: pbdecimal.Decimal{Value: "0"}}
	}
	// Simple string-based calculation for decimal
	priceFloat, _ := strconv.ParseFloat(oi.Price.Internal.Value, 64)
	subtotal := priceFloat * float64(oi.Quantity)
	return &utils.GormDecimal{Internal: pbdecimal.Decimal{Value: strconv.FormatFloat(subtotal, 'f', 2, 64)}}
}

// CanBeCancelled checks if the order can be cancelled based on its current status
func (o *Order) CanBeCancelled() bool {
	return o.Status == StatusPending || o.Status == StatusConfirmed
}

// CanBeUpdated checks if the order can be updated based on its current status
func (o *Order) CanBeUpdated() bool {
	return o.Status != StatusDelivered && o.Status != StatusCancelled
}