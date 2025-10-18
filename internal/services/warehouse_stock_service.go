package services

import (
	"errors"
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type WarehouseStockService struct {
	repo repository.WarehouseStockRepository
}

func NewWarehouseStockService(repo repository.WarehouseStockRepository) *WarehouseStockService {
	return &WarehouseStockService{
		repo: repo,
	}
}

func (s *WarehouseStockService) Create(warehouseStock *models.WarehouseStock) error {
	return s.repo.Create(warehouseStock)
}

func (s *WarehouseStockService) GetByID(id string) (*models.WarehouseStock, error) {
	return s.repo.GetByID(id)
}

func (s *WarehouseStockService) GetByWarehouseAndVariant(warehouseID, variantID string) (*models.WarehouseStock, error) {
	return s.repo.GetByWarehouseAndVariant(warehouseID, variantID)
}

func (s *WarehouseStockService) GetByVariant(variantID string) ([]*models.WarehouseStock, error) {
	return s.repo.GetByVariant(variantID)
}

func (s *WarehouseStockService) GetByWarehouse(warehouseID string) ([]*models.WarehouseStock, error) {
	return s.repo.GetByWarehouse(warehouseID)
}

func (s *WarehouseStockService) GetAvailableStock(variantID string) (int, error) {
	return s.repo.GetAvailableStock(variantID)
}

func (s *WarehouseStockService) GetAvailableStockByWarehouse(warehouseID, variantID string) (int, error) {
	return s.repo.GetAvailableStockByWarehouse(warehouseID, variantID)
}

func (s *WarehouseStockService) UpdateStock(id string, stock, reservedStock int) (*models.WarehouseStock, error) {
	return s.repo.UpdateStock(id, stock, reservedStock)
}

func (s *WarehouseStockService) ReserveStock(warehouseID, variantID string, quantity int) error {
	return s.repo.ReserveStock(warehouseID, variantID, quantity)
}

func (s *WarehouseStockService) ReleaseReservedStock(warehouseID, variantID string, quantity int) error {
	return s.repo.ReleaseReservedStock(warehouseID, variantID, quantity)
}

func (s *WarehouseStockService) ConsumeStock(warehouseID, variantID string, quantity int) error {
	return s.repo.ConsumeStock(warehouseID, variantID, quantity)
}

func (s *WarehouseStockService) Delete(id string) error {
	return s.repo.Delete(id)
}

// GetAvailabilityInfo возвращает информацию о доступности товара по всем складам
func (s *WarehouseStockService) GetAvailabilityInfo(variantID string) ([]models.WarehouseStock, error) {
	stocks, err := s.repo.GetByVariant(variantID)
	if err != nil {
		return nil, err
	}

	// Конвертируем []*models.WarehouseStock в []models.WarehouseStock
	result := make([]models.WarehouseStock, len(stocks))
	for i, stock := range stocks {
		result[i] = *stock
	}

	return result, nil
}

// CheckAvailability проверяет, доступен ли товар в нужном количестве
func (s *WarehouseStockService) CheckAvailability(variantID string, quantity int) (bool, error) {
	availableStock, err := s.repo.GetAvailableStock(variantID)
	if err != nil {
		return false, err
	}
	return availableStock >= quantity, nil
}

// CheckAvailabilityByWarehouse проверяет доступность товара на конкретном складе
func (s *WarehouseStockService) CheckAvailabilityByWarehouse(warehouseID, variantID string, quantity int) (bool, error) {
	availableStock, err := s.repo.GetAvailableStockByWarehouse(warehouseID, variantID)
	if err != nil {
		return false, err
	}
	return availableStock >= quantity, nil
}

// Дополнительные методы для работы с slug и SKU
func (s *WarehouseStockService) GetByWarehouseSlug(warehouseSlug string) ([]*models.WarehouseStock, error) {
	// Нужно получить ID склада по slug
	// Пока возвращаем ошибку, так как нужен доступ к WarehouseService
	return nil, errors.New("not implemented - need WarehouseService dependency")
}

func (s *WarehouseStockService) GetByVariantSKU(sku string) ([]*models.WarehouseStock, error) {
	// Нужно получить ID варианта по SKU
	// Пока возвращаем ошибку, так как нужен доступ к ProductVariantService
	return nil, errors.New("not implemented - need ProductVariantService dependency")
}

func (s *WarehouseStockService) GetAvailabilityInfoBySKU(sku string) ([]models.WarehouseStock, error) {
	// Нужно получить ID варианта по SKU
	return nil, errors.New("not implemented - need ProductVariantService dependency")
}

func (s *WarehouseStockService) CheckAvailabilityBySKU(sku string, quantity int) (bool, error) {
	// Нужно получить ID варианта по SKU
	return false, errors.New("not implemented - need ProductVariantService dependency")
}

func (s *WarehouseStockService) CheckAvailabilityByWarehouseSlug(warehouseSlug, sku string, quantity int) (bool, error) {
	// Нужно получить ID склада по slug и ID варианта по SKU
	return false, errors.New("not implemented - need WarehouseService and ProductVariantService dependencies")
}

func (s *WarehouseStockService) GetAvailableStockBySKU(sku string) (int, error) {
	// Нужно получить ID варианта по SKU
	return 0, errors.New("not implemented - need ProductVariantService dependency")
}
