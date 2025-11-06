package handlers

import (
	"mobile-store-back/internal/services"
	"mobile-store-back/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			ProductID uuid.UUID `json:"product_id" validate:"required"`
			Quantity  int       `json:"quantity" validate:"required,min=1"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		item, err := cartService.AddItem(userID.(string), req.ProductID.String(), req.Quantity)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

// UpdateCartItem - обновление количества товара в корзине (только для авторизованных)
func UpdateCartItem(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
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

		item, err := cartService.UpdateItem(id, userID.(string), req.Quantity)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

// RemoveFromCart - удаление товара из корзины (только для авторизованных)
func RemoveFromCart(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		
		err := cartService.RemoveItem(id, userID.(string))
		utils.HandleError(c, err)
		if err != nil {
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
