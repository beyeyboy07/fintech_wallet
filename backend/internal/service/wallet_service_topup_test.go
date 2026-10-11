package service

import (
	"sync"
	"testing"
	"time"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestProcessTopUpConcurrent(t *testing.T) {

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "3307")
	t.Setenv("DB_USER", "root")
	t.Setenv("DB_PASSWORD", "root")
	t.Setenv("DB_NAME", "fintech_wallet_test")

	cfg := config.Load()

	db, err := database.Connect(cfg)

	if err != nil {
		t.Fatal(err)
	}

	// Siapkan schema database test
	err = database.Migrate(db)
	require.NoError(t, err)

	// Siapkan user test
	user := &model.User{
		Name:     "Test User",
		Email:    "topup-concurrent-test@example.com",
		Password: "test-password",
		Role:     model.RoleUser,
	}

	err = db.Where("email = ?", user.Email).First(user).Error

	if err != nil {
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)

		err = db.Create(user).Error
		require.NoError(t, err)
	}

	// Siapkan wallet test
	wallet := &model.Wallet{
		UserID:        user.ID,
		AccountNumber: "8801000000000001",
		Balance:       100000,
	}

	err = db.Where("user_id = ?", user.ID).First(wallet).Error

	if err != nil {
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)

		err = db.Create(wallet).Error
		require.NoError(t, err)
	}

	walletRepository := repository.NewWalletRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	outboxRepository := repository.NewOutboxRepository(db)
	walletService := NewWalletService(
		walletRepository,
		transactionRepository,
		outboxRepository,
		nil,
		nil,
		db,
	)

	// Buat wallet test
	// Ambil wallet milik user test
	wallet, err = walletRepository.FindByUserID(user.ID)
	require.NoError(t, err)

	wallet.Balance = 100000

	err = db.Save(wallet).Error
	require.NoError(t, err)

	// Buat transaction PENDING
	transaction := &model.Transaction{
		UserID:    user.ID,
		Type:      "TOPUP",
		Amount:    50000,
		Status:    "PENDING",
		Reference: "TOPUP-CONCURRENT-TEST-" + time.Now().Format("20060102150405.000000"),
	}

	err = transactionRepository.Create(transaction)
	require.NoError(t, err)

	var wg sync.WaitGroup

	results := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()

		err := walletService.ProcessTopUp(
			transaction.ID,
			user.ID,
			50000,
		)

		results <- err
	}()

	go func() {
		defer wg.Done()

		err := walletService.ProcessTopUp(
			transaction.ID,
			user.ID,
			50000,
		)

		results <- err
	}()

	wg.Wait()

	close(results)

	for err := range results {
		assert.NoError(t, err)
	}

	// Cek saldo
	result, err := walletRepository.FindByUserID(user.ID)
	require.NoError(t, err)

	assert.Equal(
		t,
		int64(150000),
		result.Balance,
	)

	// Cek status transaction
	transactionResult, err :=
		transactionRepository.FindByID(transaction.ID)

	require.NoError(t, err)

	assert.Equal(
		t,
		"SUCCESS",
		transactionResult.Status,
	)
}
