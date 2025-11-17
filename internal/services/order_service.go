package services

import (
	"errors"
	"fmt"
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type OrderService struct {
	repo        repository.OrderRepository
	productRepo repository.ProductRepository
	variantRepo repository.ProductVariantRepository
}

type OrderItemInput struct {
	ProductID         *uuid.UUID
	ProductSlug       *string
	ProductVariantID  *uuid.UUID
	ProductVariantSKU *string
	Quantity          int
}

func NewOrderService(repo repository.OrderRepository, productRepo repository.ProductRepository, variantRepo repository.ProductVariantRepository) *OrderService {
	return &OrderService{
		repo:        repo,
		productRepo: productRepo,
		variantRepo: variantRepo,
	}
}

func (s *OrderService) Create(userID string, items []OrderItemInput, shippingMethod string, shippingAddress string, pickupPoint string, paymentMethod string, customerNotes string) (*models.Order, error) {
	itemsStr := make([]struct {
		ProductID        string
		ProductVariantID *string
		Quantity         int
	}, len(items))

	for i, item := range items {
		productID, err := s.resolveProductIdentifier(item.ProductID, item.ProductSlug)
		if err != nil {
			return nil, err
		}

		variantID, err := s.resolveVariantIdentifier(item.ProductVariantID, item.ProductVariantSKU, productID)
		if err != nil {
			return nil, err
		}

		itemsStr[i] = struct {
			ProductID        string
			ProductVariantID *string
			Quantity         int
		}{
			ProductID:        productID.String(),
			ProductVariantID: nil,
			Quantity:         item.Quantity,
		}
		if variantID != nil {
			idStr := variantID.String()
			itemsStr[i].ProductVariantID = &idStr
		}
	}

	return s.repo.Create(userID, itemsStr, shippingMethod, shippingAddress, pickupPoint, paymentMethod, customerNotes)
}

func (s *OrderService) GetByID(id string) (*models.Order, error) {
	return s.repo.GetByID(id)
}

func (s *OrderService) GetByUserID(userID string) ([]*models.Order, error) {
	return s.repo.GetByUserID(userID)
}

func (s *OrderService) Update(id string, userID string, status *string, paymentStatus *string, trackingNumber *string, customerNotes *string, shippingMethod *string, shippingAddress *string, pickupPoint *string) (*models.Order, error) {
	return s.repo.Update(id, userID, status, paymentStatus, trackingNumber, customerNotes, shippingMethod, shippingAddress, pickupPoint)
}

func (s *OrderService) UpdateStatus(id string, status string, trackingNumber *string) (*models.Order, error) {
	return s.repo.UpdateStatus(id, status, trackingNumber)
}

func (s *OrderService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *OrderService) List() ([]*models.Order, error) {
	return s.repo.List()
}

func (s *OrderService) resolveProductIdentifier(id *uuid.UUID, slug *string) (uuid.UUID, error) {
	if id != nil {
		return *id, nil
	}

	if slug == nil || strings.TrimSpace(*slug) == "" {
		return uuid.UUID{}, errors.New("product identifier is required")
	}

	product, err := s.productRepo.GetBySlug(strings.TrimSpace(*slug))
	if err != nil {
		return uuid.UUID{}, err
	}

	return product.ID, nil
}

func (s *OrderService) resolveVariantIdentifier(id *uuid.UUID, sku *string, productID uuid.UUID) (*uuid.UUID, error) {
	if id != nil {
		return id, nil
	}

	if sku == nil || strings.TrimSpace(*sku) == "" {
		return nil, nil
	}

	variant, err := s.variantRepo.GetBySKU(strings.TrimSpace(*sku))
	if err != nil {
		return nil, err
	}

	if variant.ProductID != productID {
		return nil, fmt.Errorf("variant %s does not belong to product %s", variant.ID.String(), productID.String())
	}

	return &variant.ID, nil
}
