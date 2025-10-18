package repository

import (
	"mobile-store-back/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ImageRepository interface {
	Create(image *models.Image) error
	GetByID(id string) (*models.Image, error)
	GetByProductID(productID string) ([]*models.Image, error)
	GetByProductSlug(productSlug string) ([]*models.Image, error)
	GetPrimaryByProductID(productID string) (*models.Image, error)
	Update(id string, cloudinaryPublicID *string, url *string, isPrimary *bool) (*models.Image, error)
	SetPrimary(id string) error
	UnsetPrimaryForProduct(productID string) error
	Delete(id string) error
}

type imageRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewImageRepository(db *gorm.DB, redis *redis.Client) ImageRepository {
	return &imageRepository{
		db:    db,
		redis: redis,
	}
}

func (r *imageRepository) Create(image *models.Image) error {
	return r.db.Create(image).Error
}

func (r *imageRepository) GetByID(id string) (*models.Image, error) {
	var image models.Image
	err := r.db.Preload("Product").First(&image, "id = ?", id).Error
	return &image, err
}

func (r *imageRepository) GetByProductID(productID string) ([]*models.Image, error) {
	var images []*models.Image
	err := r.db.Where("product_id = ?", productID).Order("is_primary DESC, created_at ASC").Find(&images).Error
	return images, err
}

func (r *imageRepository) GetByProductSlug(productSlug string) ([]*models.Image, error) {
	var images []*models.Image
	err := r.db.Joins("JOIN products ON images.product_id = products.id").
		Where("products.slug = ?", productSlug).
		Order("images.is_primary DESC, images.created_at ASC").
		Find(&images).Error
	return images, err
}

func (r *imageRepository) GetPrimaryByProductID(productID string) (*models.Image, error) {
	var image models.Image
	err := r.db.Where("product_id = ? AND is_primary = ?", productID, true).First(&image).Error
	return &image, err
}

func (r *imageRepository) Update(id string, cloudinaryPublicID *string, url *string, isPrimary *bool) (*models.Image, error) {
	var image models.Image
	if err := r.db.First(&image, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if cloudinaryPublicID != nil {
		image.CloudinaryPublicID = *cloudinaryPublicID
	}
	if url != nil {
		image.URL = *url
	}
	if isPrimary != nil {
		image.IsPrimary = *isPrimary
	}

	if err := r.db.Save(&image).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *imageRepository) SetPrimary(id string) error {
	// Сначала получаем изображение
	var image models.Image
	if err := r.db.First(&image, "id = ?", id).Error; err != nil {
		return err
	}

	// Убираем primary с других изображений этого товара
	if err := r.db.Model(&models.Image{}).
		Where("product_id = ? AND id != ?", image.ProductID, id).
		Update("is_primary", false).Error; err != nil {
		return err
	}

	// Устанавливаем primary для выбранного изображения
	return r.db.Model(&image).Update("is_primary", true).Error
}

func (r *imageRepository) UnsetPrimaryForProduct(productID string) error {
	return r.db.Model(&models.Image{}).
		Where("product_id = ?", productID).
		Update("is_primary", false).Error
}

func (r *imageRepository) Delete(id string) error {
	return r.db.Delete(&models.Image{}, "id = ?", id).Error
}
