package handlers

import (
	"net/http"

	"mobile-store-back/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func CreateProductVariant(productVariantService *services.ProductVariantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ProductID string  `json:"product_id" validate:"required"`
			SKU       string  `json:"sku" validate:"required"`
			Name      string  `json:"name" validate:"required,min=2"`
			Color     string  `json:"color"`
			Size      string  `json:"size"`
			Price     float64 `json:"price" validate:"required,min=0"`
			IsActive  bool    `json:"is_active"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		variant, err := productVariantService.Create(req.ProductID, req.SKU, req.Name, req.Color, req.Size, req.Price, req.IsActive)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, variant)
	}
}

func GetProductVariant(productVariantService *services.ProductVariantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		variant, err := productVariantService.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product variant not found"})
			return
		}

		c.JSON(http.StatusOK, variant)
	}
}

func GetProductVariants(productVariantService *services.ProductVariantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Query("product_id")
		if productID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product_id parameter is required"})
			return
		}

		variants, err := productVariantService.GetByProductID(productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"variants": variants})
	}
}

func GetProductVariantsByProductID(productVariantService *services.ProductVariantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		if productID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
			return
		}

		variants, err := productVariantService.GetByProductID(productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"variants": variants})
	}
}

func UpdateProductVariant(productVariantService *services.ProductVariantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			SKU      *string  `json:"sku" validate:"omitempty,min=2"`
			Name     *string  `json:"name" validate:"omitempty,min=2"`
			Color    *string  `json:"color"`
			Size     *string  `json:"size"`
			Price    *float64 `json:"price" validate:"omitempty,min=0"`
			IsActive *bool    `json:"is_active"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		variant, err := productVariantService.Update(id, req.SKU, req.Name, req.Color, req.Size, req.Price, req.IsActive)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, variant)
	}
}

func DeleteProductVariant(productVariantService *services.ProductVariantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := productVariantService.Delete(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product variant deleted successfully"})
	}
}

