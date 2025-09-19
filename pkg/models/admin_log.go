package models

import (
	"time"

	"gorm.io/gorm"
)

type AdminActionLog struct {
	gorm.Model
	AdminID    uint      `gorm:"not null;index" json:"admin_id"` // User ID of the admin
	Admin      User      `gorm:"foreignKey:AdminID" json:"admin"`
	Action     string    `gorm:"size:255;not null" json:"action"`    // e.g., "blocked user", "approved withdrawal"
	TargetType string    `gorm:"size:100" json:"target_type"`        // e.g., "User", "WithdrawRequest"
	TargetID   *uint     `json:"target_id,omitempty"`                // ID of the entity affected
	Details    string    `gorm:"type:text" json:"details,omitempty"` // JSON string or text for more details
	Timestamp  time.Time `gorm:"not null" json:"timestamp"`
	IPAddress  string    `gorm:"size:50" json:"ip_address,omitempty"`
}
