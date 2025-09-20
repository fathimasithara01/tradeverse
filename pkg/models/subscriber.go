package models

import "time"

type Subscriber struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TraderID   uint      `json:"trader_id"`
	UserID     uint      `json:"user_id"`
	Allocation float64   `json:"allocation"`
	Risk       string    `json:"risk"`   // e.g., "low", "medium", "high"
	Status     string    `json:"status"` // e.g., "active", "inactive"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
