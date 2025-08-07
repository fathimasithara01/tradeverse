package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	UserID        uint    `json:"user_id"`
	TraderID      uint    `json:"trader_id"`
	PlanID        uint    `json:"plan_id"`
	Amount        float64 `json:"amount"`
	OrderID       string  `json:"order_id"`
	PaymentID     string  `json:"payment_id"` // Razorpay payment ID
	PaymentMethod string  `json:"payment_method"`
	Status        string  `json:"status"` // success
	Method        string  `json:"method"` // razorpay, stripe
}
