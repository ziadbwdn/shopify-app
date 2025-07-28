package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"shopify-app/internal/contract"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
)

// menuRepository implements the contract.MenuRepository interface
type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository creates a new instance of the menu repository
func NewMenuRepository(db *gorm.DB) contract.MenuRepository {
	return &menuRepository{db: db}
}

// CreateMenu creates a new menu item in the database
func (r *menuRepository) CreateMenu(ctx context.Context, menu *entities.Menu) *exception.AppError {
	if err := r.db.WithContext(ctx).Create(menu).Error; err != nil {
		return exception.NewAppError(err, "failed to create menu")
	}
	return nil
}

// GetMenuByID retrieves a menu item by its ID
func (r *menuRepository) GetMenuByID(ctx context.Context, id utils.BinaryUUID) (*entities.Menu, *exception.AppError) {
	var menu entities.Menu
	if err := r.db.WithContext(ctx).First(&menu, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(err, "menu not found")
		}
		return nil, exception.NewAppError(err, "failed to get menu by id")
	}
	return &menu, nil
}

// GetAllMenus retrieves all menu items with optional filtering and pagination
func (r *menuRepository) GetAllMenus(ctx context.Context, offset, limit int, search, category string, activeOnly bool) ([]entities.Menu, int64, *exception.AppError) {
	var menus []entities.Menu
	var count int64

	query := r.db.WithContext(ctx).Model(&entities.Menu{})

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to count menus")
	}

	if err := query.Offset(offset).Limit(limit).Find(&menus).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to get all menus")
	}

	return menus, count, nil
}

// UpdateMenu updates an existing menu item
func (r *menuRepository) UpdateMenu(ctx context.Context, menu *entities.Menu) *exception.AppError {
	if err := r.db.WithContext(ctx).Save(menu).Error; err != nil {
		return exception.NewAppError(err, "failed to update menu")
	}
	return nil
}

// DeleteMenu soft deletes a menu item
func (r *menuRepository) DeleteMenu(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	if err := r.db.WithContext(ctx).Delete(&entities.Menu{}, "id = ?", id).Error; err != nil {
		return exception.NewAppError(err, "failed to delete menu")
	}
	return nil
}

// UpdateMenuStock updates the stock quantity of a menu item
func (r *menuRepository) UpdateMenuStock(ctx context.Context, id utils.BinaryUUID, newStock int) *exception.AppError {
	if err := r.db.WithContext(ctx).Model(&entities.Menu{}).Where("id = ?", id).Update("stock", newStock).Error; err != nil {
		return exception.NewAppError(err, "failed to update menu stock")
	}
	return nil
}

// ReduceMenuStock reduces the stock quantity of a menu item
func (r *menuRepository) ReduceMenuStock(ctx context.Context, id utils.BinaryUUID, quantity int) *exception.AppError {
	if err := r.db.WithContext(ctx).Model(&entities.Menu{}).Where("id = ?", id).Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
		return exception.NewAppError(err, "failed to reduce menu stock")
	}
	return nil
}

// GetMenusByIDs retrieves multiple menu items by their IDs
func (r *menuRepository) GetMenusByIDs(ctx context.Context, ids []utils.BinaryUUID) ([]entities.Menu, *exception.AppError) {
	var menus []entities.Menu
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&menus).Error; err != nil {
		return nil, exception.NewAppError(err, "failed to get menus by ids")
	}
	return menus, nil
}

// GetMenusByCategory retrieves menu items by category
func (r *menuRepository) GetMenusByCategory(ctx context.Context, category string, offset, limit int) ([]entities.Menu, int64, *exception.AppError) {
	var menus []entities.Menu
	var count int64

	query := r.db.WithContext(ctx).Model(&entities.Menu{}).Where("category = ?", category)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to count menus by category")
	}

	if err := query.Offset(offset).Limit(limit).Find(&menus).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to get menus by category")
	}

	return menus, count, nil
}

// SearchMenus searches menu items by name or description
func (r *menuRepository) SearchMenus(ctx context.Context, query string, offset, limit int) ([]entities.Menu, int64, *exception.AppError) {
	var menus []entities.Menu
	var count int64

	dbQuery := r.db.WithContext(ctx).Model(&entities.Menu{}).Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")

	if err := dbQuery.Count(&count).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to count menus by search")
	}

	if err := dbQuery.Offset(offset).Limit(limit).Find(&menus).Error; err != nil {
		return nil, 0, exception.NewAppError(err, "failed to get menus by search")
	}

	return menus, count, nil
}

// GetCategories retrieves all distinct categories
func (r *menuRepository) GetCategories(ctx context.Context) ([]string, *exception.AppError) {
	var categories []string
	if err := r.db.WithContext(ctx).Model(&entities.Menu{}).Distinct().Pluck("category", &categories).Error; err != nil {
		return nil, exception.NewAppError(err, "failed to get categories")
	}
	return categories, nil
}
