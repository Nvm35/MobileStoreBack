package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Глобальный валидатор для переиспользования
var validate = validator.New()

// ValidateRequest валидирует структуру запроса и возвращает ошибку в JSON если валидация не прошла
func ValidateRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	return true
}

// HandleError обрабатывает ошибки сервисов и возвращает соответствующий HTTP статус
func HandleError(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

// HandleInternalError обрабатывает внутренние ошибки сервера
func HandleInternalError(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// HandleNotFound обрабатывает ошибки "не найдено"
func HandleNotFound(c *gin.Context, err error, message string) {
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": message})
		return
	}
}
