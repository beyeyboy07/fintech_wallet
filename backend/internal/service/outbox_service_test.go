package service

import (
	"testing"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/redis"
	"fintech-wallet/backend/internal/repository"
	"fintech-wallet/backend/internal/worker"
)

func TestOutboxDispatch(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_USER", "root")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "fintech_wallet")
	t.Setenv("REDIS_HOST", "localhost")
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "0")

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatal(err)
	}

	redisClient := redis.NewClient(cfg)

	if err := redis.Ping(redisClient); err != nil {
		t.Fatal(err)
	}

	queueClient := worker.NewClient(
		cfg.RedisHost + ":" + cfg.RedisPort,
	)
	defer queueClient.Close()

	outboxRepository :=
		repository.NewOutboxRepository(db)

	defer db.
		Where("event_type = ?", "test:outbox-dispatch").
		Delete(&model.OutboxEvent{})

	outboxService := NewOutboxService(
		outboxRepository,
		queueClient,
	)

	err = db.
		Where("event_type = ?", "test:outbox-dispatch").
		Delete(&model.OutboxEvent{}).
		Error

	if err != nil {
		t.Fatal(err)
	}

	event := &model.OutboxEvent{
		EventType:   "test:outbox-dispatch",
		AggregateID: 999999,
		Payload:     `{"test":true}`,
		Status:      "PENDING",
	}

	err = outboxRepository.Create(event)
	if err != nil {
		t.Fatal(err)
	}

	err = outboxService.Dispatch(10)
	if err != nil {
		t.Fatal(err)
	}

	var updatedEvent model.OutboxEvent

	err = db.First(
		&updatedEvent,
		event.ID,
	).Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedEvent.Status != "PROCESSED" {
		t.Fatalf(
			"expected status PROCESSED, got %s",
			updatedEvent.Status,
		)
	}

}
