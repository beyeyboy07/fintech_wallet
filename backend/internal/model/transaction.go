package model

import "time"

type Transaction struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	SenderID      *uint     `gorm:"index" json:"sender_id"`
	ReceiverID    *uint     `gorm:"index" json:"-"`
	AccountNumber string    `gorm:"index" json:"accountnumber"`
	Type          string    `gorm:"size:50;not null" json:"type"`
	Amount        int64     `gorm:"not null" json:"amount"`
	Status        string    `gorm:"size:30;not null" json:"status"`
	Reference     string    `gorm:"size:100;uniqueIndex;not null" json:"reference"`
	Description   string    `gorm:"size:255" json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

const (
	MaxTransferAmount int64 = 100_000_000
)
