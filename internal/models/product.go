package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null" validate:"required,min=2"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"not null" validate:"required,min=0"`
	SKU         string    `json:"sku" gorm:"uniqueIndex;not null" validate:"required"`
	Stock       int       `json:"stock" gorm:"not null;default:0" validate:"min=0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Weight      float64   `json:"weight" validate:"min=0"`
	Dimensions  string    `json:"dimensions"` // "10x5x2 cm"
	Brand       string    `json:"brand" validate:"required,min=2"`
	Model       string    `json:"model" validate:"required,min=2"`
	Color       string    `json:"color"`
	Material    string    `json:"material"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Связи
	CategoryID uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`
	Category   Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Images     []Image   `json:"images,omitempty" gorm:"foreignKey:ProductID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:ProductID"`
}

type Category struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex" validate:"required,min=2"`
	Description string    `json:"description" gorm:"type:text"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	ParentID    *uuid.UUID `json:"parent_id" gorm:"type:uuid"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Связи
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Products []Product  `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

type Image struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	URL       string    `json:"url" gorm:"not null" validate:"required,url"`
	Alt       string    `json:"alt" validate:"required,min=2"`
	IsPrimary bool      `json:"is_primary" gorm:"default:false"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`

	// Связи
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}
