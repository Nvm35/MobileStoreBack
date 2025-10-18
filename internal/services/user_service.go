package services

import (
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *UserService) UpdateProfile(userID string, firstName *string, lastName *string, phone *string) (*models.User, error) {
	return s.repo.UpdateProfile(userID, firstName, lastName, phone)
}

func (s *UserService) Update(id string, firstName *string, lastName *string, phone *string, isActive *bool, isAdmin *bool) (*models.User, error) {
	return s.repo.Update(id, firstName, lastName, phone, isActive, isAdmin)
}

func (s *UserService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *UserService) List(limit, offset int) ([]*models.User, error) {
	return s.repo.List(limit, offset)
}
