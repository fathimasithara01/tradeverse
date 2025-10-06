package models

import (
	"gorm.io/gorm"
)

type CustomerWallet struct {
	gorm.Model

	UserID   uint    `gorm:"uniqueIndex" json:"user_id"`
	Balance  float64 `gorm:"not null;default:0" json:"balance"`
	Currency string  `gorm:"size:10;default:'INR'" json:"currency"`
}
