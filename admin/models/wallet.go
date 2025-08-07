package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserID  uint    `json:"user_id"`
	Balance float64 `json:"balance"`
}

type WalletTransaction struct {
	gorm.Model
	UserID      uint    `json:"user_id"`
	Type        string  `json:"type"` // credit / debit
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}
