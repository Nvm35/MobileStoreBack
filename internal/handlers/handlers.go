package handlers

import (
	"mobile-store-back/internal/config"
	"mobile-store-back/internal/middleware"
	"mobile-store-back/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, services *services.Services, cfg *config.Config) {
	// Health check (без /api префикса)
	router.GET("/health", HealthCheck())
	
	// CORS test endpoint
	router.GET("/cors-test", CORSTest())
	
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
		
		// Корзина (доступна для неавторизованных пользователей через сессии)
		setupPublicCartRoutes(public, services)
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
		categories.GET("/:slug", GetCategory(services.Category)) // поддерживает и slug, и ID
		categories.GET("/:slug/products", GetCategoryProducts(services.Category))
	}

	// Продукты (публичные) - основной эндпоинт с поиском и фильтрацией
	products := router.Group("/products")
	{
		products.GET("/", GetProducts(services.Product)) // поддерживает поиск и фильтрацию через query параметры
		products.GET("/featured", GetFeaturedProducts(services.Product)) // товары с feature=true
		products.GET("/:slug", GetProduct(services.Product)) // поддерживает и slug, и ID
		products.GET("/:slug/reviews", GetProductReviews(services.Review))
		products.GET("/:slug/variants", GetProductVariantsByProductID(services.ProductVariant))
	}

	// Склады (публичные)
	warehouses := router.Group("/warehouses")
	{
		warehouses.GET("/", GetWarehouses(services.Warehouse))
		warehouses.GET("/main", GetMainWarehouse(services.Warehouse))
		warehouses.GET("/:slug", GetWarehouse(services.Warehouse)) // поддерживает и slug, и ID
		warehouses.GET("/city/:city", GetWarehousesByCity(services.Warehouse))
	}

	// Остатки товаров (публичные) - упрощенные URL
	stocks := router.Group("/stocks")
	{
		// Остатки по складу
		stocks.GET("/warehouse/:warehouse_slug", GetWarehouseStocks(services.WarehouseStock))
		
		// Остатки по варианту товара (используем SKU вместо ID)
		stocks.GET("/variant/:sku", GetVariantStocks(services.WarehouseStock))
	}

	// Изображения товаров (публичные)
	images := router.Group("/images")
	{
		images.GET("/product/:slug", GetProductImages(services.Image))
	}
}

func setupPublicCartRoutes(router *gin.RouterGroup, services *services.Services) {
	// Корзина ТОЛЬКО для авторизованных пользователей (JWT токен обязателен)
	// Неавторизованные пользователи работают с корзиной только в localStorage на фронте
	cart := router.Group("/cart")
	cart.Use(middleware.AuthRequired(services.Auth)) // Требуем авторизацию
	{
		cart.GET("/", GetCart(services.Cart))
		cart.POST("/", AddToCart(services.Cart))
		cart.PUT("/:id", UpdateCartItem(services.Cart))
		cart.DELETE("/:id", RemoveFromCart(services.Cart))
		cart.DELETE("/", ClearCart(services.Cart))
		cart.GET("/count", GetCartCount(services.Cart))
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
		
		// Покупки (только заказы, корзина доступна публично)
		setupOrderRoutes(protected, services)
		
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

func setupOrderRoutes(router *gin.RouterGroup, services *services.Services) {
	// Заказы (только для авторизованных пользователей)
	orders := router.Group("/orders")
	{
		orders.POST("/", CreateOrder(services.Order))
		orders.GET("/", GetUserOrders(services.Order))
		orders.GET("/:id", GetOrder(services.Order))
		orders.PUT("/:id", UpdateOrder(services.Order))
	}

	// Избранное (только для авторизованных пользователей)
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

	// Управление вариантами товаров
	variants := router.Group("/product-variants")
	{
		variants.POST("/", CreateProductVariant(services.ProductVariant))
		variants.GET("/:id", GetProductVariant(services.ProductVariant))
		variants.PUT("/:id", UpdateProductVariant(services.ProductVariant))
		variants.DELETE("/:id", DeleteProductVariant(services.ProductVariant))
	}

	// Управление категориями
	categories := router.Group("/categories")
	{
		categories.GET("/", GetCategories(services.Category))
		categories.POST("/", CreateCategory(services.Category))
		categories.GET("/:id", GetCategory(services.Category))
		categories.PUT("/:id", UpdateCategory(services.Category))
		categories.DELETE("/:id", DeleteCategory(services.Category))
	}

	// Управление складами
	warehouses := router.Group("/warehouses")
	{
		warehouses.POST("/", CreateWarehouse(services.Warehouse))
		warehouses.GET("/:id", GetWarehouse(services.Warehouse))
		warehouses.PUT("/:id", UpdateWarehouse(services.Warehouse))
		warehouses.DELETE("/:id", DeleteWarehouse(services.Warehouse))
	}

	// Управление остатками товаров
	warehouseStocks := router.Group("/warehouse-stocks")
	{
		warehouseStocks.POST("/", CreateWarehouseStock(services.WarehouseStock))
		warehouseStocks.PUT("/:id", UpdateWarehouseStock(services.WarehouseStock))
		warehouseStocks.DELETE("/:id", DeleteWarehouseStock(services.WarehouseStock))
	}

	// Управление изображениями
	images := router.Group("/images")
	{
		images.POST("/product/:id", UploadProductImage(services.Image))
		images.DELETE("/:id", DeleteImage(services.Image))
		images.PUT("/:id/primary", SetPrimaryImage(services.Image))
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
