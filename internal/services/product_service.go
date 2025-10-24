package services

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
	"mobile-store-back/internal/utils"

	"github.com/google/uuid"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) Create(name string, description string, basePrice float64, sku string, isActive bool, feature bool, brand string, model string, material string, categoryID uuid.UUID, tags []string, videoURL *string) (*models.Product, error) {
	// Генерируем slug из названия товара
	slug := utils.GenerateSlug(name)

	// Проверяем уникальность slug
	uniqueSlug := utils.GenerateUniqueSlug(slug, func(slugToCheck string) bool {
		_, err := s.repo.GetBySlug(slugToCheck)
		return err != nil // Если ошибка, значит slug уникален
	})

	return s.repo.Create(name, uniqueSlug, description, basePrice, sku, isActive, feature, brand, model, material, categoryID.String(), tags, videoURL)
}

func (s *ProductService) GetByID(id string) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) GetBySKU(sku string) (*models.Product, error) {
	return s.repo.GetBySKU(sku)
}

func (s *ProductService) Update(id string, name *string, description *string, basePrice *float64, isActive *bool, feature *bool, brand *string, model *string, material *string, categoryID *uuid.UUID, tags []string, videoURL *string) (*models.Product, error) {
	var categoryIDStr *string
	if categoryID != nil {
		s := categoryID.String()
		categoryIDStr = &s
	}

	return s.repo.Update(id, name, description, basePrice, isActive, feature, brand, model, material, categoryIDStr, tags, videoURL)
}

func (s *ProductService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *ProductService) List() ([]*models.Product, error) {
	return s.repo.List()
}

func (s *ProductService) Search(query string) ([]*models.Product, error) {
	return s.repo.Search(query)
}

func (s *ProductService) GetByCategory(categoryID string) ([]*models.Product, error) {
	return s.repo.GetByCategory(categoryID)
}

func (s *ProductService) ListWithFilters(brand, minPrice, maxPrice string) ([]*models.Product, error) {
	return s.repo.ListWithFilters(brand, minPrice, maxPrice)
}

// GetBySlug получает товар по slug
func (s *ProductService) GetBySlug(slug string) (*models.Product, error) {
	return s.repo.GetBySlug(slug)
}

// GetBySlugOrID получает товар по slug или ID (универсальный метод)
func (s *ProductService) GetBySlugOrID(identifier string) (*models.Product, error) {
	// Сначала пробуем найти по slug
	product, err := s.repo.GetBySlug(identifier)
	if err == nil {
		return product, nil
	}
	
	// Если не найден по slug, пробуем по ID
	return s.repo.GetByID(identifier)
}

// GetFeatured получает все товары с feature=true
func (s *ProductService) GetFeatured() ([]*models.Product, error) {
	return s.repo.GetFeatured()
}
