package handlers

import (
	"mobile-store-back/internal/config"
	"mobile-store-back/internal/middleware"
	"mobile-store-back/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, services *services.Services, cfg *config.Config) {
	api := router.Group("/api")
	{
		setupPublicRoutes(api, services)
		setupProtectedRoutes(api, services)
		setupAdminRoutes(api, services)
	}
}

// ============================================================================
// ПУБЛИЧНЫЕ МАРШРУТЫ (без аутентификации)
// ============================================================================
func setupPublicRoutes(api *gin.RouterGroup, services *services.Services) {
	public := api.Group("/")
	{
		// Аутентификация
		setupAuthRoutes(public, services)
		
		// Каталог товаров
		setupCatalogRoutes(public, services)
	}
}

func setupAuthRoutes(router *gin.RouterGroup, services *services.Services) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", Register(services.Auth))
		auth.POST("/login", Login(services.Auth))
	}
}

func setupCatalogRoutes(router *gin.RouterGroup, services *services.Services) {
	// Категории (публичные)
	categories := router.Group("/categories")
	{
		categories.GET("/", GetCategories(services.Category))
		categories.GET("/:slug", GetCategoryBySlug(services.Category))
		categories.GET("/:slug/products", GetCategoryProductsBySlug(services.Category))
	}

	// Продукты (публичные)
	products := router.Group("/products")
	{
		products.GET("/", GetProducts(services.Product))
		products.GET("/:id", GetProduct(services.Product))
		products.GET("/search", SearchProducts(services.Product))
		products.GET("/:id/reviews", GetProductReviews(services.Review))
	}
}

// ============================================================================
// ЗАЩИЩЕННЫЕ МАРШРУТЫ (требуют аутентификации)
// ============================================================================
func setupProtectedRoutes(api *gin.RouterGroup, services *services.Services) {
	protected := api.Group("/")
	protected.Use(middleware.AuthRequired(services.Auth))
	{
		// Профиль пользователя
		setupUserRoutes(protected, services)
		
		// Покупки
		setupShoppingRoutes(protected, services)
		
		// Отзывы и рейтинги
		setupReviewRoutes(protected, services)
		
	}
}

func setupUserRoutes(router *gin.RouterGroup, services *services.Services) {
	users := router.Group("/users")
	{
		users.GET("/profile", GetProfile(services.User))
		users.PUT("/profile", UpdateProfile(services.User))
	}
}

func setupShoppingRoutes(router *gin.RouterGroup, services *services.Services) {
	// Заказы
	orders := router.Group("/orders")
	{
		orders.POST("/", CreateOrder(services.Order))
		orders.GET("/", GetUserOrders(services.Order))
		orders.GET("/:id", GetOrder(services.Order))
		orders.PUT("/:id", UpdateOrder(services.Order))
	}

	// Корзина
	cart := router.Group("/cart")
	{
		cart.GET("/", GetCart(services.Cart))
		cart.POST("/", AddToCart(services.Cart))
		cart.PUT("/:id", UpdateCartItem(services.Cart))
		cart.DELETE("/:id", RemoveFromCart(services.Cart))
		cart.DELETE("/", ClearCart(services.Cart))
		cart.GET("/count", GetCartCount(services.Cart))
	}

	// Избранное
	wishlist := router.Group("/wishlist")
	{
		wishlist.GET("/", GetWishlist(services.Wishlist))
		wishlist.POST("/", AddToWishlist(services.Wishlist))
		wishlist.DELETE("/:id", RemoveFromWishlist(services.Wishlist))
		wishlist.DELETE("/", ClearWishlist(services.Wishlist))
		wishlist.GET("/check/:product_id", IsInWishlist(services.Wishlist))
	}
}

func setupReviewRoutes(router *gin.RouterGroup, services *services.Services) {
	reviews := router.Group("/reviews")
	{
		reviews.POST("/", CreateReview(services.Review))
		reviews.GET("/my", GetUserReviews(services.Review))
		reviews.PUT("/:id", UpdateReview(services.Review))
		reviews.DELETE("/:id", DeleteReview(services.Review))
		reviews.POST("/:id/vote", VoteReview(services.Review))
	}
}


// ============================================================================
// АДМИНСКИЕ МАРШРУТЫ (требуют админских прав)
// ============================================================================
func setupAdminRoutes(api *gin.RouterGroup, services *services.Services) {
	admin := api.Group("/admin")
	admin.Use(middleware.AdminRequired(services.Auth))
	{
		// Управление пользователями
		setupAdminUserRoutes(admin, services)
		
		// Управление каталогом
		setupAdminCatalogRoutes(admin, services)
		
		// Управление заказами
		setupAdminOrderRoutes(admin, services)
		
		// Управление контентом
		setupAdminContentRoutes(admin, services)
	}
}

func setupAdminUserRoutes(router *gin.RouterGroup, services *services.Services) {
	users := router.Group("/users")
	{
		users.GET("/", GetUsers(services.User))
		users.GET("/:id", GetUser(services.User))
		users.PUT("/:id", UpdateUser(services.User))
		users.DELETE("/:id", DeleteUser(services.User))
	}
}

func setupAdminCatalogRoutes(router *gin.RouterGroup, services *services.Services) {
	// Управление продуктами
	products := router.Group("/products")
	{
		products.POST("/", CreateProduct(services.Product))
		products.PUT("/:id", UpdateProduct(services.Product))
		products.DELETE("/:id", DeleteProduct(services.Product))
	}

	// Управление категориями
	categories := router.Group("/categories")
	{
		categories.POST("/", CreateCategory(services.Category))
		categories.GET("/:id", GetCategory(services.Category))
		categories.PUT("/:id", UpdateCategory(services.Category))
		categories.DELETE("/:id", DeleteCategory(services.Category))
	}
}

func setupAdminOrderRoutes(router *gin.RouterGroup, services *services.Services) {
	orders := router.Group("/orders")
	{
		orders.GET("/", GetAllOrders(services.Order))
		orders.PUT("/:id/status", UpdateOrderStatus(services.Order))
	}
}

func setupAdminContentRoutes(router *gin.RouterGroup, services *services.Services) {
	// Управление отзывами
	reviews := router.Group("/reviews")
	{
		reviews.GET("/", GetAllReviews(services.Review))
		reviews.PUT("/:id/approve", ApproveReview(services.Review))
	}

}
