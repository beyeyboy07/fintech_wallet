package repository

import (
	"fintech-wallet/backend/internal/model"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) Create(
	refreshToken *model.RefreshToken,
) error {
	return r.db.Create(refreshToken).Error
}

func (r *RefreshTokenRepository) FindByTokenHash(
	tokenHash string,
) (*model.RefreshToken, error) {

	var refreshToken model.RefreshToken

	err := r.db.
		Where("token_hash = ?", tokenHash).
		First(&refreshToken).
		Error

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) Revoke(
	refreshToken *model.RefreshToken,
) error {

	now := gorm.Expr("NOW()")

	return r.db.
		Model(refreshToken).
		Update("revoked_at", now).
		Error
}
