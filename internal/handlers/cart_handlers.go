package handlers

import (
	"mobile-store-back/internal/services"
	"mobile-store-back/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetCart - получение корзины пользователя (только для авторизованных)
func GetCart(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// user_id устанавливается в AuthRequired middleware
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		items, err := cartService.GetByUserID(userID.(string))
		utils.HandleInternalError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items})
	}
}

// AddToCart - добавление товара в корзину (только для авторизованных)
func AddToCart(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// user_id устанавливается в AuthRequired middleware
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req struct {
			Product  string `json:"product" validate:"required"`
			Quantity int    `json:"quantity" validate:"required,min=1"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		item, err := cartService.AddItem(userID.(string), req.Product, req.Quantity)
		if err != nil {
			// Обрабатываем разные типы ошибок
			errMsg := err.Error()
			if errMsg == "record not found" || strings.Contains(errMsg, "not found") ||
				strings.Contains(errMsg, "Product not found") || strings.Contains(errMsg, "product not found or inactive") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found or inactive"})
			} else if strings.Contains(errMsg, "invalid product identifier") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product identifier"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			}
			return
		}

		// Если товар был обновлен (уже существовал), возвращаем 200 OK
		// Если товар был создан, возвращаем 201 Created
		// Для простоты всегда возвращаем 200 OK, так как теперь используется upsert логика
		c.JSON(http.StatusOK, item)
	}
}

// UpdateCartItem - обновление количества товара в корзине (только для авторизованных)
func UpdateCartItem(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.Param("id")
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req struct {
			Quantity int `json:"quantity" validate:"required,min=1"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		item, err := cartService.UpdateItem(identifier, userID.(string), req.Quantity)
		if err != nil {
			if err.Error() == "record not found" || strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

// RemoveFromCart - удаление товара из корзины (только для авторизованных)
func RemoveFromCart(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.Param("id")
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		err := cartService.RemoveItem(identifier, userID.(string))
		if err != nil {
			if err.Error() == "record not found" || strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
	}
}

// ClearCart - очистка корзины (только для авторизованных)
func ClearCart(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		err := cartService.Clear(userID.(string))
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cart cleared"})
	}
}

// GetCartCount - получение количества товаров в корзине (только для авторизованных)
func GetCartCount(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		count, err := cartService.GetCount(userID.(string))
		utils.HandleInternalError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"count": count})
	}
}
