package models

import (
	"time"

	"gorm.io/gorm"
)

type SubscriptionPlan struct {
	gorm.Model
	Name            string  `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Description     string  `gorm:"size:255" json:"description"`
	Price           float64 `gorm:"not null" json:"price"`
	Duration        int     `gorm:"not null" json:"duration"`         // Duration in days, months, etc.
	Interval        string  `gorm:"size:50;not null" json:"interval"` // e.g., "monthly", "yearly", "days"
	IsActive        bool    `gorm:"default:true" json:"is_active"`
	Features        string  `gorm:"type:text" json:"features"` // JSON string or comma-separated list of features
	MaxFollowers    int     `json:"max_followers"`
	CommissionRate  float64 `gorm:"type:decimal(5,4)" json:"commission_rate"` // e.g., 0.05 for 5%
	AnalyticsAccess string  `gorm:"size:100" json:"analytics_access"`         // e.g., "basic", "premium"
}

type TraderSubscription struct {
	gorm.Model
	UserID                   uint                   `gorm:"not null" json:"user_id"`
	User                     User                   `gorm:"foreignKey:UserID" json:"user"`
	TraderSubscriptionPlanID uint                   `gorm:"not null" json:"trader_subscription_plan_id"`
	TraderSubscriptionPlan   TraderSubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID" json:"trader_subscription_plan"`
	StartDate                time.Time              `gorm:"not null" json:"start_date"`
	EndDate                  time.Time              `gorm:"not null" json:"end_date"`
	IsActive                 bool                   `gorm:"default:true" json:"is_active"`
	PaymentStatus            string                 `gorm:"size:50;not null" json:"payment_status"` // e.g., "paid", "pending", "failed"
	AmountPaid               float64                `gorm:"not null" json:"amount_paid"`
	TransactionID            string                 `gorm:"size:255" json:"transaction_id"` // For payment gateway reference
	TraderProfileID          *uint                  `json:"trader_profile_id"`
}
