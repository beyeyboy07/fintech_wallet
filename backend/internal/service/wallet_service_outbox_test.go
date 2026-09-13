package service

import (
	"errors"
	"testing"
	"time"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"

	"gorm.io/gorm"
)

type FailingOutboxRepository struct {
	captured *model.OutboxEvent
}

func (r *FailingOutboxRepository) Create(
	event *model.OutboxEvent,
) error {
	r.captured = event

	return errors.New("forced outbox error")
}

func (r *FailingOutboxRepository) FindPending(
	limit int,
) ([]model.OutboxEvent, error) {
	return nil, nil
}

func (r *FailingOutboxRepository) FindByIDForUpdate(
	id uint,
) (*model.OutboxEvent, error) {
	return nil, errors.New("not implemented")
}

func (r *FailingOutboxRepository) UpdateStatus(
	id uint,
	status string,
) error {
	return nil
}

func (r *FailingOutboxRepository) WithTx(
	tx *gorm.DB,
) repository.OutboxRepositoryInterface {
	return r
}

func TestTopUpAsyncRollbackWhenOutboxFails(t *testing.T) {
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

	walletRepository :=
		repository.NewWalletRepository(db)

	transactionRepository :=
		repository.NewTransactionRepository(db)

	failingOutboxRepository :=
		&FailingOutboxRepository{}

	walletService := NewWalletService(
		walletRepository,
		transactionRepository,
		failingOutboxRepository,
		nil,
		nil,
		db,
	)

	userID := uint(1)
	amount := int64(50000)

	_, err = walletRepository.FindByUserID(userID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = walletService.TopUpAsync(
		userID,
		amount,
	)

	if err == nil {
		t.Fatal("expected outbox error")
	}

	if failingOutboxRepository.captured == nil {
		t.Fatal("expected outbox event to be captured")
	}

	var transaction model.Transaction

	err = db.First(
		&transaction,
		failingOutboxRepository.captured.AggregateID,
	).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf(
			"expected transaction rollback, got error: %v",
			err,
		)
	}
}

func (r *FailingOutboxRepository) ClaimPending(
	id uint,
) (*model.OutboxEvent, error) {
	return nil, errors.New("not implemented")
}

func (r *FailingOutboxRepository) ResetStaleProcessing(
	olderThan time.Duration,
) error {
	return nil
}
