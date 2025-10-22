package handlers

import (
	"net/http"

	"mobile-store-back/internal/models"
	"mobile-store-back/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCategories(categoryService *services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := categoryService.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"categories": categories})
	}
}

func GetCategory(categoryService *services.CategoryService) gin.HandlerFunc {
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
		slug := c.Param("slug")

		category, err := categoryService.GetBySlugWithProducts(slug)
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


// Админские хендлеры
func CreateCategory(categoryService *services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string `json:"name" validate:"required,min=2"`
			Description string `json:"description"`
			Slug        string `json:"slug" validate:"required"`
			ImageURL    string `json:"image_url"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		category := &models.Category{
			ID:          uuid.New(),
			Name:        req.Name,
			Description: req.Description,
			Slug:        req.Slug,
			ImageURL:    req.ImageURL,
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
			Name        *string `json:"name"`
			Description *string `json:"description"`
			Slug        *string `json:"slug"`
			ImageURL    *string `json:"image_url"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		category, err := categoryService.Update(id, req.Name, req.Description, req.Slug, req.ImageURL)
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
