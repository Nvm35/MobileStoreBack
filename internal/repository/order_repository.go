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

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id string) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("User").Preload("OrderItems").Preload("OrderItems.Product").Preload("Address").
		First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(userID string, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.Preload("OrderItems").Preload("OrderItems.Product").Preload("Address").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) Delete(id string) error {
	return r.db.Delete(&models.Order{}, "id = ?", id).Error
}

func (r *orderRepository) List(limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.Preload("User").Preload("OrderItems").Preload("OrderItems.Product").Preload("Address").
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
