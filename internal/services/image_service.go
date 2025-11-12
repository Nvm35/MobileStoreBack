package services

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"

	"github.com/google/uuid"
)

type ImageService struct {
	repo repository.ImageRepository
}

func NewImageService(repo repository.ImageRepository) *ImageService {
	return &ImageService{repo: repo}
}

func (s *ImageService) Create(image *models.Image) error {
	image.ID = uuid.New()
	return s.repo.Create(image)
}

func (s *ImageService) GetByID(id string) (*models.Image, error) {
	return s.repo.GetByID(id)
}

func (s *ImageService) GetByProductID(productID string) ([]*models.Image, error) {
	return s.repo.GetByProductID(productID)
}

func (s *ImageService) GetByProductSlug(productSlug string) ([]*models.Image, error) {
	return s.repo.GetByProductSlug(productSlug)
}

func (s *ImageService) GetPrimaryByProductID(productID string) (*models.Image, error) {
	return s.repo.GetPrimaryByProductID(productID)
}

func (s *ImageService) Update(id string, cloudinaryPublicID *string, url *string, isPrimary *bool) (*models.Image, error) {
	image, err := s.repo.Update(id, cloudinaryPublicID, url, isPrimary)
	if err != nil {
		return nil, err
	}

	if isPrimary != nil && *isPrimary {
		if err := s.repo.SetPrimary(id); err != nil {
			return nil, err
		}
		return s.repo.GetByID(id)
	}

	return image, nil
}

func (s *ImageService) SetPrimary(id string) error {
	return s.repo.SetPrimary(id)
}

func (s *ImageService) Delete(id string) error {
	return s.repo.Delete(id)
}
