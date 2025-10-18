package repository

import (
	"mobile-store-back/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type warehouseStockRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewWarehouseStockRepository(db *gorm.DB, redis *redis.Client) WarehouseStockRepository {
	return &warehouseStockRepository{
		db:    db,
		redis: redis,
	}
}

func (r *warehouseStockRepository) Create(warehouseStock *models.WarehouseStock) error {
	return r.db.Create(warehouseStock).Error
}

func (r *warehouseStockRepository) GetByID(id string) (*models.WarehouseStock, error) {
	var stock models.WarehouseStock
	err := r.db.Preload("Warehouse").Preload("ProductVariant").
		First(&stock, "id = ?", id).Error
	return &stock, err
}

func (r *warehouseStockRepository) GetByWarehouseAndVariant(warehouseID, variantID string) (*models.WarehouseStock, error) {
	var stock models.WarehouseStock
	err := r.db.Preload("Warehouse").Preload("ProductVariant").
		Where("warehouse_id = ? AND product_variant_id = ?", warehouseID, variantID).
		First(&stock).Error
	return &stock, err
}

func (r *warehouseStockRepository) GetByVariant(variantID string) ([]*models.WarehouseStock, error) {
	var stocks []*models.WarehouseStock
	err := r.db.Preload("Warehouse").
		Where("product_variant_id = ?", variantID).
		Find(&stocks).Error
	return stocks, err
}

func (r *warehouseStockRepository) GetByWarehouse(warehouseID string) ([]*models.WarehouseStock, error) {
	var stocks []*models.WarehouseStock
	err := r.db.Preload("ProductVariant").
		Where("warehouse_id = ?", warehouseID).
		Find(&stocks).Error
	return stocks, err
}

func (r *warehouseStockRepository) GetAvailableStock(variantID string) (int, error) {
	var totalStock int
	err := r.db.Model(&models.WarehouseStock{}).
		Where("product_variant_id = ?", variantID).
		Select("COALESCE(SUM(stock - reserved_stock), 0)").
		Scan(&totalStock).Error
	return totalStock, err
}

func (r *warehouseStockRepository) GetAvailableStockByWarehouse(warehouseID, variantID string) (int, error) {
	var stock int
	err := r.db.Model(&models.WarehouseStock{}).
		Where("warehouse_id = ? AND product_variant_id = ?", warehouseID, variantID).
		Select("COALESCE(stock - reserved_stock, 0)").
		Scan(&stock).Error
	return stock, err
}

func (r *warehouseStockRepository) UpdateStock(id string, stock, reservedStock int) (*models.WarehouseStock, error) {
	var warehouseStock models.WarehouseStock
	if err := r.db.First(&warehouseStock, "id = ?", id).Error; err != nil {
		return nil, err
	}

	warehouseStock.Stock = stock
	warehouseStock.ReservedStock = reservedStock

	err := r.db.Save(&warehouseStock).Error
	return &warehouseStock, err
}

func (r *warehouseStockRepository) ReserveStock(warehouseID, variantID string, quantity int) error {
	return r.db.Model(&models.WarehouseStock{}).
		Where("warehouse_id = ? AND product_variant_id = ? AND (stock - reserved_stock) >= ?", 
			warehouseID, variantID, quantity).
		Update("reserved_stock", gorm.Expr("reserved_stock + ?", quantity)).Error
}

func (r *warehouseStockRepository) ReleaseReservedStock(warehouseID, variantID string, quantity int) error {
	return r.db.Model(&models.WarehouseStock{}).
		Where("warehouse_id = ? AND product_variant_id = ? AND reserved_stock >= ?", 
			warehouseID, variantID, quantity).
		Update("reserved_stock", gorm.Expr("reserved_stock - ?", quantity)).Error
}

func (r *warehouseStockRepository) ConsumeStock(warehouseID, variantID string, quantity int) error {
	return r.db.Model(&models.WarehouseStock{}).
		Where("warehouse_id = ? AND product_variant_id = ? AND reserved_stock >= ?", 
			warehouseID, variantID, quantity).
		Updates(map[string]interface{}{
			"stock":         gorm.Expr("stock - ?", quantity),
			"reserved_stock": gorm.Expr("reserved_stock - ?", quantity),
		}).Error
}

func (r *warehouseStockRepository) Delete(id string) error {
	return r.db.Delete(&models.WarehouseStock{}, "id = ?", id).Error
}

