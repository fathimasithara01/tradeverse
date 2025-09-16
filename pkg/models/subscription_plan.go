package models

import "gorm.io/gorm"

type SubscriptionPlan struct {
	gorm.Model
	Name        string  `gorm:"size:100;not null;unique" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:numeric(18,4);not null" json:"price"`
	Duration    int     `gorm:"not null" json:"duration"`         // e.g., 1, 3, 12
	Interval    string  `gorm:"size:20;not null" json:"interval"` // e.g., "day", "month", "year"
	IsActive    bool    `gorm:"default:true" json:"is_active"`

	IsTraderPlan     bool    `gorm:"default:true" json:"is_trader_plan"`                              // True if this plan upgrades a user to a Trader
	Features         string  `gorm:"type:text" json:"features"`                                       // JSON or comma-separated list of features
	MaxFollowers     int     `json:"max_followers,omitempty"`                                         // Max followers for traders on this plan
	CommissionRate   float64 `gorm:"type:numeric(5,4);default:0.00" json:"commission_rate,omitempty"` // Commission rate for traders
	AnalyticsAccess  string  `gorm:"size:50" json:"analytics_access,omitempty"`
	CreatedByAdminID uint    `json:"created_by_admin_id"`
}
