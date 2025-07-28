
package entities
/*
import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Role      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Menu struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"not null"`
	Description string
	Price       float64   `gorm:"not null"`
	Category    string    `gorm:"not null"`
	Stock       int       `gorm:"not null"`
	ImageURL    string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type Cart struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null"`
	User      User
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type CartItem struct {
	ID        uint `gorm:"primaryKey"`
	CartID    uint `gorm:"not null"`
	Cart      Cart
	MenuID    uint `gorm:"not null"`
	Menu      Menu
	Quantity  int       `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Order struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	User        User
	TotalAmount float64   `gorm:"not null"`
	Status      string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint `gorm:"not null"`
	Order     Order
	MenuID    uint `gorm:"not null"`
	Menu      Menu
	Quantity  int       `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	MenuName  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

*/