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
	items, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	
	// Заполняем product_slug и variant_sku для каждого элемента
	for i := range items {
		if items[i].Product.Slug != "" {
			items[i].ProductSlug = items[i].Product.Slug
		}
		if items[i].ProductVariant != nil && items[i].ProductVariant.SKU != "" {
			items[i].VariantSKU = items[i].ProductVariant.SKU
		}
	}
	
	return items, nil
}

func (s *CartService) AddItem(userID string, productID string, variantID *string, quantity int) (*models.CartItem, error) {
	item, err := s.repo.AddItem(userID, productID, variantID, quantity)
	if err != nil {
		return nil, err
	}
	
	// Заполняем product_slug и variant_sku
	if item.Product.Slug != "" {
		item.ProductSlug = item.Product.Slug
	}
	if item.ProductVariant != nil && item.ProductVariant.SKU != "" {
		item.VariantSKU = item.ProductVariant.SKU
	}
	
	return item, nil
}

func (s *CartService) UpdateItem(id string, userID string, quantity int) (*models.CartItem, error) {
	item, err := s.repo.UpdateItem(id, userID, quantity)
	if err != nil {
		return nil, err
	}
	
	// Заполняем product_slug и variant_sku
	if item.Product.Slug != "" {
		item.ProductSlug = item.Product.Slug
	}
	if item.ProductVariant != nil && item.ProductVariant.SKU != "" {
		item.VariantSKU = item.ProductVariant.SKU
	}
	
	return item, nil
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

func (s *CartService) MergeCart(userID string, sessionID string) error {
	return s.repo.MergeCart(userID, sessionID)
}
