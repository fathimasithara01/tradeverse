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
	Duration    time.Duration `gorm:"type:integer;not null" json:"duration"`
	Interval    string        `gorm:"size:20;not null" json:"interval"`
	IsActive    bool          `gorm:"default:true" json:"is_active"`

	IsTraderPlan     bool    `gorm:"default:false" json:"is_trader_plan"`
	TraderID         *uint   `gorm:"index" json:"trader_id,omitempty"` 
	Features         string  `gorm:"type:text" json:"features"`
	MaxFollowers     int     `json:"max_followers,omitempty"`
	CommissionRate   float64 `gorm:"type:numeric(5,4);default:0.10" json:"commission_rate,omitempty"` 
	AnalyticsAccess  string  `gorm:"size:50" json:"analytics_access,omitempty"`
	CreatedByAdminID uint    `json:"created_by_admin_id"` 

	IsUpgradeToTrader bool `gorm:"default:false" json:"is_upgrade_to_trader"`
}
