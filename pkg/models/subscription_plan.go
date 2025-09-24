package models

import "gorm.io/gorm"

type SubscriptionPlan struct {
	gorm.Model
	Name        string  `gorm:"size:100;not null;unique" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:numeric(18,4);not null" json:"price"`
	Duration    int     `gorm:"not null" json:"duration"`
	Interval    string  `gorm:"size:20;not null" json:"interval"`
	IsActive    bool    `gorm:"default:true" json:"is_active"`

	Currency string `gorm:"size:10;not null;default:'USD'" json:"currency"`

	IsTraderPlan     bool    `gorm:"default:false" json:"is_trader_plan"`
	Features         string  `gorm:"type:text" json:"features"`
	MaxFollowers     int     `json:"max_followers,omitempty"`
	CommissionRate   float64 `gorm:"type:numeric(5,4);default:0.00" json:"commission_rate,omitempty"`
	AnalyticsAccess  string  `gorm:"size:50" json:"analytics_access,omitempty"`
	CreatedByAdminID uint    `json:"created_by_admin_id"`
}
