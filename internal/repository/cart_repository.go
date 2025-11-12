package repository

import (
	"errors"
	"mobile-store-back/internal/models"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// UniqueConstraintError - специальная ошибка для обработки unique constraint вне транзакции
type UniqueConstraintError struct {
	Err error
}

func (e *UniqueConstraintError) Error() string {
	return e.Err.Error()
}

type cartRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewCartRepository(db *gorm.DB, redis *redis.Client) CartRepository {
	return &cartRepository{
		db:    db,
		redis: redis,
	}
}

func (r *cartRepository) GetByUserID(userID string) ([]models.CartItem, error) {
	var items []models.CartItem
	err := r.db.Where("user_id = ?", userID).Preload("Product").Find(&items).Error
	return items, err
}

func (r *cartRepository) AddItem(userID string, productID string, quantity int) (*models.CartItem, error) {
	// Парсим userID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	
	// Проверяем, существует ли пользователь
	var user models.User
	err = r.db.First(&user, "id = ?", userUUID).Error
	if err != nil {
		return nil, err
	}
	
	// Получаем товар для получения цены
	var product models.Product
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}
	
	// Проверяем, что товар существует и активен
	err = r.db.Where("id = ? AND is_active = ?", productUUID, true).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found or inactive")
		}
		return nil, err
	}
	
	// Используем транзакцию для атомарности и предотвращения race condition
	var item models.CartItem
	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Проверяем, есть ли уже такой товар в корзине
		err := tx.Where("user_id = ? AND product_id = ?", userUUID, productUUID).First(&item).Error
		
		if err == nil {
			// Товар уже есть в корзине - обновляем количество (upsert логика)
			item.Quantity += quantity
			item.Price = product.BasePrice // Обновляем цену на актуальную
			return tx.Save(&item).Error
		}
		
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Другая ошибка при поиске
			return err
		}
		
		// Товара нет, создаем новый
		// Не устанавливаем SessionID, чтобы избежать конфликтов с уникальным индексом
		item = models.CartItem{
			UserID:    userUUID,
			ProductID: productUUID,
			Quantity:  quantity,
			Price:     product.BasePrice,
			// SessionID не устанавливаем - оставляем пустым (NULL в БД)
		}
		
		// Создаем запись
		err = tx.Create(&item).Error
		if err != nil {
			// Проверяем, не ошибка ли это unique constraint (race condition)
			if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
				// Транзакция помечена как aborted, нужно откатить и начать новую
				// Возвращаем специальную ошибку, чтобы обработать её вне транзакции
				return &UniqueConstraintError{Err: err}
			}
			return err
		}
		
		return nil
	})
	
	// Обрабатываем ошибку unique constraint вне транзакции
	if err != nil {
		var uniqueErr *UniqueConstraintError
		if errors.As(err, &uniqueErr) {
			// Товар был добавлен другим запросом - получаем его и обновляем количество в новой транзакции
			err = r.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Where("user_id = ? AND product_id = ?", userUUID, productUUID).First(&item).Error
				if err != nil {
					return err
				}
				item.Quantity += quantity
				item.Price = product.BasePrice
				return tx.Save(&item).Error
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	
	// Загружаем связанные данные
	err = r.db.Preload("Product").First(&item, item.ID).Error
	return &item, err
}

func (r *cartRepository) UpdateItem(id string, userID string, quantity int) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error
	if err != nil {
		return nil, err
	}
	
	item.Quantity = quantity
	err = r.db.Save(&item).Error
	if err != nil {
		return nil, err
	}
	
	// Загружаем связанные данные
	err = r.db.Preload("Product").First(&item, item.ID).Error
	return &item, err
}

func (r *cartRepository) RemoveItem(id string, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.CartItem{}).Error
}

func (r *cartRepository) Clear(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error
}

func (r *cartRepository) GetCount(userID string) (int, error) {
	var count int64
	err := r.db.Model(&models.CartItem{}).Where("user_id = ?", userID).Count(&count).Error
	return int(count), err
}

// MergeCart объединяет корзину сессии с корзиной пользователя после логина
// В БД user_id NOT NULL, поэтому записи с session_id и без user_id не могут существовать
// Этот метод используется для очистки старых записей и объединения данных
func (r *cartRepository) MergeCart(userID string, sessionID string) error {
	// Поскольку в БД user_id NOT NULL, мы не можем иметь записи только с session_id
	// Но на всякий случай проверяем и удаляем любые ошибочные записи с session_id без валидного user_id
	// Это не должно происходить, но если миграция не была выполнена, такие записи могут существовать
	
	// Удаляем записи с session_id, которые не имеют валидного user_id (если такие есть)
	// В нормальной работе этого не должно быть, так как user_id NOT NULL
	err := r.db.Exec("DELETE FROM cart_items WHERE session_id = ? AND (user_id IS NULL OR user_id NOT IN (SELECT id FROM users))", sessionID).Error
	if err != nil {
		return err
	}
	
	// Если frontend отправляет товары после логина, они должны быть добавлены через AddItem
	// Этот метод просто очищает старые/некорректные записи
	return nil
}
