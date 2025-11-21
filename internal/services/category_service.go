package services

import (
	"errors"

	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

var ErrCategoryHasProducts = errors.New("category has related products")

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) GetAll() ([]*models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) GetByID(id string) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) GetBySlug(slug string) (*models.Category, error) {
	return s.repo.GetBySlug(slug)
}

func (s *CategoryService) Create(category *models.Category) error {
	return s.repo.Create(category)
}

func (s *CategoryService) Update(id string, name *string, description *string, slug *string, imageURL *string) (*models.Category, error) {
	return s.repo.Update(id, name, description, slug, imageURL)
}

func (s *CategoryService) Delete(id string) error {
	hasProducts, err := s.repo.HasProducts(id)
	if err != nil {
		return err
	}

	if hasProducts {
		return ErrCategoryHasProducts
	}

	return s.repo.Delete(id)
}

func (s *CategoryService) GetWithProducts(id string) (*models.Category, error) {
	return s.repo.GetWithProducts(id)
}

func (s *CategoryService) GetBySlugWithProducts(slug string) (*models.Category, error) {
	return s.repo.GetBySlugWithProducts(slug)
}
