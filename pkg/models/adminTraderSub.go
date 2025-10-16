package models

import (
	"time"

	"gorm.io/gorm"
)

// this is admin created customer to trader upgrade subscription

type CustomerToTraderSub struct {
	gorm.Model
	UserID             uint                        `gorm:"not null;index" json:"user_id"`
	User               User                        `gorm:"foreignKey:UserID" json:"user"`
	SubscriptionPlanID uint                        `gorm:"not null;index" json:"subscription_plan_id"`
	SubscriptionPlan   AdminTraderSubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID" json:"subscription_plan"`

	TraderID *uint `gorm:"index" json:"trader_id,omitempty"`

	StartDate     time.Time `gorm:"not null" json:"start_date"`
	EndDate       time.Time `gorm:"not null" json:"end_date"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	PaymentStatus string    `gorm:"size:50;not null" json:"payment_status"`
	AmountPaid    float64   `gorm:"type:numeric(18,4);not null" json:"amount_paid"`
	TransactionID string    `gorm:"size:255" json:"transaction_id"`
}
