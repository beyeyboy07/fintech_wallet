package service

import (
	"testing"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"
)

func TestWebhookDuplicate(t *testing.T) {
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

	webhookRepository := repository.NewWebhookRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	walletRepository := repository.NewWalletRepository(db)

	service := NewWebhookService(
		db,
		webhookRepository,
		transactionRepository,
		walletRepository,
	)
	eventID := "test-webhook-duplicate"

	db.
		Where("event_id = ?", eventID).
		Delete(&model.WebhookEvent{})

	transaction := &model.Transaction{
		UserID:      1,
		Type:        "TOPUP",
		Amount:      50000,
		Status:      "PENDING",
		Reference:   "TEST-WEBHOOK-DUPLICATE",
		Description: "Webhook duplicate test",
	}

	if err := db.Create(transaction).Error; err != nil {
		t.Fatal(err)
	}

	request := model.PaymentWebhookRequest{
		EventID:       eventID,
		TransactionID: transaction.ID,
		Status:        "SUCCESS",
		Amount:        50000,
	}

	err = service.Receive(request)
	if err != nil {
		t.Fatal(err)
	}

	err = service.Receive(request)

	if err != ErrDuplicateWebhook {
		t.Fatalf(
			"expected ErrDuplicateWebhook, got %v",
			err,
		)
	}

	db.
		Where("event_id = ?", eventID).
		Delete(&model.WebhookEvent{})
	db.Delete(&model.Transaction{}, transaction.ID)
}

func TestWebhookTransactionNotFound(t *testing.T) {
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

	webhookRepository := repository.NewWebhookRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	walletRepository := repository.NewWalletRepository(db)

	service := NewWebhookService(
		db,
		webhookRepository,
		transactionRepository,
		walletRepository,
	)

	request := model.PaymentWebhookRequest{
		EventID:       "test-webhook-not-found",
		TransactionID: 999999,
		Status:        "SUCCESS",
		Amount:        50000,
	}

	err = service.Receive(request)

	if err != ErrTransactionNotFound {
		t.Fatalf(
			"expected ErrTransactionNotFound, got %v",
			err,
		)
	}
}

func TestWebhookInvalidTransactionType(t *testing.T) {
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

	webhookRepository := repository.NewWebhookRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	walletRepository := repository.NewWalletRepository(db)

	service := NewWebhookService(
		db,
		webhookRepository,
		transactionRepository,
		walletRepository,
	)

	transaction := &model.Transaction{
		UserID:      1,
		Type:        "TRANSFER",
		Amount:      50000,
		Status:      "PENDING",
		Reference:   "TEST-WEBHOOK-INVALID-TYPE",
		Description: "Webhook invalid type test",
	}

	if err := db.Create(transaction).Error; err != nil {
		t.Fatal(err)
	}

	defer db.Delete(&model.Transaction{}, transaction.ID)

	request := model.PaymentWebhookRequest{
		EventID:       "test-webhook-invalid-type",
		TransactionID: transaction.ID,
		Status:        "SUCCESS",
		Amount:        50000,
	}

	err = service.Receive(request)

	if err != ErrInvalidTransactionType {
		t.Fatalf(
			"expected ErrInvalidTransactionType, got %v",
			err,
		)
	}
}

func TestWebhookInvalidAmount(t *testing.T) {
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

	webhookRepository := repository.NewWebhookRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	walletRepository := repository.NewWalletRepository(db)

	service := NewWebhookService(
		db,
		webhookRepository,
		transactionRepository,
		walletRepository,
	)

	transaction := &model.Transaction{
		UserID:      1,
		Type:        "TOPUP",
		Amount:      50000,
		Status:      "PENDING",
		Reference:   "TEST-WEBHOOK-INVALID-AMOUNT",
		Description: "Webhook invalid amount test",
	}

	if err := db.Create(transaction).Error; err != nil {
		t.Fatal(err)
	}

	defer db.Delete(&model.Transaction{}, transaction.ID)

	request := model.PaymentWebhookRequest{
		EventID:       "test-webhook-invalid-amount",
		TransactionID: transaction.ID,
		Status:        "SUCCESS",
		Amount:        100000,
	}

	err = service.Receive(request)

	if err != ErrInvalidWebhookAmount {
		t.Fatalf(
			"expected ErrInvalidWebhookAmount, got %v",
			err,
		)
	}
}

func TestWebhookTransactionAlreadyProcessed(t *testing.T) {
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

	webhookRepository := repository.NewWebhookRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	walletRepository := repository.NewWalletRepository(db)

	service := NewWebhookService(
		db,
		webhookRepository,
		transactionRepository,
		walletRepository,
	)

	transaction := &model.Transaction{
		UserID:      1,
		Type:        "TOPUP",
		Amount:      50000,
		Status:      "SUCCESS",
		Reference:   "TEST-WEBHOOK-PROCESSED",
		Description: "Webhook processed test",
	}

	if err := db.Create(transaction).Error; err != nil {
		t.Fatal(err)
	}

	defer db.Delete(&model.Transaction{}, transaction.ID)

	request := model.PaymentWebhookRequest{
		EventID:       "test-webhook-processed",
		TransactionID: transaction.ID,
		Status:        "SUCCESS",
		Amount:        50000,
	}

	err = service.Receive(request)

	if err != ErrTransactionProcessed {
		t.Fatalf(
			"expected ErrTransactionProcessed, got %v",
			err,
		)
	}
}

func TestWebhookSuccessUpdatesWalletAndTransaction(t *testing.T) {
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

	webhookRepository := repository.NewWebhookRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	walletRepository := repository.NewWalletRepository(db)

	service := NewWebhookService(
		db,
		webhookRepository,
		transactionRepository,
		walletRepository,
	)

	userID := uint(9999)

	// Pastikan wallet test sebelumnya tidak bentrok.
	// db.
	// 	Where("user_id = ?", userID).
	// 	Delete(&model.Wallet{})

	wallet := &model.Wallet{
		UserID:  userID,
		Balance: 100000,
	}

	err = walletRepository.Create(wallet)
	if err != nil {
		t.Fatal(err)
	}

	transaction := &model.Transaction{
		UserID:      userID,
		Type:        "TOPUP",
		Amount:      50000,
		Status:      "PENDING",
		Reference:   "TOPUP-WEBHOOK-001",
		Description: "Webhook success test",
	}

	err = transactionRepository.Create(transaction)
	if err != nil {
		t.Fatal(err)
	}

	request := model.PaymentWebhookRequest{
		EventID:       "evt-success-001",
		TransactionID: transaction.ID,
		Status:        "SUCCESS",
		Amount:        50000,
	}

	err = service.Receive(request)
	if err != nil {
		t.Fatal(err)
	}

	updatedWallet, err := walletRepository.FindByUserID(userID)
	if err != nil {
		t.Fatal(err)
	}

	if updatedWallet.Balance != 150000 {
		t.Fatalf(
			"expected balance 150000, got %d",
			updatedWallet.Balance,
		)
	}

	updatedTransaction, err := transactionRepository.FindByID(transaction.ID)
	if err != nil {
		t.Fatal(err)
	}

	if updatedTransaction.Status != "SUCCESS" {
		t.Fatalf(
			"expected transaction status SUCCESS, got %s",
			updatedTransaction.Status,
		)
	}

	event, err := webhookRepository.FindByEventID(request.EventID)
	if err != nil {
		t.Fatal(err)
	}

	if event.Status != "PROCESSED" {
		t.Fatalf(
			"expected webhook status PROCESSED, got %s",
			event.Status,
		)
	}

	// Cleanup
	db.
		Where("event_id = ?", request.EventID).
		Delete(&model.WebhookEvent{})

	db.Delete(&model.Transaction{}, transaction.ID)
	db.Delete(&model.Wallet{}, wallet.ID)
}
