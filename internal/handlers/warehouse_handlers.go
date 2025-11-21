package handlers

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetWarehouses - получение списка складов
func GetWarehouses(warehouseService *services.WarehouseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		warehouses, err := warehouseService.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"warehouses": warehouses})
	}
}

// GetWarehouse - получение склада по slug или ID
func GetWarehouse(warehouseService *services.WarehouseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пробуем получить параметр как "id" (для админских роутов) или "slug" (для публичных)
		identifier := c.Param("id")
		if identifier == "" {
			identifier = c.Param("slug")
		}

		warehouse, err := warehouseService.GetBySlugOrID(identifier)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
			return
		}

		c.JSON(http.StatusOK, warehouse)
	}
}

// GetWarehousesByCity - получение складов по городу
func GetWarehousesByCity(warehouseService *services.WarehouseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		city := c.Param("city")

		warehouses, err := warehouseService.GetByCity(city)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"warehouses": warehouses})
	}
}

// GetMainWarehouse - получение главного склада
func GetMainWarehouse(warehouseService *services.WarehouseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		warehouse, err := warehouseService.GetMain()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Main warehouse not found"})
			return
		}

		c.JSON(http.StatusOK, warehouse)
	}
}

// CreateWarehouse - создание склада (админ)
func CreateWarehouse(warehouseService *services.WarehouseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string `json:"name" validate:"required,min=2"`
			Address     string `json:"address" validate:"required"`
			City        string `json:"city" validate:"required"`
			Phone       string `json:"phone"`
			Email       string `json:"email"`
			IsActive    bool   `json:"is_active"`
			IsMain      bool   `json:"is_main"`
			ManagerName string `json:"manager_name"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		warehouse := &models.Warehouse{
			Name:        req.Name,
			Address:     req.Address,
			City:        req.City,
			Phone:       req.Phone,
			Email:       req.Email,
			IsActive:    req.IsActive,
			IsMain:      req.IsMain,
			ManagerName: req.ManagerName,
		}

		if err := warehouseService.Create(warehouse); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, warehouse)
	}
}

// UpdateWarehouse - обновление склада (админ)
func UpdateWarehouse(warehouseService *services.WarehouseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			Name        *string `json:"name"`
			Address     *string `json:"address"`
			City        *string `json:"city"`
			Phone       *string `json:"phone"`
			Email       *string `json:"email"`
			IsActive    *bool   `json:"is_active"`
			ManagerName *string `json:"manager_name"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		warehouse, err := warehouseService.Update(id, req.Name, req.Address, req.City, req.Phone, req.Email, req.IsActive, req.ManagerName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, warehouse)
	}
}

// DeleteWarehouse - удаление склада (админ)
func DeleteWarehouse(warehouseService *services.WarehouseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := warehouseService.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Warehouse deleted successfully"})
	}
}
