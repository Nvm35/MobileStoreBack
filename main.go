package main

import (
	"log"
	"mobile-store-back/internal/config"
	"mobile-store-back/internal/database"
	"mobile-store-back/internal/handlers"
	"mobile-store-back/internal/middleware"
	"mobile-store-back/internal/repository"
	"mobile-store-back/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title Mobile Store API
// @version 1.0
// @description API для магазина мобильных аксессуаров
// @host localhost:8080
// @BasePath /api
func main() {
	// Инициализация конфигурации
	cfg := config.Load()

	// Устанавливаем режим Gin в зависимости от окружения
	if cfg.Env == "production" || cfg.Env == "render" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	// Инициализация базы данных
	db, err := database.Connect(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Инициализация Redis
	redisClient := database.ConnectRedis(cfg)

	// Инициализация репозиториев
	repos := repository.New(db, redisClient)

	// Инициализация сервисов
	services := services.New(repos, cfg)

	// Инициализация роутера
	router := gin.Default()

	// Middleware - CORS должен быть первым
	router.Use(middleware.CORS())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.Recovery(logger))

	// Инициализация обработчиков
	handlers.SetupRoutes(router, services, cfg)

	// Запуск фоновой задачи очистки истекших сессий (каждые 24 часа)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		logger.Info("Session cleanup worker started")

		for range ticker.C {
			if err := repos.Auth.DeleteExpiredSessions(); err != nil {
				logger.Error("Failed to delete expired sessions", zap.Error(err))
			} else {
				logger.Info("Expired sessions cleaned up successfully")
			}
		}
	}()

	// Запуск сервера
	logger.Info("Starting server",
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port))

	if err := router.Run(cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
