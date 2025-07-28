// internal/contract/menu_contract.go
package contract

import (
	"context"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
)

// MenuRepository defines the contract for menu data access operations
type MenuRepository interface {
	// CreateMenu creates a new menu item in the database
	CreateMenu(ctx context.Context, menu *entities.Menu) *exception.AppError
	
	// GetMenuByID retrieves a menu item by its ID
	GetMenuByID(ctx context.Context, id utils.BinaryUUID) (*entities.Menu, *exception.AppError)
	
	// GetAllMenus retrieves all menu items with optional filtering and pagination
	GetAllMenus(ctx context.Context, offset, limit int, search, category string, activeOnly bool) ([]entities.Menu, int64, *exception.AppError)
	
	// UpdateMenu updates an existing menu item
	UpdateMenu(ctx context.Context, menu *entities.Menu) *exception.AppError
	
	// DeleteMenu soft deletes a menu item
	DeleteMenu(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	
	// UpdateMenuStock updates the stock quantity of a menu item
	UpdateMenuStock(ctx context.Context, id utils.BinaryUUID, newStock int) *exception.AppError
	
	// ReduceMenuStock reduces the stock quantity of a menu item (for orders)
	ReduceMenuStock(ctx context.Context, id utils.BinaryUUID, quantity int) *exception.AppError
	
	// GetMenusByIDs retrieves multiple menu items by their IDs
	GetMenusByIDs(ctx context.Context, ids []utils.BinaryUUID) ([]entities.Menu, *exception.AppError)
	
	// GetMenusByCategory retrieves menu items by category
	GetMenusByCategory(ctx context.Context, category string, offset, limit int) ([]entities.Menu, int64, *exception.AppError)
	
	// SearchMenus searches menu items by name or description
	SearchMenus(ctx context.Context, query string, offset, limit int) ([]entities.Menu, int64, *exception.AppError)
	
	// GetCategories retrieves all distinct categories
	GetCategories(ctx context.Context) ([]string, *exception.AppError)
}

// MenuService defines the contract for menu business logic operations
type MenuService interface {
	// AddMenu handles adding a new menu item with validation
	AddMenu(ctx context.Context, name, description string, price *utils.GormDecimal, category string, stock int, imageURL string) (*entities.Menu, *exception.AppError)
	
	// GetMenus retrieves menu list with filtering and pagination
	GetMenus(ctx context.Context, offset, limit int, search, category string, activeOnly bool) ([]entities.Menu, int64, *exception.AppError)
	
	// GetMenuByID retrieves a specific menu item
	GetMenuByID(ctx context.Context, id utils.BinaryUUID) (*entities.Menu, *exception.AppError)
	
	// UpdateMenu handles updating menu item with validation
	UpdateMenu(ctx context.Context, id utils.BinaryUUID, name, description string, price *utils.GormDecimal, category string, stock int, imageURL string) (*entities.Menu, *exception.AppError)
	
	// DeleteMenu handles menu item deletion
	DeleteMenu(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	
	// UpdateMenuStock updates menu stock with validation
	UpdateMenuStock(ctx context.Context, id utils.BinaryUUID, newStock int) (*entities.Menu, *exception.AppError)
	
	// CheckMenuAvailability checks if menu items are available for given quantities
	CheckMenuAvailability(ctx context.Context, items map[utils.BinaryUUID]int) (map[utils.BinaryUUID]bool, *exception.AppError)
	
	// ReserveMenuStock reserves stock for multiple menu items (for order processing)
	ReserveMenuStock(ctx context.Context, items map[utils.BinaryUUID]int) *exception.AppError
	
	// GetMenusByCategory retrieves menu items by category
	GetMenusByCategory(ctx context.Context, category string, offset, limit int) ([]entities.Menu, int64, *exception.AppError)
	
	// SearchMenus searches menu items
	SearchMenus(ctx context.Context, query string, offset, limit int) ([]entities.Menu, int64, *exception.AppError)
	
	// GetCategories retrieves all available categories
	GetCategories(ctx context.Context) ([]string, *exception.AppError)
	
	// ToggleMenuStatus toggles menu active/inactive status
	ToggleMenuStatus(ctx context.Context, id utils.BinaryUUID) (*entities.Menu, *exception.AppError)
}