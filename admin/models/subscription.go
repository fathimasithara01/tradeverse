package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	UserID   uint      `json:"user_id"`
	TraderID uint      `json:"trader_id"`
	PlanID   uint      `json:"plan_id"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
	Amount   float64   `json:"amount"`
	Status   string    `json:"status"` // active, expired, cancelled
}
