package repository

import (
	"mobile-store-back/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type warehouseRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewWarehouseRepository(db *gorm.DB, redis *redis.Client) WarehouseRepository {
	return &warehouseRepository{
		db:    db,
		redis: redis,
	}
}

func (r *warehouseRepository) Create(warehouse *models.Warehouse) error {
	return r.db.Create(warehouse).Error
}

func (r *warehouseRepository) GetByID(id string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := r.db.First(&warehouse, "id = ?", id).Error
	return &warehouse, err
}

func (r *warehouseRepository) GetBySlug(slug string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := r.db.First(&warehouse, "slug = ?", slug).Error
	return &warehouse, err
}

func (r *warehouseRepository) GetBySlugOrID(identifier string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	// Сначала пробуем найти по slug
	err := r.db.First(&warehouse, "slug = ?", identifier).Error
	if err == nil {
		return &warehouse, nil
	}
	// Если не найден по slug, пробуем по ID
	err = r.db.First(&warehouse, "id = ?", identifier).Error
	return &warehouse, err
}

func (r *warehouseRepository) GetByCity(city string) ([]*models.Warehouse, error) {
	var warehouses []*models.Warehouse
	err := r.db.Where("city = ? AND is_active = ?", city, true).Find(&warehouses).Error
	return warehouses, err
}

func (r *warehouseRepository) GetMain() (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := r.db.Where("is_main = ? AND is_active = ?", true, true).First(&warehouse).Error
	return &warehouse, err
}

func (r *warehouseRepository) List(limit, offset int) ([]*models.Warehouse, error) {
	var warehouses []*models.Warehouse
	err := r.db.Where("is_active = ?", true).
		Limit(limit).Offset(offset).
		Order("is_main DESC, name ASC").
		Find(&warehouses).Error
	return warehouses, err
}

func (r *warehouseRepository) Update(id string, name *string, address *string, city *string, phone *string, email *string, isActive *bool, managerName *string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	if err := r.db.First(&warehouse, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if name != nil {
		warehouse.Name = *name
	}
	if address != nil {
		warehouse.Address = *address
	}
	if city != nil {
		warehouse.City = *city
	}
	if phone != nil {
		warehouse.Phone = *phone
	}
	if email != nil {
		warehouse.Email = *email
	}
	if isActive != nil {
		warehouse.IsActive = *isActive
	}
	if managerName != nil {
		warehouse.ManagerName = *managerName
	}

	err := r.db.Save(&warehouse).Error
	return &warehouse, err
}

func (r *warehouseRepository) Delete(id string) error {
	return r.db.Delete(&models.Warehouse{}, "id = ?", id).Error
}
