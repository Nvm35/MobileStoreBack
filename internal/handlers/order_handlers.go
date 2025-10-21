package handlers

import (
	"mobile-store-back/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func CreateOrder(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		
		var req struct {
			Items             []struct {
				ProductID        uuid.UUID  `json:"product_id" validate:"required"`
				ProductVariantID *uuid.UUID `json:"product_variant_id"`
				Quantity         int        `json:"quantity" validate:"required,min=1"`
			} `json:"items" validate:"required,min=1"`
			// Способ доставки
			ShippingMethod    string `json:"shipping_method" validate:"required,oneof=delivery pickup"`
			// Адрес доставки (если нужен другой адрес, чем у пользователя)
			ShippingAddress   string `json:"shipping_address"`
			// Пункт самовывоза (если выбран pickup)
			PickupPoint       string `json:"pickup_point"`
			PaymentMethod     string     `json:"payment_method" validate:"required,oneof=cash card transfer"`
			CustomerNotes     string     `json:"customer_notes"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидация
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Конвертируем структуры для service
		items := make([]struct {
			ProductID        uuid.UUID
			ProductVariantID *uuid.UUID
			Quantity         int
		}, len(req.Items))
		
		for i, item := range req.Items {
			items[i] = struct {
				ProductID        uuid.UUID
				ProductVariantID *uuid.UUID
				Quantity         int
			}{
				ProductID:        item.ProductID,
				ProductVariantID: item.ProductVariantID,
				Quantity:         item.Quantity,
			}
		}
		
		order, err := orderService.Create(userID.(string), items, req.ShippingMethod, req.ShippingAddress, req.PickupPoint, req.PaymentMethod, req.CustomerNotes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, order)
	}
}

func GetUserOrders(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		orders, err := orderService.GetByUserID(userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}

func GetOrder(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		order, err := orderService.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func UpdateOrder(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := c.Get("user_id")
		
		var req struct {
			Status            *string    `json:"status" validate:"omitempty,oneof=pending confirmed processing shipped delivered cancelled returned"`
			PaymentStatus     *string    `json:"payment_status" validate:"omitempty,oneof=pending paid failed refunded cancelled"`
			TrackingNumber    *string    `json:"tracking_number"`
			CustomerNotes     *string    `json:"customer_notes"`
			// Способ доставки
			ShippingMethod    *string `json:"shipping_method" validate:"omitempty,oneof=delivery pickup"`
			// Адрес доставки (если нужен другой адрес, чем у пользователя)
			ShippingAddress   *string `json:"shipping_address"`
			// Пункт самовывоза (если выбран pickup)
			PickupPoint       *string `json:"pickup_point"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидация
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order, err := orderService.Update(id, userID.(string), req.Status, req.PaymentStatus, req.TrackingNumber, req.CustomerNotes, req.ShippingMethod, req.ShippingAddress, req.PickupPoint)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func GetAllOrders(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orders, err := orderService.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}

func UpdateOrderStatus(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		var req struct {
			Status         string  `json:"status" validate:"required,oneof=pending confirmed processing shipped delivered cancelled returned"`
			TrackingNumber *string `json:"tracking_number"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидация
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order, err := orderService.UpdateStatus(id, req.Status, req.TrackingNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}
