package handlers

import (
	"mobile-store-back/internal/config"
	"mobile-store-back/internal/middleware"
	"mobile-store-back/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, services *services.Services, cfg *config.Config) {
	// API группа
	api := router.Group("/api")
	{
		// Публичные маршруты
		public := api.Group("/")
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
				products.GET("/:id/reviews", GetProductReviews(services.Review))
			}
		}

		// Защищенные маршруты
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired(services.Auth))
		{
			// Пользователи
			users := protected.Group("/users")
			{
				users.GET("/profile", GetProfile(services.User))
				users.PUT("/profile", UpdateProfile(services.User))
			}

			// Адреса теперь встроены в профиль пользователя

			// Заказы
			orders := protected.Group("/orders")
			{
				orders.POST("/", CreateOrder(services.Order))
				orders.GET("/", GetUserOrders(services.Order))
				orders.GET("/:id", GetOrder(services.Order))
				orders.PUT("/:id", UpdateOrder(services.Order))
			}

			// Корзина
			cart := protected.Group("/cart")
			{
				cart.GET("/", GetCart(services.Cart))
				cart.POST("/", AddToCart(services.Cart))
				cart.PUT("/:id", UpdateCartItem(services.Cart))
				cart.DELETE("/:id", RemoveFromCart(services.Cart))
				cart.DELETE("/", ClearCart(services.Cart))
				cart.GET("/count", GetCartCount(services.Cart))
			}

			// Избранное
			wishlist := protected.Group("/wishlist")
			{
				wishlist.GET("/", GetWishlist(services.Wishlist))
				wishlist.POST("/", AddToWishlist(services.Wishlist))
				wishlist.DELETE("/:id", RemoveFromWishlist(services.Wishlist))
				wishlist.DELETE("/", ClearWishlist(services.Wishlist))
				wishlist.GET("/check/:product_id", IsInWishlist(services.Wishlist))
			}

			// Отзывы
			reviews := protected.Group("/reviews")
			{
				reviews.POST("/", CreateReview(services.Review))
				reviews.GET("/my", GetUserReviews(services.Review))
				reviews.PUT("/:id", UpdateReview(services.Review))
				reviews.DELETE("/:id", DeleteReview(services.Review))
				reviews.POST("/:id/vote", VoteReview(services.Review))
			}

			// Промокоды
			coupons := protected.Group("/coupons")
			{
				coupons.GET("/", GetCoupons(services.Coupon))
				coupons.GET("/:id", GetCoupon(services.Coupon))
				coupons.POST("/validate", ValidateCoupon(services.Coupon))
			}
		}

		// Админские маршруты
		admin := api.Group("/admin")
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
				adminOrders.PUT("/:id/status", UpdateOrderStatus(services.Order))
			}

			// Управление отзывами
			adminReviews := admin.Group("/reviews")
			{
				adminReviews.GET("/", GetAllReviews(services.Review))
				adminReviews.PUT("/:id/approve", ApproveReview(services.Review))
			}

			// Управление промокодами
			adminCoupons := admin.Group("/coupons")
			{
				adminCoupons.POST("/", CreateCoupon(services.Coupon))
				adminCoupons.PUT("/:id", UpdateCoupon(services.Coupon))
				adminCoupons.DELETE("/:id", DeleteCoupon(services.Coupon))
				adminCoupons.GET("/:id/usage", GetCouponUsage(services.Coupon))
			}
		}
	}
}
