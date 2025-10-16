package repository

import (
	"mobile-store-back/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type couponRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewCouponRepository(db *gorm.DB, redis *redis.Client) CouponRepository {
	return &couponRepository{
		db:    db,
		redis: redis,
	}
}

func (r *couponRepository) GetAll(limit, offset int) ([]models.Coupon, error) {
	var coupons []models.Coupon
	err := r.db.Where("is_active = ?", true).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&coupons).Error
	return coupons, err
}

func (r *couponRepository) GetByID(id string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.db.Where("id = ?", id).First(&coupon).Error
	return &coupon, err
}

func (r *couponRepository) Validate(code string, userID string, orderAmount float64) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.db.Where("code = ? AND is_active = ?", code, true).First(&coupon).Error
	if err != nil {
		return nil, err
	}
	
	// Проверяем срок действия
	now := time.Now()
	if coupon.StartsAt != nil && now.Before(*coupon.StartsAt) {
		return nil, gorm.ErrRecordNotFound
	}
	if coupon.ExpiresAt != nil && now.After(*coupon.ExpiresAt) {
		return nil, gorm.ErrRecordNotFound
	}
	
	// Проверяем минимальную сумму заказа
	if orderAmount < coupon.MinimumAmount {
		return nil, gorm.ErrRecordNotFound
	}
	
	// Проверяем лимит использований
	if coupon.UsageLimit != nil && coupon.UsedCount >= *coupon.UsageLimit {
		return nil, gorm.ErrRecordNotFound
	}
	
	return &coupon, nil
}

func (r *couponRepository) Create(code string, name string, description string, couponType string, value float64, minimumAmount float64, maximumDiscount *float64, usageLimit *int, startsAt *string, expiresAt *string) (*models.Coupon, error) {
	coupon := models.Coupon{
		Code:            code,
		Name:            name,
		Description:     description,
		Type:            couponType,
		Value:           value,
		MinimumAmount:   minimumAmount,
		MaximumDiscount: maximumDiscount,
		UsageLimit:      usageLimit,
		IsActive:        true,
	}
	
	// Парсим даты если они переданы
	if startsAt != nil && *startsAt != "" {
		if t, err := time.Parse(time.RFC3339, *startsAt); err == nil {
			coupon.StartsAt = &t
		}
	}
	if expiresAt != nil && *expiresAt != "" {
		if t, err := time.Parse(time.RFC3339, *expiresAt); err == nil {
			coupon.ExpiresAt = &t
		}
	}
	
	err := r.db.Create(&coupon).Error
	return &coupon, err
}

func (r *couponRepository) Update(id string, name *string, description *string, value *float64, minimumAmount *float64, maximumDiscount *float64, usageLimit *int, isActive *bool, startsAt *string, expiresAt *string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.db.Where("id = ?", id).First(&coupon).Error
	if err != nil {
		return nil, err
	}
	
	if name != nil {
		coupon.Name = *name
	}
	if description != nil {
		coupon.Description = *description
	}
	if value != nil {
		coupon.Value = *value
	}
	if minimumAmount != nil {
		coupon.MinimumAmount = *minimumAmount
	}
	if maximumDiscount != nil {
		coupon.MaximumDiscount = maximumDiscount
	}
	if usageLimit != nil {
		coupon.UsageLimit = usageLimit
	}
	if isActive != nil {
		coupon.IsActive = *isActive
	}
	
	// Парсим даты если они переданы
	if startsAt != nil && *startsAt != "" {
		if t, err := time.Parse(time.RFC3339, *startsAt); err == nil {
			coupon.StartsAt = &t
		}
	}
	if expiresAt != nil && *expiresAt != "" {
		if t, err := time.Parse(time.RFC3339, *expiresAt); err == nil {
			coupon.ExpiresAt = &t
		}
	}
	
	err = r.db.Save(&coupon).Error
	return &coupon, err
}

func (r *couponRepository) Delete(id string) error {
	return r.db.Delete(&models.Coupon{}, "id = ?", id).Error
}

func (r *couponRepository) GetUsage(id string) ([]models.CouponUsage, error) {
	var usages []models.CouponUsage
	err := r.db.Where("coupon_id = ?", id).
		Preload("User").
		Preload("Order").
		Order("used_at DESC").
		Find(&usages).Error
	return usages, err
}
