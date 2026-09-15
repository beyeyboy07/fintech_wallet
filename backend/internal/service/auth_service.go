package service

import (
	"errors"
	"time"

	"fintech-wallet/backend/internal/auth"
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"

	"gorm.io/gorm"
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	userRepository         *repository.UserRepository
	walletRepository       *repository.WalletRepository
	refreshTokenRepository *repository.RefreshTokenRepository
	jwtManager             *auth.JWTManager
	db                     *gorm.DB
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthService(
	userRepository *repository.UserRepository,
	walletRepository *repository.WalletRepository,
	refreshTokenRepository *repository.RefreshTokenRepository,
	jwtManager *auth.JWTManager,
	db *gorm.DB,
) *AuthService {
	return &AuthService{
		userRepository:         userRepository,
		walletRepository:       walletRepository,
		refreshTokenRepository: refreshTokenRepository,
		jwtManager:             jwtManager,
		db:                     db,
	}
}

func (s *AuthService) Register(
	name string,
	email string,
	password string,
) (*model.User, error) {

	var createdUser *model.User

	err := s.db.Transaction(func(tx *gorm.DB) error {

		userRepository := s.userRepository.WithTx(tx)
		walletRepository := s.walletRepository.WithTx(tx)

		// Cek email
		existingUser, err := userRepository.FindByEmail(email)

		if err == nil && existingUser != nil {
			return ErrEmailAlreadyExists
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Hash password
		hashedPassword, err := auth.HashPassword(password)

		if err != nil {
			return err
		}

		// Create User
		user := &model.User{
			Name:     name,
			Email:    email,
			Password: hashedPassword,
			Role:     model.RoleUser,
		}

		if err := userRepository.Create(user); err != nil {
			return err
		}

		// Create Wallet
		wallet := &model.Wallet{
			UserID:  user.ID,
			Balance: 0,
		}

		if err := walletRepository.Create(wallet); err != nil {
			return err
		}

		createdUser = user

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *AuthService) Login(
	email string,
	password string,
) (*LoginResult, error) {

	user, err := s.userRepository.FindByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	err = auth.CheckPassword(password, user.Password)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.jwtManager.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
	)

	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken()

	if err != nil {
		return nil, err
	}

	refreshTokenHash := auth.HashRefreshToken(refreshToken)

	refreshTokenModel := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	err = s.refreshTokenRepository.Create(refreshTokenModel)

	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GetCurrentUser(
	userID uint,
) (*model.User, error) {

	return s.userRepository.FindByID(userID)
}

func (s *AuthService) Refresh(
	refreshToken string,
) (*RefreshResult, error) {

	tokenHash := auth.HashRefreshToken(refreshToken)

	var result RefreshResult

	err := s.db.Transaction(func(tx *gorm.DB) error {

		// Cari refresh token lama
		var refreshTokenModel model.RefreshToken

		err := tx.
			Where("token_hash = ?", tokenHash).
			First(&refreshTokenModel).
			Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidCredentials
			}

			return err
		}

		// Cek revoked
		if refreshTokenModel.RevokedAt != nil {
			return ErrInvalidCredentials
		}

		// Cek expired
		if time.Now().After(refreshTokenModel.ExpiresAt) {
			return ErrInvalidCredentials
		}

		// Cari user
		user, err := s.userRepository.FindByID(refreshTokenModel.UserID)

		if err != nil {
			return err
		}

		// Revoke token lama
		now := time.Now()

		err = tx.
			Model(&refreshTokenModel).
			Update("revoked_at", now).
			Error

		if err != nil {
			return err
		}

		// Generate refresh token baru
		newRefreshToken, err := auth.GenerateRefreshToken()

		if err != nil {
			return err
		}

		newRefreshTokenHash := auth.HashRefreshToken(
			newRefreshToken,
		)

		newRefreshTokenModel := &model.RefreshToken{
			UserID:    user.ID,
			TokenHash: newRefreshTokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}

		// Simpan refresh token baru
		err = tx.Create(newRefreshTokenModel).Error

		if err != nil {
			return err
		}

		// Generate access token baru
		newAccessToken, err := s.jwtManager.GenerateToken(
			user.ID,
			user.Email,
			user.Role,
		)

		if err != nil {
			return err
		}

		result = RefreshResult{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	tokenHash := auth.HashRefreshToken(refreshToken)

	refreshTokenModel, err := s.refreshTokenRepository.FindByTokenHash(tokenHash)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidCredentials
		}

		return err
	}

	if refreshTokenModel.RevokedAt != nil {
		return ErrInvalidCredentials
	}

	if time.Now().After(refreshTokenModel.ExpiresAt) {
		return ErrInvalidCredentials
	}

	return s.refreshTokenRepository.Revoke(refreshTokenModel)
}
