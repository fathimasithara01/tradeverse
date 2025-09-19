// package models

// import (
// 	"time"

// 	"gorm.io/gorm"
// )

// // TradeType defines the type of trade order
// type TradeType string

// const (
// 	TradeTypeMarket TradeType = "MARKET"
// 	TradeTypeLimit  TradeType = "LIMIT"
// 	TradeTypeStop   TradeType = "STOP"
// )

// // TradeSide defines whether it's a buy or sell order
// type TradeSide string

// const (
// 	TradeSideBuy  TradeSide = "BUY"
// 	TradeSideSell TradeSide = "SELL"
// )

// // TradeStatus defines the current status of a trade
// type TradeStatus string

// const (
// 	TradeStatusPending   TradeStatus = "PENDING"
// 	TradeStatusOpen      TradeStatus = "OPEN"
// 	TradeStatusClosed    TradeStatus = "CLOSED"
// 	TradeStatusCancelled TradeStatus = "CANCELLED"
// 	TradeStatusFailed    TradeStatus = "FAILED"
// )

// // Trade represents a single trade executed by a trader or copied by a customer.
// type Trade struct {
// 	gorm.Model
// 	TraderID        uint        `gorm:"index;not null" json:"trader_id"` // The user (trader) who initiated or is associated with this trade
// 	Symbol          string      `gorm:"size:20;not null" json:"symbol"`  // e.g., "BTC/USD", "ETH/INR"
// 	TradeType       TradeType   `gorm:"size:10;not null" json:"trade_type"`
// 	Side            TradeSide   `gorm:"size:5;not null" json:"side"` // BUY or SELL
// 	EntryPrice      float64     `gorm:"type:numeric(18,4);not null" json:"entry_price"`
// 	ExecutedPrice   *float64    `gorm:"type:numeric(18,4)" json:"executed_price,omitempty"` // Actual price trade was executed, can differ for market orders
// 	Quantity        float64     `gorm:"type:numeric(18,8);not null" json:"quantity"`
// 	Leverage        uint        `gorm:"default:1" json:"leverage"` // e.g., 1x, 5x, 10x
// 	StopLossPrice   *float64    `gorm:"type:numeric(18,4)" json:"stop_loss_price,omitempty"`
// 	TakeProfitPrice *float64    `gorm:"type:numeric(18,4)" json:"take_profit_price,omitempty"`
// 	Status          TradeStatus `gorm:"size:20;not null" json:"status"`
// 	ClosePrice      *float64    `gorm:"type:numeric(18,4)" json:"close_price,omitempty"`
// 	ClosedAt        *time.Time  `json:"closed_at,omitempty"`
// 	Pnl             *float64    `gorm:"type:numeric(18,4)" json:"pnl,omitempty"` // Profit and Loss for the trade
// 	Fees            float64     `gorm:"type:numeric(18,4);default:0.00" json:"fees"`

// 	// Copy Trading specific fields
// 	IsCopyTrade     bool  `gorm:"default:false" json:"is_copy_trade"`
// 	OriginalTradeID *uint `gorm:"index" json:"original_trade_id,omitempty"` // ID of the master trade if this is a copy
// 	CopyProfileID   *uint `gorm:"index" json:"copy_profile_id,omitempty"`   // Link to the specific copy profile used
// 	CustomerID      *uint `gorm:"index" json:"customer_id,omitempty"`       // If it's a copy trade, this is the customer's ID
// }

// // TradeInput for creating new trades
// type TradeInput struct {
// 	Symbol          string    `json:"symbol" binding:"required"`
// 	TradeType       TradeType `json:"trade_type" binding:"required,oneof=MARKET LIMIT STOP"`
// 	Side            TradeSide `json:"side" binding:"required,oneof=BUY SELL"`
// 	EntryPrice      float64   `json:"entry_price" binding:"required_if=TradeType LIMIT STOP,omitempty,gt=0"` // Required for LIMIT/STOP orders
// 	Quantity        float64   `json:"quantity" binding:"required,gt=0"`
// 	Leverage        uint      `json:"leverage" binding:"omitempty,gt=0"`
// 	StopLossPrice   *float64  `json:"stop_loss_price,omitempty"`
// 	TakeProfitPrice *float64  `json:"take_profit_price,omitempty"`
// 	// For copy trading, these would be set internally, not by direct input
// 	// IsCopyTrade bool `json:"is_copy_trade"`
// 	// OriginalTradeID *uint `json:"original_trade_id"`
// }

// // TradeUpdateInput for updating existing trades (e.g., modifying stop loss, closing)
// type TradeUpdateInput struct {
// 	StopLossPrice   *float64    `json:"stop_loss_price,omitempty"`
// 	TakeProfitPrice *float64    `json:"take_profit_price,omitempty"`
// 	Action          string      `json:"action,omitempty" binding:"omitempty,oneof=CLOSE CANCEL"` // e.g., "CLOSE", "CANCEL"
// 	ClosePrice      *float64    `json:"close_price,omitempty"`                                   // Required if action is CLOSE
// 	Status          TradeStatus `json:"status,omitempty"`                                        // Admin-only perhaps, for manual status change
// }

// // TradeListResponse for listing trades with pagination
// type TradeListResponse struct {
// 	Trades []Trade `json:"trades"`
// 	Total  int64   `json:"total"`
// 	Page   int     `json:"page"`
// 	Limit  int     `json:"limit"`
// }

// // PaginationParams already defined in wallet.go, can be reused or re-defined if specific
// // type PaginationParams struct {
// // 	Page  int `form:"page,default=1"`
// // 	Limit int `form:"limit,default=10"`
// // }

// pkg/models/models.go (simplified example)
package models

import (
	"time"

	"gorm.io/gorm"
)

// Trader represents a user of the trading platform
type Trader struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	// Other trader specific fields
}

// Trade represents an individual trading position
type Trade struct {
	gorm.Model
	TraderID        uint
	Symbol          string    `gorm:"not null"`                  // e.g., "BTC/USD"
	TradeType       TradeType `gorm:"type:varchar(20);not null"` // MARKET, LIMIT, STOP
	Side            TradeSide `gorm:"type:varchar(10);not null"` // BUY, SELL
	EntryPrice      float64   `gorm:"not null"`
	ExecutedPrice   *float64  // Actual price at which the order was filled
	Quantity        float64   `gorm:"not null"`
	Leverage        uint      `gorm:"not null;default:1"`
	StopLossPrice   *float64
	TakeProfitPrice *float64
	ClosePrice      *float64    // Price at which the trade was closed
	Pnl             *float64    // Profit and Loss
	Fees            float64     `gorm:"not null;default:0"`
	Status          TradeStatus `gorm:"type:varchar(20);not null"` // PENDING, OPEN, CLOSED, CANCELLED, FAILED
	OpenedAt        *time.Time
	ClosedAt        *time.Time
	Trader          Trader `gorm:"foreignKey:TraderID"`
	// WalletTransactionID *uint // Optional: Link to the transaction that funded/affected this trade
}

// Enums
type TradeType string

const (
	TradeTypeMarket TradeType = "MARKET"
	TradeTypeLimit  TradeType = "LIMIT"
	TradeTypeStop   TradeType = "STOP"
)

type TradeSide string

const (
	TradeSideBuy  TradeSide = "BUY"
	TradeSideSell TradeSide = "SELL"
)

type TradeStatus string

const (
	TradeStatusPending   TradeStatus = "PENDING"
	TradeStatusOpen      TradeStatus = "OPEN"
	TradeStatusClosed    TradeStatus = "CLOSED"
	TradeStatusCancelled TradeStatus = "CANCELLED"
	TradeStatusFailed    TradeStatus = "FAILED"
)

// Input DTOs
type TradeInput struct {
	Symbol          string    `json:"symbol" binding:"required"`
	TradeType       TradeType `json:"trade_type" binding:"required"`
	Side            TradeSide `json:"side" binding:"required"`
	EntryPrice      float64   `json:"entry_price"` // Required for LIMIT/STOP, optional for MARKET
	Quantity        float64   `json:"quantity" binding:"required,gt=0"`
	Leverage        uint      `json:"leverage"`
	StopLossPrice   *float64  `json:"stop_loss_price"`
	TakeProfitPrice *float64  `json:"take_profit_price"`
}

type TradeUpdateInput struct {
	StopLossPrice   *float64 `json:"stop_loss_price"`
	TakeProfitPrice *float64 `json:"take_profit_price"`
	Action          string   `json:"action"`      // e.g., "CLOSE", "CANCEL"
	ClosePrice      *float64 `json:"close_price"` // Required if Action is "CLOSE"
}

type TradeListResponse struct {
	Trades []Trade `json:"trades"`
	Total  int64   `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}

// Helper function for *time.Time
func TimePtr(t time.Time) *time.Time {
	return &t
}
