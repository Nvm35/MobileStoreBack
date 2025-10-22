package handlers

import (
	"mobile-store-back/internal/middleware"
	"mobile-store-back/internal/services"
	"mobile-store-back/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetCart - получение корзины пользователя (авторизованного или по сессии)
func GetCart(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userOrSessionID, _ := middleware.GetUserOrSessionID(c)
		
		items, err := cartService.GetByUserID(userOrSessionID)
		utils.HandleInternalError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items})
	}
}

// AddToCart - добавление товара в корзину (для авторизованных и неавторизованных пользователей)
func AddToCart(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userOrSessionID, _ := middleware.GetUserOrSessionID(c)
		
		var req struct {
			ProductID uuid.UUID `json:"product_id" validate:"required"`
			Quantity  int       `json:"quantity" validate:"required,min=1"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		item, err := cartService.AddItem(userOrSessionID, req.ProductID.String(), req.Quantity)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

// UpdateCartItem - обновление количества товара в корзине
func UpdateCartItem(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userOrSessionID, _ := middleware.GetUserOrSessionID(c)
		
		var req struct {
			Quantity int `json:"quantity" validate:"required,min=1"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		item, err := cartService.UpdateItem(id, userOrSessionID, req.Quantity)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

// RemoveFromCart - удаление товара из корзины
func RemoveFromCart(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userOrSessionID, _ := middleware.GetUserOrSessionID(c)
		
		err := cartService.RemoveItem(id, userOrSessionID)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
	}
}

// ClearCart - очистка корзины
func ClearCart(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userOrSessionID, _ := middleware.GetUserOrSessionID(c)
		
		err := cartService.Clear(userOrSessionID)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cart cleared"})
	}
}

// GetCartCount - получение количества товаров в корзине
func GetCartCount(cartService *services.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userOrSessionID, _ := middleware.GetUserOrSessionID(c)
		
		count, err := cartService.GetCount(userOrSessionID)
		utils.HandleInternalError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"count": count})
	}
}
