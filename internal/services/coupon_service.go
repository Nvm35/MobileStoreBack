package services

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type CouponService struct {
	repo repository.CouponRepository
}

func NewCouponService(repo repository.CouponRepository) *CouponService {
	return &CouponService{repo: repo}
}

func (s *CouponService) GetAll(limit, offset int) ([]models.Coupon, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *CouponService) GetByID(id string) (*models.Coupon, error) {
	return s.repo.GetByID(id)
}

func (s *CouponService) Validate(code string, userID string, orderAmount float64) (*models.Coupon, error) {
	return s.repo.Validate(code, userID, orderAmount)
}

func (s *CouponService) Create(code string, name string, description string, couponType string, value float64, minimumAmount float64, maximumDiscount *float64, usageLimit *int, startsAt *string, expiresAt *string) (*models.Coupon, error) {
	return s.repo.Create(code, name, description, couponType, value, minimumAmount, maximumDiscount, usageLimit, startsAt, expiresAt)
}

func (s *CouponService) Update(id string, name *string, description *string, value *float64, minimumAmount *float64, maximumDiscount *float64, usageLimit *int, isActive *bool, startsAt *string, expiresAt *string) (*models.Coupon, error) {
	return s.repo.Update(id, name, description, value, minimumAmount, maximumDiscount, usageLimit, isActive, startsAt, expiresAt)
}

func (s *CouponService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *CouponService) GetUsage(id string) ([]models.CouponUsage, error) {
	return s.repo.GetUsage(id)
}
