package models

import "gorm.io/gorm"

type SubscriptionPlan struct {
	gorm.Model
	Name        string  `gorm:"size:100;not null;unique" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:numeric(18,4);not null" json:"price"`
	Duration    int     `gorm:"not null" json:"duration"`         // Duration in days/months
	Interval    string  `gorm:"size:20;not null" json:"interval"` // e.g., "day", "month", "year"
	IsActive    bool    `gorm:"default:true" json:"is_active"`

	Currency string `gorm:"size:10;not null;default:'USD'" json:"currency"`

	IsTraderPlan     bool    `gorm:"default:false" json:"is_trader_plan"` // True if this plan is offered by a trader
	TraderID         *uint   `gorm:"index" json:"trader_id,omitempty"`    // Link to the Trader (User.ID) if IsTraderPlan is true
	Features         string  `gorm:"type:text" json:"features"`
	MaxFollowers     int     `json:"max_followers,omitempty"`
	CommissionRate   float64 `gorm:"type:numeric(5,4);default:0.10" json:"commission_rate,omitempty"` // Commission for the admin (e.g., 0.10 for 10%)
	AnalyticsAccess  string  `gorm:"size:50" json:"analytics_access,omitempty"`
	CreatedByAdminID uint    `json:"created_by_admin_id"` // Who created this plan (Admin UserID)
}
