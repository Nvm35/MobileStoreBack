package handlers

import (
	"mobile-store-back/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetWishlist - получение избранного пользователя
func GetWishlist(wishlistService *services.WishlistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		items, err := wishlistService.GetByUserID(userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items})
	}
}

// AddToWishlist - добавление товара в избранное
func AddToWishlist(wishlistService *services.WishlistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		
		var req struct {
			ProductID uuid.UUID `json:"product_id" validate:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item, err := wishlistService.AddItem(userID.(string), req.ProductID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

// RemoveFromWishlist - удаление товара из избранного
func RemoveFromWishlist(wishlistService *services.WishlistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := c.Get("user_id")
		
		err := wishlistService.RemoveItem(id, userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item removed from wishlist"})
	}
}

// IsInWishlist - проверка, есть ли товар в избранном
func IsInWishlist(wishlistService *services.WishlistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		productID := c.Param("product_id")
		
		isInWishlist, err := wishlistService.IsInWishlist(userID.(string), productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"is_in_wishlist": isInWishlist})
	}
}

// ClearWishlist - очистка избранного
func ClearWishlist(wishlistService *services.WishlistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		
		err := wishlistService.Clear(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wishlist cleared"})
	}
}
