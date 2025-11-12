package models

import (
	"time"

	"github.com/google/uuid"
)

// Модель ShippingMethod удалена - способы доставки теперь встроены в Order

// CartItem - элементы корзины
type CartItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	SessionID *string   `json:"session_id,omitempty" gorm:"type:varchar(255)"` // Используем указатель, чтобы NULL в БД был пустой строкой
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Quantity  int       `json:"quantity" gorm:"not null" validate:"required,min=1"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null" validate:"min=0"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Связи
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"`
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
