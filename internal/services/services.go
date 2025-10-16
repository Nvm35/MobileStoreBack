package services

import (
	"mobile-store-back/internal/config"
	"mobile-store-back/internal/repository"
)

type Services struct {
	Auth      *AuthService
	User      *UserService
	Product   *ProductService
	Order     *OrderService
	Cart      *CartService
	Wishlist  *WishlistService
	Review    *ReviewService
	Coupon    *CouponService
	// AddressService удален - адреса теперь встроены в User
}

func New(repos *repository.Repository, cfg *config.Config) *Services {
	return &Services{
		Auth:     NewAuthService(repos.Auth, cfg),
		User:     NewUserService(repos.User),
		Product:  NewProductService(repos.Product),
		Order:    NewOrderService(repos.Order),
		Cart:     NewCartService(repos.Cart),
		Wishlist: NewWishlistService(repos.Wishlist),
		Review:   NewReviewService(repos.Review),
		Coupon:   NewCouponService(repos.Coupon),
		// Address:  NewAddressService(repos.Address), // удален
	}
}
