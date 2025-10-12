package handlers

import (
	"mobile-store-back/internal/config"
	"mobile-store-back/internal/middleware"
	"mobile-store-back/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, services *services.Services, cfg *config.Config) {
	// API v1 группа
	v1 := router.Group("/api/v1")
	{
		// Публичные маршруты
		public := v1.Group("/")
		{
			// Аутентификация
			auth := public.Group("/auth")
			{
				auth.POST("/register", Register(services.Auth))
				auth.POST("/login", Login(services.Auth))
			}

			// Продукты (публичные)
			products := public.Group("/products")
			{
				products.GET("/", GetProducts(services.Product))
				products.GET("/:id", GetProduct(services.Product))
				products.GET("/search", SearchProducts(services.Product))
				products.GET("/category/:category_id", GetProductsByCategory(services.Product))
			}
		}

		// Защищенные маршруты
		protected := v1.Group("/")
		protected.Use(middleware.AuthRequired(services.Auth))
		{
			// Пользователи
			users := protected.Group("/users")
			{
				users.GET("/profile", GetProfile(services.User))
				users.PUT("/profile", UpdateProfile(services.User))
			}

			// Заказы
			orders := protected.Group("/orders")
			{
				orders.POST("/", CreateOrder(services.Order))
				orders.GET("/", GetUserOrders(services.Order))
				orders.GET("/:id", GetOrder(services.Order))
				orders.PUT("/:id", UpdateOrder(services.Order))
			}
		}

		// Админские маршруты
		admin := v1.Group("/admin")
		admin.Use(middleware.AdminRequired(services.Auth))
		{
			// Управление пользователями
			adminUsers := admin.Group("/users")
			{
				adminUsers.GET("/", GetUsers(services.User))
				adminUsers.GET("/:id", GetUser(services.User))
				adminUsers.PUT("/:id", UpdateUser(services.User))
				adminUsers.DELETE("/:id", DeleteUser(services.User))
			}

			// Управление продуктами
			adminProducts := admin.Group("/products")
			{
				adminProducts.POST("/", CreateProduct(services.Product))
				adminProducts.PUT("/:id", UpdateProduct(services.Product))
				adminProducts.DELETE("/:id", DeleteProduct(services.Product))
			}

			// Управление заказами
			adminOrders := admin.Group("/orders")
			{
				adminOrders.GET("/", GetAllOrders(services.Order))
				adminOrders.PUT("/:id", UpdateOrderStatus(services.Order))
			}
		}
	}
}
