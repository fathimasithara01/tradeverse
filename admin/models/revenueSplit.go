package models

import "gorm.io/gorm"

type RevenueSplit struct {
	gorm.Model
	PaymentID   uint    `json:"payment_id"`
	UserID      uint    `json:"user_id"`
	TraderID    uint    `json:"trader_id"`
	AdminShare  float64 `json:"admin_share"`
	TraderShare float64 `json:"trader_share"`
	TotalAmount float64 `json:"total_amount"`
}
