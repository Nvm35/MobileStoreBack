package repository

import (
	"mobile-store-back/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB, redis *redis.Client) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(userID string, items []struct {
	ProductID        string
	ProductVariantID *string
	Quantity         int
}, shippingMethod string, shippingAddress string, pickupPoint string, paymentMethod string, customerNotes string) (*models.Order, error) {
	// TODO: Реализовать полную логику создания заказа
	// Это требует сложной бизнес-логики: расчет цен, создание order_items и т.д.
	// Пока возвращаем заглушку
	return nil, gorm.ErrNotImplemented
}

func (r *orderRepository) GetByID(id string) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("User").Preload("OrderItems").Preload("OrderItems.Product").Preload("OrderItems.ProductVariant").
		First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(userID string) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.Preload("OrderItems").Preload("OrderItems.Product").Preload("OrderItems.ProductVariant").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) Update(id string, userID string, status *string, paymentStatus *string, trackingNumber *string, customerNotes *string, shippingMethod *string, shippingAddress *string, pickupPoint *string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&order).Error
	if err != nil {
		return nil, err
	}
	
	if status != nil {
		order.Status = models.OrderStatus(*status)
	}
	if paymentStatus != nil {
		order.PaymentStatus = models.PaymentStatus(*paymentStatus)
	}
	if trackingNumber != nil {
		order.TrackingNumber = *trackingNumber
	}
	// AdminNotes удален из модели
	if customerNotes != nil {
		order.CustomerNotes = *customerNotes
	}
	if shippingMethod != nil {
		order.ShippingMethod = *shippingMethod
	}
	if shippingAddress != nil {
		order.ShippingAddress = *shippingAddress
	}
	if pickupPoint != nil {
		order.PickupPoint = *pickupPoint
	}
	
	err = r.db.Save(&order).Error
	return &order, err
}

func (r *orderRepository) UpdateStatus(id string, status string, trackingNumber *string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	
	order.Status = models.OrderStatus(status)
	if trackingNumber != nil {
		order.TrackingNumber = *trackingNumber
	}
	// AdminNotes удален из модели
	
	err = r.db.Save(&order).Error
	return &order, err
}

func (r *orderRepository) Delete(id string) error {
	return r.db.Delete(&models.Order{}, "id = ?", id).Error
}

func (r *orderRepository) List() ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.Preload("User").Preload("OrderItems").Preload("OrderItems.Product").Preload("OrderItems.ProductVariant").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
