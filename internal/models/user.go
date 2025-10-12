package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password  string    `json:"-" gorm:"not null" validate:"required,min=6"`
	FirstName string    `json:"first_name" gorm:"not null" validate:"required,min=2"`
	LastName  string    `json:"last_name" gorm:"not null" validate:"required,min=2"`
	Phone     string    `json:"phone" validate:"omitempty,e164"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Связи
	Addresses []Address `json:"addresses,omitempty" gorm:"foreignKey:UserID"`
	Orders    []Order   `json:"orders,omitempty" gorm:"foreignKey:UserID"`
}

type Address struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Title    string    `json:"title" gorm:"not null" validate:"required,min=2"`
	Address  string    `json:"address" gorm:"not null" validate:"required,min=10"`
	City     string    `json:"city" gorm:"not null" validate:"required,min=2"`
	PostalCode string  `json:"postal_code" gorm:"not null" validate:"required,min=5"`
	Country  string    `json:"country" gorm:"not null" validate:"required,min=2"`
	IsDefault bool     `json:"is_default" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Связи
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
