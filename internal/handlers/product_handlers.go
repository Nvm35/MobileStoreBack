package handlers

import (
	"mobile-store-back/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func GetProducts(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		products, err := productService.List(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	}
}

func GetProduct(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		product, err := productService.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func SearchProducts(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
			return
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		products, err := productService.Search(query, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	}
}

func GetProductsByCategory(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("category_id")
		
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		products, err := productService.GetByCategory(categoryID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	}
}

func CreateProduct(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name             string    `json:"name" validate:"required,min=2"`
			Description      string    `json:"description"`
			ShortDescription string    `json:"short_description"`
			Price            float64   `json:"price" validate:"required,min=0"`
			ComparePrice     *float64  `json:"compare_price" validate:"omitempty,min=0"`
			SKU              string    `json:"sku" validate:"required"`
			Stock            int       `json:"stock" validate:"min=0"`
			IsActive         bool      `json:"is_active"`
			IsFeatured       bool      `json:"is_featured"`
			IsNew            bool      `json:"is_new"`
			Weight           *float64  `json:"weight" validate:"omitempty,min=0"`
			Dimensions       string    `json:"dimensions"`
			Brand            string    `json:"brand" validate:"required,min=2"`
			Model            string    `json:"model"`
			Color            string    `json:"color"`
			Material         string    `json:"material"`
			CategoryID       uuid.UUID `json:"category_id" validate:"required"`
			Tags             []string  `json:"tags"`
			MetaTitle        string    `json:"meta_title"`
			MetaDescription  string    `json:"meta_description"`
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

		product, err := productService.Create(req.Name, req.Description, req.ShortDescription, req.Price, req.ComparePrice, req.SKU, req.Stock, req.IsActive, req.IsFeatured, req.IsNew, req.Weight, req.Dimensions, req.Brand, req.Model, req.Color, req.Material, req.CategoryID, req.Tags, req.MetaTitle, req.MetaDescription)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	}
}

func UpdateProduct(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		var req struct {
			Name             *string    `json:"name" validate:"omitempty,min=2"`
			Description      *string    `json:"description"`
			ShortDescription *string    `json:"short_description"`
			Price            *float64   `json:"price" validate:"omitempty,min=0"`
			ComparePrice     *float64   `json:"compare_price" validate:"omitempty,min=0"`
			Stock            *int       `json:"stock" validate:"omitempty,min=0"`
			IsActive         *bool      `json:"is_active"`
			IsFeatured       *bool      `json:"is_featured"`
			IsNew            *bool      `json:"is_new"`
			Weight           *float64   `json:"weight" validate:"omitempty,min=0"`
			Dimensions       *string    `json:"dimensions"`
			Brand            *string    `json:"brand" validate:"omitempty,min=2"`
			Model            *string    `json:"model"`
			Color            *string    `json:"color"`
			Material         *string    `json:"material"`
			CategoryID       *uuid.UUID `json:"category_id"`
			Tags             []string   `json:"tags"`
			MetaTitle        *string    `json:"meta_title"`
			MetaDescription  *string    `json:"meta_description"`
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

		product, err := productService.Update(id, req.Name, req.Description, req.ShortDescription, req.Price, req.ComparePrice, req.Stock, req.IsActive, req.IsFeatured, req.IsNew, req.Weight, req.Dimensions, req.Brand, req.Model, req.Color, req.Material, req.CategoryID, req.Tags, req.MetaTitle, req.MetaDescription)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteProduct(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		err := productService.Delete(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}
