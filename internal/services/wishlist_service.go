package services

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type WishlistService struct {
	repo repository.WishlistRepository
}

func NewWishlistService(repo repository.WishlistRepository) *WishlistService {
	return &WishlistService{repo: repo}
}

func (s *WishlistService) GetByUserID(userID string, limit, offset int) ([]models.WishlistItem, error) {
	return s.repo.GetByUserID(userID, limit, offset)
}

func (s *WishlistService) AddItem(userID string, productID string) (*models.WishlistItem, error) {
	return s.repo.AddItem(userID, productID)
}

func (s *WishlistService) RemoveItem(id string, userID string) error {
	return s.repo.RemoveItem(id, userID)
}

func (s *WishlistService) IsInWishlist(userID string, productID string) (bool, error) {
	return s.repo.IsInWishlist(userID, productID)
}

func (s *WishlistService) Clear(userID string) error {
	return s.repo.Clear(userID)
}
