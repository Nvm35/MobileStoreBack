package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Разрешенные домены для CORS
		allowedOrigins := GetCORSOrigins()
		
		// Проверяем, разрешен ли origin
		allowedOrigin := ""
		if origin != "" {
			// Проверяем точные совпадения
			for _, allowed := range allowedOrigins {
				if strings.EqualFold(origin, allowed) {
					allowedOrigin = origin
					break
				}
			}
			
			// Если точного совпадения нет, проверяем Vercel поддомены
			if allowedOrigin == "" && strings.Contains(origin, "vercel.app") {
				allowedOrigin = origin
			}
		} else {
			// Если нет Origin заголовка, разрешаем все (для тестирования)
			allowedOrigin = "*"
		}
		
		// ВАЖНО: Устанавливаем CORS заголовки ДО обработки запроса
		// Это гарантирует, что заголовки будут в ответе даже при ошибках авторизации
		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Header("Access-Control-Max-Age", "86400") // 24 часа

		// Обрабатываем preflight запросы (OPTIONS)
		// Preflight запросы должны обрабатываться ДО проверки авторизации
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func GetCORSOrigins() []string {
	// Базовые домены для локальной разработки
	origins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3001",
		"http://localhost:5173", // Vite dev server
		"http://127.0.0.1:5173",
	}
	
	// Добавляем домены из переменной окружения
	if corsOrigins := os.Getenv("CORS_ORIGINS"); corsOrigins != "" {
		envOrigins := strings.Split(corsOrigins, ",")
		for _, origin := range envOrigins {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				origins = append(origins, origin)
			}
		}
	}
	
	return origins
}
