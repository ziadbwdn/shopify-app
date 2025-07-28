package repository

/*
import (
	"gorm.io/gorm"
	"shopify-app/internal/entities"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) GetUserByID(id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) UpdateUser(user *entities.User) error {
	return r.db.Save(user).Error
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *menuRepository {
	return &menuRepository{db}
}

func (r *menuRepository) CreateMenu(menu *entities.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) GetMenuByID(id uint) (*entities.Menu, error) {
	var menu entities.Menu
	err := r.db.First(&menu, id).Error
	return &menu, err
}

func (r *menuRepository) GetAllMenus(offset, limit int, search, category string) ([]*entities.Menu, int64, error) {
	var menus []*entities.Menu
	var count int64

	query := r.db.Model(&entities.Menu{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&count)
	query.Offset(offset).Limit(limit).Find(&menus)

	return menus, count, query.Error
}

func (r *menuRepository) UpdateMenu(menu *entities.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) DeleteMenu(id uint) error {
	return r.db.Delete(&entities.Menu{}, id).Error
}

func (r *menuRepository) UpdateMenuStock(id uint, stock int) error {
	return r.db.Model(&entities.Menu{}).Where("id = ?", id).Update("stock", stock).Error
}
*/