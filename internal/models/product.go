package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Product struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string          `json:"name" gorm:"not null" validate:"required,min=2"`
	Slug        string          `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	Description string          `json:"description" gorm:"type:text"`
	BasePrice   float64         `json:"base_price" gorm:"not null" validate:"required,min=0"`
	SKU         string          `json:"sku" gorm:"uniqueIndex;not null" validate:"required"`
	Stock       int             `json:"stock" gorm:"not null;default:0" validate:"min=0"`
	IsActive    bool            `json:"is_active" gorm:"default:true"`
	Brand       string          `json:"brand" gorm:"not null" validate:"required,min=2"`
	Model       string          `json:"model"`
	Material    string          `json:"material" gorm:"type:varchar(255)"`
	CategoryID  uuid.UUID       `json:"category_id" gorm:"type:uuid;not null"`
	Tags        pq.StringArray  `json:"tags" gorm:"type:text[]"`
	ViewCount   int             `json:"view_count" gorm:"default:0"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`

	// Связи
	Category   Category         `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Variants   []ProductVariant `json:"variants,omitempty" gorm:"foreignKey:ProductID"`
	Images     []Image          `json:"images,omitempty" gorm:"foreignKey:ProductID"`
	OrderItems []OrderItem      `json:"order_items,omitempty" gorm:"foreignKey:ProductID"`
}

type Category struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name          string     `json:"name" gorm:"not null;uniqueIndex" validate:"required,min=2"`
	Description   string     `json:"description" gorm:"type:text"`
	Slug          string     `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	ImageURL      string     `json:"image_url" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at"`

	// Связи
	Products []Product  `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

type ProductVariant struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	SKU       string    `json:"sku" gorm:"uniqueIndex;not null" validate:"required"`
	Name      string    `json:"name" gorm:"not null" validate:"required,min=2"`
	Color     string    `json:"color" gorm:"type:varchar(100)"`
	Size      string    `json:"size" gorm:"type:varchar(50)"`
	Price     float64   `json:"price" gorm:"not null" validate:"required,min=0"`
	Stock     int       `json:"stock" gorm:"not null;default:0" validate:"min=0"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Связи
	Product   Product    `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:ProductVariantID"`
}

type Image struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID         uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	CloudinaryPublicID string  `json:"cloudinary_public_id" gorm:"type:varchar(255);not null"`
	URL               string   `json:"url" gorm:"not null" validate:"required,url"`
	IsPrimary         bool     `json:"is_primary" gorm:"default:false"`
	CreatedAt         time.Time `json:"created_at"`

	// Связи
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}
