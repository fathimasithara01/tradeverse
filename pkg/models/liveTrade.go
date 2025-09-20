package models

import "time"

type LiveTrade struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	TraderID   uint       `json:"trader_id"`
	Symbol     string     `json:"symbol"`
	TradeType  string     `json:"trade_type"` // e.g., spot, futures
	Side       string     `json:"side"`       // buy or sell
	EntryPrice float64    `json:"entry_price"`
	ClosePrice *float64   `json:"close_price,omitempty"`
	Status     string     `json:"status"` // OPEN or CLOSED
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
