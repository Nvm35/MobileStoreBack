package services

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type CartService struct {
	repo repository.CartRepository
}

func NewCartService(repo repository.CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) GetByUserID(userID string) ([]models.CartItem, error) {
	return s.repo.GetByUserID(userID)
}

func (s *CartService) AddItem(userID string, productID string, quantity int) (*models.CartItem, error) {
	return s.repo.AddItem(userID, productID, quantity)
}

func (s *CartService) UpdateItem(id string, userID string, quantity int) (*models.CartItem, error) {
	return s.repo.UpdateItem(id, userID, quantity)
}

func (s *CartService) RemoveItem(id string, userID string) error {
	return s.repo.RemoveItem(id, userID)
}

func (s *CartService) Clear(userID string) error {
	return s.repo.Clear(userID)
}

func (s *CartService) GetCount(userID string) (int, error) {
	return s.repo.GetCount(userID)
}
