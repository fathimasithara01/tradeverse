package models

import (
	"time"

	"gorm.io/gorm"
)

// type TraderSubscription struct {
// 	gorm.Model
// 	UserID                   uint             `gorm:"not null" json:"user_id"`
// 	User                     User             `gorm:"foreignKey:UserID" json:"user"`
// 	TraderSubscriptionPlanID uint             `gorm:"not null" json:"trader_subscription_plan_id"`
// 	TraderSubscriptionPlan   SubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID" json:"trader_subscription_plan"`
// 	StartDate                time.Time        `gorm:"not null" json:"start_date"`
// 	EndDate                  time.Time        `gorm:"not null" json:"end_date"`
// 	IsActive                 bool             `gorm:"default:true" json:"is_active"`
// 	PaymentStatus            string           `gorm:"size:50;not null" json:"payment_status"` // e.g., "paid", "pending", "failed"
// 	AmountPaid               float64          `gorm:"not null" json:"amount_paid"`
// 	TransactionID            string           `gorm:"size:255" json:"transaction_id"` // For payment gateway reference
// 	TraderProfileID          *uint            `json:"trader_profile_id"`

// 	Allocation     float64 `gorm:"type:decimal(5,2);default:1.0" json:"allocation"`      // e.g., 1.0 for 100% of trader's trade size
// 	RiskMultiplier float64 `gorm:"type:decimal(5,2);default:1.0" json:"risk_multiplier"` // e.g., 1.0 for standard risk

// 	IsPaused       bool       `gorm:"default:false" json:"is_paused"`
// 	LastPauseDate  *time.Time `json:"last_pause_date"`
// 	LastResumeDate *time.Time `json:"last_resume_date"`
// }

type TraderSubscription struct {
	gorm.Model
	UserID                   uint             `gorm:"not null;index" json:"user_id"` // User who subscribed
	User                     User             `gorm:"foreignKey:UserID" json:"user"`
	TraderSubscriptionPlanID uint             `gorm:"not null;index" json:"trader_subscription_plan_id"` // The specific plan subscribed to
	TraderSubscriptionPlan   SubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID" json:"trader_subscription_plan"`
	StartDate                time.Time        `gorm:"not null" json:"start_date"`
	EndDate                  time.Time        `gorm:"not null" json:"end_date"`
	IsActive                 bool             `gorm:"default:true" json:"is_active"`          // Whether the subscription is currently active
	PaymentStatus            string           `gorm:"size:50;not null" json:"payment_status"` // e.g., "paid", "pending", "failed"
	AmountPaid               float64          `gorm:"type:numeric(18,4);not null" json:"amount_paid"`
	TransactionID            string           `gorm:"size:255" json:"transaction_id"` // Reference ID for the payment transaction
}
