package models

import (
	"time"

	"github.com/google/uuid"
)

// Модель ShippingMethod удалена - способы доставки теперь встроены в Order

// CartItem - элементы корзины
type CartItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	SessionID string    `json:"session_id" gorm:"type:varchar(255)"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Quantity  int       `json:"quantity" gorm:"not null" validate:"required,min=1"`
	Price     float64   `json:"price" gorm:"not null" validate:"min=0"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Связи
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// WishlistItem - элементы избранного
type WishlistItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`

	// Связи
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// Coupon - промокоды
type Coupon struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code          string     `json:"code" gorm:"uniqueIndex;not null" validate:"required,min=3"`
	Name          string     `json:"name" gorm:"not null" validate:"required,min=2"`
	Description   string     `json:"description" gorm:"type:text"`
	Type          string     `json:"type" gorm:"not null" validate:"required,oneof=percentage fixed"`
	Value         float64    `json:"value" gorm:"not null" validate:"required,min=0"`
	MinimumAmount float64    `json:"minimum_amount" gorm:"default:0" validate:"min=0"`
	MaximumDiscount *float64 `json:"maximum_discount" gorm:"type:decimal(10,2)"`
	UsageLimit    *int       `json:"usage_limit"`
	UsedCount     int        `json:"used_count" gorm:"default:0"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`
	StartsAt      *time.Time `json:"starts_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Связи
	CouponUsages []CouponUsage `json:"coupon_usages,omitempty" gorm:"foreignKey:CouponID"`
}

// CouponUsage - использование промокодов
type CouponUsage struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CouponID      uuid.UUID `json:"coupon_id" gorm:"type:uuid;not null"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	OrderID       uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	DiscountAmount float64  `json:"discount_amount" gorm:"not null" validate:"min=0"`
	UsedAt        time.Time `json:"used_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Связи
	Coupon Coupon `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
	User   User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Order  Order  `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// Review - отзывы
type Review struct {
	ID             uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         uuid.UUID     `json:"user_id" gorm:"type:uuid;not null"`
	ProductID      uuid.UUID     `json:"product_id" gorm:"type:uuid;not null"`
	OrderID        *uuid.UUID    `json:"order_id" gorm:"type:uuid"`
	Rating         int           `json:"rating" gorm:"not null" validate:"required,min=1,max=5"`
	Title          string        `json:"title" gorm:"type:varchar(255)"`
	Comment        string        `json:"comment" gorm:"type:text"`
	IsVerified     bool          `json:"is_verified" gorm:"default:false"`
	IsApproved     bool          `json:"is_approved" gorm:"default:true"`
	HelpfulCount   int           `json:"helpful_count" gorm:"default:0"`
	UnhelpfulCount int           `json:"unhelpful_count" gorm:"default:0"`
	HelpfulVotes   string `json:"helpful_votes" gorm:"type:jsonb;default:'[]'"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`

	// Связи
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Order   *Order  `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}
