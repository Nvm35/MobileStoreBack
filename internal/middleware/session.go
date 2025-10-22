package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	SessionCookieName = "session_id"
	SessionDuration   = 30 * 24 * time.Hour // 30 дней
)

// SessionMiddleware создает или получает сессию для пользователя
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем session_id из cookie
		sessionID, err := c.Cookie(SessionCookieName)
		if err != nil || sessionID == "" {
			// Создаем новую сессию
			sessionID = generateSessionID()
			c.SetCookie(SessionCookieName, sessionID, int(SessionDuration.Seconds()), "/", "", false, true)
		}

		// Сохраняем session_id в контексте
		c.Set("session_id", sessionID)
		c.Next()
	}
}

// generateSessionID генерирует уникальный ID сессии
func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetSessionID получает session_id из контекста
func GetSessionID(c *gin.Context) string {
	if sessionID, exists := c.Get("session_id"); exists {
		return sessionID.(string)
	}
	return ""
}

// GetUserOrSessionID возвращает user_id если пользователь авторизован, иначе session_id
func GetUserOrSessionID(c *gin.Context) (string, bool) {
	// Сначала проверяем, авторизован ли пользователь
	if userID, exists := c.Get("user_id"); exists {
		return userID.(string), true // true = авторизованный пользователь
	}
	
	// Если не авторизован, используем session_id
	if sessionID, exists := c.Get("session_id"); exists {
		return sessionID.(string), false // false = неавторизованный пользователь
	}
	
	return "", false
}
