package repository

import (
	"fintech-wallet/backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IdempotencyRepository struct {
	db *gorm.DB
}

func NewIdempotencyRepository(db *gorm.DB) *IdempotencyRepository {
	return &IdempotencyRepository{
		db: db,
	}
}

func (r *IdempotencyRepository) WithTx(
	tx *gorm.DB,
) *IdempotencyRepository {
	return &IdempotencyRepository{
		db: tx,
	}
}

func (r *IdempotencyRepository) Create(
	idempotencyKey *model.IdempotencyKey,
) error {
	return r.db.Create(idempotencyKey).Error
}

func (r *IdempotencyRepository) FindByUserIDAndKey(
	userID uint,
	key string,
) (*model.IdempotencyKey, error) {

	var idempotencyKey model.IdempotencyKey

	err := r.db.
		Where("user_id = ? AND `key` = ?", userID, key).
		First(&idempotencyKey).
		Error

	if err != nil {
		return nil, err
	}

	return &idempotencyKey, nil
}

func (r *IdempotencyRepository) FindByUserIDAndKeyForUpdate(
	userID uint,
	key string,
) (*model.IdempotencyKey, error) {

	var idempotencyKey model.IdempotencyKey

	err := r.db.
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("user_id = ? AND `key` = ?", userID, key).
		First(&idempotencyKey).
		Error

	if err != nil {
		return nil, err
	}

	return &idempotencyKey, nil
}

func (r *IdempotencyRepository) Update(
	idempotencyKey *model.IdempotencyKey,
) error {
	return r.db.Save(idempotencyKey).Error
}
