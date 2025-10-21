package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck - проверка состояния сервиса
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"service":   "mobile-store-back",
		})
	}
}

// CORSTest - тест CORS
func CORSTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "CORS test successful",
			"origin":  c.Request.Header.Get("Origin"),
			"method":  c.Request.Method,
		})
	}
}

