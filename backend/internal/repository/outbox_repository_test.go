package repository

import (
	"testing"
	"time"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/model"
)

func TestResetStaleProcessing(t *testing.T) {
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

	repo := NewOutboxRepository(db)

	event := &model.OutboxEvent{
		EventType:   "test:stale-recovery",
		AggregateID: 999999,
		Payload:     `{"test":true}`,
		Status:      "PROCESSING",
	}

	err = repo.Create(event)
	if err != nil {
		t.Fatal(err)
	}

	// Simulasikan event yang sudah stale.
	oldTime := time.Now().Add(-20 * time.Minute)

	err = db.
		Model(&model.OutboxEvent{}).
		Where("id = ?", event.ID).
		Update("updated_at", oldTime).
		Error

	if err != nil {
		t.Fatal(err)
	}

	err = repo.ResetStaleProcessing(
		10 * time.Minute,
	)

	if err != nil {
		t.Fatal(err)
	}

	var updatedEvent model.OutboxEvent

	err = db.
		First(&updatedEvent, event.ID).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedEvent.Status != "PENDING" {
		t.Fatalf(
			"expected status PENDING, got %s",
			updatedEvent.Status,
		)
	}

	// Cleanup
	db.Delete(&model.OutboxEvent{}, event.ID)
}
