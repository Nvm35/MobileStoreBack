package services

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type WarehouseService struct {
	repo repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) *WarehouseService {
	return &WarehouseService{
		repo: repo,
	}
}

func (s *WarehouseService) Create(warehouse *models.Warehouse) error {
	return s.repo.Create(warehouse)
}

func (s *WarehouseService) GetByID(id string) (*models.Warehouse, error) {
	return s.repo.GetByID(id)
}

func (s *WarehouseService) GetBySlug(slug string) (*models.Warehouse, error) {
	return s.repo.GetBySlug(slug)
}

func (s *WarehouseService) GetBySlugOrID(identifier string) (*models.Warehouse, error) {
	return s.repo.GetBySlugOrID(identifier)
}

func (s *WarehouseService) GetByCity(city string) ([]*models.Warehouse, error) {
	return s.repo.GetByCity(city)
}

func (s *WarehouseService) GetMain() (*models.Warehouse, error) {
	return s.repo.GetMain()
}

func (s *WarehouseService) List() ([]*models.Warehouse, error) {
	return s.repo.List()
}

func (s *WarehouseService) Update(id string, name *string, address *string, city *string, phone *string, email *string, isActive *bool, managerName *string) (*models.Warehouse, error) {
	return s.repo.Update(id, name, address, city, phone, email, isActive, managerName)
}

func (s *WarehouseService) Delete(id string) error {
	return s.repo.Delete(id)
}
