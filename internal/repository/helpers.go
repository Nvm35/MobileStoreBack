package repository

import (
	"errors"
	"mobile-store-back/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func findProductByIdentifier(db *gorm.DB, identifier string, onlyActive bool) (*models.Product, error) {
	var product models.Product
	query := db.Model(&models.Product{})
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	if id, err := uuid.Parse(identifier); err == nil {
		if err := query.Where("id = ?", id).First(&product).Error; err == nil {
			return &product, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// Если это не UUID или запись по UUID не найдена, пробуем slug
	if err := query.Where("slug = ?", identifier).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func findProductIDByIdentifier(db *gorm.DB, identifier string) (uuid.UUID, error) {
	product, err := findProductByIdentifier(db, identifier, false)
	if err != nil {
		return uuid.UUID{}, err
	}
	return product.ID, nil
}

