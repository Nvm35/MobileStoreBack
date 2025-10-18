package services

import (
	"errors"
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type ProductVariantService struct {
	repo repository.ProductVariantRepository
}

func NewProductVariantService(repo repository.ProductVariantRepository) *ProductVariantService {
	return &ProductVariantService{
		repo: repo,
	}
}

func (s *ProductVariantService) Create(productID string, sku string, name string, color string, size string, price float64, isActive bool) (*models.ProductVariant, error) {
	return s.repo.Create(productID, sku, name, color, size, price, isActive)
}

func (s *ProductVariantService) GetByID(id string) (*models.ProductVariant, error) {
	return s.repo.GetByID(id)
}

func (s *ProductVariantService) GetBySKU(sku string) (*models.ProductVariant, error) {
	return s.repo.GetBySKU(sku)
}

func (s *ProductVariantService) GetByProductID(productID string) ([]*models.ProductVariant, error) {
	return s.repo.GetByProductID(productID)
}

func (s *ProductVariantService) Update(id string, sku *string, name *string, color *string, size *string, price *float64, isActive *bool) (*models.ProductVariant, error) {
	return s.repo.Update(id, sku, name, color, size, price, isActive)
}

func (s *ProductVariantService) Delete(id string) error {
	return s.repo.Delete(id)
}

// GetByProductSlugOrID - получение вариантов товара по slug или ID продукта
func (s *ProductVariantService) GetByProductSlugOrID(identifier string) ([]*models.ProductVariant, error) {
	// Пока возвращаем ошибку, так как нужен доступ к ProductService
	return nil, errors.New("not implemented - need ProductService dependency")
}

