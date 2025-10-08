package models

import (
	"time"

	"gorm.io/gorm"
)

// TraderSubscription represents a subscription of a customer to a trader's signals
type TraderSubscription struct {
	gorm.Model

	// Customer
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	// Trader
	TraderID uint `gorm:"not null;index" json:"trader_id"`
	Trader   User `gorm:"foreignKey:TraderID" json:"trader"`

	// Subscription Plan
	TraderSubscriptionPlanID uint              `gorm:"not null;index" json:"trader_subscription_plan_id"`
	TraderSubscriptionPlan   *SubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID" json:"trader_subscription_plan"`

	// Duration & Status
	StartDate     time.Time `gorm:"not null" json:"start_date"`
	EndDate       time.Time `gorm:"not null" json:"end_date"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	PaymentStatus string    `gorm:"size:50;not null" json:"payment_status"`

	// Financial Details
	AmountPaid      float64 `gorm:"type:numeric(18,4);not null;default:0" json:"amount_paid"`
	TraderShare     float64 `gorm:"type:numeric(18,4);not null;default:0" json:"trader_share"`
	AdminCommission float64 `gorm:"type:numeric(18,4);not null;default:0" json:"admin_commission"`
	TransactionID   string  `gorm:"size:255;not null" json:"transaction_id"`
}

// Request payload
type TraderSubscriptionRequest struct {
	CustomerID               uint `json:"customer_id" binding:"required"`
	TraderID                 uint `json:"trader_id" binding:"required"`
	TraderSubscriptionPlanID uint `json:"trader_subscription_plan_id" binding:"required"`
}

// Response payload
type TraderSubscriptionResponse struct {
	TraderName      string  `json:"trader_name"`
	PlanName        string  `json:"plan_name"`
	AmountPaid      float64 `json:"amount_paid"`
	TraderShare     float64 `json:"trader_share"`
	AdminCommission float64 `json:"admin_commission"`
	PaymentStatus   string  `json:"payment_status"`
	TransactionID   string  `json:"transaction_id"`
	StartDate       string  `json:"start_date"`
	EndDate         string  `json:"end_date"`
	IsActive        bool    `json:"is_active"`
}
