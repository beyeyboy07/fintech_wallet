package repository

import (
	"errors"
	"fintech-wallet/backend/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OutboxRepositoryInterface interface {
	Create(event *model.OutboxEvent) error
	FindPending(limit int) ([]model.OutboxEvent, error)
	FindByIDForUpdate(id uint) (*model.OutboxEvent, error)
	ClaimPending(id uint) (*model.OutboxEvent, error)
	UpdateStatus(id uint, status string) error
	ResetStaleProcessing(olderThan time.Duration) error
	WithTx(tx *gorm.DB) OutboxRepositoryInterface
}

type OutboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) *OutboxRepository {
	return &OutboxRepository{
		db: db,
	}
}

func (r *OutboxRepository) WithTx(tx *gorm.DB) OutboxRepositoryInterface {
	return &OutboxRepository{
		db: tx,
	}
}

func (r *OutboxRepository) Create(
	event *model.OutboxEvent,
) error {
	return r.db.Create(event).Error
}

func (r *OutboxRepository) FindPending(
	limit int,
) ([]model.OutboxEvent, error) {

	var events []model.OutboxEvent

	err := r.db.
		Where("status = ?", "PENDING").
		Order("id ASC").
		Limit(limit).
		Find(&events).
		Error

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *OutboxRepository) FindByIDForUpdate(
	id uint,
) (*model.OutboxEvent, error) {

	var event model.OutboxEvent

	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&event, id).
		Error

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *OutboxRepository) UpdateStatus(
	id uint,
	status string,
) error {

	result := r.db.
		Model(&model.OutboxEvent{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("outbox event not found when updating status")
	}

	return nil
}

func (r *OutboxRepository) ClaimPending(
	id uint,
) (*model.OutboxEvent, error) {

	var event model.OutboxEvent

	err := r.db.Transaction(func(tx *gorm.DB) error {

		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&event, id).
			Error

		if err != nil {
			return err
		}

		if event.Status != "PENDING" {
			return nil
		}

		err = tx.
			Model(&model.OutboxEvent{}).
			Where(
				"id = ? AND status = ?",
				id,
				"PENDING",
			).
			Update(
				"status",
				"PROCESSING",
			).
			Error

		if err != nil {
			return err
		}

		event.Status = "PROCESSING"

		return nil
	})

	if err != nil {
		return nil, err
	}

	if event.Status != "PROCESSING" {
		return nil, nil
	}

	return &event, nil
}

func (r *OutboxRepository) ResetStaleProcessing(
	olderThan time.Duration,
) error {

	cutoff := time.Now().Add(-olderThan)

	return r.db.
		Model(&model.OutboxEvent{}).
		Where(
			"status = ? AND updated_at < ?",
			"PROCESSING",
			cutoff,
		).
		Update(
			"status",
			"PENDING",
		).
		Error
}
