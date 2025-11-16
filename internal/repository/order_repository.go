package repository

import (
	"fmt"
	"mobile-store-back/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB, redis *redis.Client) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(userID string, items []struct {
	ProductID        string
	ProductVariantID *string
	Quantity         int
}, shippingMethod string, shippingAddress string, pickupPoint string, paymentMethod string, customerNotes string) (*models.Order, error) {
	var createdOrder *models.Order
	
	// Начинаем транзакцию
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Парсим userID
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return fmt.Errorf("invalid user_id: %w", err)
		}

		// Получаем главный склад
		var mainWarehouse models.Warehouse
		if err := tx.Where("is_main = ? AND is_active = ?", true, true).First(&mainWarehouse).Error; err != nil {
			return fmt.Errorf("main warehouse not found: %w", err)
		}

		// Генерируем номер заказа
		orderNumber := fmt.Sprintf("ORD-%d", time.Now().UnixNano())

		// Подготавливаем данные для заказа
		var totalAmount float64
		var orderItems []models.OrderItem

		// Обрабатываем каждый товар
		for _, item := range items {
			// Получаем товар
			var product models.Product
			productUUID, err := uuid.Parse(item.ProductID)
			if err != nil {
				return fmt.Errorf("invalid product_id: %w", err)
			}

			if err := tx.Where("id = ? AND is_active = ?", productUUID, true).First(&product).Error; err != nil {
				return fmt.Errorf("product not found or inactive: %w", err)
			}

			var variant *models.ProductVariant
			var variantUUID *uuid.UUID
			var price float64

			// Если указан вариант товара
			if item.ProductVariantID != nil {
				parsedVariantUUID, err := uuid.Parse(*item.ProductVariantID)
				if err != nil {
					return fmt.Errorf("invalid product_variant_id: %w", err)
				}
				variantUUID = &parsedVariantUUID

				var v models.ProductVariant
				if err := tx.Where("id = ? AND product_id = ? AND is_active = ?", variantUUID, productUUID, true).First(&v).Error; err != nil {
					return fmt.Errorf("product variant not found or inactive: %w", err)
				}
				variant = &v
				price = variant.Price
			} else {
				// Используем базовую цену товара
				price = product.BasePrice
				// Если нет варианта, проверяем наличие через варианты товара
				// Для упрощения, если нет варианта, считаем что товар доступен
				// В реальной системе может потребоваться другая логика
			}

			// Проверяем и резервируем товар на складе (если есть вариант)
			if variantUUID != nil {
				var warehouseStock models.WarehouseStock
				err := tx.Where("warehouse_id = ? AND product_variant_id = ?", mainWarehouse.ID, variantUUID).First(&warehouseStock).Error
				
				if err != nil {
					return fmt.Errorf("stock not found for variant %s on warehouse", variantUUID.String())
				}

				availableStock := warehouseStock.Stock - warehouseStock.ReservedStock
				if availableStock < item.Quantity {
					return fmt.Errorf("insufficient stock for variant %s: available %d, requested %d", variantUUID.String(), availableStock, item.Quantity)
				}

				// Резервируем товар
				if err := tx.Model(&warehouseStock).
					Where("warehouse_id = ? AND product_variant_id = ? AND (stock - reserved_stock) >= ?",
						mainWarehouse.ID, variantUUID, item.Quantity).
					Update("reserved_stock", gorm.Expr("reserved_stock + ?", item.Quantity)).Error; err != nil {
					return fmt.Errorf("failed to reserve stock: %w", err)
				}
			}

			// Рассчитываем сумму для этого товара
			itemTotal := price * float64(item.Quantity)
			totalAmount += itemTotal

			// Создаем OrderItem
			orderItem := models.OrderItem{
				ProductID:        productUUID,
				ProductVariantID: variantUUID,
				Quantity:         item.Quantity,
				Price:            price,
			}
			orderItems = append(orderItems, orderItem)
		}

		// Создаем заказ
		order := models.Order{
			UserID:          userUUID,
			WarehouseID:     &mainWarehouse.ID,
			OrderNumber:     orderNumber,
			Status:          models.OrderStatusPending,
			TotalAmount:     totalAmount,
			PaymentMethod:   paymentMethod,
			PaymentStatus:   models.PaymentStatusPending,
			ShippingMethod:  shippingMethod,
			ShippingAddress: shippingAddress,
			PickupPoint:     pickupPoint,
			CustomerNotes:   customerNotes,
		}

		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Создаем OrderItems
		for i := range orderItems {
			orderItems[i].OrderID = order.ID
			if err := tx.Create(&orderItems[i]).Error; err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
		}

		// Загружаем связанные данные для ответа
		if err := tx.Preload("User").
			Preload("Warehouse").
			Preload("OrderItems").
			Preload("OrderItems.Product").
			Preload("OrderItems.ProductVariant").
			First(&order, order.ID).Error; err != nil {
			return fmt.Errorf("failed to load order data: %w", err)
		}

		createdOrder = &order
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return createdOrder, nil
}

func (r *orderRepository) GetByID(id string) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("User").Preload("OrderItems").Preload("OrderItems.Product").Preload("OrderItems.ProductVariant").
		First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(userID string) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.Preload("OrderItems").Preload("OrderItems.Product").Preload("OrderItems.ProductVariant").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) Update(id string, userID string, status *string, paymentStatus *string, trackingNumber *string, customerNotes *string, shippingMethod *string, shippingAddress *string, pickupPoint *string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&order).Error
	if err != nil {
		return nil, err
	}
	
	if status != nil {
		order.Status = models.OrderStatus(*status)
	}
	if paymentStatus != nil {
		order.PaymentStatus = models.PaymentStatus(*paymentStatus)
	}
	if trackingNumber != nil {
		order.TrackingNumber = *trackingNumber
	}
	// AdminNotes удален из модели
	if customerNotes != nil {
		order.CustomerNotes = *customerNotes
	}
	if shippingMethod != nil {
		order.ShippingMethod = *shippingMethod
	}
	if shippingAddress != nil {
		order.ShippingAddress = *shippingAddress
	}
	if pickupPoint != nil {
		order.PickupPoint = *pickupPoint
	}
	
	err = r.db.Save(&order).Error
	return &order, err
}

func (r *orderRepository) UpdateStatus(id string, status string, trackingNumber *string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	
	order.Status = models.OrderStatus(status)
	if trackingNumber != nil {
		order.TrackingNumber = *trackingNumber
	}
	// AdminNotes удален из модели
	
	err = r.db.Save(&order).Error
	return &order, err
}

func (r *orderRepository) Delete(id string) error {
	return r.db.Delete(&models.Order{}, "id = ?", id).Error
}

func (r *orderRepository) List() ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.Preload("User").Preload("OrderItems").Preload("OrderItems.Product").Preload("OrderItems.ProductVariant").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
