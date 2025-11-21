package services

import (
	"errors"

	"github.com/google/uuid"

	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type ProductVariantService struct {
	repo        repository.ProductVariantRepository
	productRepo repository.ProductRepository
}

func NewProductVariantService(repo repository.ProductVariantRepository, productRepo repository.ProductRepository) *ProductVariantService {
	return &ProductVariantService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *ProductVariantService) Create(productID string, sku string, name string, color string, size string, price float64, isActive bool) (*models.ProductVariant, error) {
	variant, err := s.repo.Create(productID, sku, name, color, size, price, isActive)
	if err != nil {
		return nil, err
	}
	
	// Заполняем product_slug если нужно
	if s.productRepo != nil {
		product, err := s.productRepo.GetByID(productID)
		if err == nil {
			variant.ProductSlug = product.Slug
		}
	}
	
	return variant, nil
}

func (s *ProductVariantService) GetByID(id string) (*models.ProductVariant, error) {
	variant, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// Заполняем product_slug если нужно
	if s.productRepo != nil {
		product, err := s.productRepo.GetByID(variant.ProductID.String())
		if err == nil {
			variant.ProductSlug = product.Slug
		}
	}
	
	return variant, nil
}

func (s *ProductVariantService) GetBySKU(sku string) (*models.ProductVariant, error) {
	variant, err := s.repo.GetBySKU(sku)
	if err != nil {
		return nil, err
	}
	
	// Заполняем product_slug если нужно
	if s.productRepo != nil {
		product, err := s.productRepo.GetByID(variant.ProductID.String())
		if err == nil {
			variant.ProductSlug = product.Slug
		}
	}
	
	return variant, nil
}

func (s *ProductVariantService) GetByProductID(productID string) ([]*models.ProductVariant, error) {
	variants, err := s.repo.GetByProductID(productID)
	if err != nil {
		return nil, err
	}
	
	// Заполняем product_slug для каждого варианта
	if s.productRepo != nil {
		product, err := s.productRepo.GetByID(productID)
		if err == nil {
			for _, variant := range variants {
				variant.ProductSlug = product.Slug
			}
		}
	}
	
	return variants, nil
}

func (s *ProductVariantService) Update(id string, sku *string, name *string, color *string, size *string, price *float64, isActive *bool) (*models.ProductVariant, error) {
	variant, err := s.repo.Update(id, sku, name, color, size, price, isActive)
	if err != nil {
		return nil, err
	}
	
	// Заполняем product_slug если нужно
	if s.productRepo != nil {
		product, err := s.productRepo.GetByID(variant.ProductID.String())
		if err == nil {
			variant.ProductSlug = product.Slug
		}
	}
	
	return variant, nil
}

func (s *ProductVariantService) Delete(id string) error {
	return s.repo.Delete(id)
}

// GetByProductSlugOrID - получение вариантов товара по slug или ID продукта
func (s *ProductVariantService) GetByProductSlugOrID(identifier string) ([]*models.ProductVariant, error) {
	if s.productRepo == nil {
		return nil, errors.New("product repository dependency is not configured")
	}

	var product *models.Product
	var err error

	// Если пришел UUID, работаем напрямую
	if _, err := uuid.Parse(identifier); err == nil {
		product, err = s.productRepo.GetByID(identifier)
		if err != nil {
			return nil, err
		}
	} else {
		product, err = s.productRepo.GetBySlug(identifier)
		if err != nil {
			return nil, err
		}
	}

	variants, err := s.repo.GetByProductID(product.ID.String())
	if err != nil {
		return nil, err
	}

	// Заполняем product_slug для каждого варианта
	for _, variant := range variants {
		variant.ProductSlug = product.Slug
	}

	return variants, nil
}
