package handlers

import (
	"mobile-store-back/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		// TODO: Implement product creation
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
	}
}

func UpdateProduct(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement product update
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
	}
}

func DeleteProduct(productService *services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement product deletion
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
	}
}
