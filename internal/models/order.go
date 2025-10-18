package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID     `json:"user_id" gorm:"type:uuid;not null"`
	OrderNumber     string        `json:"order_number" gorm:"uniqueIndex;not null"`
	Status          OrderStatus   `json:"status" gorm:"not null;default:'pending'"`
	TotalAmount     float64       `json:"total_amount" gorm:"not null" validate:"min=0"`
	PaymentMethod   string        `json:"payment_method" gorm:"not null" validate:"required,oneof=cash card transfer"`
	PaymentStatus   PaymentStatus `json:"payment_status" gorm:"not null;default:'pending'"`
	// Способ доставки
	ShippingMethod  string        `json:"shipping_method" gorm:"not null;default:'delivery'" validate:"required,oneof=delivery pickup"`
	// Адрес доставки (если нужен другой адрес, чем у пользователя)
	ShippingAddress string        `json:"shipping_address" gorm:"type:text"`
	// Пункт самовывоза (если выбран pickup)
	PickupPoint     string        `json:"pickup_point" gorm:"type:text"`
	TrackingNumber  string        `json:"tracking_number"`
	Notes           string        `json:"notes" gorm:"type:text"`
	CustomerNotes   string        `json:"customer_notes" gorm:"type:text"`
	ShippedAt       *time.Time    `json:"shipped_at"`
	DeliveredAt     *time.Time    `json:"delivered_at"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`

	// Связи
	User       User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID               uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID          uuid.UUID        `json:"order_id" gorm:"type:uuid;not null"`
	ProductID        uuid.UUID        `json:"product_id" gorm:"type:uuid;not null"`
	ProductVariantID *uuid.UUID       `json:"product_variant_id" gorm:"type:uuid"`
	Quantity         int              `json:"quantity" gorm:"not null" validate:"required,min=1"`
	Price            float64          `json:"price" gorm:"not null" validate:"min=0"`
	CreatedAt        time.Time        `json:"created_at"`

	// Связи
	Order          Order          `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Product        Product        `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	ProductVariant *ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID"`
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusReturned   OrderStatus = "returned"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)
