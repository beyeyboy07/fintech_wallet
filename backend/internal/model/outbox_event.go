package model

import "time"

type OutboxEvent struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	EventType   string     `gorm:"size:100;not null" json:"event_type"`
	AggregateID uint       `gorm:"not null;index" json:"aggregate_id"`
	Payload     string     `gorm:"type:json;not null" json:"payload"`
	Status      string     `gorm:"size:30;not null;index" json:"status"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
