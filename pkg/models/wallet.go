package models

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	UserID  uint    `gorm:"uniqueIndex;not null"`
	Balance float64 `gorm:"type:numeric(18,4);default:0.00"`
	// Currency    string  `gorm:"size:3;default:'INR'"`
	Currency string `gorm:"size:10;not null;default:'USD'" json:"currency"`

	LastUpdated time.Time
}

type TransactionType string

const (
	TxTypeDeposit      TransactionType = "DEPOSIT"
	TxTypeWithdraw     TransactionType = "WITHDRAW"
	TxTypeFee          TransactionType = "FEE"
	TxTypeTransfer     TransactionType = "TRANSFER"
	TxTypeReversal     TransactionType = "REVERSAL"
	TxTypeSubscription TransactionType = "SUBSCRIPTION_PAYMENT"

	TxTypeTradeProfit        TransactionType = "TRADE_PROFIT"   // NEW
	TxTypeTradeLoss          TransactionType = "TRADE_LOSS"     // NEW
	TxTypeCopyTradeFee       TransactionType = "COPY_TRADE_FEE" // NEW (paid by customer to trader)
	TxTypeReferralCommission TransactionType = "REFERRAL_COMMISSION"
)

type TransactionStatus string

const (
	TxStatusPending   TransactionStatus = "PENDING"
	TxStatusSuccess   TransactionStatus = "SUCCESS"
	TxStatusFailed    TransactionStatus = "FAILED"
	TxStatusCancelled TransactionStatus = "CANCELLED"
	TxStatusReversed  TransactionStatus = "REVERSED"
	TxStatusRejected  TransactionStatus = "REJECTED" // <--- ADD THIS LINE
)

type WalletTransaction struct {
	gorm.Model
	WalletID           uint              `gorm:"index;not null"`
	UserID             uint              `gorm:"index;not null"`
	TransactionType    TransactionType   `gorm:"size:20;not null"`
	Amount             float64           `gorm:"type:numeric(18,4);not null"`
	Currency           string            `gorm:"size:3;not null"`
	Status             TransactionStatus `gorm:"size:20;not null"`
	ReferenceID        string            `gorm:"size:100"` // General reference (e.g., related trade ID, subscription ID)
	PaymentGatewayTxID string            `gorm:"size:100"`
	Description        string            `gorm:"type:text"`
	BalanceBefore      float64           `gorm:"type:numeric(18,4)"`
	BalanceAfter       float64           `gorm:"type:numeric(18,4)"`
	// New field for associating with a specific trade or copy trade
	TradeID              *uint `gorm:"index" json:"trade_id,omitempty"`               // For TxTypeTradeProfit/Loss
	CopyTradeID          *uint `gorm:"index" json:"copy_trade_id,omitempty"`          // For TxTypeCopyTradeFee
	ReferralID           *uint `gorm:"index" json:"referral_id,omitempty"`            // For TxTypeReferralCommission
	SubscriptionID       *uint `gorm:"index" json:"subscription_id,omitempty"`        // For TxTypeSubscription
	TraderSubscriptionID *uint `gorm:"index" json:"trader_subscription_id,omitempty"` // For TxTypeSubscription
}
type DepositRequest struct {
	gorm.Model
	UserID              uint              `gorm:"index;not null"`
	Amount              float64           `gorm:"type:numeric(18,4);not null"`
	Currency            string            `gorm:"size:3;not null"`
	Status              TransactionStatus `gorm:"size:20;not null"`
	PaymentGateway      string            `gorm:"size:50"`
	PaymentGatewayTxID  string            `gorm:"size:100"`
	RedirectURL         string            `gorm:"size:255"`
	WalletTransactionID *uint             `gorm:"index"`
	AdminNotes          string            `gorm:"type:text" json:"admin_notes,omitempty"` // NEW
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
	AdminNotes          string            `gorm:"type:text" json:"admin_notes,omitempty"` // NEW
}

type WalletSummaryResponse struct {
	UserID      uint      `json:"user_id"`
	WalletID    uint      `json:"wallet_id"` // Added WalletID
	Balance     float64   `json:"balance"`
	Currency    string    `json:"currency"`
	LastUpdated time.Time `json:"last_updated"`
}

type DepositRequestInput struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required,oneof=INR USD"`
}

type DepositResponse struct {
	DepositID          uint              `json:"deposit_id"`
	Amount             float64           `json:"amount"`
	Currency           string            `json:"currency"`
	Status             TransactionStatus `json:"status"`
	RedirectURL        string            `json:"redirect_url,omitempty"`
	PaymentGatewayTxID string            `json:"payment_gateway_tx_id,omitempty"`
	Message            string            `json:"message"`
}

type DepositVerifyInput struct {
	PaymentGatewayTxID string  `json:"payment_gateway_tx_id" binding:"required"`
	Amount             float64 `json:"amount"` // Can be optional if gateway webhook provides it reliably
	Status             string  `json:"status"` // E.g., "success", "failed"
	WebhookSignature   string  `json:"webhook_signature,omitempty"`
}

type WithdrawalRequestInput struct {
	Amount             float64 `json:"amount" binding:"required,gt=0"`
	Currency           string  `json:"currency" binding:"required,oneof=INR USD"`
	BeneficiaryAccount string  `json:"beneficiary_account" binding:"required"` // Can be a JSON object for more details
}

type WithdrawalResponse struct {
	WithdrawalID       uint              `json:"withdrawal_id"`
	Amount             float64           `json:"amount"`
	Currency           string            `json:"currency"`
	Status             TransactionStatus `json:"status"`
	PaymentGatewayTxID string            `json:"payment_gateway_tx_id,omitempty"`
	Message            string            `json:"message"`
}

type TransactionListResponse struct {
	Transactions []WalletTransaction `json:"transactions"`
	Total        int64               `json:"total"`
	Page         int                 `json:"page"`
	Limit        int                 `json:"limit"`
}

type PaginationParams struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}
