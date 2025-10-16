package models

import "gorm.io/gorm"

type TraderSignalSubscriptionPlan struct {
	gorm.Model

	TraderID        uint    `gorm:"index;not null" json:"trader_id"`
	Name            string  `gorm:"not null" json:"name"`
	Description     string  `gorm:"type:text" json:"description"`
	Price           float64 `gorm:"type:numeric(18,4);not null" json:"price"`
	Currency        string  `gorm:"size:10;not null;default:'USD'" json:"currency"`
	DurationDays    uint    `json:"duration_days"`
	IsActive        bool    `gorm:"default:true" json:"is_active"`
	AdminCommission float64 `gorm:"column:admin_commission;type:numeric(5,2);default:0.00" json:"admin_commission_percentage"`
	TraderShare     float64 `gorm:"type:numeric(18,4);not null;default:0.00" json:"trader_share"`

	Trader User `gorm:"foreignKey:TraderID"`
}
