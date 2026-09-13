package model

import "time"

type WebhookEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventID   string    `gorm:"size:100;uniqueIndex;not null" json:"event_id"`
	EventType string    `gorm:"size:100;not null" json:"event_type"`
	Status    string    `gorm:"size:30;not null" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
