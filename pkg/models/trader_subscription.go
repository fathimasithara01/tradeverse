package models

import (
	"time"

	"gorm.io/gorm"
)

type TraderSubscription struct {
	gorm.Model
	UserID                   uint              `gorm:"not null;index" json:"user_id"`     // This is the Customer's UserID
	User                     User              `gorm:"foreignKey:UserID" json:"user"`     // The Customer
	TraderID                 uint              `gorm:"not null;index" json:"trader_id"`   // The Trader's UserID
	Trader                   User              `gorm:"foreignKey:TraderID" json:"trader"` // The Trader
	TraderSubscriptionPlanID uint              `gorm:"not null;index" json:"trader_subscription_plan_id"`
	TraderSubscriptionPlan   *SubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID" json:"trader_subscription_plan"`

	StartDate     time.Time `gorm:"not null" json:"start_date"`
	EndDate       time.Time `gorm:"not null" json:"end_date"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	PaymentStatus string    `gorm:"size:50;not null" json:"payment_status"`
	AmountPaid    float64   `gorm:"type:numeric(18,4);not null" json:"amount_paid"`
	TransactionID string    `gorm:"size:255" json:"transaction_id"` // Reference to the wallet transaction
}
