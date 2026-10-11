package model

import "time"

type Wallet struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"-"`
	AccountNumber string    `gorm:"size:20;not null;uniqueIndex" json:"account_number"`
	Balance       int64     `gorm:"not null;default:0" json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
