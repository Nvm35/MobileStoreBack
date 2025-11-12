package handlers

import (
	"mobile-store-back/internal/services"
	"mobile-store-back/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req services.RegisterRequest
		if !utils.ValidateRequest(c, &req) {
			return
		}

		response, err := authService.Register(&req)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusCreated, response)
	}
}

func Login(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req services.LoginRequest
		if !utils.ValidateRequest(c, &req) {
			return
		}

		response, err := authService.Login(&req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// После успешного логина корзина синхронизируется на фронте
		// Фронт должен отправить товары из localStorage на бэк после логина
		// Используется JWT токен через заголовок Authorization: Bearer <token>

		c.JSON(http.StatusOK, response)
	}
}

// Refresh обновляет JWT токен
func Refresh(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
				"code":  "TOKEN_MISSING",
			})
			return
		}

		tokenString := authHeader
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		response, err := authService.RefreshToken(tokenString)
		if err != nil {
			// Если токен истек слишком давно, возвращаем специальный код
			if err.Error() == "token expired too long ago, please login again" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
					"code":  "TOKEN_EXPIRED_TOO_LONG",
					"redirect": true,
				})
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
				"code":  "TOKEN_INVALID",
			})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
