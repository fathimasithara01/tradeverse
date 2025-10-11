package models

import (
	"time"

	"gorm.io/gorm"
)

type UserSubscription struct {
	gorm.Model
	UserID             uint             `gorm:"index;not null" json:"user_id"`
	User               User             `gorm:"foreignKey:UserID"`
	SubscriptionPlanID uint             `gorm:"index;not null" json:"subscription_plan_id"`
	Plan               SubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID"`
	StartDate          time.Time        `gorm:"not null" json:"start_date"`
	EndDate            time.Time        `gorm:"not null" json:"end_date"`
	IsActive           bool             `gorm:"default:true" json:"is_active"`
	TransactionID      uint             `gorm:"index" json:"transaction_id"`
}

type TraderSubscriptionPlan struct {
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

type CustomerTraderSubscription struct {
	gorm.Model
	CustomerID uint `gorm:"index;not null" json:"customer_id"`
	Customer   User `gorm:"foreignKey:CustomerID"`

	TraderID uint `gorm:"index;not null" json:"trader_id"`
	Trader   User `gorm:"foreignKey:TraderID"`

	TraderSubscriptionPlanID uint                   `gorm:"index;not null" json:"trader_subscription_plan_id"`
	Plan                     TraderSubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID"`

	StartDate           time.Time `gorm:"not null" json:"start_date"`
	EndDate             time.Time `gorm:"not null" json:"end_date"`
	IsActive            bool      `gorm:"default:true" json:"is_active"`
	WalletTransactionID *uint     `gorm:"index" json:"wallet_transaction_id"`
	TransactionID       uint

	PaymentStatus          string  `gorm:"size:50;not null" json:"payment_status"`
	AmountPaid             float64 `gorm:"type:numeric(18,4);not null" json:"amount_paid"`
	TraderShare            float64 `gorm:"type:numeric(18,4);not null" json:"trader_share"`
	AdminCommission        float64 `gorm:"type:numeric(18,4);not null" json:"admin_commission"`
	TransactionReferenceID string  `gorm:"size:255;not null" json:"transaction_reference_id"`
}

type TraderSubscriptionResponse struct {
	TraderSubscriptionID uint    `json:"trader_subscription_id"`
	TraderName           string  `json:"trader_name"`
	PlanName             string  `json:"plan_name"`
	AmountPaid           float64 `json:"amount_paid"`
	TraderShare          float64 `json:"trader_share"`
	AdminCommission      float64 `json:"admin_commission"`
	PaymentStatus        string  `json:"payment_status"`
	TransactionID        string  `json:"transaction_id"`
	StartDate            string  `json:"start_date"`
	EndDate              string  `json:"end_date"`
	IsActive             bool    `json:"is_active"`
	Message              string  `json:"message"`
	Status               string  `json:"status"`
}

// Request
type TraderSubscriptionRequest struct {
	CustomerID               uint `json:"customer_id" binding:"required"`
	TraderID                 uint `json:"trader_id" binding:"required"`
	TraderSubscriptionPlanID uint `json:"trader_subscription_plan_id" binding:"required"`
}

type CreateTraderSubscriptionPlanInput struct {
	Name                      string  `json:"name" binding:"required"`
	Description               string  `json:"description"`
	Price                     float64 `json:"price" binding:"required,gt=0"`
	Currency                  string  `json:"currency" binding:"required,oneof=INR USD"`
	DurationDays              uint    `json:"duration_days" binding:"required,gt=0"`
	AdminCommissionPercentage float64 `json:"admin_commission_percentage" binding:"required,gte=0,lte=100"`
}
type SubscribeToTraderInput struct {
	TraderSubscriptionPlanID uint `json:"trader_subscription_plan_id" binding:"required"`
}
