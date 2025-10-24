package repository

import (
	"mobile-store-back/internal/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(name string, slug string, description string, basePrice float64, sku string, isActive bool, feature bool, brand string, model string, material string, categoryID string, tags []string, videoURL *string) (*models.Product, error) {
	categoryUUID, _ := uuid.Parse(categoryID)

	product := models.Product{
		Name:        name,
		Slug:        slug,
		Description: description,
		BasePrice:   basePrice,
		SKU:         sku,
		IsActive:    isActive,
		Feature:     feature,
		Brand:       brand,
		Model:       model,
		Material:    material,
		CategoryID:  categoryUUID,
		Tags:        tags,
		VideoURL:    videoURL,
	}

	err := r.db.Create(&product).Error
	if err != nil {
		return nil, err
	}
	
	return &product, nil
}

func (r *productRepository) GetByID(id string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Category").Preload("Images").First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) GetBySKU(sku string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Category").Preload("Images").First(&product, "sku = ?", sku).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(id string, name *string, description *string, basePrice *float64, isActive *bool, feature *bool, brand *string, model *string, material *string, categoryID *string, tags []string, videoURL *string) (*models.Product, error) {
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
	if isActive != nil {
		product.IsActive = *isActive
	}
	if feature != nil {
		product.Feature = *feature
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
	if videoURL != nil {
		product.VideoURL = videoURL
	}

	// Обновляем в базе данных
	err = r.db.Save(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Delete(id string) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}

func (r *productRepository) List() ([]*models.Product, error) {
	var products []*models.Product
	if err := r.db.Preload("Category").Preload("Images").Where("is_active = ?", true).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Search(query string) ([]*models.Product, error) {
	var products []*models.Product
	searchQuery := "%" + strings.ToLower(query) + "%"
	
	if err := r.db.Preload("Category").Preload("Images").
		Where("is_active = ? AND (LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(brand) LIKE ? OR LOWER(model) LIKE ?)", 
			true, searchQuery, searchQuery, searchQuery, searchQuery).
		Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetByCategory(categoryID string) ([]*models.Product, error) {
	var products []*models.Product
	if err := r.db.Preload("Category").Preload("Images").
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) ListWithFilters(brand, minPrice, maxPrice string) ([]*models.Product, error) {
	var products []*models.Product
	query := r.db.Preload("Category").Preload("Images").Where("is_active = ?", true)
	
	if brand != "" {
		query = query.Where("brand = ?", brand)
	}
	if minPrice != "" {
		query = query.Where("base_price >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("base_price <= ?", maxPrice)
	}
	
	if err := query.Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetBySlug получает товар по slug
func (r *productRepository) GetBySlug(slug string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Category").Preload("Images").Where("slug = ? AND is_active = ?", slug, true).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// GetFeatured получает все товары с feature=true
func (r *productRepository) GetFeatured() ([]*models.Product, error) {
	var products []*models.Product
	if err := r.db.Preload("Category").Preload("Images").
		Where("feature = ? AND is_active = ?", true, true).
		Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

