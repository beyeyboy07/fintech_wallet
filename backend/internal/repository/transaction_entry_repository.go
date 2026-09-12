package repository

import (
	"fintech-wallet/backend/internal/model"

	"gorm.io/gorm"
)

type TransactionEntryRepository struct {
	db *gorm.DB
}

func NewTransactionEntryRepository(
	db *gorm.DB,
) *TransactionEntryRepository {
	return &TransactionEntryRepository{
		db: db,
	}
}

func (r *TransactionEntryRepository) WithTx(
	tx *gorm.DB,
) *TransactionEntryRepository {
	return &TransactionEntryRepository{
		db: tx,
	}
}

func (r *TransactionEntryRepository) Create(
	entry *model.TransactionEntry,
) error {
	return r.db.Create(entry).Error
}
