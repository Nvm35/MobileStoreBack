package repository

import (
	"mobile-store-back/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewCategoryRepository(db *gorm.DB, redis *redis.Client) CategoryRepository {
	return &categoryRepository{
		db:    db,
		redis: redis,
	}
}

func (r *categoryRepository) GetAll(limit, offset int) ([]*models.Category, error) {
	var categories []*models.Category
	
	query := r.db.Order("name ASC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetByID(id string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(id string, name *string, description *string, slug *string, imageURL *string) (*models.Category, error) {
	var category models.Category
	
	updates := make(map[string]interface{})
	if name != nil {
		updates["name"] = *name
	}
	if description != nil {
		updates["description"] = *description
	}
	if slug != nil {
		updates["slug"] = *slug
	}
	if imageURL != nil {
		updates["image_url"] = *imageURL
	}
	
	err := r.db.Model(&category).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	
	return r.GetByID(id)
}

func (r *categoryRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Category{}).Error
}

func (r *categoryRepository) GetWithProducts(id string, limit, offset int) (*models.Category, error) {
	var category models.Category
	
	// Загружаем категорию с продуктами
	query := r.db.Preload("Products", func(db *gorm.DB) *gorm.DB {
		productQuery := db.Where("is_active = ?", true)
		if limit > 0 {
			productQuery = productQuery.Limit(limit)
		}
		if offset > 0 {
			productQuery = productQuery.Offset(offset)
		}
		return productQuery.Order("created_at DESC")
	}).Where("id = ?", id)
	
	err := query.First(&category).Error
	if err != nil {
		return nil, err
	}
	
	return &category, nil
}

