package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomerWallet struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex" json:"user_id"`
	Balance   float64        `gorm:"not null;default:0" json:"balance"`
	Currency  string         `gorm:"size:10;default:'INR'" json:"currency"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
