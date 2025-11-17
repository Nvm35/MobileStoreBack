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

func (r *wishlistRepository) GetByUserID(userID string) ([]models.WishlistItem, error) {
	var items []models.WishlistItem
	err := r.db.Where("user_id = ?", userID).
		Preload("Product").
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *wishlistRepository) AddItem(userID string, productIdentifier string) (*models.WishlistItem, error) {
	product, err := findProductByIdentifier(r.db, productIdentifier, true)
	if err != nil {
		return nil, err
	}

	var existingItem models.WishlistItem
	err = r.db.Where("user_id = ? AND product_id = ?", userID, product.ID).First(&existingItem).Error
	if err == nil {
		return &existingItem, nil
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	userUUID, _ := uuid.Parse(userID)
	item := models.WishlistItem{
		UserID:    userUUID,
		ProductID: product.ID,
	}

	if err := r.db.Create(&item).Error; err != nil {
		return nil, err
	}

	err = r.db.Preload("Product").First(&item, item.ID).Error
	return &item, err
}

func (r *wishlistRepository) RemoveItem(identifier string, userID string) error {
	if id, err := uuid.Parse(identifier); err == nil {
		if err := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.WishlistItem{}).Error; err == nil {
			return nil
		} else if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	productID, err := findProductIDByIdentifier(r.db, identifier)
	if err != nil {
		return err
	}

	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.WishlistItem{}).Error
}

func (r *wishlistRepository) IsInWishlist(userID string, productIdentifier string) (bool, error) {
	productID, err := findProductIDByIdentifier(r.db, productIdentifier)
	if err != nil {
		return false, err
	}
	var count int64
	err = r.db.Model(&models.WishlistItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	return count > 0, err
}

func (r *wishlistRepository) Clear(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.WishlistItem{}).Error
}
