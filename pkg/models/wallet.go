package models

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	UserID      uint    `gorm:"uniqueIndex;not null"`
	Balance     float64 `gorm:"type:numeric(18,4);default:0.00"`
	Currency    string  `gorm:"size:3;default:'INR'"`
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
	ReferenceID        string            `gorm:"size:100"`
	PaymentGatewayTxID string            `gorm:"size:100"`
	Description        string            `gorm:"type:text"`
	BalanceBefore      float64           `gorm:"type:numeric(18,4)"`
	BalanceAfter       float64           `gorm:"type:numeric(18,4)"`
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
	Amount             float64 `json:"amount"`
	Status             string  `json:"status"`
	WebhookSignature   string  `json:"webhook_signature,omitempty"`
}

type WithdrawalRequestInput struct {
	Amount             float64 `json:"amount" binding:"required,gt=0"`
	Currency           string  `json:"currency" binding:"required,oneof=INR USD"`
	BeneficiaryAccount string  `json:"beneficiary_account" binding:"required"`
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
