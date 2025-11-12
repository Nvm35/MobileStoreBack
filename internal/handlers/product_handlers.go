package handlers

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/services"
	"mobile-store-back/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetProducts(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Поддержка фильтрации и поиска в одном эндпоинте
		query := c.Query("q")
		categoryID := c.Query("category_id")
		brand := c.Query("brand")
		minPrice := c.Query("min_price")
		maxPrice := c.Query("max_price")

		var products []*models.Product
		var err error

		if query != "" {
			// Поиск товаров
			products, err = productService.Search(query)
		} else if categoryID != "" {
			// Фильтр по категории
			products, err = productService.GetByCategory(categoryID)
		} else {
			// Получить все товары с фильтрами
			products, err = productService.ListWithFilters(brand, minPrice, maxPrice)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	}
}

func GetProduct(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.Param("slug") // Может быть как ID, так и slug

		// Пробуем найти по slug или ID
		product, err := productService.GetBySlugOrID(identifier)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func CreateProduct(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string    `json:"name" validate:"required,min=2"`
			Description string    `json:"description"`
			BasePrice   float64   `json:"base_price" validate:"required,min=0"`
			SKU         string    `json:"sku" validate:"required"`
			IsActive    bool      `json:"is_active"`
			Feature     bool      `json:"feature"`
			Brand       string    `json:"brand" validate:"required,min=2"`
			Model       string    `json:"model"`
			Material    string    `json:"material"`
			CategoryID  uuid.UUID `json:"category_id" validate:"required"`
			Tags        []string  `json:"tags"`
			VideoURL    *string   `json:"video_url" validate:"omitempty,url"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		product, err := productService.Create(req.Name, req.Description, req.BasePrice, req.SKU, req.IsActive, req.Feature, req.Brand, req.Model, req.Material, req.CategoryID, req.Tags, req.VideoURL)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusCreated, product)
	}
}

func UpdateProduct(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			Name        *string    `json:"name" validate:"omitempty,min=2"`
			Description *string    `json:"description"`
			BasePrice   *float64   `json:"base_price" validate:"omitempty,min=0"`
			IsActive    *bool      `json:"is_active"`
			Feature     *bool      `json:"feature"`
			Brand       *string    `json:"brand" validate:"omitempty,min=2"`
			Model       *string    `json:"model"`
			Material    *string    `json:"material"`
			CategoryID  *uuid.UUID `json:"category_id"`
			Tags        *[]string  `json:"tags"`
			VideoURL    *string    `json:"video_url" validate:"omitempty,url"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		product, err := productService.Update(id, req.Name, req.Description, req.BasePrice, req.IsActive, req.Feature, req.Brand, req.Model, req.Material, req.CategoryID, req.Tags, req.VideoURL)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteProduct(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := productService.Delete(id)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}

func GetFeaturedProducts(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := productService.GetFeatured()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	}
}
