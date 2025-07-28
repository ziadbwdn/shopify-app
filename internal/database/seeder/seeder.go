package seeder

import (
	"fmt"
	"log"
	"shopify-app/internal/entities"
	"shopify-app/internal/utils"

	"gorm.io/gorm"
)

// Seed populates the database with initial data.
// It's idempotent and will not create duplicate entries.
func Seed(db *gorm.DB) error {
	if err := seedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := seedMenus(db); err != nil {
		return fmt.Errorf("failed to seed menus: %w", err)
	}

	log.Println("Database seeding completed successfully.")
	return nil
}

// seedUsers creates an admin and a customer if they don't already exist.
func seedUsers(db *gorm.DB) error {
	// Check if admin user already exists
	var adminCount int64
	db.Model(&entities.User{}).Where("email = ?", "admin@example.com").Count(&adminCount)

	if adminCount == 0 {
		log.Println("Seeding admin user...")
		hashedPassword, _ := utils.HashPassword("AdminPassword123!")
		admin := entities.User{
			Email:    "admin@example.com",
			Password: hashedPassword,
			Role:     entities.RoleAdmin,
		}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	}

	// Check if customer user already exists
	var customerCount int64
	db.Model(&entities.User{}).Where("email = ?", "customer@example.com").Count(&customerCount)

	if customerCount == 0 {
		log.Println("Seeding customer user...")
		hashedPassword, _ := utils.HashPassword("CustomerPassword123!")
		customer := entities.User{
			Email:    "customer@example.com",
			Password: hashedPassword,
			Role:     entities.RoleCustomer,
		}
		if err := db.Create(&customer).Error; err != nil {
			return err
		}
	}

	return nil
}

// seedMenus creates some sample menu items if none exist.
func seedMenus(db *gorm.DB) error {
	var menuCount int64
	db.Model(&entities.Menu{}).Count(&menuCount)

	if menuCount == 0 {
		log.Println("Seeding menu items...")
		menus := []entities.Menu{
			{
				Name:        "Classic Burger",
				Description: "A juicy beef patty with lettuce, tomato, and our special sauce.",
				Price:       utils.MustNewGormDecimal("12.99"),
				Category:    "Burgers",
				Stock:       100,
				IsActive:    true,
			},
			{
				Name:        "Margherita Pizza",
				Description: "Classic pizza with fresh mozzarella, tomatoes, and basil.",
				Price:       utils.MustNewGormDecimal("15.50"),
				Category:    "Pizzas",
				Stock:       50,
				IsActive:    true,
			},
			{
				Name:        "Caesar Salad",
				Description: "Crisp romaine lettuce with Caesar dressing, croutons, and parmesan cheese.",
				Price:       utils.MustNewGormDecimal("9.75"),
				Category:    "Salads",
				Stock:       75,
				IsActive:    true,
			},
			{
				Name:        "Chocolate Lava Cake",
				Description: "Warm chocolate cake with a gooey molten center.",
				Price:       utils.MustNewGormDecimal("7.00"),
				Category:    "Desserts",
				Stock:       40,
				IsActive:    true,
			},
		}

		if err := db.Create(&menus).Error; err != nil {
			return err
		}
	}

	return nil
}