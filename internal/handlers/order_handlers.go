package handlers

import (
	"mobile-store-back/internal/services"
	"mobile-store-back/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateOrder(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		var req struct {
			Items []struct {
				ProductID         *uuid.UUID `json:"product_id"`
				ProductSlug       *string    `json:"product_slug"`
				ProductVariantID  *uuid.UUID `json:"product_variant_id"`
				ProductVariantSKU *string    `json:"product_variant_sku"`
				Quantity          int        `json:"quantity" validate:"required,min=1"`
			} `json:"items" validate:"required,min=1"`
			// Способ доставки
			ShippingMethod string `json:"shipping_method" validate:"required,oneof=delivery pickup"`
			// Адрес доставки (если нужен другой адрес, чем у пользователя)
			ShippingAddress string `json:"shipping_address"`
			// Пункт самовывоза (если выбран pickup)
			PickupPoint   string `json:"pickup_point"`
			PaymentMethod string `json:"payment_method" validate:"required,oneof=cash card transfer"`
			CustomerNotes string `json:"customer_notes"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		items := make([]services.OrderItemInput, len(req.Items))

		for i, item := range req.Items {
			productSlug := normalizePointer(item.ProductSlug)
			variantSKU := normalizePointer(item.ProductVariantSKU)

			if item.ProductID == nil && productSlug == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Each item must include product_id or product_slug"})
				return
			}

			items[i] = services.OrderItemInput{
				ProductID:         item.ProductID,
				ProductSlug:       productSlug,
				ProductVariantID:  item.ProductVariantID,
				ProductVariantSKU: variantSKU,
				Quantity:          item.Quantity,
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
		utils.HandleInternalError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}

func GetOrder(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.Param("identifier")

		order, err := orderService.GetByID(identifier)
		utils.HandleNotFound(c, err, "Order not found")
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func UpdateOrder(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.Param("identifier")
		userID, _ := c.Get("user_id")

		var req struct {
			Status         *string `json:"status" validate:"omitempty,oneof=pending confirmed processing shipped delivered cancelled returned"`
			PaymentStatus  *string `json:"payment_status" validate:"omitempty,oneof=pending paid failed refunded cancelled"`
			TrackingNumber *string `json:"tracking_number"`
			CustomerNotes  *string `json:"customer_notes"`
			// Способ доставки
			ShippingMethod *string `json:"shipping_method" validate:"omitempty,oneof=delivery pickup"`
			// Адрес доставки (если нужен другой адрес, чем у пользователя)
			ShippingAddress *string `json:"shipping_address"`
			// Пункт самовывоза (если выбран pickup)
			PickupPoint *string `json:"pickup_point"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		order, err := orderService.Update(identifier, userID.(string), req.Status, req.PaymentStatus, req.TrackingNumber, req.CustomerNotes, req.ShippingMethod, req.ShippingAddress, req.PickupPoint)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func GetAllOrders(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orders, err := orderService.List()
		utils.HandleInternalError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}

func UpdateOrderStatus(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.Param("identifier")

		var req struct {
			Status         string  `json:"status" validate:"required,oneof=pending confirmed processing shipped delivered cancelled returned"`
			TrackingNumber *string `json:"tracking_number"`
		}

		// Валидация (ValidateRequest сам делает ShouldBindJSON)
		if !utils.ValidateRequest(c, &req) {
			return
		}

		order, err := orderService.UpdateStatus(identifier, req.Status, req.TrackingNumber)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func normalizePointer(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
