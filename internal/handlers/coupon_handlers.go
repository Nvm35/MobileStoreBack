package handlers

import (
	"mobile-store-back/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// GetCoupons - получение списка промокодов
func GetCoupons(couponService *services.CouponService) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		coupons, err := couponService.GetAll(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"coupons": coupons})
	}
}

// GetCoupon - получение конкретного промокода
func GetCoupon(couponService *services.CouponService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		coupon, err := couponService.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
			return
		}

		c.JSON(http.StatusOK, coupon)
	}
}

// ValidateCoupon - валидация промокода
func ValidateCoupon(couponService *services.CouponService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		
		var req struct {
			Code        string  `json:"code" validate:"required"`
			OrderAmount float64 `json:"order_amount" validate:"required,min=0"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидация
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validation, err := couponService.Validate(req.Code, userID.(string), req.OrderAmount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, validation)
	}
}

// CreateCoupon - создание промокода (админ)
func CreateCoupon(couponService *services.CouponService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Code            string    `json:"code" validate:"required,min=3"`
			Name            string    `json:"name" validate:"required,min=2"`
			Description     string    `json:"description"`
			Type            string    `json:"type" validate:"required,oneof=percentage fixed"`
			Value           float64   `json:"value" validate:"required,min=0"`
			MinimumAmount   float64   `json:"minimum_amount" validate:"min=0"`
			MaximumDiscount *float64  `json:"maximum_discount" validate:"omitempty,min=0"`
			UsageLimit      *int      `json:"usage_limit" validate:"omitempty,min=1"`
			StartsAt        *string   `json:"starts_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
			ExpiresAt       *string   `json:"expires_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидация
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		coupon, err := couponService.Create(req.Code, req.Name, req.Description, req.Type, req.Value, req.MinimumAmount, req.MaximumDiscount, req.UsageLimit, req.StartsAt, req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, coupon)
	}
}

// UpdateCoupon - обновление промокода (админ)
func UpdateCoupon(couponService *services.CouponService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		var req struct {
			Name            *string   `json:"name" validate:"omitempty,min=2"`
			Description     *string   `json:"description"`
			Value           *float64  `json:"value" validate:"omitempty,min=0"`
			MinimumAmount   *float64  `json:"minimum_amount" validate:"omitempty,min=0"`
			MaximumDiscount *float64  `json:"maximum_discount" validate:"omitempty,min=0"`
			UsageLimit      *int      `json:"usage_limit" validate:"omitempty,min=1"`
			IsActive        *bool     `json:"is_active"`
			StartsAt        *string   `json:"starts_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
			ExpiresAt       *string   `json:"expires_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидация
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		coupon, err := couponService.Update(id, req.Name, req.Description, req.Value, req.MinimumAmount, req.MaximumDiscount, req.UsageLimit, req.IsActive, req.StartsAt, req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, coupon)
	}
}

// DeleteCoupon - удаление промокода (админ)
func DeleteCoupon(couponService *services.CouponService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		err := couponService.Delete(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Coupon deleted successfully"})
	}
}

// GetCouponUsage - получение статистики использования промокода (админ)
func GetCouponUsage(couponService *services.CouponService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		usage, err := couponService.GetUsage(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, usage)
	}
}
