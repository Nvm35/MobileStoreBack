package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"mobile-store-back/internal/models"
	"strings"
	"time"

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

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
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

func (r *productRepository) Update(product *models.Product) error {
	// Обновляем в базе данных
	if err := r.db.Save(product).Error; err != nil {
		return err
	}

	// Удаляем из кэша
	cacheKey := fmt.Sprintf("product:%s", product.ID.String())
	r.redis.Del(context.Background(), cacheKey)

	return nil
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
