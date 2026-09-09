package service

import (
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(
	userRepository *repository.UserRepository,
) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUsers() ([]model.User, error) {
	return s.userRepository.FindAll()
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.userRepository.FindByID(id)
}
