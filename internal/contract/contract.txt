/*
package contract

import "shopify-app/internal/entities"

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByID(id uint) (*entities.User, error)
	UpdateUser(user *entities.User) error
}

type MenuRepository interface {
	CreateMenu(menu *entities.Menu) error
	GetMenuByID(id uint) (*entities.Menu, error)
	GetAllMenus(offset, limit int, search, category string) ([]*entities.Menu, int64, error)
	UpdateMenu(menu *entities.Menu) error
	DeleteMenu(id uint) error
	UpdateMenuStock(id uint, stock int) error
}
*/