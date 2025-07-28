package database

import (
	"log"
	"shopify-app/internal/entities"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	return db.AutoMigrate(
		&entities.User{},
		&entities.Menu{},
		&entities.Cart{},
		&entities.CartItem{},
		&entities.Order{},
		&entities.OrderItem{},
	)
}
