package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type WarehouseStockService struct {
	repo          repository.WarehouseStockRepository
	warehouseRepo repository.WarehouseRepository
	variantRepo   repository.ProductVariantRepository
}

func NewWarehouseStockService(repo repository.WarehouseStockRepository, warehouseRepo repository.WarehouseRepository, variantRepo repository.ProductVariantRepository) *WarehouseStockService {
	return &WarehouseStockService{
		repo:          repo,
		warehouseRepo: warehouseRepo,
		variantRepo:   variantRepo,
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
	if s.warehouseRepo == nil {
		return nil, errors.New("warehouse repository dependency is not configured")
	}
	warehouse, err := s.warehouseRepo.GetBySlug(warehouseSlug)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByWarehouse(warehouse.ID.String())
}

func (s *WarehouseStockService) GetByVariantSKU(sku string) ([]*models.WarehouseStock, error) {
	if s.variantRepo == nil {
		return nil, errors.New("product variant repository dependency is not configured")
	}
	variant, err := s.variantRepo.GetBySKU(sku)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByVariant(variant.ID.String())
}

func (s *WarehouseStockService) GetAvailabilityInfoBySKU(sku string) ([]models.WarehouseStock, error) {
	if s.variantRepo == nil {
		return nil, errors.New("product variant repository dependency is not configured")
	}
	variant, err := s.variantRepo.GetBySKU(sku)
	if err != nil {
		return nil, err
	}
	return s.GetAvailabilityInfo(variant.ID.String())
}

func (s *WarehouseStockService) CheckAvailabilityBySKU(sku string, quantity int) (bool, error) {
	if s.variantRepo == nil {
		return false, errors.New("product variant repository dependency is not configured")
	}

	variant, err := s.variantRepo.GetBySKU(sku)
	if err != nil {
		return false, err
	}
	return s.CheckAvailability(variant.ID.String(), quantity)
}

func (s *WarehouseStockService) CheckAvailabilityByWarehouseSlug(warehouseSlug, sku string, quantity int) (bool, error) {
	if s.warehouseRepo == nil || s.variantRepo == nil {
		return false, errors.New("warehouse or product variant repository dependency is not configured")
	}

	warehouse, err := s.warehouseRepo.GetBySlug(warehouseSlug)
	if err != nil {
		return false, err
	}

	variant, err := s.variantRepo.GetBySKU(sku)
	if err != nil {
		return false, err
	}

	return s.CheckAvailabilityByWarehouse(warehouse.ID.String(), variant.ID.String(), quantity)
}

func (s *WarehouseStockService) GetAvailableStockBySKU(sku string) (int, error) {
	if s.variantRepo == nil {
		return 0, errors.New("product variant repository dependency is not configured")
	}

	variant, err := s.variantRepo.GetBySKU(sku)
	if err != nil {
		return 0, err
	}

	return s.repo.GetAvailableStock(variant.ID.String())
}

// TransferStock перемещает остатки между складами по идентификаторам (UUID или slug для складов и UUID или SKU для варианта)
func (s *WarehouseStockService) TransferStock(sourceWarehouse, destinationWarehouse, variantIdentifier string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	sourceWarehouse = strings.TrimSpace(sourceWarehouse)
	destinationWarehouse = strings.TrimSpace(destinationWarehouse)
	variantIdentifier = strings.TrimSpace(variantIdentifier)

	if sourceWarehouse == "" || destinationWarehouse == "" || variantIdentifier == "" {
		return errors.New("source warehouse, destination warehouse and variant identifiers are required")
	}

	if sourceWarehouse == destinationWarehouse {
		return errors.New("source and destination warehouses must be different")
	}

	sourceID, err := s.resolveWarehouseID(sourceWarehouse)
	if err != nil {
		return err
	}

	destID, err := s.resolveWarehouseID(destinationWarehouse)
	if err != nil {
		return err
	}

	variantID, err := s.resolveVariantID(variantIdentifier)
	if err != nil {
		return err
	}

	return s.repo.TransferStock(sourceID, destID, variantID, quantity)
}

func (s *WarehouseStockService) resolveWarehouseID(identifier string) (string, error) {
	if _, err := uuid.Parse(identifier); err == nil {
		return identifier, nil
	}
	if s.warehouseRepo == nil {
		return "", errors.New("warehouse repository dependency is not configured")
	}
	warehouse, err := s.warehouseRepo.GetBySlugOrID(identifier)
	if err != nil {
		return "", err
	}
	return warehouse.ID.String(), nil
}

func (s *WarehouseStockService) resolveVariantID(identifier string) (string, error) {
	if _, err := uuid.Parse(identifier); err == nil {
		return identifier, nil
	}
	if s.variantRepo == nil {
		return "", errors.New("product variant repository dependency is not configured")
	}
	variant, err := s.variantRepo.GetBySKU(identifier)
	if err != nil {
		return "", err
	}
	return variant.ID.String(), nil
}

func (s *WarehouseStockService) List() ([]*models.WarehouseStock, error) {
	return s.repo.List()
}
