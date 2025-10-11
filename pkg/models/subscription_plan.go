package models

import (
	"time"

	"gorm.io/gorm"
)

type SubscriptionPlan struct {
	gorm.Model
	Name        string        `gorm:"size:100;not null;unique" json:"name"`
	Description string        `gorm:"type:text" json:"description"`
	Price       float64       `gorm:"type:numeric(18,4);not null" json:"price"`
	Currency    string        `gorm:"size:10;not null;default:'USD'" json:"currency"`
	Duration    time.Duration `gorm:"type:integer;not null" json:"duration"` // Duration (e.g., in seconds)
	Interval    string        `gorm:"size:20;not null" json:"interval"`      // e.g., "day", "month", "year"
	IsActive    bool          `gorm:"default:true" json:"is_active"`

	// Fields specific to general plans, potentially created by admin for general users
	IsTraderPlan     bool    `gorm:"default:false" json:"is_trader_plan"` // True if this plan is offered by a trader (perhaps not needed if TraderSubscriptionPlan handles this)
	TraderID         *uint   `gorm:"index" json:"trader_id,omitempty"`    // Link to the Trader (User.ID) if IsTraderPlan is true (Can be removed if TraderSubscriptionPlan is used instead)
	Features         string  `gorm:"type:text" json:"features"`
	MaxFollowers     int     `json:"max_followers,omitempty"`
	CommissionRate   float64 `gorm:"type:numeric(5,4);default:0.10" json:"commission_rate,omitempty"` // Commission for the admin
	AnalyticsAccess  string  `gorm:"size:50" json:"analytics_access,omitempty"`
	CreatedByAdminID uint    `json:"created_by_admin_id"` // User ID of the admin who created it

	IsUpgradeToTrader bool `gorm:"default:false" json:"is_upgrade_to_trader"`
}
type UserSubscription struct {
	gorm.Model
	UserID             uint             `gorm:"index;not null" json:"user_id"`
	User               User             `gorm:"foreignKey:UserID"` // Association to User
	SubscriptionPlanID uint             `gorm:"index;not null" json:"subscription_plan_id"`
	Plan               SubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID"` // Association to SubscriptionPlan
	StartDate          time.Time        `gorm:"not null" json:"start_date"`
	EndDate            time.Time        `gorm:"not null" json:"end_date"`
	IsActive           bool             `gorm:"default:true" json:"is_active"`
	TransactionID      uint             `gorm:"index" json:"transaction_id"` // Link to WalletTransaction
}

// type TraderSubscriptionPlan struct {
// 	gorm.Model
// 	TraderID                  uint    `gorm:"index;not null" json:"trader_id"` // The trader (User.ID) offering this plan
// 	Name                      string  `gorm:"not null" json:"name"`            // e.g., "Monthly Access to My Signals"
// 	Description               string  `gorm:"type:text" json:"description"`
// 	Price                     float64 `gorm:"type:numeric(18,4);not null" json:"price"`
// 	Currency                  string  `gorm:"size:10;not null;default:'USD'" json:"currency"`
// 	DurationDays              uint    `json:"duration_days"` // How long the customer gets access for
// 	IsActive                  bool    `gorm:"default:true" json:"is_active"`
// 	AdminCommissionPercentage float64 `gorm:"type:numeric(5,2);default:0.00" json:"admin_commission_percentage"` // e.g., 10.00 for 10%

//		Trader User `gorm:"foreignKey:TraderID"`
//	}
type TraderSubscriptionPlan struct {
	gorm.Model
	User   User `gorm:"foreignKey:UserID" json:"user"`
	UserID uint `gorm:"not null;index" json:"user_id"` // Customer

	TraderID                  uint              `gorm:"index;not null" json:"trader_id"` // The trader (User.ID) offering this plan
	Name                      string            `gorm:"not null" json:"name"`            // e.g., "Monthly Access to My Signals"
	Description               string            `gorm:"type:text" json:"description"`
	Price                     float64           `gorm:"type:numeric(18,4);not null" json:"price"`
	Currency                  string            `gorm:"size:10;not null;default:'USD'" json:"currency"`
	DurationDays              uint              `json:"duration_days"` // How long the customer gets access for
	IsActive                  bool              `gorm:"default:true" json:"is_active"`
	AdminCommissionPercentage float64           `gorm:"type:numeric(5,2);default:0.00" json:"admin_commission_percentage"` // e.g., 10.00 for 10%
	AmountPaid                float64           `gorm:"type:numeric(18,4);not null" json:"amount_paid"`
	StartDate                 time.Time         `gorm:"not null" json:"start_date"`
	EndDate                   time.Time         `gorm:"not null" json:"end_date"`
	Trader                    User              `gorm:"foreignKey:TraderID"`
	PaymentStatus             string            `gorm:"size:50;not null" json:"payment_status"`
	TraderSubscriptionPlan    *SubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID" json:"trader_subscription_plan"`
	TraderSubscriptionPlanID  uint              `gorm:"not null;index" json:"trader_subscription_plan_id"`
}
type CustomerTraderSubscription struct {
	gorm.Model
	CustomerID uint `gorm:"index;not null" json:"customer_id"`
	Customer   User `gorm:"foreignKey:CustomerID"`

	TraderID uint `gorm:"index;not null" json:"trader_id"` // Foreign key to the User who is the trader
	Trader   User `gorm:"foreignKey:TraderID"`             // Association to the Trader User

	TraderSubscriptionPlanID uint                   `gorm:"index;not null" json:"trader_subscription_plan_id"` // Foreign key to the specific plan
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
	TransactionReferenceID string  `gorm:"size:255;not null" json:"transaction_reference_id"` // Renamed from TransactionID to avoid conflict/clarify
}
