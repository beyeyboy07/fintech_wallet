package service

import (
	"errors"

	"fintech-wallet/backend/internal/auth"
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	userRepository         *repository.UserRepository
	refreshTokenRepository *repository.RefreshTokenRepository
	jwtManager             *auth.JWTManager
}

func NewAuthService(
	userRepository *repository.UserRepository,
	refreshTokenRepository *repository.RefreshTokenRepository,
	jwtManager *auth.JWTManager,
) *AuthService {
	return &AuthService{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		jwtManager:             jwtManager,
	}
}

func (s *AuthService) Register(
	name string,
	email string,
	password string,
) (*model.User, error) {

	existingUser, err := s.userRepository.FindByEmail(email)

	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "USER",
	}

	err = s.userRepository.Create(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(
	email string,
	password string,
) (string, error) {

	user, err := s.userRepository.FindByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}

		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GetCurrentUser(
	userID uint,
) (*model.User, error) {

	return s.userRepository.FindByID(userID)
}
