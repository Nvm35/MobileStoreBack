package repository

import (
	"mobile-store-back/internal/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type reviewRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewReviewRepository(db *gorm.DB, redis *redis.Client) ReviewRepository {
	return &reviewRepository{
		db:    db,
		redis: redis,
	}
}

func (r *reviewRepository) GetByProductID(productID string) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("product_id = ? AND is_approved = ?", productID, true).
		Preload("User").
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) Create(userID string, productID string, orderID *string, rating int, title string, comment string) (*models.Review, error) {
	userUUID, _ := uuid.Parse(userID)
	productUUID, _ := uuid.Parse(productID)
	
	var orderUUID *uuid.UUID
	if orderID != nil {
		if parsed, err := uuid.Parse(*orderID); err == nil {
			orderUUID = &parsed
		}
	}
	
	review := models.Review{
		UserID:    userUUID,
		ProductID: productUUID,
		OrderID:   orderUUID,
		Rating:    rating,
		Title:     title,
		Comment:   comment,
	}
	
	err := r.db.Create(&review).Error
	if err != nil {
		return nil, err
	}
	
	// Загружаем связанные данные
	err = r.db.Preload("User").Preload("Product").First(&review, review.ID).Error
	return &review, err
}

func (r *reviewRepository) Update(id string, userID string, rating *int, title *string, comment *string) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&review).Error
	if err != nil {
		return nil, err
	}
	
	if rating != nil {
		review.Rating = *rating
	}
	if title != nil {
		review.Title = *title
	}
	if comment != nil {
		review.Comment = *comment
	}
	
	err = r.db.Save(&review).Error
	if err != nil {
		return nil, err
	}
	
	// Загружаем связанные данные
	err = r.db.Preload("User").Preload("Product").First(&review, review.ID).Error
	return &review, err
}

func (r *reviewRepository) Delete(id string, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Review{}).Error
}

func (r *reviewRepository) Vote(id string, userID string, helpful bool) error {
	// TODO: Реализовать голосование за полезность отзыва
	// Это требует более сложной логики для работы с JSON полем helpful_votes
	return nil
}

func (r *reviewRepository) GetByUserID(userID string) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("user_id = ?", userID).
		Preload("Product").
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetAll() ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Preload("User").Preload("Product").
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) Approve(id string, approved bool) error {
	return r.db.Model(&models.Review{}).Where("id = ?", id).Update("is_approved", approved).Error
}
