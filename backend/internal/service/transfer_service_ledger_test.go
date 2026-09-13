package service

import (
	"testing"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"
)

func TestTransferCreatesLedgerEntries(t *testing.T) {

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
	transactionRepository := repository.NewTransactionRepository(db)
	transactionEntryRepository := repository.NewTransactionEntryRepository(db)
	idempotencyRepository := repository.NewIdempotencyRepository(db)

	transferService := NewTransferService(
		walletRepository,
		transactionRepository,
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

	// Transfer
	transaction, err := transferService.Transfer(
		2,
		3,
		80000,
		"Ledger test",
	)

	if err != nil {
		t.Fatal(err)
	}

	// Ambil semua ledger entry untuk transaction tersebut
	var entries []model.TransactionEntry

	err = db.
		Where("transaction_id = ?", transaction.ID).
		Order("id ASC").
		Find(&entries).
		Error

	if err != nil {
		t.Fatal(err)
	}

	// Harus tepat 2 entry
	if len(entries) != 2 {
		t.Fatalf(
			"expected 2 ledger entries, got %d",
			len(entries),
		)
	}

	// Entry pertama harus DEBIT sender
	if entries[0].WalletID != 2 {
		t.Errorf(
			"expected first entry wallet_id 2, got %d",
			entries[0].WalletID,
		)
	}

	if entries[0].EntryType != "DEBIT" {
		t.Errorf(
			"expected first entry type DEBIT, got %s",
			entries[0].EntryType,
		)
	}

	if entries[0].Amount != 80000 {
		t.Errorf(
			"expected first entry amount 80000, got %d",
			entries[0].Amount,
		)
	}

	// Entry kedua harus CREDIT receiver
	if entries[1].WalletID != 3 {
		t.Errorf(
			"expected second entry wallet_id 3, got %d",
			entries[1].WalletID,
		)
	}

	if entries[1].EntryType != "CREDIT" {
		t.Errorf(
			"expected second entry type CREDIT, got %s",
			entries[1].EntryType,
		)
	}

	if entries[1].Amount != 80000 {
		t.Errorf(
			"expected second entry amount 80000, got %d",
			entries[1].Amount,
		)
	}

	// Validasi double-entry:
	// total debit harus sama dengan total credit
	var debitTotal int64
	var creditTotal int64

	err = db.Model(&model.TransactionEntry{}).
		Where(
			"transaction_id = ? AND entry_type = ?",
			transaction.ID,
			"DEBIT",
		).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&debitTotal).
		Error

	if err != nil {
		t.Fatal(err)
	}

	err = db.Model(&model.TransactionEntry{}).
		Where(
			"transaction_id = ? AND entry_type = ?",
			transaction.ID,
			"CREDIT",
		).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&creditTotal).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if debitTotal != creditTotal {
		t.Errorf(
			"ledger is unbalanced: debit=%d credit=%d",
			debitTotal,
			creditTotal,
		)
	}

	if debitTotal != 80000 {
		t.Errorf(
			"expected debit total 80000, got %d",
			debitTotal,
		)
	}

	if creditTotal != 80000 {
		t.Errorf(
			"expected credit total 80000, got %d",
			creditTotal,
		)
	}
}
