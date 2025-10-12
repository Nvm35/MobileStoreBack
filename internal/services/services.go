package services

import (
	"mobile-store-back/internal/config"
	"mobile-store-back/internal/repository"
)

type Services struct {
	Auth    *AuthService
	User    *UserService
	Product *ProductService
	Order   *OrderService
}

func New(repos *repository.Repository, cfg *config.Config) *Services {
	return &Services{
		Auth:    NewAuthService(repos.Auth, cfg),
		User:    NewUserService(repos.User),
		Product: NewProductService(repos.Product),
		Order:   NewOrderService(repos.Order),
	}
}
