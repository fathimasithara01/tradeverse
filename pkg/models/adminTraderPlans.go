package models

import (
	"time"

	"gorm.io/gorm"
)

type AdminTraderSubscriptionPlan struct {
	gorm.Model
	Name        string        `gorm:"size:100;not null;unique" json:"name"`
	Description string        `gorm:"type:text" json:"description"`
	Price       float64       `gorm:"type:numeric(18,4);not null" json:"price"`
	Currency    string        `gorm:"size:10;not null;default:'USD'" json:"currency"`
	Duration    time.Duration `gorm:"type:integer;not null" json:"duration"` // Duration (e.g., in seconds)
	Interval    string        `gorm:"size:20;not null" json:"interval"`      // e.g., "day", "month", "year"
	IsActive    bool          `gorm:"default:true" json:"is_active"`

	IsTraderPlan     bool    `gorm:"default:false" json:"is_trader_plan"` // True if this plan is offered by a trader (perhaps not needed if TraderSubscriptionPlan handles this)
	TraderID         *uint   `gorm:"index" json:"trader_id,omitempty"`    // Link to the Trader (User.ID) if IsTraderPlan is true (Can be removed if TraderSubscriptionPlan is used instead)
	Features         string  `gorm:"type:text" json:"features"`
	MaxFollowers     int     `json:"max_followers,omitempty"`
	CommissionRate   float64 `gorm:"type:numeric(5,4);default:0.10" json:"commission_rate,omitempty"` // Commission for the admin
	AnalyticsAccess  string  `gorm:"size:50" json:"analytics_access,omitempty"`
	CreatedByAdminID uint    `json:"created_by_admin_id"` // User ID of the admin who created it

	IsUpgradeToTrader bool `gorm:"default:false" json:"is_upgrade_to_trader"`
}
