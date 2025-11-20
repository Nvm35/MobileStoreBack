package handlers

import (
	"mobile-store-back/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCloudinaryImages - получение списка изображений из Cloudinary (админ)
func GetCloudinaryImages(cloudinaryService *services.CloudinaryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		folder := c.Query("folder")
		maxResults := 50 // по умолчанию
		if maxStr := c.Query("max_results"); maxStr != "" {
			if parsed, err := strconv.Atoi(maxStr); err == nil && parsed > 0 && parsed <= 500 {
				maxResults = parsed
			}
		}
		nextCursor := c.Query("next_cursor")

		images, err := cloudinaryService.GetImages(folder, maxResults, nextCursor)
		if err != nil {
			// Возвращаем более детальную ошибку для отладки
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
				"message": "Failed to get images from Cloudinary. Please check your Cloudinary credentials in .env file.",
			})
			return
		}

		c.JSON(http.StatusOK, images)
	}
}

// GetCloudinaryFolders - получение списка папок из Cloudinary (админ)
func GetCloudinaryFolders(cloudinaryService *services.CloudinaryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		folders, err := cloudinaryService.GetFolders()
		if err != nil {
			// Возвращаем более детальную ошибку для отладки
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
				"message": "Failed to get folders from Cloudinary. Please check your Cloudinary credentials in .env file.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"folders": folders})
	}
}

// UploadToCloudinary - загрузка изображения в Cloudinary (админ)
func UploadToCloudinary(cloudinaryService *services.CloudinaryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		folder := c.PostForm("folder")

		// Получаем файл из формы
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required: " + err.Error()})
			return
		}

		// Открываем файл
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open file: " + err.Error()})
			return
		}
		defer src.Close()

		// Загружаем в Cloudinary
		uploadedImage, err := cloudinaryService.UploadImage(src, file.Filename, folder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload to cloudinary: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, uploadedImage)
	}
}

// DeleteCloudinaryImage - удаление изображения из Cloudinary (админ)
func DeleteCloudinaryImage(cloudinaryService *services.CloudinaryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		publicID := c.Param("public_id")
		if publicID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "public_id is required"})
			return
		}

		err := cloudinaryService.DeleteImage(publicID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
	}
}

// CheckCloudinaryConfig - проверка конфигурации Cloudinary (админ, для отладки)
func CheckCloudinaryConfig(cloudinaryService *services.CloudinaryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем конфигурацию без вызова API
		cfg := cloudinaryService.GetConfig()
		configStatus := map[string]interface{}{
			"cloud_name_set": cfg != nil && cfg.CloudName != "",
			"api_key_set":    cfg != nil && cfg.APIKey != "",
			"api_secret_set": cfg != nil && cfg.APISecret != "",
		}

		// Показываем частично скрытые значения для отладки (первые 4 символа)
		debugInfo := map[string]interface{}{}
		if cfg != nil {
			if cfg.CloudName != "" {
				debugInfo["cloud_name"] = cfg.CloudName[:min(4, len(cfg.CloudName))] + "..."
			}
			if cfg.APIKey != "" {
				debugInfo["api_key"] = cfg.APIKey[:min(4, len(cfg.APIKey))] + "..."
			}
			if cfg.APISecret != "" {
				debugInfo["api_secret"] = "***" + cfg.APISecret[max(0, len(cfg.APISecret)-4):]
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"config":     configStatus,
			"debug_info": debugInfo,
			"message":     "Check your .env file for CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY, CLOUDINARY_API_SECRET",
			"note":        "If all values are set but you still get 401, check: 1) Correct values from Cloudinary Console, 2) No extra spaces/quotes in .env, 3) Server restarted after .env changes",
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

