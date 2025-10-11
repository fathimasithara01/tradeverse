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
	TxStatusProcessing TransactionStatus = "PROCESSING" // Added for clarity, especially for withdrawals
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
	Status   TransactionStatus `gorm:"type:varchar(20);default:'PENDING'"` // Changed to TransactionStatus

	PaymentGateway      string `gorm:"size:50"`
	PaymentGatewayTxID  string `gorm:"size:100"`
	RedirectURL         string `gorm:"size:255"`
	WalletTransactionID *uint  `gorm:"index"`
	AdminNotes          string `gorm:"type:text" json:"admin_notes,omitempty"`
	// PaymentMethod       string `gorm:"type:varchar(50);not null"` // e.g., "razorpay", "bank_transfer"
	RequestTime time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"` // Or default:now()

	// RequestTime    time.Time `gorm:"not null"`
	CompletionTime *time.Time
	PaymentMethod  string `gorm:"type:varchar(50);not null;default:'unknown'"` // Add this line or modify existing

}

type DepositRequestInput struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Currency      string  `json:"currency" binding:"required,oneof=INR USD"` // Assuming INR and USD
}

type DepositResponse struct {
	DepositID          uint              `json:"deposit_id"`
	Message            string            `json:"message"`
	RedirectURL        string            `json:"redirect_url,omitempty"` // URL to payment gateway
	Amount             float64           `json:"amount"`
	Currency           string            `json:"currency"`
	Status             TransactionStatus `json:"status"`
	PaymentGatewayTxID string            `json:"payment_gateway_tx_id,omitempty"`
}

type DepositVerifyInput struct {
	PaymentStatus      string  `json:"payment_status" binding:"required"`
	TransactionID      string  `json:"transaction_id"`
	PaymentGatewayTxID string  `json:"payment_gateway_tx_id" binding:"required"`
	Amount             float64 `json:"amount"`
	Status             string  `json:"status"`
	WebhookSignature   string  `json:"webhook_signature,omitempty"`
} // DepositVerifyResponse is the response structure after verifying a deposit.
type DepositVerifyResponse struct {
	DepositID     uint              `json:"deposit_id"`
	Status        TransactionStatus `json:"status"`
	TransactionID string            `json:"transaction_id,omitempty"` // Internal WalletTransaction ID
	Message       string            `json:"message"`
}

// WithdrawalRequest represents a user's request to withdraw funds.
type WithdrawalRequest struct {
	gorm.Model
	UserID             uint              `gorm:"not null"`
	Amount             float64           `gorm:"type:decimal(10,2);not null"`
	Currency           string            `gorm:"size:3;not null"` // <--- ADDED THIS FIELD
	BankAccountNumber  string            `gorm:"type:varchar(50);not null"`
	BankAccountHolder  string            `gorm:"type:varchar(100);not null"`
	IFSCCode           string            `gorm:"type:varchar(20);not null"`
	Status             TransactionStatus `gorm:"type:varchar(20);default:'PENDING'"` // Changed to TransactionStatus
	RequestTime        time.Time         `gorm:"not null"`
	ProcessingTime     *time.Time
	CompletionTime     *time.Time
	AdminNotes         string `gorm:"type:text"`
	PaymentGatewayTxID string `json:"payment_gateway_tx_id"` // <--- ADD THIS LINE
}

// WithdrawalRequestInput is the input structure for creating a withdrawal request.
type WithdrawalRequestInput struct {
	Amount             float64 `json:"amount" binding:"required,gt=0"`
	BankAccountNumber  string  `json:"bank_account_number" binding:"required"`
	BankAccountHolder  string  `json:"bank_account_holder" binding:"required"`
	IFSCCode           string  `json:"ifsc_code" binding:"required"`
	Currency           string  `json:"currency" binding:"required,oneof=INR USD"`
	BeneficiaryAccount string  `json:"beneficiary_account" binding:"required"`
}

// WithdrawalResponse is the output structure after a withdrawal request.
type WithdrawalResponse struct {
	WithdrawalID       uint              `json:"withdrawal_id"`
	Amount             float64           `json:"amount"`
	Currency           string            `json:"currency"`
	Status             TransactionStatus `json:"status"` // Changed to TransactionStatus
	PaymentGatewayTxID string            `json:"payment_gateway_tx_id,omitempty"`
	Message            string            `json:"message"`
}

type WalletTransaction struct {
	gorm.Model
	WalletID uint            `gorm:"index;not null" json:"wallet_id"`
	Type     TransactionType `gorm:"type:varchar(20);not null" json:"type"`
	UserID   uint            `gorm:"not null;index" json:"user_id"`
	Name     string          `gorm:"size:100;not null" json:"name"`

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

type WithdrawRequest struct {
	gorm.Model
	UserID              uint              `gorm:"index;not null"`
	Amount              float64           `gorm:"type:numeric(18,4);not null"`
	Currency            string            `gorm:"size:3;not null"`
	Status              TransactionStatus `gorm:"size:20;not null"`
	BeneficiaryAccount  string            `gorm:"size:100;not null"`
	PaymentGateway      string            `gorm:"size:50"`
	PaymentGatewayTxID  string            `gorm:"size:100"`
	WalletTransactionID *uint             `gorm:"index"`
	AdminNotes          string            `gorm:"type:text" json:"admin_notes,omitempty"`
}

// WalletSummaryResponse provides a summary of a user's wallet.
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
	Page        int    `form:"page,default=1"`
	Limit       int    `form:"limit,default=10"`
	SearchQuery string `form:"search"`
}
