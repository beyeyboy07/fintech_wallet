package repository

import (
	"fintech-wallet/backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WebhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{
		db: db,
	}
}

func (r *WebhookRepository) WithTx(tx *gorm.DB) *WebhookRepository {
	return &WebhookRepository{
		db: tx,
	}
}

func (r *WebhookRepository) Create(
	event *model.WebhookEvent,
) error {
	return r.db.Create(event).Error
}

func (r *WebhookRepository) FindByEventID(
	eventID string,
) (*model.WebhookEvent, error) {

	var event model.WebhookEvent

	err := r.db.
		Where("event_id = ?", eventID).
		First(&event).
		Error

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *WebhookRepository) FindByEventIDForUpdate(
	eventID string,
) (*model.WebhookEvent, error) {

	var event model.WebhookEvent

	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("event_id = ?", eventID).
		First(&event).
		Error

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *WebhookRepository) UpdateStatus(
	id uint,
	status string,
) error {
	return r.db.
		Model(&model.WebhookEvent{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}
