package models

import (
	"gorm.io/gorm"
)

type Subscriber struct {
	gorm.Model
	TraderID   uint    `json:"trader_id"`
	UserID     uint    `json:"user_id"`
	Allocation float64 `json:"allocation"`
	Risk       string  `json:"risk"`
	Status     string  `json:"status"`
}
