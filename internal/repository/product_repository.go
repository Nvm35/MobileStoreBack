package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"mobile-store-back/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type productRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProductRepository(db *gorm.DB, redis *redis.Client) ProductRepository {
	return &productRepository{
		db:    db,
		redis: redis,
	}
}

func (r *productRepository) Create(name string, slug string, description string, basePrice float64, sku string, stock int, isActive bool, brand string, model string, material string, categoryID string, tags []string) (*models.Product, error) {
	categoryUUID, _ := uuid.Parse(categoryID)

	product := models.Product{
		Name:        name,
		Slug:        slug,
		Description: description,
		BasePrice:   basePrice,
		SKU:         sku,
		Stock:       stock,
		IsActive:    isActive,
		Brand:       brand,
		Model:       model,
		Material:    material,
		CategoryID:  categoryUUID,
		Tags:        tags,
	}

	err := r.db.Create(&product).Error
	return &product, err
}

func (r *productRepository) GetByID(id string) (*models.Product, error) {
	// Попробуем получить из кэша
	cacheKey := fmt.Sprintf("product:%s", id)
	cached, err := r.redis.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var product models.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return &product, nil
		}
	}

	// Получаем из базы данных
	var product models.Product
	if err := r.db.Preload("Category").Preload("Images").First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	productJSON, _ := json.Marshal(product)
	r.redis.Set(context.Background(), cacheKey, productJSON, time.Hour)

	return &product, nil
}

func (r *productRepository) GetBySKU(sku string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Category").Preload("Images").First(&product, "sku = ?", sku).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(id string, name *string, description *string, basePrice *float64, stock *int, isActive *bool, brand *string, model *string, material *string, categoryID *string, tags []string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}

	if name != nil {
		product.Name = *name
	}
	if description != nil {
		product.Description = *description
	}
	if basePrice != nil {
		product.BasePrice = *basePrice
	}
	if stock != nil {
		product.Stock = *stock
	}
	if isActive != nil {
		product.IsActive = *isActive
	}
	if brand != nil {
		product.Brand = *brand
	}
	if model != nil {
		product.Model = *model
	}
	if material != nil {
		product.Material = *material
	}
	if categoryID != nil {
		if categoryUUID, err := uuid.Parse(*categoryID); err == nil {
			product.CategoryID = categoryUUID
		}
	}
	if len(tags) > 0 {
		product.Tags = tags
	}

	// Обновляем в базе данных
	err = r.db.Save(&product).Error
	if err != nil {
		return nil, err
	}

	// Удаляем из кэша
	cacheKey := fmt.Sprintf("product:%s", product.ID.String())
	r.redis.Del(context.Background(), cacheKey)

	return &product, nil
}

func (r *productRepository) Delete(id string) error {
	// Удаляем из кэша
	cacheKey := fmt.Sprintf("product:%s", id)
	r.redis.Del(context.Background(), cacheKey)

	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}

func (r *productRepository) List(limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	if err := r.db.Preload("Category").Preload("Images").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Search(query string, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	searchQuery := "%" + strings.ToLower(query) + "%"
	
	if err := r.db.Preload("Category").Preload("Images").
		Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(brand) LIKE ? OR LOWER(model) LIKE ?", 
			searchQuery, searchQuery, searchQuery, searchQuery).
		Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetByCategory(categoryID string, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	if err := r.db.Preload("Category").Preload("Images").
		Where("category_id = ?", categoryID).
		Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetBySlug получает товар по slug
func (r *productRepository) GetBySlug(slug string) (*models.Product, error) {
	// Попробуем получить из кэша
	cacheKey := fmt.Sprintf("product:slug:%s", slug)
	cached, err := r.redis.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var product models.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return &product, nil
		}
	}

	// Получаем из базы данных
	var product models.Product
	if err := r.db.Preload("Category").Preload("Images").First(&product, "slug = ?", slug).Error; err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	productJSON, _ := json.Marshal(product)
	r.redis.Set(context.Background(), cacheKey, productJSON, time.Hour)

	return &product, nil
}
