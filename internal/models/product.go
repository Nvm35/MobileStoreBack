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
	IsActive    bool            `json:"is_active" gorm:"default:true"`
	Feature     bool            `json:"feature" gorm:"default:false"` // Флаг особенного товара для витрины
	Brand       string          `json:"brand" gorm:"not null" validate:"required,min=2"`
	Model       string          `json:"model"`
	Material    string          `json:"material" gorm:"type:varchar(255)"`
	CategoryID  uuid.UUID       `json:"-" gorm:"type:uuid;not null"` // Скрываем UUID категории
	Tags        pq.StringArray  `json:"tags" gorm:"type:text[]"`
	VideoURL    *string         `json:"video_url" gorm:"type:text" validate:"omitempty,url"` // Ссылка на видео товара
	ViewCount   int             `json:"view_count" gorm:"default:0"`
	CreatedAt   time.Time       `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"type:timestamp"`

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
	CreatedAt     time.Time  `json:"created_at" gorm:"type:timestamp"`

	// Связи
	Products []Product  `json:"-" gorm:"foreignKey:CategoryID"` // Скрываем связи
}

type ProductVariant struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"-" gorm:"type:uuid;not null"` // Скрываем UUID в JSON
	SKU       string    `json:"sku" gorm:"uniqueIndex;not null" validate:"required"` // SKU как основной идентификатор
	Name      string    `json:"name" gorm:"not null" validate:"required,min=2"`
	Color     string    `json:"color" gorm:"type:varchar(100)"`
	Size      string    `json:"size" gorm:"type:varchar(50)"`
	Price     float64   `json:"price" gorm:"not null" validate:"required,min=0"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`

	// Связи (не загружаем полный Product, только slug если нужно)
	Product     Product           `json:"-" gorm:"foreignKey:ProductID"` // Скрываем полный объект Product
	WarehouseStocks []WarehouseStock `json:"-" gorm:"foreignKey:ProductVariantID"`
	OrderItems  []OrderItem       `json:"-" gorm:"foreignKey:ProductVariantID"`
	
	// Вычисляемое поле для API (заполняется в сервисе/обработчике)
	ProductSlug string `json:"product_slug,omitempty" gorm:"-"`
}

// Warehouse - склад/филиал
type Warehouse struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null" validate:"required,min=2"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	Address     string    `json:"address" gorm:"not null" validate:"required"`
	City        string    `json:"city" gorm:"not null" validate:"required"`
	Phone       string    `json:"phone" gorm:"type:varchar(20)"`
	Email       string    `json:"email" gorm:"type:varchar(255)"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	IsMain      bool      `json:"is_main" gorm:"default:false"` // главный склад
	ManagerName string    `json:"manager_name" gorm:"type:varchar(255)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Связи
	WarehouseStocks []WarehouseStock `json:"warehouse_stocks,omitempty" gorm:"foreignKey:WarehouseID"`
	Orders          []Order          `json:"orders,omitempty" gorm:"foreignKey:WarehouseID"`
}

// WarehouseStock - остатки товаров по складам
type WarehouseStock struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WarehouseID      uuid.UUID      `json:"warehouse_id" gorm:"type:uuid;not null"`
	ProductVariantID uuid.UUID      `json:"product_variant_id" gorm:"type:uuid;not null"`
	Stock            int            `json:"stock" gorm:"not null;default:0" validate:"min=0"`
	ReservedStock    int            `json:"reserved_stock" gorm:"not null;default:0" validate:"min=0"` // зарезервированный товар
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`

	// Связи
	Warehouse      Warehouse      `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	ProductVariant ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID"`

	// Уникальный индекс для комбинации склада и варианта товара
	// UNIQUE(warehouse_id, product_variant_id)
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
