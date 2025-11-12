package handlers

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetProductImages - получение изображений товара по slug
func GetProductImages(imageService *services.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		productSlug := c.Param("slug")
		if productSlug == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product slug is required"})
			return
		}

		images, err := imageService.GetByProductSlug(productSlug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"images": images})
	}
}

// UploadProductImage - загрузка изображения товара (админ)
func UploadProductImage(imageService *services.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		if productID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
			return
		}

		// Валидация UUID
		if _, err := uuid.Parse(productID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id format"})
			return
		}

		var req struct {
			CloudinaryPublicID string `json:"cloudinary_public_id" validate:"required"`
			URL                string `json:"url" validate:"required,url"`
			IsPrimary          bool   `json:"is_primary"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Создаем изображение
		productUUID, _ := uuid.Parse(productID)
		image := &models.Image{
			ProductID:          productUUID,
			CloudinaryPublicID: req.CloudinaryPublicID,
			URL:                req.URL,
			IsPrimary:          req.IsPrimary,
		}

		if err := imageService.Create(image); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Если это главное изображение, убираем primary с других
		if req.IsPrimary {
			if err := imageService.SetPrimary(image.ID.String()); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			image.IsPrimary = true
		}

		c.JSON(http.StatusCreated, image)
	}
}

// DeleteImage - удаление изображения (админ)
func DeleteImage(imageService *services.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageID := c.Param("id")
		if imageID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image id is required"})
			return
		}

		// Валидация UUID
		if _, err := uuid.Parse(imageID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image id format"})
			return
		}

		if err := imageService.Delete(imageID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// SetPrimaryImage - установка главного изображения (админ)
func SetPrimaryImage(imageService *services.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageID := c.Param("id")
		if imageID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image id is required"})
			return
		}

		// Валидация UUID
		if _, err := uuid.Parse(imageID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image id format"})
			return
		}

		if err := imageService.SetPrimary(imageID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Primary image updated successfully"})
	}
}

// UpdateImage - обновление изображения товара (админ)
func UpdateImage(imageService *services.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageID := c.Param("id")
		if imageID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image id is required"})
			return
		}

		if _, err := uuid.Parse(imageID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image id format"})
			return
		}

		var req struct {
			CloudinaryPublicID *string `json:"cloudinary_public_id"`
			URL                *string `json:"url"`
			IsPrimary          *bool   `json:"is_primary"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		image, err := imageService.Update(imageID, req.CloudinaryPublicID, req.URL, req.IsPrimary)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, image)
	}
}
