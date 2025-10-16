package repository

import (
	"mobile-store-back/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	User     UserRepository
	Product  ProductRepository
	Order    OrderRepository
	Auth     AuthRepository
	Cart     CartRepository
	Wishlist WishlistRepository
	Review   ReviewRepository
	Coupon   CouponRepository
	// AddressRepository удален - адреса теперь встроены в User
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	UpdateProfile(userID string, firstName *string, lastName *string, phone *string, dateOfBirth *string, gender *string) (*models.User, error)
	Update(id string, firstName *string, lastName *string, phone *string, dateOfBirth *string, gender *string, isActive *bool, isAdmin *bool) (*models.User, error)
	Delete(id string) error
	List(limit, offset int) ([]*models.User, error)
}

type ProductRepository interface {
	Create(name string, description string, shortDescription string, price float64, comparePrice *float64, sku string, stock int, isActive bool, isFeatured bool, isNew bool, weight *float64, dimensions string, brand string, model string, color string, material string, categoryID string, tags []string, metaTitle string, metaDescription string) (*models.Product, error)
	GetByID(id string) (*models.Product, error)
	GetBySKU(sku string) (*models.Product, error)
	Update(id string, name *string, description *string, shortDescription *string, price *float64, comparePrice *float64, stock *int, isActive *bool, isFeatured *bool, isNew *bool, weight *float64, dimensions *string, brand *string, model *string, color *string, material *string, categoryID *string, tags []string, metaTitle *string, metaDescription *string) (*models.Product, error)
	Delete(id string) error
	List(limit, offset int) ([]*models.Product, error)
	Search(query string, limit, offset int) ([]*models.Product, error)
	GetByCategory(categoryID string, limit, offset int) ([]*models.Product, error)
}

type OrderRepository interface {
	Create(userID string, items []struct {
		ProductID string
		Quantity  int
	}, shippingMethod string, shippingAddress string, pickupPoint string, paymentMethod string, customerNotes string, couponCode string) (*models.Order, error)
	GetByID(id string) (*models.Order, error)
	GetByUserID(userID string, limit, offset int) ([]*models.Order, error)
	Update(id string, userID string, status *string, paymentStatus *string, trackingNumber *string, customerNotes *string, shippingMethod *string, shippingAddress *string, pickupPoint *string) (*models.Order, error)
	UpdateStatus(id string, status string, trackingNumber *string) (*models.Order, error)
	Delete(id string) error
	List(limit, offset int) ([]*models.Order, error)
}

type AuthRepository interface {
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
}

type CartRepository interface {
	GetByUserID(userID string) ([]models.CartItem, error)
	AddItem(userID string, productID string, quantity int) (*models.CartItem, error)
	UpdateItem(id string, userID string, quantity int) (*models.CartItem, error)
	RemoveItem(id string, userID string) error
	Clear(userID string) error
	GetCount(userID string) (int, error)
}

type WishlistRepository interface {
	GetByUserID(userID string, limit, offset int) ([]models.WishlistItem, error)
	AddItem(userID string, productID string) (*models.WishlistItem, error)
	RemoveItem(id string, userID string) error
	IsInWishlist(userID string, productID string) (bool, error)
	Clear(userID string) error
}

type ReviewRepository interface {
	GetByProductID(productID string, limit, offset int) ([]models.Review, error)
	Create(userID string, productID string, orderID *string, rating int, title string, comment string) (*models.Review, error)
	Update(id string, userID string, rating *int, title *string, comment *string) (*models.Review, error)
	Delete(id string, userID string) error
	Vote(id string, userID string, helpful bool) error
	GetByUserID(userID string, limit, offset int) ([]models.Review, error)
	GetAll(limit, offset int) ([]models.Review, error)
	Approve(id string, approved bool) error
}

type CouponRepository interface {
	GetAll(limit, offset int) ([]models.Coupon, error)
	GetByID(id string) (*models.Coupon, error)
	Validate(code string, userID string, orderAmount float64) (*models.Coupon, error)
	Create(code string, name string, description string, couponType string, value float64, minimumAmount float64, maximumDiscount *float64, usageLimit *int, startsAt *string, expiresAt *string) (*models.Coupon, error)
	Update(id string, name *string, description *string, value *float64, minimumAmount *float64, maximumDiscount *float64, usageLimit *int, isActive *bool, startsAt *string, expiresAt *string) (*models.Coupon, error)
	Delete(id string) error
	GetUsage(id string) ([]models.CouponUsage, error)
}

// AddressRepository удален - адреса теперь встроены в User

func New(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		User:     NewUserRepository(db, redis),
		Product:  NewProductRepository(db, redis),
		Order:    NewOrderRepository(db, redis),
		Auth:     NewAuthRepository(db, redis),
		Cart:     NewCartRepository(db, redis),
		Wishlist: NewWishlistRepository(db, redis),
		Review:   NewReviewRepository(db, redis),
		Coupon:   NewCouponRepository(db, redis),
		// Address:  NewAddressRepository(db, redis), // удален
	}
}
