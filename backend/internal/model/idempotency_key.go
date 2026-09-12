package model

import "time"

type IdempotencyKey struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	Key           string    `gorm:"size:100;not null;uniqueIndex" json:"key"`
	RequestHash   string    `gorm:"size:64;not null" json:"request_hash"`
	TransactionID *uint     `gorm:"index" json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
