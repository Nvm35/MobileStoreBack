package repository

import (
	"mobile-store-back/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB, redis *redis.Client) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *authRepository) CreateSession(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *authRepository) UpdateSession(session *models.Session) error {
	return r.db.Save(session).Error
}

func (r *authRepository) GetSessionByID(id string) (*models.Session, error) {
	var session models.Session
	if err := r.db.First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *authRepository) DeleteSessionByID(id string) error {
	return r.db.Delete(&models.Session{}, "id = ?", id).Error
}

func (r *authRepository) DeleteExpiredSessions() error {
	return r.db.Where("expires_at < ?", gorm.Expr("NOW()")).Delete(&models.Session{}).Error
}

func (r *authRepository) DeleteOldSessionsForUser(userID string, keepLast int) error {
	// Находим ID сессий, которые нужно оставить (последние N по created_at)
	var keepIDs []string
	if err := r.db.Model(&models.Session{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(keepLast).
		Pluck("id", &keepIDs).Error; err != nil {
		return err
	}

	// Если нет сессий для удаления
	if len(keepIDs) == 0 {
		return nil
	}

	// Удаляем все сессии пользователя, кроме последних N
	return r.db.Where("user_id = ? AND id NOT IN ?", userID, keepIDs).
		Delete(&models.Session{}).Error
}
