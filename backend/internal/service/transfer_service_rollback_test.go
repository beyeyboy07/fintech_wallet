package service

import (
	"errors"
	"testing"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"

	"gorm.io/gorm"
)

type FailingTransactionRepository struct {
	err error
}

func (r *FailingTransactionRepository) UpdateStatus(
	id uint,
	status string,
) error {
	return r.err
}

func (r *FailingTransactionRepository) FindByIDForUpdate(
	id uint,
) (*model.Transaction, error) {
	return nil, r.err
}

func (r *FailingTransactionRepository) Create(
	transaction *model.Transaction,
) error {
	return r.err
}

func (r *FailingTransactionRepository) WithTx(
	tx *gorm.DB,
) repository.TransactionRepositoryInterface {
	return r
}

func TestTransferRollback(t *testing.T) {

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_USER", "root")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "fintech_wallet")

	cfg := config.Load()

	db, err := database.Connect(cfg)

	if err != nil {
		t.Fatal(err)
	}

	walletRepository := repository.NewWalletRepository(db)
	// transactionRepository := repository.NewTransactionRepository(db)

	failingRepository := &FailingTransactionRepository{
		err: errors.New("failed to create transaction"),
	}

	transactionEntryRepository := repository.NewTransactionEntryRepository(db)
	idempotencyRepository := repository.NewIdempotencyRepository(db)

	transferService := NewTransferService(
		walletRepository,
		failingRepository,
		transactionEntryRepository,
		idempotencyRepository,
		nil,
		db,
	)

	// Reset balance
	err = db.Exec(`
		UPDATE wallets
		SET balance = CASE
			WHEN user_id = 2 THEN 100000
			WHEN user_id = 3 THEN 0
		END
		WHERE user_id IN (2, 3)
	`).Error

	if err != nil {
		t.Fatal(err)
	}

	// Count transaction sebelum transfer
	var transactionCountBefore int64

	err = db.Model(&model.Transaction{}).
		Where("user_id = ?", 2).
		Count(&transactionCountBefore).
		Error

	if err != nil {
		t.Fatal(err)
	}

	// Jalankan transfer yang sengaja dibuat gagal
	_, err = transferService.Transfer(
		2,
		3,
		"999",
		80000,
		"Rollback test",
	)

	if err == nil {
		t.Fatal("expected transfer to fail")
	}

	// Cek balance sender
	senderWallet, err := walletRepository.FindByUserID(2)

	if err != nil {
		t.Fatal(err)
	}

	// Cek balance receiver
	receiverWallet, err := walletRepository.FindByUserID(3)

	if err != nil {
		t.Fatal(err)
	}

	// Balance harus kembali seperti sebelum transaksi
	if senderWallet.Balance != 100000 {
		t.Errorf(
			"expected sender balance 100000, got %d",
			senderWallet.Balance,
		)
	}

	if receiverWallet.Balance != 0 {
		t.Errorf(
			"expected receiver balance 0, got %d",
			receiverWallet.Balance,
		)
	}

	// Pastikan transaction juga tidak tersimpan
	var transactionCountAfter int64

	err = db.Model(&model.Transaction{}).
		Where("user_id = ?", 2).
		Count(&transactionCountAfter).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if transactionCountAfter != transactionCountBefore {
		t.Errorf(
			"expected transaction count unchanged: %d, got %d",
			transactionCountBefore,
			transactionCountAfter,
		)
	}
}

func (r *FailingTransactionRepository) FindByID(
	id uint,
) (*model.Transaction, error) {
	return nil, r.err
}

func (r *FailingTransactionRepository) FindByUserID(
	userID uint,
	limit int,
	offset int,
) ([]model.Transaction, error) {
	return nil, r.err
}

func (r *FailingTransactionRepository) CountByUserID(
	userID uint,
) (int64, error) {
	return 0, r.err
}
