package services

import (
	"mobile-store-back/internal/config"
	"mobile-store-back/internal/repository"
)

type Services struct {
	Auth           *AuthService
	User           *UserService
	Product        *ProductService
	ProductVariant *ProductVariantService
	Order          *OrderService
	Cart           *CartService
	Wishlist       *WishlistService
	Review         *ReviewService
	Category       *CategoryService
	Warehouse      *WarehouseService
	WarehouseStock *WarehouseStockService
	Image          *ImageService
}

func New(repos *repository.Repository, cfg *config.Config) *Services {
	return &Services{
		Auth:           NewAuthService(repos.Auth, cfg),
		User:           NewUserService(repos.User),
		Product:        NewProductService(repos.Product),
		ProductVariant: NewProductVariantService(repos.ProductVariant, repos.Product),
		Order:          NewOrderService(repos.Order),
		Cart:           NewCartService(repos.Cart),
		Wishlist:       NewWishlistService(repos.Wishlist),
		Review:         NewReviewService(repos.Review),
		Category:       NewCategoryService(repos.Category),
		Warehouse:      NewWarehouseService(repos.Warehouse),
		WarehouseStock: NewWarehouseStockService(repos.WarehouseStock, repos.Warehouse, repos.ProductVariant),
		Image:          NewImageService(repos.Image),
	}
}
