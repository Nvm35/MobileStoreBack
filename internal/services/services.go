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
	// AddressService удален - адреса теперь встроены в User
	// CouponService удален - купоны больше не используются
}

func New(repos *repository.Repository, cfg *config.Config) *Services {
	return &Services{
		Auth:           NewAuthService(repos.Auth, cfg),
		User:           NewUserService(repos.User),
		Product:        NewProductService(repos.Product),
		ProductVariant: NewProductVariantService(repos.ProductVariant),
		Order:          NewOrderService(repos.Order),
		Cart:           NewCartService(repos.Cart),
		Wishlist:       NewWishlistService(repos.Wishlist),
		Review:         NewReviewService(repos.Review),
		Category:       NewCategoryService(repos.Category),
		// Address:  NewAddressService(repos.Address), // удален
		// Coupon:   NewCouponService(repos.Coupon), // удален
	}
}
