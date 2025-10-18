package repository

import (
	"mobile-store-back/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	User            UserRepository
	Product         ProductRepository
	ProductVariant  ProductVariantRepository
	Order           OrderRepository
	Auth            AuthRepository
	Cart            CartRepository
	Wishlist        WishlistRepository
	Review          ReviewRepository
	Category        CategoryRepository
	Warehouse       WarehouseRepository
	WarehouseStock  WarehouseStockRepository
	// AddressRepository удален - адреса теперь встроены в User
	// CouponRepository удален - купоны больше не используются
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	UpdateProfile(userID string, firstName *string, lastName *string, phone *string) (*models.User, error)
	Update(id string, firstName *string, lastName *string, phone *string, isActive *bool, isAdmin *bool) (*models.User, error)
	Delete(id string) error
	List(limit, offset int) ([]*models.User, error)
}

type ProductRepository interface {
	Create(name string, slug string, description string, basePrice float64, sku string, isActive bool, brand string, model string, material string, categoryID string, tags []string) (*models.Product, error)
	GetByID(id string) (*models.Product, error)
	GetBySlug(slug string) (*models.Product, error)
	GetBySKU(sku string) (*models.Product, error)
	Update(id string, name *string, description *string, basePrice *float64, isActive *bool, brand *string, model *string, material *string, categoryID *string, tags []string) (*models.Product, error)
	Delete(id string) error
	List(limit, offset int) ([]*models.Product, error)
	Search(query string, limit, offset int) ([]*models.Product, error)
	GetByCategory(categoryID string, limit, offset int) ([]*models.Product, error)
}

type ProductVariantRepository interface {
	Create(productID string, sku string, name string, color string, size string, price float64, isActive bool) (*models.ProductVariant, error)
	GetByID(id string) (*models.ProductVariant, error)
	GetBySKU(sku string) (*models.ProductVariant, error)
	GetByProductID(productID string) ([]*models.ProductVariant, error)
	Update(id string, sku *string, name *string, color *string, size *string, price *float64, isActive *bool) (*models.ProductVariant, error)
	Delete(id string) error
}

type OrderRepository interface {
	Create(userID string, items []struct {
		ProductID        string
		ProductVariantID *string
		Quantity         int
	}, shippingMethod string, shippingAddress string, pickupPoint string, paymentMethod string, customerNotes string) (*models.Order, error)
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


type CategoryRepository interface {
	GetAll(limit, offset int) ([]*models.Category, error)
	GetByID(id string) (*models.Category, error)
	GetBySlug(slug string) (*models.Category, error)
	Create(category *models.Category) error
	Update(id string, name *string, description *string, slug *string, imageURL *string) (*models.Category, error)
	Delete(id string) error
	GetWithProducts(id string, limit, offset int) (*models.Category, error)
}

type WarehouseRepository interface {
	Create(warehouse *models.Warehouse) error
	GetByID(id string) (*models.Warehouse, error)
	GetBySlug(slug string) (*models.Warehouse, error)
	GetBySlugOrID(identifier string) (*models.Warehouse, error)
	GetByCity(city string) ([]*models.Warehouse, error)
	GetMain() (*models.Warehouse, error)
	List(limit, offset int) ([]*models.Warehouse, error)
	Update(id string, name *string, address *string, city *string, phone *string, email *string, isActive *bool, managerName *string) (*models.Warehouse, error)
	Delete(id string) error
}

type WarehouseStockRepository interface {
	Create(warehouseStock *models.WarehouseStock) error
	GetByID(id string) (*models.WarehouseStock, error)
	GetByWarehouseAndVariant(warehouseID, variantID string) (*models.WarehouseStock, error)
	GetByVariant(variantID string) ([]*models.WarehouseStock, error)
	GetByWarehouse(warehouseID string) ([]*models.WarehouseStock, error)
	GetAvailableStock(variantID string) (int, error)
	GetAvailableStockByWarehouse(warehouseID, variantID string) (int, error)
	UpdateStock(id string, stock, reservedStock int) (*models.WarehouseStock, error)
	ReserveStock(warehouseID, variantID string, quantity int) error
	ReleaseReservedStock(warehouseID, variantID string, quantity int) error
	ConsumeStock(warehouseID, variantID string, quantity int) error
	Delete(id string) error
}

// AddressRepository удален - адреса теперь встроены в User

func New(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		User:           NewUserRepository(db, redis),
		Product:        NewProductRepository(db, redis),
		ProductVariant: NewProductVariantRepository(db, redis),
		Order:          NewOrderRepository(db, redis),
		Auth:           NewAuthRepository(db, redis),
		Cart:           NewCartRepository(db, redis),
		Wishlist:       NewWishlistRepository(db, redis),
		Review:         NewReviewRepository(db, redis),
		Category:       NewCategoryRepository(db, redis),
		Warehouse:      NewWarehouseRepository(db, redis),
		WarehouseStock: NewWarehouseStockRepository(db, redis),
		// Address:  NewAddressRepository(db, redis), // удален
		// Coupon:   NewCouponRepository(db, redis), // удален
	}
}
