package repository

import (
	"mobile-store-back/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	User    UserRepository
	Product ProductRepository
	Order   OrderRepository
	Auth    AuthRepository
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
	List(limit, offset int) ([]*models.User, error)
}

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id string) (*models.Product, error)
	GetBySKU(sku string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id string) error
	List(limit, offset int) ([]*models.Product, error)
	Search(query string, limit, offset int) ([]*models.Product, error)
	GetByCategory(categoryID string, limit, offset int) ([]*models.Product, error)
}

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id string) (*models.Order, error)
	GetByUserID(userID string, limit, offset int) ([]*models.Order, error)
	Update(order *models.Order) error
	Delete(id string) error
	List(limit, offset int) ([]*models.Order, error)
}

type AuthRepository interface {
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
}

func New(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		User:    NewUserRepository(db, redis),
		Product: NewProductRepository(db, redis),
		Order:   NewOrderRepository(db, redis),
		Auth:    NewAuthRepository(db, redis),
	}
}
