package repository

import (
	"fintech-wallet/backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (r *WalletRepository) WithTx(tx *gorm.DB) *WalletRepository {
	return &WalletRepository{
		db: tx,
	}
}

func (r *WalletRepository) Create(
	wallet *model.Wallet,
) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepository) FindByUserID(
	userID uint,
) (*model.Wallet, error) {

	var wallet model.Wallet

	err := r.db.
		Where("user_id = ?", userID).
		First(&wallet).
		Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepository) UpdateBalance(
	userID uint,
	amount int64,
) error {

	return r.db.
		Model(&model.Wallet{}).
		Where("user_id = ?", userID).
		UpdateColumn(
			"balance",
			gorm.Expr("balance + ?", amount),
		).
		Error
}

func (r *WalletRepository) FindByUserIDForUpdate(
	userID uint,
) (*model.Wallet, error) {

	var wallet model.Wallet

	err := r.db.
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("user_id = ?", userID).
		First(&wallet).
		Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
