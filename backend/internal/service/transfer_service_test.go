package service

import (
	"sync"
	"testing"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/repository"
)

func TestTransferConcurrent(t *testing.T) {

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

	for i := 1; i <= 20; i++ {

		t.Log("Iteration:", i)

		// Reset saldo
		err := db.Exec(`
			UPDATE wallets
			SET balance = CASE
				WHEN user_id = 2 THEN 100000
				WHEN user_id = 3 THEN 0
				WHEN user_id = 4 THEN 0
			END
			WHERE user_id IN (2, 3, 4)
		`).Error

		if err != nil {
			t.Fatal(err)
		}

		var wg sync.WaitGroup

		results := make(chan error, 2)

		wg.Add(2)

		go func() {
			defer wg.Done()

			_, err := transferService.Transfer(
				2,
				3,
				"999",
				80000,
				"Concurrent transfer 1",
			)

			results <- err
		}()

		go func() {
			defer wg.Done()

			_, err := transferService.Transfer(
				2,
				4,
				"999",
				80000,
				"Concurrent transfer 2",
			)

			results <- err
		}()

		wg.Wait()

		close(results)

		successCount := 0
		failedCount := 0

		for err := range results {

			if err == nil {
				successCount++
			} else {
				failedCount++
			}
		}

		t.Log(
			"Success:",
			successCount,
			"Failed:",
			failedCount,
		)

		if successCount != 1 {
			t.Errorf(
				"iteration %d: expected 1 success, got %d",
				i,
				successCount,
			)
		}

		if failedCount != 1 {
			t.Errorf(
				"iteration %d: expected 1 failed transfer, got %d",
				i,
				failedCount,
			)
		}

		senderWallet, err := walletRepository.FindByUserID(2)
		if err != nil {
			t.Fatal(err)
		}

		receiver3Wallet, err := walletRepository.FindByUserID(3)
		if err != nil {
			t.Fatal(err)
		}

		receiver4Wallet, err := walletRepository.FindByUserID(4)
		if err != nil {
			t.Fatal(err)
		}

		if senderWallet.Balance != 20000 {
			t.Errorf(
				"iteration %d: expected sender balance 20000, got %d",
				i,
				senderWallet.Balance,
			)
		}

		if receiver3Wallet.Balance+receiver4Wallet.Balance != 80000 {
			t.Errorf(
				"iteration %d: expected receiver total balance 80000, got %d",
				i,
				receiver3Wallet.Balance+receiver4Wallet.Balance,
			)
		}
	}
}
