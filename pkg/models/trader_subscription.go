package models

import (
	"time"

	"gorm.io/gorm"
)

type UserSubscription struct {
	gorm.Model
	UserID             uint             `gorm:"index;not null" json:"user_id"`
	User               User             `gorm:"foreignKey:UserID"` // Association to User
	SubscriptionPlanID uint             `gorm:"index;not null" json:"subscription_plan_id"`
	Plan               SubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID"` // Association to SubscriptionPlan
	StartDate          time.Time        `gorm:"not null" json:"start_date"`
	EndDate            time.Time        `gorm:"not null" json:"end_date"`
	IsActive           bool             `gorm:"default:true" json:"is_active"`
	TransactionID      uint             `gorm:"index" json:"transaction_id"` // Link to
}

// type TraderSubscriptionPlan struct {
// 	gorm.Model
// 	User   User `gorm:"foreignKey:UserID" json:"user"`
// 	UserID uint `gorm:"not null;index" json:"user_id"` // Customer

//		TraderID                  uint              `gorm:"index;not null" json:"trader_id"` // The trader (User.ID) offering this plan
//		Name                      string            `gorm:"not null" json:"name"`            // e.g., "Monthly Access to My Signals"
//		Description               string            `gorm:"type:text" json:"description"`
//		Price                     float64           `gorm:"type:numeric(18,4);not null" json:"price"`
//		Currency                  string            `gorm:"size:10;not null;default:'USD'" json:"currency"`
//		DurationDays              uint              `json:"duration_days"` // How long the customer gets access for
//		IsActive                  bool              `gorm:"default:true" json:"is_active"`
//		AdminCommissionPercentage float64           `gorm:"type:numeric(5,2);default:0.00" json:"admin_commission_percentage"` // e.g., 10.00 for 10%
//		AmountPaid                float64           `gorm:"type:numeric(18,4);not null" json:"amount_paid"`
//		StartDate                 time.Time         `gorm:"not null" json:"start_date"`
//		EndDate                   time.Time         `gorm:"not null" json:"end_date"`
//		Trader                    User              `gorm:"foreignKey:TraderID"`
//		PaymentStatus             string            `gorm:"size:50;not null" json:"payment_status"`
//		TraderSubscriptionPlan    *SubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID" json:"trader_subscription_plan"`
//		TraderSubscriptionPlanID  uint              `gorm:"not null;index" json:"trader_subscription_plan_id"`
//	}

type TraderSubscriptionPlan struct {
	gorm.Model

	TraderID        uint    `gorm:"index;not null" json:"trader_id"` // The trader (User.ID) offering this plan
	Name            string  `gorm:"not null" json:"name"`            // e.g., "My Premium Signals Monthly"
	Description     string  `gorm:"type:text" json:"description"`
	Price           float64 `gorm:"type:numeric(18,4);not null" json:"price"`
	Currency        string  `gorm:"size:10;not null;default:'USD'" json:"currency"`
	DurationDays    uint    `json:"duration_days"` // How long the customer gets access for
	IsActive        bool    `gorm:"default:true" json:"is_active"`
	AdminCommission float64 `gorm:"column:admin_commission;type:numeric(5,2);default:0.00" json:"admin_commission_percentage"`
	TraderShare     float64 `gorm:"type:numeric(18,4);not null;default:0.00" json:"trader_share"` // Keep this if you want it as a plan property

	// REMOVE THESE LINES:
	// TraderSubscriptionPlan    *SubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID" json:"trader_subscription_plan"`
	// TraderSubscriptionPlanID  uint              `gorm:"not null;index" json:"trader_subscription_plan_id"` // <-- THIS IS THE PROBLEM FIELD

	// The fields below (AmountPaid, StartDate, EndDate, PaymentStatus)
	// were also likely misplaced from CustomerTraderSubscription and should be removed if they are still present.
	// AmountPaid                float64           `gorm:"type:numeric(18,4);not null" json:"amount_paid"`
	// StartDate                 time.Time         `gorm:"not null" json:"start_date"`
	// EndDate                   time.Time         `gorm:"not null" json:"end_date"`
	// PaymentStatus             string            `gorm:"size:50;not null" json:"payment_status"`

	Trader User `gorm:"foreignKey:TraderID"` // Association to the Trader who owns this plan
}

type CustomerTraderSubscription struct {
	gorm.Model
	CustomerID uint `gorm:"index;not null" json:"customer_id"`
	Customer   User `gorm:"foreignKey:CustomerID"`

	TraderID uint `gorm:"index;not null" json:"trader_id"` // Foreign key to the User who is the trader
	Trader   User `gorm:"foreignKey:TraderID"`             // Association to the Trader User

	// THIS IS WHERE TraderSubscriptionPlanID should be used, linking to the PLAN DEFINITION
	TraderSubscriptionPlanID uint                   `gorm:"index;not null" json:"trader_subscription_plan_id"` // Correctly links to the TraderSubscriptionPlan
	Plan                     TraderSubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID"`               // Association to the TraderSubscriptionPlan

	StartDate           time.Time `gorm:"not null" json:"start_date"`
	EndDate             time.Time `gorm:"not null" json:"end_date"`
	IsActive            bool      `gorm:"default:true" json:"is_active"`
	WalletTransactionID *uint     `gorm:"index" json:"wallet_transaction_id"`
	TransactionID       uint      // ID of the WalletTransaction that paid for this

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
