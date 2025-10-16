package repository

import (
	"mobile-store-back/internal/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type wishlistRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewWishlistRepository(db *gorm.DB, redis *redis.Client) WishlistRepository {
	return &wishlistRepository{
		db:    db,
		redis: redis,
	}
}

func (r *wishlistRepository) GetByUserID(userID string, limit, offset int) ([]models.WishlistItem, error) {
	var items []models.WishlistItem
	err := r.db.Where("user_id = ?", userID).
		Preload("Product").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *wishlistRepository) AddItem(userID string, productID string) (*models.WishlistItem, error) {
	// Проверяем, есть ли уже такой товар в избранном
	var existingItem models.WishlistItem
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&existingItem).Error
	
	if err == nil {
		// Товар уже в избранном
		return &existingItem, nil
	}
	
	// Товара нет, создаем новый
	userUUID, _ := uuid.Parse(userID)
	productUUID, _ := uuid.Parse(productID)
	item := models.WishlistItem{
		UserID:    userUUID,
		ProductID: productUUID,
	}
	
	err = r.db.Create(&item).Error
	if err != nil {
		return nil, err
	}
	
	// Загружаем связанные данные
	err = r.db.Preload("Product").First(&item, item.ID).Error
	return &item, err
}

func (r *wishlistRepository) RemoveItem(id string, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.WishlistItem{}).Error
}

func (r *wishlistRepository) IsInWishlist(userID string, productID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.WishlistItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	return count > 0, err
}

func (r *wishlistRepository) Clear(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.WishlistItem{}).Error
}
