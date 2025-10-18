package handlers

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/services"
	"net/http"
	"strconv"

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

// GetAvailabilityInfo - получение информации о доступности товара
func GetAvailabilityInfo(warehouseStockService *services.WarehouseStockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sku := c.Param("sku")

		availability, err := warehouseStockService.GetAvailabilityInfoBySKU(sku)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"availability": availability})
	}
}

// CheckAvailability - проверка доступности товара
func CheckAvailability(warehouseStockService *services.WarehouseStockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sku := c.Param("sku")
		quantity, _ := strconv.Atoi(c.DefaultQuery("quantity", "1"))

		available, err := warehouseStockService.CheckAvailabilityBySKU(sku, quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"available": available,
			"quantity":  quantity,
			"sku":       sku,
		})
	}
}

// CheckAvailabilityByWarehouse - проверка доступности товара на конкретном складе
func CheckAvailabilityByWarehouse(warehouseStockService *services.WarehouseStockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		warehouseSlug := c.Param("warehouse_slug")
		sku := c.Param("sku")
		quantity, _ := strconv.Atoi(c.DefaultQuery("quantity", "1"))

		available, err := warehouseStockService.CheckAvailabilityByWarehouseSlug(warehouseSlug, sku, quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"available":       available,
			"quantity":        quantity,
			"warehouse_slug":  warehouseSlug,
			"sku":            sku,
		})
	}
}

// GetTotalAvailableStock - получение общего доступного остатка товара
func GetTotalAvailableStock(warehouseStockService *services.WarehouseStockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sku := c.Param("sku")

		totalStock, err := warehouseStockService.GetAvailableStockBySKU(sku)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"total_available_stock": totalStock,
			"sku":                  sku,
		})
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
