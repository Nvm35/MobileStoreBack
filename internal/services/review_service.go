package services

import (
	"errors"
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type ReviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(repo repository.ReviewRepository) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) GetByProductID(productID string) ([]models.Review, error) {
	return s.repo.GetByProductID(productID)
}

func (s *ReviewService) Create(userID string, productID string, orderID *string, rating int, title string, comment string) (*models.Review, error) {
	return s.repo.Create(userID, productID, orderID, rating, title, comment)
}

func (s *ReviewService) Update(id string, userID string, rating *int, title *string, comment *string) (*models.Review, error) {
	return s.repo.Update(id, userID, rating, title, comment)
}

func (s *ReviewService) Delete(id string, userID string) error {
	return s.repo.Delete(id, userID)
}

func (s *ReviewService) Vote(id string, userID string, helpful bool) error {
	return s.repo.Vote(id, userID, helpful)
}

func (s *ReviewService) GetByUserID(userID string) ([]models.Review, error) {
	return s.repo.GetByUserID(userID)
}

func (s *ReviewService) GetAll() ([]models.Review, error) {
	return s.repo.GetAll()
}

func (s *ReviewService) Approve(id string, approved bool) error {
	return s.repo.Approve(id, approved)
}

// GetByProductSlugOrID - получение отзывов по slug или ID продукта
func (s *ReviewService) GetByProductSlugOrID(identifier string) ([]models.Review, error) {
	// Пока возвращаем ошибку, так как нужен доступ к ProductService
	return nil, errors.New("not implemented - need ProductService dependency")
}
