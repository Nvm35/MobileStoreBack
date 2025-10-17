package repository

import (
	"mobile-store-back/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewCategoryRepository(db *gorm.DB, redis *redis.Client) CategoryRepository {
	return CategoryRepository{
		db:    db,
		redis: redis,
	}
}

func (r *CategoryRepository) GetAll(limit, offset int) ([]*models.Category, error) {
	var categories []*models.Category
	
	query := r.db.Where("is_active = ?", true).Order("sort_order ASC, name ASC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetByID(id string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("slug = ? AND is_active = ?", slug, true).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) Update(id string, name *string, description *string, slug *string, isActive *bool, sortOrder *int, imageURL *string, metaTitle *string, metaDescription *string) (*models.Category, error) {
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
	if isActive != nil {
		updates["is_active"] = *isActive
	}
	if sortOrder != nil {
		updates["sort_order"] = *sortOrder
	}
	if imageURL != nil {
		updates["image_url"] = *imageURL
	}
	if metaTitle != nil {
		updates["meta_title"] = *metaTitle
	}
	if metaDescription != nil {
		updates["meta_description"] = *metaDescription
	}
	
	err := r.db.Model(&category).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	
	return r.GetByID(id)
}

func (r *CategoryRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Category{}).Error
}

func (r *CategoryRepository) GetWithProducts(id string, limit, offset int) (*models.Category, error) {
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
	}).Where("id = ? AND is_active = ?", id, true)
	
	err := query.First(&category).Error
	if err != nil {
		return nil, err
	}
	
	return &category, nil
}

