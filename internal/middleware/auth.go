package middleware

import (
	"errors"
	"mobile-store-back/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":    "Authorization header required",
				"code":     "TOKEN_MISSING",
				"redirect": true,
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":    "Bearer token required",
				"code":      "TOKEN_MALFORMED",
				"redirect":  true,
			})
			c.Abort()
			return
		}

		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			// Проверяем тип ошибки токена
			var tokenErr *services.TokenError
			if errors.As(err, &tokenErr) {
				response := gin.H{
					"error": err.Error(),
					"code":  "TOKEN_" + strings.ToUpper(tokenErr.Type),
				}
				// Если токен истек, предлагаем фронту сделать refresh или редирект
				if tokenErr.Type == "expired" {
					response["redirect"] = false // Не редиректим сразу, фронт может попробовать refresh
					response["can_refresh"] = true
				} else {
					response["redirect"] = true
				}
				c.JSON(http.StatusUnauthorized, response)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":    "Invalid token",
					"code":     "TOKEN_INVALID",
					"redirect": true,
				})
			}
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func AdminRequired(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			// Проверяем тип ошибки токена
			var tokenErr *services.TokenError
			if errors.As(err, &tokenErr) {
				response := gin.H{
					"error": err.Error(),
					"code":  "TOKEN_" + strings.ToUpper(tokenErr.Type),
				}
				if tokenErr.Type == "expired" {
					response["redirect"] = false
					response["can_refresh"] = true
				} else {
					response["redirect"] = true
				}
				c.JSON(http.StatusUnauthorized, response)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":    "Invalid token",
					"code":     "TOKEN_INVALID",
					"redirect": true,
				})
			}
			c.Abort()
			return
		}

		// Проверяем, является ли пользователь администратором
		user, err := authService.GetUserByID(userID)
		if err != nil || user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

// OptionalAuth проверяет авторизацию, но не требует её (не прерывает запрос, если токен отсутствует)
// Используется для endpoints, которые могут работать как для авторизованных, так и для неавторизованных пользователей
func OptionalAuth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			// Если токен невалиден, просто продолжаем без user_id
			c.Next()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

// ManagerRequired проверяет, что пользователь имеет роль manager или admin
func ManagerRequired(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Проверяем, является ли пользователь менеджером или администратором
		user, err := authService.GetUserByID(userID)
		if err != nil || (user.Role != "manager" && user.Role != "admin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Manager or admin access required"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
