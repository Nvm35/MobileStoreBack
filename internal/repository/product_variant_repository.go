package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"mobile-store-back/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type productVariantRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProductVariantRepository(db *gorm.DB, redis *redis.Client) ProductVariantRepository {
	return &productVariantRepository{
		db:    db,
		redis: redis,
	}
}

func (r *productVariantRepository) Create(productID string, sku string, name string, color string, size string, price float64, isActive bool) (*models.ProductVariant, error) {
	productUUID, _ := uuid.Parse(productID)

	variant := models.ProductVariant{
		ProductID: productUUID,
		SKU:       sku,
		Name:      name,
		Color:     color,
		Size:      size,
		Price:     price,
		IsActive:  isActive,
	}

	err := r.db.Create(&variant).Error
	return &variant, err
}

func (r *productVariantRepository) GetByID(id string) (*models.ProductVariant, error) {
	// Попробуем получить из кэша
	cacheKey := fmt.Sprintf("product_variant:%s", id)
	cached, err := r.redis.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var variant models.ProductVariant
		if json.Unmarshal([]byte(cached), &variant) == nil {
			return &variant, nil
		}
	}

	var variant models.ProductVariant
	err = r.db.Where("id = ?", id).First(&variant).Error
	if err != nil {
		return nil, err
	}

	// Кэшируем на 1 час
	if data, err := json.Marshal(variant); err == nil {
		r.redis.Set(context.Background(), cacheKey, data, time.Hour)
	}

	return &variant, nil
}

func (r *productVariantRepository) GetBySKU(sku string) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.Where("sku = ?", sku).First(&variant).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *productVariantRepository) GetByProductID(productID string) ([]*models.ProductVariant, error) {
	var variants []*models.ProductVariant
	err := r.db.Where("product_id = ?", productID).Find(&variants).Error
	return variants, err
}

func (r *productVariantRepository) Update(id string, sku *string, name *string, color *string, size *string, price *float64, isActive *bool) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.Where("id = ?", id).First(&variant).Error
	if err != nil {
		return nil, err
	}

	if sku != nil {
		variant.SKU = *sku
	}
	if name != nil {
		variant.Name = *name
	}
	if color != nil {
		variant.Color = *color
	}
	if size != nil {
		variant.Size = *size
	}
	if price != nil {
		variant.Price = *price
	}
	if isActive != nil {
		variant.IsActive = *isActive
	}

	// Обновляем в базе данных
	err = r.db.Save(&variant).Error
	if err != nil {
		return nil, err
	}

	// Удаляем из кэша
	cacheKey := fmt.Sprintf("product_variant:%s", variant.ID.String())
	r.redis.Del(context.Background(), cacheKey)

	return &variant, nil
}

func (r *productVariantRepository) Delete(id string) error {
	// Удаляем из кэша
	cacheKey := fmt.Sprintf("product_variant:%s", id)
	r.redis.Del(context.Background(), cacheKey)

	return r.db.Delete(&models.ProductVariant{}, "id = ?", id).Error
}

