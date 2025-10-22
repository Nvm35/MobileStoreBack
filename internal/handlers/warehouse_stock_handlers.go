package handlers

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetWarehouseStocks - получение остатков по складу
func GetWarehouseStocks(warehouseStockService *services.WarehouseStockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		warehouseSlug := c.Param("warehouse_slug")

		// Нужно получить ID склада по slug
		// Это будет реализовано в сервисе
		stocks, err := warehouseStockService.GetByWarehouseSlug(warehouseSlug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"stocks": stocks})
	}
}

// GetVariantStocks - получение остатков по варианту товара
func GetVariantStocks(warehouseStockService *services.WarehouseStockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sku := c.Param("sku")

		stocks, err := warehouseStockService.GetByVariantSKU(sku)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"stocks": stocks})
	}
}


// CreateWarehouseStock - создание остатка на складе (админ)
func CreateWarehouseStock(warehouseStockService *services.WarehouseStockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			WarehouseID      string `json:"warehouse_id" validate:"required"`
			ProductVariantID string `json:"product_variant_id" validate:"required"`
			Stock            int    `json:"stock" validate:"min=0"`
			ReservedStock    int    `json:"reserved_stock" validate:"min=0"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		warehouseID, _ := uuid.Parse(req.WarehouseID)
		variantID, _ := uuid.Parse(req.ProductVariantID)

		warehouseStock := &models.WarehouseStock{
			WarehouseID:      warehouseID,
			ProductVariantID: variantID,
			Stock:            req.Stock,
			ReservedStock:    req.ReservedStock,
		}

		if err := warehouseStockService.Create(warehouseStock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, warehouseStock)
	}
}

// UpdateWarehouseStock - обновление остатка на складе (админ)
func UpdateWarehouseStock(warehouseStockService *services.WarehouseStockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			Stock         int `json:"stock" validate:"min=0"`
			ReservedStock int `json:"reserved_stock" validate:"min=0"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		warehouseStock, err := warehouseStockService.UpdateStock(id, req.Stock, req.ReservedStock)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, warehouseStock)
	}
}

// DeleteWarehouseStock - удаление остатка на складе (админ)
func DeleteWarehouseStock(warehouseStockService *services.WarehouseStockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := warehouseStockService.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Warehouse stock deleted successfully"})
	}
}
