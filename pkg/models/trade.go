package models

import (
	"time"

	"gorm.io/gorm"
)

type TradeType string

const (
	TradeTypeMarket TradeType = "MARKET"
	TradeTypeLimit  TradeType = "LIMIT"
	TradeTypeStop   TradeType = "STOP"
)

type TradeSide string

const (
	TradeTypeSpot   TradeType = "spot"
	TradeTypeFuture TradeType = "future"
	TradeSideBuy    TradeSide = "BUY"
	TradeSideSell   TradeSide = "SELL"
)

type TradeStatus string

const (
	TradeStatusPending   TradeStatus = "PENDING"
	TradeStatusOpen      TradeStatus = "OPEN"
	TradeStatusClosed    TradeStatus = "CLOSED"
	TradeStatusCancelled TradeStatus = "CANCELLED"
	TradeStatusFailed    TradeStatus = "FAILED"
)

type Trader struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type Trade struct {
	gorm.Model

	TraderID        uint        `gorm:"index;not null" json:"trader_id"` // This is the column that stores the FK
	Symbol          string      `gorm:"size:20;not null" json:"symbol"`
	TradeType       TradeType   `gorm:"size:10;not null" json:"trade_type"`
	Side            TradeSide   `gorm:"size:5;not null" json:"side"`
	EntryPrice      float64     `gorm:"type:numeric(18,4);not null" json:"entry_price"`
	ExecutedPrice   *float64    `gorm:"type:numeric(18,4)" json:"executed_price,omitempty"`
	Quantity        float64     `gorm:"type:numeric(18,8);not null" json:"quantity"`
	Leverage        uint        `gorm:"default:1" json:"leverage"`
	StopLossPrice   *float64    `gorm:"type:numeric(18,4)" json:"stop_loss_price,omitempty"`
	TakeProfitPrice *float64    `gorm:"type:numeric(18,4)" json:"take_profit_price,omitempty"`
	Status          TradeStatus `gorm:"size:20;not null" json:"status"`
	ClosePrice      *float64    `gorm:"type:numeric(18,4)" json:"close_price,omitempty"`
	OpenedAt        *time.Time  `json:"opened_at,omitempty"`
	ClosedAt        *time.Time  `json:"closed_at,omitempty"`
	Pnl             *float64    `gorm:"type:numeric(18,4)" json:"pnl,omitempty"`
	Fees            float64     `gorm:"type:numeric(18,4);default:0.00" json:"fees"`

	// THIS IS THE IMPORTANT CHANGE:
	// We explicitly define the foreign key relationship here.
	// `Trader` is the associated model. `foreignKey:TraderID` tells GORM
	// that the `TraderID` field in *this* `Trade` struct holds the foreign key.
	// `references:ID` tells GORM that it refers to the `ID` column of the `User` table.
	// `constraint:OnUpdate:CASCADE,OnDelete:SET NULL` are good practices for referential integrity.
	Trader User `gorm:"foreignKey:TraderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	IsCopyTrade     bool  `gorm:"default:false" json:"is_copy_trade"`
	OriginalTradeID *uint `gorm:"index" json:"original_trade_id,omitempty"`
	CopyProfileID   *uint `gorm:"index" json:"copy_profile_id,omitempty"`
	CustomerID      *uint `gorm:"index" json:"customer_id,omitempty"`
}

// type Trade struct {
// 	gorm.Model

// 	TraderID        uint        `gorm:"index;not null" json:"trader_id"`
// 	Symbol          string      `gorm:"size:20;not null" json:"symbol"`
// 	TradeType       TradeType   `gorm:"size:10;not null" json:"trade_type"`
// 	Side            TradeSide   `gorm:"size:5;not null" json:"side"`
// 	EntryPrice      float64     `gorm:"type:numeric(18,4);not null" json:"entry_price"`
// 	ExecutedPrice   *float64    `gorm:"type:numeric(18,4)" json:"executed_price,omitempty"`
// 	Quantity        float64     `gorm:"type:numeric(18,8);not null" json:"quantity"`
// 	Leverage        uint        `gorm:"default:1" json:"leverage"`
// 	StopLossPrice   *float64    `gorm:"type:numeric(18,4)" json:"stop_loss_price,omitempty"`
// 	TakeProfitPrice *float64    `gorm:"type:numeric(18,4)" json:"take_profit_price,omitempty"`
// 	Status          TradeStatus `gorm:"size:20;not null" json:"status"`
// 	ClosePrice      *float64    `gorm:"type:numeric(18,4)" json:"close_price,omitempty"`
// 	OpenedAt        *time.Time  `json:"opened_at,omitempty"`
// 	ClosedAt        *time.Time  `json:"closed_at,omitempty"`
// 	Pnl             *float64    `gorm:"type:numeric(18,4)" json:"pnl,omitempty"`
// 	Fees            float64     `gorm:"type:numeric(18,4);default:0.00" json:"fees"`

// 	Trader User `gorm:"foreignKey:TraderID" json:"-"`

// 	IsCopyTrade     bool  `gorm:"default:false" json:"is_copy_trade"`
// 	OriginalTradeID *uint `gorm:"index" json:"original_trade_id,omitempty"`
// 	CopyProfileID   *uint `gorm:"index" json:"copy_profile_id,omitempty"`
// 	CustomerID      *uint `gorm:"index" json:"customer_id,omitempty"`
// }

type TradeInput struct {
	Symbol          string    `json:"symbol" binding:"required"`
	TradeType       TradeType `json:"trade_type" binding:"required"`
	Side            TradeSide `json:"side" binding:"required"`
	EntryPrice      float64   `json:"entry_price"`
	Quantity        float64   `json:"quantity" binding:"required,gt=0"`
	Leverage        uint      `json:"leverage"`
	StopLossPrice   *float64  `json:"stop_loss_price"`
	TakeProfitPrice *float64  `json:"take_profit_price"`
}

type TradeUpdateInput struct {
	StopLossPrice   *float64    `json:"stop_loss_price,omitempty"`
	TakeProfitPrice *float64    `json:"take_profit_price,omitempty"`
	Action          string      `json:"action,omitempty" binding:"omitempty,oneof=CLOSE CANCEL"` // e.g., "CLOSE", "CANCEL"
	ClosePrice      *float64    `json:"close_price,omitempty"`                                   // Required if action is CLOSE
	Status          TradeStatus `json:"status,omitempty"`                                        // Admin-only perhaps, for manual status change
}

type TradeListResponse struct {
	Trades []Trade `json:"trades"`
	Total  int64   `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

type TradeRequest struct {
	Symbol          string  `json:"symbol"`
	TradeType       string  `json:"trade_type"`
	Side            string  `json:"side"`
	EntryPrice      float64 `json:"entry_price"`
	Quantity        float64 `json:"quantity"`
	Leverage        int     `json:"leverage"`
	StopLossPrice   float64 `json:"stop_loss_price"`
	TakeProfitPrice float64 `json:"take_profit_price"`
	TraderID        uint    `json:"-"`
}
