package handlers

import (
	"net/http"

	"mobile-store-back/internal/services"
	"mobile-store-back/internal/utils"

	"github.com/gin-gonic/gin"
)


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

		if !utils.ValidateRequest(c, &req) {
			return
		}

		variant, err := productVariantService.Create(req.ProductID, req.SKU, req.Name, req.Color, req.Size, req.Price, req.IsActive)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusCreated, variant)
	}
}

func GetProductVariant(productVariantService *services.ProductVariantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		variant, err := productVariantService.GetByID(id)
		utils.HandleNotFound(c, err, "Product variant not found")
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, variant)
	}
}


func GetProductVariantsByProductID(productVariantService *services.ProductVariantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.Param("slug") // Может быть как slug, так и ID
		if identifier == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product identifier is required"})
			return
		}

		variants, err := productVariantService.GetByProductSlugOrID(identifier)
		utils.HandleInternalError(c, err)
		if err != nil {
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

		if !utils.ValidateRequest(c, &req) {
			return
		}

		variant, err := productVariantService.Update(id, req.SKU, req.Name, req.Color, req.Size, req.Price, req.IsActive)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, variant)
	}
}

func DeleteProductVariant(productVariantService *services.ProductVariantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := productVariantService.Delete(id)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product variant deleted successfully"})
	}
}

