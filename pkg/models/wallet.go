// package models

// import (
// 	"time"

// 	"gorm.io/gorm"
// )

// type TransactionStatus string

// const (
// 	TxStatusPending    TransactionStatus = "PENDING"
// 	TxStatusSuccess    TransactionStatus = "SUCCESS"
// 	TxStatusFailed     TransactionStatus = "FAILED"
// 	TxStatusCancelled  TransactionStatus = "CANCELLED"
// 	TxStatusReversed   TransactionStatus = "REVERSED"
// 	TxStatusRejected   TransactionStatus = "REJECTED"
// 	TxStatusProcessing TransactionStatus = "PROCESSING"
// )

// type TransactionType string

// const (
// 	TxTypeDeposit            TransactionType = "DEPOSIT"
// 	TxTypeWithdrawal         TransactionType = "WITHDRAWAL"
// 	TxTypeTraderRevenue      TransactionType = "trader_revenue"
// 	TxTypeAdminCommission    TransactionType = "admin_commission"
// 	TxTypeFee                TransactionType = "FEE"
// 	TxTypeTransfer           TransactionType = "TRANSFER"
// 	TxTypeReversal           TransactionType = "REVERSAL"
// 	TxTypeSubscription       TransactionType = "SUBSCRIPTION_PAYMENT"
// 	TxTypeTradeOpeningFunds  TransactionType = "TRADE_OPENING_FUNDS"
// 	TxTypeTradeClosingFunds  TransactionType = "TRADE_CLOSING_FUNDS"
// 	TxTypeTradeProfit        TransactionType = "TRADE_PROFIT"
// 	TxTypeTradeLoss          TransactionType = "TRADE_LOSS"
// 	TxTypeCopyTradeFee       TransactionType = "COPY_TRADE_FEE"
// 	TxTypeReferralCommission TransactionType = "REFERRAL_COMMISSION"
// 	TxTypeCommission         TransactionType = "commission"
// 	TxTypeCredit             TransactionType = "credit"
// 	TxTypeDebit              TransactionType = "debit"
// 	TxTypeSignalPayment      TransactionType = "SIGNAL_PAYMENT"
// )

// type Wallet struct {
// 	gorm.Model
// 	WalletID    uint    `json:"wallet_id"`
// 	UserID      uint    `gorm:"uniqueIndex;not null" json:"user_id"`
// 	Balance     float64 `gorm:"type:numeric(18,4);default:0.00" json:"balance"`
// 	Currency    string  `gorm:"size:10;not null;default:'USD'" json:"currency"`
// 	LastUpdated time.Time

// 	Transactions []WalletTransaction `gorm:"foreignKey:WalletID" json:"transactions,omitempty"`
// }

// type DepositRequest struct {
// 	gorm.Model
// 	UserID   uint              `gorm:"index;not null"`
// 	Amount   float64           `gorm:"type:numeric(18,4);not null"`
// 	Currency string            `gorm:"size:3;not null"`
// 	Status   TransactionStatus `gorm:"type:varchar(20);default:'PENDING'"`

// 	PaymentGateway      string    `gorm:"size:50"`
// 	PaymentGatewayTxID  string    `gorm:"size:100"`
// 	RedirectURL         string    `gorm:"size:255"`
// 	WalletTransactionID *uint     `gorm:"index"`
// 	AdminNotes          string    `gorm:"type:text" json:"admin_notes,omitempty"`
// 	RequestTime         time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"` // Or default:now()

// 	CompletionTime *time.Time
// 	PaymentMethod  string `gorm:"type:varchar(50);not null;default:'unknown'"`
// }

// type DepositRequestInput struct {
// 	Amount        float64 `json:"amount" binding:"required,gt=0"`
// 	PaymentMethod string  `json:"payment_method" binding:"required"`
// 	Currency      string  `json:"currency" binding:"required,oneof=INR USD"`
// }

// type DepositResponse struct {
// 	DepositID          uint              `json:"deposit_id"`
// 	Message            string            `json:"message"`
// 	RedirectURL        string            `json:"redirect_url,omitempty"`
// 	Amount             float64           `json:"amount"`
// 	Currency           string            `json:"currency"`
// 	Status             TransactionStatus `json:"status"`
// 	PaymentGatewayTxID string            `json:"payment_gateway_tx_id,omitempty"`
// }

// type DepositVerifyInput struct {
// 	PaymentStatus      string  `json:"payment_status" binding:"required"`
// 	TransactionID      string  `json:"transaction_id"`
// 	PaymentGatewayTxID string  `json:"payment_gateway_tx_id" binding:"required"`
// 	Amount             float64 `json:"amount"`
// 	Status             string  `json:"status"`
// 	WebhookSignature   string  `json:"webhook_signature,omitempty"`
// }
// type DepositVerifyResponse struct {
// 	DepositID     uint              `json:"deposit_id"`
// 	Status        TransactionStatus `json:"status"`
// 	TransactionID string            `json:"transaction_id,omitempty"`
// 	Message       string            `json:"message"`
// }

// type WithdrawalRequest struct {
// 	ID                 uint              `gorm:"primaryKey" json:"id"`
// 	UserID             uint              `gorm:"not null;index;comment:ID of the user requesting withdrawal" json:"user_id"`
// 	User               User              `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
// 	Amount             float64           `gorm:"type:decimal(18,4);not null" json:"amount"`
// 	Currency           string            `gorm:"size:3;not null;default:'USD'" json:"currency"`
// 	BankAccountNumber  string            `gorm:"size:50;not null" json:"bank_account_number"`
// 	BankAccountHolder  string            `gorm:"size:100;not null" json:"bank_account_holder"`
// 	IFSCCode           string            `gorm:"size:20;not null" json:"ifsc_code"`
// 	Status             TransactionStatus `gorm:"type:varchar(20);default:'PENDING';index" json:"status"`
// 	RequestTime        time.Time         `gorm:"not null;index" json:"request_time"`
// 	ProcessingTime     *time.Time        `json:"processing_time,omitempty"`
// 	CompletionTime     *time.Time        `json:"completion_time,omitempty"`
// 	AdminNotes         string            `gorm:"type:text" json:"admin_notes,omitempty"`
// 	PaymentGatewayTxID string            `gorm:"size:100" json:"payment_gateway_tx_id,omitempty"`
// 	CreatedAt          time.Time         `json:"created_at"`
// 	UpdatedAt          time.Time         `json:"updated_at"`
// }

// type WithdrawalRequestInput struct {
// 	Amount             float64 `json:"amount" binding:"required,gt=0"`
// 	BankAccountNumber  string  `json:"bank_account_number" binding:"required"`
// 	BankAccountHolder  string  `json:"bank_account_holder" binding:"required"`
// 	IFSCCode           string  `json:"ifsc_code" binding:"required"`
// 	Currency           string  `json:"currency" binding:"required,oneof=INR USD"`
// 	BeneficiaryAccount string  `json:"beneficiary_account" binding:"required"`
// }

// type WithdrawalResponse struct {
// 	WithdrawalID       uint              `json:"withdrawal_id"`
// 	Amount             float64           `json:"amount"`
// 	Currency           string            `json:"currency"`
// 	Status             TransactionStatus `json:"status"`
// 	PaymentGatewayTxID string            `json:"payment_gateway_tx_id,omitempty"`
// 	Message            string            `json:"message"`
// }

// type WalletTransaction struct {
// 	gorm.Model
// 	WalletID uint `gorm:"index;not null" json:"wallet_id"`

// 	Type   TransactionType `gorm:"type:varchar(20);not null" json:"type"`
// 	UserID uint            `gorm:"not null;index" json:"user_id"`
// 	Name   string          `gorm:"size:100;not null" json:"name"`

// 	User User `gorm:"foreignKey:UserID"`

// 	TransactionType    TransactionType   `gorm:"size:30;not null"`
// 	Amount             float64           `gorm:"type:numeric(18,4);not null"`
// 	Currency           string            `gorm:"size:3;not null"`
// 	Status             TransactionStatus `gorm:"size:20;not null"`
// 	Notes              string            `json:"notes,omitempty"`
// 	ReferenceID        string            `gorm:"size:100"`
// 	PaymentGatewayTxID string            `gorm:"size:100"`
// 	Description        string            `gorm:"type:text"`
// 	BalanceBefore      float64           `gorm:"type:numeric(18,4)"`
// 	BalanceAfter       float64           `gorm:"type:numeric(18,4)"`

// 	TransactionID string `gorm:"size:255" json:"transaction_id"`

// 	TradeID              *uint `gorm:"index" json:"trade_id,omitempty"`
// 	CopyTradeID          *uint `gorm:"index" json:"copy_trade_id,omitempty"`
// 	ReferralID           *uint `gorm:"index" json:"referral_id,omitempty"`
// 	SubscriptionID       *uint `gorm:"index" json:"subscription_id,omitempty"`
// 	TraderSubscriptionID *uint `gorm:"index" json:"trader_subscription_id,omitempty"`
// }

// type WithdrawRequest struct {
// 	gorm.Model
// 	UserID              uint              `gorm:"index;not null"`
// 	Amount              float64           `gorm:"type:numeric(18,4);not null"`
// 	Currency            string            `gorm:"size:3;not null"`
// 	Status              TransactionStatus `gorm:"size:20;not null"`
// 	BeneficiaryAccount  string            `gorm:"size:100;not null"`
// 	PaymentGateway      string            `gorm:"size:50"`
// 	PaymentGatewayTxID  string            `gorm:"size:100"`
// 	WalletTransactionID *uint             `gorm:"index"`
// 	AdminNotes          string            `gorm:"type:text" json:"admin_notes,omitempty"`
// }

// type WalletSummaryResponse struct {
// 	UserID      uint      `json:"user_id"`
// 	WalletID    uint      `json:"wallet_id"`
// 	Balance     float64   `json:"balance"`
// 	Currency    string    `json:"currency"`
// 	LastUpdated time.Time `json:"last_updated"`
// }

// type TransactionListResponse struct {
// 	Transactions []WalletTransaction `json:"transactions"`
// 	Total        int64               `json:"total"`
// 	Page         int                 `json:"page"`
// 	Limit        int                 `json:"limit"`
// }

// type PaginationParams struct {
// 	Page   int    `form:"page"`
// 	Limit  int    `form:"limit"`
// 	Search string `form:"search"`
// }

// // type PaginationParams struct {
// // 	Page        int    `form:"page,default=1"`
// // 	Limit       int    `form:"limit,default=10"`
// // 	SearchQuery string `form:"search"`
// // }

// type AdminTransactionDisplayDTO struct {
// 	ID              uint              `json:"id"`
// 	UserID          uint              `json:"user_id"`
// 	UserName        string            `json:"user_name"`
// 	UserEmail       string            `json:"user_email"`
// 	UserPhone       string            `json:"user_phone"`
// 	TransactionType TransactionType   `json:"transaction_type"`
// 	Amount          float64           `json:"amount"`
// 	Currency        string            `json:"currency"`
// 	Status          TransactionStatus `json:"status"`
// 	ReferenceID     string            `json:"reference_id"`
// 	Description     string            `json:"description"`
// 	CreatedAt       time.Time         `json:"created_at"`
// 	BalanceBefore   float64           `json:"balance_before"`
// 	BalanceAfter    float64           `json:"balance_after"`
// }

//	type AllTransactionsListResponse struct {
//		Transactions []AdminTransactionDisplayDTO `json:"transactions"`
//		Total        int64                        `json:"total"`
//		Page         int                          `json:"page"`
//		Limit        int                          `json:"limit"`
//	}
package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionStatus string

const (
	TxStatusPending    TransactionStatus = "PENDING"
	TxStatusSuccess    TransactionStatus = "SUCCESS"
	TxStatusFailed     TransactionStatus = "FAILED"
	TxStatusCancelled  TransactionStatus = "CANCELLED"
	TxStatusReversed   TransactionStatus = "REVERSED"
	TxStatusRejected   TransactionStatus = "REJECTED"
	TxStatusProcessing TransactionStatus = "PROCESSING"
)

type TransactionType string

const (
	TxTypeDeposit            TransactionType = "DEPOSIT"
	TxTypeWithdrawal         TransactionType = "WITHDRAWAL"
	TxTypeTraderRevenue      TransactionType = "trader_revenue"
	TxTypeAdminCommission    TransactionType = "admin_commission"
	TxTypeFee                TransactionType = "FEE"
	TxTypeTransfer           TransactionType = "TRANSFER"
	TxTypeReversal           TransactionType = "REVERSAL"
	TxTypeSubscription       TransactionType = "SUBSCRIPTION_PAYMENT"
	TxTypeTradeOpeningFunds  TransactionType = "TRADE_OPENING_FUNDS"
	TxTypeTradeClosingFunds  TransactionType = "TRADE_CLOSING_FUNDS"
	TxTypeTradeProfit        TransactionType = "TRADE_PROFIT"
	TxTypeTradeLoss          TransactionType = "TRADE_LOSS"
	TxTypeCopyTradeFee       TransactionType = "COPY_TRADE_FEE"
	TxTypeReferralCommission TransactionType = "REFERRAL_COMMISSION"
	TxTypeCommission         TransactionType = "commission"
	TxTypeCredit             TransactionType = "credit"
	TxTypeDebit              TransactionType = "debit"
	TxTypeSignalPayment      TransactionType = "SIGNAL_PAYMENT"
)

type Wallet struct {
	gorm.Model
	WalletID    uint    `json:"wallet_id"`
	UserID      uint    `gorm:"uniqueIndex;not null" json:"user_id"`
	Balance     float64 `gorm:"type:numeric(18,4);default:0.00" json:"balance"`
	Currency    string  `gorm:"size:10;not null;default:'USD'" json:"currency"`
	LastUpdated time.Time

	Transactions []WalletTransaction `gorm:"foreignKey:WalletID" json:"transactions,omitempty"`
}

type DepositRequest struct {
	gorm.Model
	UserID   uint              `gorm:"index;not null"`
	Amount   float64           `gorm:"type:numeric(18,4);not null"`
	Currency string            `gorm:"size:3;not null"`
	Status   TransactionStatus `gorm:"type:varchar(20);default:'PENDING'"`

	PaymentGateway      string    `gorm:"size:50"`
	PaymentGatewayTxID  string    `gorm:"size:100"`
	RedirectURL         string    `gorm:"size:255"`
	WalletTransactionID *uint     `gorm:"index"`
	AdminNotes          string    `gorm:"type:text" json:"admin_notes,omitempty"`
	RequestTime         time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"` // Or default:now()

	CompletionTime *time.Time
	PaymentMethod  string `gorm:"type:varchar(50);not null;default:'unknown'"`
}

type DepositRequestInput struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"payment_method"` // Made not required in binding, will be defaulted by controller for admin
	Currency      string  `json:"currency" binding:"required,oneof=INR USD"`
}

type DepositResponse struct {
	DepositID          uint              `json:"deposit_id"`
	Message            string            `json:"message"`
	RedirectURL        string            `json:"redirect_url,omitempty"`
	Amount             float64           `json:"amount"`
	Currency           string            `json:"currency"`
	Status             TransactionStatus `json:"status"`
	PaymentGatewayTxID string            `json:"payment_gateway_tx_id,omitempty"`
}

type DepositVerifyInput struct {
	PaymentStatus      string  `json:"payment_status"`
	TransactionID      string  `json:"transaction_id"`
	PaymentGatewayTxID string  `json:"payment_gateway_tx_id" binding:"required"`
	Amount             float64 `json:"amount"`
	Status             string  `json:"status"`
	WebhookSignature   string  `json:"webhook_signature,omitempty"`
}
type DepositVerifyResponse struct {
	DepositID     uint              `json:"deposit_id"`
	Status        TransactionStatus `json:"status"`
	TransactionID string            `json:"transaction_id,omitempty"`
	Message       string            `json:"message"`
}

// WithdrawalRequest (This seems to be a DTO or specific API response model, not the DB model)
type WithdrawalRequest struct { // Renamed from original to WithdrawalRequestDTO to avoid confusion with DB model
	ID                 uint              `gorm:"primaryKey" json:"id"`
	UserID             uint              `gorm:"not null;index;comment:ID of the user requesting withdrawal" json:"user_id"`
	User               User              `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"` // Ensure User is preloaded
	Amount             float64           `gorm:"type:decimal(18,4);not null" json:"amount"`
	Currency           string            `gorm:"size:3;not null;default:'USD'" json:"currency"`
	BankAccountNumber  string            `gorm:"size:50;not null" json:"bank_account_number"`
	BankAccountHolder  string            `gorm:"size:100;not null" json:"bank_account_holder"`
	IFSCCode           string            `gorm:"size:20;not null" json:"ifsc_code"`
	Status             TransactionStatus `gorm:"type:varchar(20);default:'PENDING';index" json:"status"`
	RequestTime        time.Time         `gorm:"not null;index" json:"request_time"`
	ProcessingTime     *time.Time        `json:"processing_time,omitempty"`
	CompletionTime     *time.Time        `json:"completion_time,omitempty"`
	AdminNotes         string            `gorm:"type:text" json:"admin_notes,omitempty"`
	PaymentGatewayTxID string            `gorm:"size:100" json:"payment_gateway_tx_id,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

type WithdrawalRequestInput struct {
	Amount             float64 `json:"amount" binding:"required,gt=0"`
	BankAccountNumber  string  `json:"bank_account_number" binding:"required"`
	BankAccountHolder  string  `json:"bank_account_holder" binding:"required"`
	IFSCCode           string  `json:"ifsc_code" binding:"required"`
	Currency           string  `json:"currency" binding:"required,oneof=INR USD"`
	BeneficiaryAccount string  `json:"beneficiary_account"` // This can be a descriptive string
}

type WithdrawalResponse struct {
	WithdrawalID       uint              `json:"withdrawal_id"`
	Amount             float64           `json:"amount"`
	Currency           string            `json:"currency"`
	Status             TransactionStatus `json:"status"`
	PaymentGatewayTxID string            `json:"payment_gateway_tx_id,omitempty"`
	Message            string            `json:"message"`
}

type WalletTransaction struct {
	gorm.Model
	WalletID uint `gorm:"index;not null" json:"wallet_id"`

	Type   TransactionType `gorm:"type:varchar(20);not null" json:"type"`
	UserID uint            `gorm:"not null;index" json:"user_id"`
	Name   string          `gorm:"size:100;not null" json:"name"`

	User User `gorm:"foreignKey:UserID"`

	TransactionType    TransactionType   `gorm:"size:30;not null"`
	Amount             float64           `gorm:"type:numeric(18,4);not null"`
	Currency           string            `gorm:"size:3;not null"`
	Status             TransactionStatus `gorm:"size:20;not null"`
	Notes              string            `json:"notes,omitempty"`
	ReferenceID        string            `gorm:"size:100"`
	PaymentGatewayTxID string            `gorm:"size:100"`
	Description        string            `gorm:"type:text"`
	BalanceBefore      float64           `gorm:"type:numeric(18,4)"`
	BalanceAfter       float64           `gorm:"type:numeric(18,4)"`

	TransactionID string `gorm:"size:255" json:"transaction_id"`

	TradeID              *uint `gorm:"index" json:"trade_id,omitempty"`
	CopyTradeID          *uint `gorm:"index" json:"copy_trade_id,omitempty"`
	ReferralID           *uint `gorm:"index" json:"referral_id,omitempty"`
	SubscriptionID       *uint `gorm:"index" json:"subscription_id,omitempty"`
	TraderSubscriptionID *uint `gorm:"index" json:"trader_subscription_id,omitempty"`
}

// WithdrawRequest is the actual DB model for withdrawal requests
type WithdrawRequest struct {
	gorm.Model
	UserID              uint              `gorm:"index;not null"`
	User                User              `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Added for GORM Preload
	Amount              float64           `gorm:"type:numeric(18,4);not null"`
	Currency            string            `gorm:"size:3;not null"`
	Status              TransactionStatus `gorm:"size:20;not null"`
	BeneficiaryAccount  string            `gorm:"type:text;not null"` // Full description
	BankAccountNumber   string            `gorm:"size:50;not null"`   // Separate field for easier handling
	BankAccountHolder   string            `gorm:"size:100;not null"`  // Separate field
	IFSCCode            string            `gorm:"size:20;not null"`   // Separate field
	PaymentGateway      string            `gorm:"size:50"`
	PaymentGatewayTxID  string            `gorm:"size:100"`
	WalletTransactionID *uint             `gorm:"index"`
	AdminNotes          string            `gorm:"type:text" json:"admin_notes,omitempty"`

	RequestTime    time.Time  `gorm:"not null;index" json:"request_time"` // Added from the other WithdrawalRequest
	ProcessingTime *time.Time `json:"processing_time,omitempty"`          // Added from the other WithdrawalRequest
	CompletionTime *time.Time `json:"completion_time,omitempty"`          // Added from the other WithdrawalRequest

}

type WalletSummaryResponse struct {
	UserID      uint      `json:"user_id"`
	WalletID    uint      `json:"wallet_id"`
	Balance     float64   `json:"balance"`
	Currency    string    `json:"currency"`
	LastUpdated time.Time `json:"last_updated"`
}

type TransactionListResponse struct {
	Transactions []WalletTransaction `json:"transactions"`
	Total        int64               `json:"total"`
	Page         int                 `json:"page"`
	Limit        int                 `json:"limit"`
}

type PaginationParams struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
}

type AdminTransactionDisplayDTO struct {
	ID              uint              `json:"id"`
	UserID          uint              `json:"user_id"`
	UserName        string            `json:"user_name"`
	UserEmail       string            `json:"user_email"`
	UserPhone       string            `json:"user_phone"`
	TransactionType TransactionType   `json:"transaction_type"`
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	Status          TransactionStatus `json:"status"`
	ReferenceID     string            `json:"reference_id"`
	Description     string            `json:"description"`
	CreatedAt       time.Time         `json:"created_at"`
	BalanceBefore   float64           `json:"balance_before"`
	BalanceAfter    float64           `json:"balance_after"`
}

type AllTransactionsListResponse struct {
	Transactions []AdminTransactionDisplayDTO `json:"transactions"`
	Total        int64                        `json:"total"`
	Page         int                          `json:"page"`
	Limit        int                          `json:"limit"`
}
