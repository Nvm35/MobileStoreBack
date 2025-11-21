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
	err := r.db.Where("user_id = ?", userID).
		Preload("Product").
		Preload("ProductVariant").
		Find(&items).Error
	return items, err
}

func (r *cartRepository) AddItem(userID string, productIdentifier string, variantIdentifier *string, quantity int) (*models.CartItem, error) {
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

	product, err := findProductByIdentifier(r.db, productIdentifier, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found or inactive")
		}
		return nil, err
	}
	productUUID := product.ID

	// Определяем вариант товара и цену
	var variantUUID *uuid.UUID
	var itemPrice float64 = product.BasePrice

	if variantIdentifier != nil && *variantIdentifier != "" {
		variant, err := findProductVariantByIdentifier(r.db, *variantIdentifier, true)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("product variant not found or inactive")
			}
			return nil, err
		}
		// Проверяем, что вариант принадлежит этому товару
		if variant.ProductID != productUUID {
			return nil, errors.New("product variant does not belong to the specified product")
		}
		variantUUID = &variant.ID
		itemPrice = variant.Price // Используем цену варианта, если он указан
	}

	// Используем транзакцию для атомарности и предотвращения race condition
	var item models.CartItem
	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Проверяем, есть ли уже такой товар с таким вариантом в корзине
		query := tx.Where("user_id = ? AND product_id = ?", userUUID, productUUID)
		if variantUUID != nil {
			query = query.Where("product_variant_id = ?", *variantUUID)
		} else {
			query = query.Where("product_variant_id IS NULL")
		}
		err := query.First(&item).Error

		if err == nil {
			// Товар уже есть в корзине - обновляем количество (upsert логика)
			item.Quantity += quantity
			item.Price = itemPrice // Обновляем цену на актуальную
			return tx.Save(&item).Error
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Другая ошибка при поиске
			return err
		}

		// Товара нет, создаем новый
		// SessionID не используется, так как корзина работает только для авторизованных пользователей
		item = models.CartItem{
			UserID:          userUUID,
			ProductID:       productUUID,
			ProductVariantID: variantUUID,
			Quantity:        quantity,
			Price:           itemPrice,
			// SessionID остается nil (NULL в БД) - не используется в текущей логике
		}

		// Создаем запись
		err = tx.Create(&item).Error
		if err != nil {
			// Проверяем, не ошибка ли это unique constraint (race condition)
			if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") || strings.Contains(err.Error(), "unique constraint") {
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
			// Товар был добавлен другим запросом или constraint нарушен - получаем существующую запись
			err = r.db.Transaction(func(tx *gorm.DB) error {
				// Пробуем найти запись с учетом варианта
				query := tx.Where("user_id = ? AND product_id = ?", userUUID, productUUID)
				if variantUUID != nil {
					query = query.Where("product_variant_id = ?", *variantUUID)
				} else {
					query = query.Where("product_variant_id IS NULL")
				}
				err := query.First(&item).Error
				
				if err == nil {
					// Запись найдена - обновляем количество
					item.Quantity += quantity
					item.Price = itemPrice
					return tx.Save(&item).Error
				}
				
				// Если запись не найдена, возможно старый constraint без варианта
				// Пробуем найти любую запись с этим товаром
				if errors.Is(err, gorm.ErrRecordNotFound) {
					var existingItem models.CartItem
					err = tx.Where("user_id = ? AND product_id = ?", userUUID, productUUID).First(&existingItem).Error
					if err == nil {
						// Нашли запись без варианта - если добавляем с вариантом, создаем новую
						// Но если constraint старый, это не сработает
						// В этом случае обновляем существующую запись
						if variantUUID != nil && existingItem.ProductVariantID == nil {
							// Обновляем существующую запись, добавляя вариант
							existingItem.ProductVariantID = variantUUID
							existingItem.Quantity = quantity
							existingItem.Price = itemPrice
							return tx.Save(&existingItem).Error
						}
						// Иначе просто обновляем количество
						existingItem.Quantity += quantity
						existingItem.Price = itemPrice
						return tx.Save(&existingItem).Error
					}
					// Если и это не сработало, создаем новую запись
					item = models.CartItem{
						UserID:          userUUID,
						ProductID:       productUUID,
						ProductVariantID: variantUUID,
						Quantity:        quantity,
						Price:           itemPrice,
					}
					return tx.Create(&item).Error
				}
				
				return err
			})
			if err != nil {
				return nil, err
			}
		} else {
			// Если это не unique constraint ошибка, возвращаем её
			return nil, err
		}
	}

	// Загружаем связанные данные
	err = r.db.Preload("Product").Preload("ProductVariant").First(&item, item.ID).Error
	return &item, err
}

func (r *cartRepository) UpdateItem(identifier string, userID string, quantity int) (*models.CartItem, error) {
	var item models.CartItem

	if id, err := uuid.Parse(identifier); err == nil {
		if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err == nil {
			item.Quantity = quantity
			if err := r.db.Save(&item).Error; err != nil {
				return nil, err
			}
			err = r.db.Preload("Product").Preload("ProductVariant").First(&item, item.ID).Error
			return &item, err
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	productID, err := findProductIDByIdentifier(r.db, identifier)
	if err != nil {
		return nil, err
	}

	if err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error; err != nil {
		return nil, err
	}

	item.Quantity = quantity
	if err := r.db.Save(&item).Error; err != nil {
		return nil, err
	}

	err = r.db.Preload("Product").Preload("ProductVariant").First(&item, item.ID).Error
	return &item, err
}

func (r *cartRepository) RemoveItem(identifier string, userID string) error {
	if id, err := uuid.Parse(identifier); err == nil {
		if err := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.CartItem{}).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	productID, err := findProductIDByIdentifier(r.db, identifier)
	if err != nil {
		return err
	}

	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.CartItem{}).Error
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
