package services

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) Create(product *models.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) GetByID(id string) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) GetBySKU(sku string) (*models.Product, error) {
	return s.repo.GetBySKU(sku)
}

func (s *ProductService) Update(product *models.Product) error {
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *ProductService) List(limit, offset int) ([]*models.Product, error) {
	return s.repo.List(limit, offset)
}

func (s *ProductService) Search(query string, limit, offset int) ([]*models.Product, error) {
	return s.repo.Search(query, limit, offset)
}

func (s *ProductService) GetByCategory(categoryID string, limit, offset int) ([]*models.Product, error) {
	return s.repo.GetByCategory(categoryID, limit, offset)
}
