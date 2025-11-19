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
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	userJSON, _ := json.Marshal(user)
	r.redis.Set(context.Background(), cacheKey, userJSON, time.Hour)

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateProfile(
	userID string,
	firstName *string,
	lastName *string,
	phone *string,
	addressStreet *string,
	addressCity *string,
	addressState *string,
	addressPostalCode *string,
) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	
	// Обновление основных полей
	if firstName != nil {
		user.FirstName = *firstName
	}
	if lastName != nil {
		user.LastName = *lastName
	}
	if phone != nil {
		user.Phone = *phone
	}
	
	// Обновление полей адреса
	if addressStreet != nil {
		user.AddressStreet = *addressStreet
	}
	if addressCity != nil {
		user.AddressCity = *addressCity
	}
	if addressState != nil {
		user.AddressState = *addressState
	}
	if addressPostalCode != nil {
		user.AddressPostalCode = *addressPostalCode
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

func (r *userRepository) Update(id string, firstName *string, lastName *string, phone *string, isActive *bool, role *string) (*models.User, error) {
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
	if isActive != nil {
		user.IsActive = *isActive
	}
	if role != nil {
		// Проверяем, что роль валидна
		validRoles := map[string]bool{"admin": true, "manager": true, "customer": true}
		if !validRoles[*role] {
			return nil, fmt.Errorf("invalid role: %s", *role)
		}
		user.Role = *role
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

func (r *userRepository) List() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
