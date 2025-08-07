package models

import "gorm.io/gorm"

type Signal struct {
	gorm.Model
	TraderID   uint    `json:"trader_id"`
	Title      string  `json:"title"`
	Pair       string  `json:"pair"`   // e.g., BTC/USDT
	Action     string  `json:"action"` // e.g., BUY or SELL
	EntryPrice string  `json:"entry_price"`
	StopLoss   string  `json:"stop_loss"`
	TakeProfit string  `json:"take_profit"`
	TimeFrame  string  `json:"time_frame"`                     // e.g., 1H, 4H
	Status     string  `json:"status" gorm:"default:'active'"` // Options: active, inactive
	Profit     float64 `json:"profit"`                         // optional
}
