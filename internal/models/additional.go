package models

import (
	"time"

	"github.com/google/uuid"
)

// Модель ShippingMethod удалена - способы доставки теперь встроены в Order

// CartItem - элементы корзины
type CartItem struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"` // ID нужен для удаления/обновления конкретного элемента
	UserID          uuid.UUID       `json:"-" gorm:"type:uuid;not null;index"` // Скрываем UUID пользователя
	SessionID       *string         `json:"-" gorm:"type:varchar(255)"` // Не используется в текущей логике
	ProductID       uuid.UUID       `json:"-" gorm:"type:uuid;not null;index"` // Скрываем UUID товара
	ProductVariantID *uuid.UUID     `json:"-" gorm:"type:uuid;index"` // Скрываем UUID варианта
	Quantity        int             `json:"quantity" gorm:"not null" validate:"required,min=1"`
	Price           float64         `json:"price" gorm:"type:decimal(10,2);not null" validate:"min=0"`
	ExpiresAt       *time.Time      `json:"-" gorm:"type:timestamp"`
	CreatedAt       time.Time       `json:"-" gorm:"type:timestamp"`
	UpdatedAt       time.Time       `json:"-" gorm:"type:timestamp"`

	// Связи (загружаем только нужные поля)
	User          User           `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Product       Product        `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"`
	ProductVariant *ProductVariant `json:"variant,omitempty" gorm:"foreignKey:ProductVariantID;references:ID"`
	
	// Вычисляемые поля для API (заполняются в сервисе/обработчике)
	ProductSlug   string  `json:"product_slug,omitempty" gorm:"-"`
	VariantSKU    string  `json:"variant_sku,omitempty" gorm:"-"`
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
