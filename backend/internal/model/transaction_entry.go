package model

import "time"

type TransactionEntry struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID uint      `gorm:"not null;index" json:"transaction_id"`
	WalletID      uint      `gorm:"not null;index" json:"wallet_id"`
	EntryType     string    `gorm:"size:20;not null" json:"entry_type"`
	Amount        int64     `gorm:"not null" json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}
