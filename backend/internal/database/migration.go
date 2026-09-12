package database

import (
	"fintech-wallet/backend/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.Wallet{},
		&model.Transaction{},
		&model.TransactionEntry{},
		&model.IdempotencyKey{},
	)
}
