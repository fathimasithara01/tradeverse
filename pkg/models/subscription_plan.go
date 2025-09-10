package models

import "gorm.io/gorm"

type SubscriptionPlan struct {
	gorm.Model
	Name         string  `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Description  string  `gorm:"size:255" json:"description"`
	Price        float64 `gorm:"not null" json:"price"`
	Duration     int     `gorm:"not null" json:"duration"`         // Duration in days, months, etc.
	Interval     string  `gorm:"size:50;not null" json:"interval"` // e.g., "monthly", "yearly", "days"
	IsActive     bool    `gorm:"default:true" json:"is_active"`
	Features     string  `gorm:"type:text" json:"features"` // JSON string or comma-separated list of features
	MaxFollowers int     `json:"max_followers"`
	Status       string  `json:"status"`
}
