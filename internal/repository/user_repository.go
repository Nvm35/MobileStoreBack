package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"mobile-store-back/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) UserRepository {
	return &userRepository{
		db:    db,
		redis: redis,
	}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	// Попробуем получить из кэша
	cacheKey := fmt.Sprintf("user:%s", id)
	cached, err := r.redis.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	// Получаем из базы данных
	var user models.User
	if err := r.db.Preload("Addresses").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	userJSON, _ := json.Marshal(user)
	r.redis.Set(context.Background(), cacheKey, userJSON, time.Hour)

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Addresses").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateProfile(userID string, firstName *string, lastName *string, phone *string, dateOfBirth *string, gender *string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	
	if firstName != nil {
		user.FirstName = *firstName
	}
	if lastName != nil {
		user.LastName = *lastName
	}
	if phone != nil {
		user.Phone = *phone
	}
	if dateOfBirth != nil {
		if t, err := time.Parse("2006-01-02", *dateOfBirth); err == nil {
			user.DateOfBirth = &t
		}
	}
	if gender != nil {
		user.Gender = *gender
	}
	
	err = r.db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	// Удаляем из кэша
	cacheKey := fmt.Sprintf("user:%s", user.ID.String())
	r.redis.Del(context.Background(), cacheKey)

	return &user, nil
}

func (r *userRepository) Update(id string, firstName *string, lastName *string, phone *string, dateOfBirth *string, gender *string, isActive *bool, isAdmin *bool) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	
	if firstName != nil {
		user.FirstName = *firstName
	}
	if lastName != nil {
		user.LastName = *lastName
	}
	if phone != nil {
		user.Phone = *phone
	}
	if dateOfBirth != nil {
		if t, err := time.Parse("2006-01-02", *dateOfBirth); err == nil {
			user.DateOfBirth = &t
		}
	}
	if gender != nil {
		user.Gender = *gender
	}
	if isActive != nil {
		user.IsActive = *isActive
	}
	if isAdmin != nil {
		user.IsAdmin = *isAdmin
	}
	
	err = r.db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	// Удаляем из кэша
	cacheKey := fmt.Sprintf("user:%s", user.ID.String())
	r.redis.Del(context.Background(), cacheKey)

	return &user, nil
}

func (r *userRepository) Delete(id string) error {
	// Удаляем из кэша
	cacheKey := fmt.Sprintf("user:%s", id)
	r.redis.Del(context.Background(), cacheKey)

	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) List(limit, offset int) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
