package models

import (
	"time"

	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	UserID    *uint     `json:"user_id"`    // Can be null (admin system logs)
	ActorRole string    `json:"actor_role"` // admin, user, trader
	Action    string    `json:"action"`     // e.g., "login", "create_signal", "ban_user"
	Details   string    `json:"details"`    // custom message
	Timestamp time.Time `json:"timestamp"`
}
