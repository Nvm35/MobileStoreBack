package repository

import (
	"errors"

	"mobile-store-back/internal/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
			"stock":          gorm.Expr("stock - ?", quantity),
			"reserved_stock": gorm.Expr("reserved_stock - ?", quantity),
		}).Error
}

func (r *warehouseStockRepository) Delete(id string) error {
	return r.db.Delete(&models.WarehouseStock{}, "id = ?", id).Error
}

func (r *warehouseStockRepository) TransferStock(warehouseFromID, warehouseToID, variantID string, quantity int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var source models.WarehouseStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND product_variant_id = ?", warehouseFromID, variantID).
			First(&source).Error; err != nil {
			return err
		}

		available := source.Stock - source.ReservedStock
		if available < quantity {
			return errors.New("insufficient available stock on source warehouse")
		}

		source.Stock -= quantity
		if err := tx.Save(&source).Error; err != nil {
			return err
		}

		var target models.WarehouseStock
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND product_variant_id = ?", warehouseToID, variantID).
			First(&target).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			toWarehouseID, parseErr := uuid.Parse(warehouseToID)
			if parseErr != nil {
				return parseErr
			}
			variantUUID, parseErr := uuid.Parse(variantID)
			if parseErr != nil {
				return parseErr
			}

			target = models.WarehouseStock{
				WarehouseID:      toWarehouseID,
				ProductVariantID: variantUUID,
				Stock:            quantity,
				ReservedStock:    0,
			}
			if err := tx.Create(&target).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			target.Stock += quantity
			if err := tx.Save(&target).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *warehouseStockRepository) List() ([]*models.WarehouseStock, error) {
	var stocks []*models.WarehouseStock
	err := r.db.Preload("Warehouse").Preload("ProductVariant").
		Find(&stocks).Error
	return stocks, err
}
