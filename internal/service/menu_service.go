package service

import (
	"context"
	"fmt"
	"shopify-app/internal/contract"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
)

type menuService struct {
	menuRepo contract.MenuRepository
}

func NewMenuService(menuRepo contract.MenuRepository) contract.MenuService {
	return &menuService{menuRepo: menuRepo}
}

func (s *menuService) AddMenu(ctx context.Context, name, description string, price *utils.GormDecimal, category string, stock int, imageURL string) (*entities.Menu, *exception.AppError) {
	menu := &entities.Menu{
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
		Stock:       stock,
		ImageURL:    imageURL,
		IsActive:    true,
	}

	if err := s.menuRepo.CreateMenu(ctx, menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *menuService) GetMenus(ctx context.Context, offset, limit int, search, category string, activeOnly bool) ([]entities.Menu, int64, *exception.AppError) {
	return s.menuRepo.GetAllMenus(ctx, offset, limit, search, category, activeOnly)
}

func (s *menuService) GetMenuByID(ctx context.Context, id utils.BinaryUUID) (*entities.Menu, *exception.AppError) {
	return s.menuRepo.GetMenuByID(ctx, id)
}

func (s *menuService) UpdateMenu(ctx context.Context, id utils.BinaryUUID, name, description string, price *utils.GormDecimal, category string, stock int, imageURL string) (*entities.Menu, *exception.AppError) {
	menu, err := s.menuRepo.GetMenuByID(ctx, id)
	if err != nil {
		return nil, err
	}

	menu.Name = name
	menu.Description = description
	menu.Price = price
	menu.Category = category
	menu.Stock = stock
	menu.ImageURL = imageURL

	if err := s.menuRepo.UpdateMenu(ctx, menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *menuService) DeleteMenu(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	return s.menuRepo.DeleteMenu(ctx, id)
}

func (s *menuService) UpdateMenuStock(ctx context.Context, id utils.BinaryUUID, newStock int) (*entities.Menu, *exception.AppError) {
	if _, err := s.menuRepo.GetMenuByID(ctx, id); err != nil {
		return nil, err
	}

	if err := s.menuRepo.UpdateMenuStock(ctx, id, newStock); err != nil {
		return nil, err
	}
	return s.menuRepo.GetMenuByID(ctx, id)
}

func (s *menuService) CheckMenuAvailability(ctx context.Context, items map[utils.BinaryUUID]int) (map[utils.BinaryUUID]bool, *exception.AppError) {
	availability := make(map[utils.BinaryUUID]bool)
	var ids []utils.BinaryUUID
	for id := range items {
		ids = append(ids, id)
	}

	menus, err := s.menuRepo.GetMenusByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	menuMap := make(map[utils.BinaryUUID]entities.Menu)
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}

	for id, quantity := range items {
		if menu, ok := menuMap[id]; ok {
			availability[id] = menu.IsInStock(quantity)
		} else {
			availability[id] = false
		}
	}
	return availability, nil
}

func (s *menuService) ReserveMenuStock(ctx context.Context, items map[utils.BinaryUUID]int) *exception.AppError {
	// This should be transactional in a real application
	for id, quantity := range items {
		if err := s.menuRepo.ReduceMenuStock(ctx, id, quantity); err != nil {
			// Here you would ideally roll back previous stock reductions
			return exception.NewAppError(err, fmt.Sprintf("failed to reserve stock for menu item %s", id))
		}
	}
	return nil
}

func (s *menuService) GetMenusByCategory(ctx context.Context, category string, offset, limit int) ([]entities.Menu, int64, *exception.AppError) {
	return s.menuRepo.GetMenusByCategory(ctx, category, offset, limit)
}

func (s *menuService) SearchMenus(ctx context.Context, query string, offset, limit int) ([]entities.Menu, int64, *exception.AppError) {
	return s.menuRepo.SearchMenus(ctx, query, offset, limit)
}

func (s *menuService) GetCategories(ctx context.Context) ([]string, *exception.AppError) {
	return s.menuRepo.GetCategories(ctx)
}

func (s *menuService) ToggleMenuStatus(ctx context.Context, id utils.BinaryUUID) (*entities.Menu, *exception.AppError) {
	menu, err := s.menuRepo.GetMenuByID(ctx, id)
	if err != nil {
		return nil, err
	}
	menu.IsActive = !menu.IsActive
	if err := s.menuRepo.UpdateMenu(ctx, menu); err != nil {
		return nil, err
	}
	return menu, nil
}
