package handlers

import (
	"mobile-store-back/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateOrder(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement order creation
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
	}
}

func GetUserOrders(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		orders, err := orderService.GetByUserID(userID.(string), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}

func GetOrder(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		order, err := orderService.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func UpdateOrder(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement order update
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
	}
}

func GetAllOrders(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		orders, err := orderService.List(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}

func UpdateOrderStatus(orderService *services.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement order status update
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
	}
}
