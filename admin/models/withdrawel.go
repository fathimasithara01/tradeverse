package models

import "gorm.io/gorm"

type Withdrawal struct {
	gorm.Model
	TraderID uint    `json:"trader_id"`
	Amount   float64 `json:"amount"`
	Status   string  `json:"status"` // pending, approved, rejected
	Note     string  `json:"note"`
}
