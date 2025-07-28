package main

import (
	"fmt"
	"log"
	"shopify-app/internal/api/router"
	"shopify-app/internal/config"
	"shopify-app/internal/database"
	"shopify-app/internal/database/seeder"
	"shopify-app/internal/repository"
	"shopify-app/internal/service"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Seed the database with initial data
	if err := seeder.Seed(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg)
	menuService := service.NewMenuService(menuRepo)
	cartService := service.NewCartService(cartRepo, menuRepo)
	orderService := service.NewOrderService(orderRepo, cartService, menuRepo)
	reportService := service.NewReportService(reportRepo)

	// Setup router
	r := router.Setup(cfg, userService, menuService, cartService, orderService, reportService)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}