package services

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) Create(order *models.Order) error {
	return s.repo.Create(order)
}

func (s *OrderService) GetByID(id string) (*models.Order, error) {
	return s.repo.GetByID(id)
}

func (s *OrderService) GetByUserID(userID string, limit, offset int) ([]*models.Order, error) {
	return s.repo.GetByUserID(userID, limit, offset)
}

func (s *OrderService) Update(order *models.Order) error {
	return s.repo.Update(order)
}

func (s *OrderService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *OrderService) List(limit, offset int) ([]*models.Order, error) {
	return s.repo.List(limit, offset)
}
