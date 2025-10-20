package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler - централизованная обработка ошибок
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("Panic recovered", zap.String("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"code":  "INTERNAL_ERROR",
			})
		} else {
			logger.Error("Panic recovered", zap.Any("error", recovered))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error", 
				"code":  "INTERNAL_ERROR",
			})
		}
		c.Abort()
	})
}

// ErrorResponse - стандартный формат ошибки
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// SendError - отправка ошибки в стандартном формате
func SendError(c *gin.Context, status int, message string, code ...string) {
	errorCode := "UNKNOWN_ERROR"
	if len(code) > 0 {
		errorCode = code[0]
	}
	
	c.JSON(status, ErrorResponse{
		Error: message,
		Code:  errorCode,
	})
}
