package models

import (
	"time"

	"gorm.io/gorm"
)

type TraderSignalSubscriptionPlan struct {
	gorm.Model
	TraderID        uint      `gorm:"index;not null" json:"trader_id"`
	Trader          User      `gorm:"foreignKey:TraderID"`
	Name            string    `gorm:"size:255;not null" json:"name"`
	Description     string    `gorm:"type:text" json:"description"`
	Price           float64   `gorm:"type:numeric(18,4);not null" json:"price"`
	Currency        string    `gorm:"size:10;not null" json:"currency"`
	DurationDays    uint      `gorm:"not null" json:"duration_days"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	AdminCommission float64   `gorm:"type:numeric(5,2);not null;default:0.0" json:"admin_commission_percentage"` 
	TraderShare     float64   `gorm:"type:numeric(18,4);not null;default:0.0" json:"trader_share"`              
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
