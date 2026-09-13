package service

import (
	"sync"
	"testing"
	"time"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/redis"
	"fintech-wallet/backend/internal/repository"
)

func TestProcessTopUpDuplicate(t *testing.T) {
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
	outboxRepository := repository.NewOutboxRepository(db)
	redisClient := redis.NewClient(cfg)
	redisService := NewRedisService(redisClient)

	walletService := NewWalletService(
		walletRepository,
		transactionRepository,
		outboxRepository,
		redisService,
		nil,
		db,
	)

	userID := uint(1)

	// Reset saldo
	err = db.Model(&model.Wallet{}).
		Where("user_id = ?", userID).
		Update("balance", int64(100000)).
		Error

	if err != nil {
		t.Fatal(err)
	}

	// Buat satu transaksi top up
	transaction := &model.Transaction{
		UserID:    userID,
		Type:      "TOPUP",
		Amount:    50000,
		Status:    "PENDING",
		Reference: "TOPUP-DUPLICATE-" + time.Now().Format("20060102150405.000000"),
	}

	err = transactionRepository.Create(transaction)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	// Simulasikan 2 worker memproses transaction yang sama
	for i := 0; i < 2; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			errCh <- walletService.ProcessTopUp(
				transaction.ID,
				userID,
				50000,
			)
		}()
	}

	wg.Wait()
	close(errCh)

	// Keduanya harus berhasil secara teknis.
	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// Cek saldo akhir
	wallet, err := walletRepository.FindByUserID(userID)
	if err != nil {
		t.Fatal(err)
	}

	if wallet.Balance != 150000 {
		t.Fatalf(
			"expected balance 150000, got %d",
			wallet.Balance,
		)
	}

	// Cek status transaction
	transaction, err = transactionRepository.FindByID(transaction.ID)
	if err != nil {
		t.Fatal(err)
	}

	if transaction.Status != "SUCCESS" {
		t.Fatalf(
			"expected status SUCCESS, got %s",
			transaction.Status,
		)
	}
}
