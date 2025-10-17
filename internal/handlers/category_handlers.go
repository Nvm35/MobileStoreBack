package handlers

import (
	"net/http"
	"strconv"

	"mobile-store-back/internal/models"
	"mobile-store-back/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCategories(categoryService *services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		categories, err := categoryService.GetAll(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"categories": categories})
	}
}

func GetCategory(categoryService *services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		category, err := categoryService.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"category": category})
	}
}

func GetCategoryBySlug(categoryService *services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		
		category, err := categoryService.GetBySlug(slug)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"category": category})
	}
}

func GetCategoryProducts(categoryService *services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		category, err := categoryService.GetWithProducts(id, limit, offset)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"category": category,
			"products": category.Products,
		})
	}
}

func GetCategoryProductsBySlug(categoryService *services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		// Сначала получаем категорию по slug
		category, err := categoryService.GetBySlug(slug)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		// Затем получаем продукты этой категории
		categoryWithProducts, err := categoryService.GetWithProducts(category.ID.String(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"category": category,
			"products": categoryWithProducts.Products,
		})
	}
}

// Админские хендлеры
func CreateCategory(categoryService *services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name            string `json:"name" validate:"required,min=2"`
			Description     string `json:"description"`
			Slug            string `json:"slug" validate:"required"`
			IsActive        bool   `json:"is_active"`
			SortOrder       int    `json:"sort_order"`
			ImageURL        string `json:"image_url"`
			MetaTitle       string `json:"meta_title"`
			MetaDescription string `json:"meta_description"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		category := &models.Category{
			ID:              uuid.New(),
			Name:            req.Name,
			Description:     req.Description,
			Slug:            req.Slug,
			IsActive:        req.IsActive,
			SortOrder:       req.SortOrder,
			ImageURL:        req.ImageURL,
			MetaTitle:       req.MetaTitle,
			MetaDescription: req.MetaDescription,
		}

		err := categoryService.Create(category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"category": category})
	}
}

func UpdateCategory(categoryService *services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		var req struct {
			Name            *string `json:"name"`
			Description     *string `json:"description"`
			Slug            *string `json:"slug"`
			IsActive        *bool   `json:"is_active"`
			SortOrder       *int    `json:"sort_order"`
			ImageURL        *string `json:"image_url"`
			MetaTitle       *string `json:"meta_title"`
			MetaDescription *string `json:"meta_description"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		category, err := categoryService.Update(id, req.Name, req.Description, req.Slug, req.IsActive, req.SortOrder, req.ImageURL, req.MetaTitle, req.MetaDescription)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"category": category})
	}
}

func DeleteCategory(categoryService *services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := categoryService.Delete(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
	}
}
