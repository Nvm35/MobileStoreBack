package repository

import (
	"mobile-store-back/internal/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type cartRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewCartRepository(db *gorm.DB, redis *redis.Client) CartRepository {
	return &cartRepository{
		db:    db,
		redis: redis,
	}
}

func (r *cartRepository) GetByUserID(userID string) ([]models.CartItem, error) {
	var items []models.CartItem
	err := r.db.Where("user_id = ?", userID).Preload("Product").Find(&items).Error
	return items, err
}

func (r *cartRepository) AddItem(userID string, productID string, quantity int) (*models.CartItem, error) {
	// Проверяем, есть ли уже такой товар в корзине
	var existingItem models.CartItem
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&existingItem).Error
	
	if err == nil {
		// Товар уже есть, обновляем количество
		existingItem.Quantity += quantity
		err = r.db.Save(&existingItem).Error
		return &existingItem, err
	}
	
	// Товара нет, создаем новый
	userUUID, _ := uuid.Parse(userID)
	productUUID, _ := uuid.Parse(productID)
	item := models.CartItem{
		UserID:    userUUID,
		ProductID: productUUID,
		Quantity:  quantity,
	}
	
	err = r.db.Create(&item).Error
	if err != nil {
		return nil, err
	}
	
	// Загружаем связанные данные
	err = r.db.Preload("Product").First(&item, item.ID).Error
	return &item, err
}

func (r *cartRepository) UpdateItem(id string, userID string, quantity int) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error
	if err != nil {
		return nil, err
	}
	
	item.Quantity = quantity
	err = r.db.Save(&item).Error
	if err != nil {
		return nil, err
	}
	
	// Загружаем связанные данные
	err = r.db.Preload("Product").First(&item, item.ID).Error
	return &item, err
}

func (r *cartRepository) RemoveItem(id string, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.CartItem{}).Error
}

func (r *cartRepository) Clear(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error
}

func (r *cartRepository) GetCount(userID string) (int, error) {
	var count int64
	err := r.db.Model(&models.CartItem{}).Where("user_id = ?", userID).Count(&count).Error
	return int(count), err
}
