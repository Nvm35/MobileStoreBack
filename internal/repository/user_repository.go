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

func (r *userRepository) Update(user *models.User) error {
	// Обновляем в базе данных
	if err := r.db.Save(user).Error; err != nil {
		return err
	}

	// Удаляем из кэша
	cacheKey := fmt.Sprintf("user:%s", user.ID.String())
	r.redis.Del(context.Background(), cacheKey)

	return nil
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
