package models

// type TraderSubscriptionPlan struct {
// 	gorm.Model
// 	TraderID uint // Foreign key to User

// 	// TraderID                  uint    `gorm:"index;not null" json:"trader_id"` // The trader offering this plan
// 	Name                      string  `gorm:"not null" json:"name"` // e.g., "Monthly Access to My Signals"
// 	Description               string  `gorm:"type:text" json:"description"`
// 	Price                     float64 `gorm:"type:numeric(18,4);not null" json:"price"`
// 	Currency                  string  `gorm:"size:10;not null;default:'USD'" json:"currency"`
// 	DurationDays              uint    `json:"duration_days"` // How long the customer gets access for
// 	IsActive                  bool    `gorm:"default:true" json:"is_active"`
// 	AdminCommissionPercentage float64 `gorm:"type:numeric(5,2);default:0.00" json:"admin_commission_percentage"` // e.g., 10.00 for 10%

// 	// Associations
// 	Trader User `gorm:"foreignKey:TraderID"`

// 	UserID uint `gorm:"not null;index" json:"user_id"` // Customer
// 	User   User `gorm:"foreignKey:UserID" json:"user"`

// 	TraderSubscriptionPlanID uint              `gorm:"not null;index" json:"trader_subscription_plan_id"`
// 	TraderSubscriptionPlan   *SubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID" json:"trader_subscription_plan"`
// 	PaymentStatus            string            `gorm:"size:50;not null" json:"payment_status"`

// 	AmountPaid      float64   `gorm:"type:numeric(18,4);not null" json:"amount_paid"`
// 	TraderShare     float64   `gorm:"type:numeric(18,4);not null" json:"trader_share"`
// 	AdminCommission float64   `gorm:"type:numeric(18,4);not null" json:"admin_commission"`
// 	TransactionID   string    `gorm:"size:255;not null" json:"transaction_id"`
// 	StartDate       time.Time `gorm:"not null" json:"start_date"`
// 	EndDate         time.Time `gorm:"not null" json:"end_date"`
// }
// // type TraderSubscriptionPlan struct {
// 	gorm.Model
// 	TraderID uint // Foreign key to User

// 	// TraderID                  uint    `gorm:"index;not null" json:"trader_id"` // The trader offering this plan
// 	Name                      string  `gorm:"not null" json:"name"` // e.g., "Monthly Access to My Signals"
// 	Description               string  `gorm:"type:text" json:"description"`
// 	Price                     float64 `gorm:"type:numeric(18,4);not null" json:"price"`
// 	Currency                  string  `gorm:"size:10;not null;default:'USD'" json:"currency"`
// 	DurationDays              uint    `json:"duration_days"` // How long the customer gets access for
// 	IsActive                  bool    `gorm:"default:true" json:"is_active"`
// 	AdminCommissionPercentage float64 `gorm:"type:numeric(5,2);default:0.00" json:"admin_commission_percentage"` // e.g., 10.00 for 10%

// 	// Associations
// 	Trader User `gorm:"foreignKey:TraderID"`

// 	UserID uint `gorm:"not null;index" json:"user_id"` // Customer
// 	User   User `gorm:"foreignKey:UserID" json:"user"`

// 	TraderSubscriptionPlanID uint              `gorm:"not null;index" json:"trader_subscription_plan_id"`
// 	TraderSubscriptionPlan   *SubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID" json:"trader_subscription_plan"`
// 	PaymentStatus            string            `gorm:"size:50;not null" json:"payment_status"`

// 	AmountPaid      float64   `gorm:"type:numeric(18,4);not null" json:"amount_paid"`
// 	TraderShare     float64   `gorm:"type:numeric(18,4);not null" json:"trader_share"`
// 	AdminCommission float64   `gorm:"type:numeric(18,4);not null" json:"admin_commission"`
// 	TransactionID   string    `gorm:"size:255;not null" json:"transaction_id"`
// 	StartDate       time.Time `gorm:"not null" json:"start_date"`
// 	EndDate         time.Time `gorm:"not null" json:"end_date"`
// }
// type TraderSubscriptionPlan struct {
// 	gorm.Model

// 	UserID uint `gorm:"not null;index" json:"user_id"` // Customer
// 	User   User `gorm:"foreignKey:UserID" json:"user"`

// 	TraderID uint `gorm:"not null;index" json:"trader_id"`
// 	Trader   User `gorm:"foreignKey:TraderID" json:"trader"`

// 	TraderSubscriptionPlanID uint              `gorm:"not null;index" json:"trader_subscription_plan_id"`
// 	TraderSubscriptionPlan   *SubscriptionPlan `gorm:"foreignKey:TraderSubscriptionPlanID" json:"trader_subscription_plan"`

// 	StartDate     time.Time `gorm:"not null" json:"start_date"`
// 	EndDate       time.Time `gorm:"not null" json:"end_date"`
// 	IsActive      bool      `gorm:"default:true" json:"is_active"`
// 	PaymentStatus string    `gorm:"size:50;not null" json:"payment_status"`

// 	AmountPaid      float64 `gorm:"type:numeric(18,4);not null" json:"amount_paid"`
// 	TraderShare     float64 `gorm:"type:numeric(18,4);not null" json:"trader_share"`
// 	AdminCommission float64 `gorm:"type:numeric(18,4);not null" json:"admin_commission"`
// 	TransactionID   string  `gorm:"size:255;not null" json:"transaction_id"`
// }

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
