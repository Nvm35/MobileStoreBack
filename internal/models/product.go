package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	ID               uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name             string          `json:"name" gorm:"not null" validate:"required,min=2"`
	Slug             string          `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	Description      string          `json:"description" gorm:"type:text"`
	ShortDescription string          `json:"short_description" gorm:"type:varchar(500)"`
	Price            float64         `json:"price" gorm:"not null" validate:"required,min=0"`
	ComparePrice     *float64        `json:"compare_price" gorm:"type:decimal(10,2)"`
	SKU              string          `json:"sku" gorm:"uniqueIndex;not null" validate:"required"`
	Stock            int             `json:"stock" gorm:"not null;default:0" validate:"min=0"`
	IsActive         bool            `json:"is_active" gorm:"default:true"`
	IsFeatured       bool            `json:"is_featured" gorm:"default:false"`
	IsNew            bool            `json:"is_new" gorm:"default:false"`
	Weight           *float64        `json:"weight" gorm:"type:decimal(8,2)" validate:"min=0"`
	Dimensions       string          `json:"dimensions" gorm:"type:varchar(100)"` // "10x5x2 cm"
	Brand            string          `json:"brand" gorm:"not null" validate:"required,min=2"`
	Model            string          `json:"model"`
	Color            string          `json:"color" gorm:"type:varchar(100)"`
	Material         string          `json:"material" gorm:"type:varchar(255)"`
	CategoryID       uuid.UUID       `json:"category_id" gorm:"type:uuid;not null"`
	Tags             pq.StringArray  `json:"tags" gorm:"type:text[]"`
	MetaTitle        string          `json:"meta_title" gorm:"type:varchar(255)"`
	MetaDescription  string          `json:"meta_description" gorm:"type:text"`
	ViewCount        int             `json:"view_count" gorm:"default:0"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `json:"-" gorm:"index"`

	// Связи
	Category   Category    `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Images     []Image     `json:"images,omitempty" gorm:"foreignKey:ProductID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:ProductID"`
}

type Category struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name          string     `json:"name" gorm:"not null;uniqueIndex" validate:"required,min=2"`
	Description   string     `json:"description" gorm:"type:text"`
	Slug          string     `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`
	SortOrder     int        `json:"sort_order" gorm:"default:0"`
	ImageURL      string     `json:"image_url" gorm:"type:text"`
	MetaTitle     string     `json:"meta_title" gorm:"type:varchar(255)"`
	MetaDescription string   `json:"meta_description" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Связи
	Products []Product  `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

type Image struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID         uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	CloudinaryPublicID string  `json:"cloudinary_public_id" gorm:"type:varchar(255);not null"`
	URL               string   `json:"url" gorm:"not null" validate:"required,url"`
	Alt               string   `json:"alt" gorm:"not null" validate:"required,min=2"`
	IsPrimary         bool     `json:"is_primary" gorm:"default:false"`
	SortOrder         int      `json:"sort_order" gorm:"default:0"`
	CreatedAt         time.Time `json:"created_at"`

	// Связи
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}
