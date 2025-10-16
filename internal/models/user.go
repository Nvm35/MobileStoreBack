package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                     uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email                  string          `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password               string          `json:"-" gorm:"not null" validate:"required,min=6"`
	FirstName              string          `json:"first_name" gorm:"not null" validate:"required,min=2"`
	LastName               string          `json:"last_name" gorm:"not null" validate:"required,min=2"`
	Phone                  string          `json:"phone" gorm:"type:varchar(20)" validate:"omitempty,e164"`
	DateOfBirth            *time.Time      `json:"date_of_birth" gorm:"type:date"`
	Gender                 string          `json:"gender" gorm:"type:varchar(10)" validate:"omitempty,oneof=male female"`
	IsActive               bool            `json:"is_active" gorm:"default:true"`
	IsAdmin                bool            `json:"is_admin" gorm:"default:false"`
	EmailVerified          bool            `json:"email_verified" gorm:"default:false"`
	EmailVerificationToken string          `json:"-" gorm:"type:varchar(255)"`
	PasswordResetToken     string          `json:"-" gorm:"type:varchar(255)"`
	PasswordResetExpires   *time.Time      `json:"-"`
	LastLogin              *time.Time      `json:"last_login"`
	Notifications          string          `json:"notifications" gorm:"type:jsonb;default:'[]'"`
	// Адрес пользователя (встроенный)
	AddressTitle           string          `json:"address_title" gorm:"type:varchar(255)"`
	AddressFirstName        string          `json:"address_first_name" gorm:"type:varchar(255)"`
	AddressLastName         string          `json:"address_last_name" gorm:"type:varchar(255)"`
	AddressCompany          string          `json:"address_company" gorm:"type:varchar(255)"`
	AddressStreet           string          `json:"address_street" gorm:"type:text"`
	AddressCity             string          `json:"address_city" gorm:"type:varchar(255)"`
	AddressState            string          `json:"address_state" gorm:"type:varchar(255)"`
	AddressPostalCode       string          `json:"address_postal_code" gorm:"type:varchar(20)"`
	AddressCountry          string          `json:"address_country" gorm:"type:varchar(255)"`
	AddressPhone            string          `json:"address_phone" gorm:"type:varchar(20)"`
	CreatedAt              time.Time       `json:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `json:"-" gorm:"index"`

	// Связи
	Orders    []Order   `json:"orders,omitempty" gorm:"foreignKey:UserID"`
}

// Модель Address удалена - адрес теперь встроен в User
