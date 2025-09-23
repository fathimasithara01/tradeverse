package models

import (
	"time"

	"gorm.io/gorm"
)

type AdminActionLog struct {
	gorm.Model
	AdminID    uint      `gorm:"not null;index" json:"admin_id"`
	Admin      User      `gorm:"foreignKey:AdminID" json:"admin"`
	Action     string    `gorm:"size:255;not null" json:"action"`
	TargetType string    `gorm:"size:100" json:"target_type"`
	TargetID   *uint     `json:"target_id,omitempty"`
	Details    string    `gorm:"type:text" json:"details,omitempty"`
	Timestamp  time.Time `gorm:"not null" json:"timestamp"`
	IPAddress  string    `gorm:"size:50" json:"ip_address,omitempty"`
}
