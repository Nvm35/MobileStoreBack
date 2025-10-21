package services

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"

	"github.com/google/uuid"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) Create(userID string, items []struct {
	ProductID        uuid.UUID
	ProductVariantID *uuid.UUID
	Quantity         int
}, shippingMethod string, shippingAddress string, pickupPoint string, paymentMethod string, customerNotes string) (*models.Order, error) {
	// Конвертируем UUID в строки для repository
	itemsStr := make([]struct {
		ProductID        string
		ProductVariantID *string
		Quantity         int
	}, len(items))
	
	for i, item := range items {
		itemsStr[i] = struct {
			ProductID        string
			ProductVariantID *string
			Quantity         int
		}{
			ProductID:        item.ProductID.String(),
			ProductVariantID: nil,
			Quantity:         item.Quantity,
		}
		if item.ProductVariantID != nil {
			variantIDStr := item.ProductVariantID.String()
			itemsStr[i].ProductVariantID = &variantIDStr
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
