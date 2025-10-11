package models

import (
	"time"

	"gorm.io/gorm"
)

// type Signal struct {
// 	gorm.Model
// 	TraderID       uint      `gorm:"not null;index" json:"trader_id"`
// 	TraderName     string    `gorm:"not null" json:"trader_name"`
// 	Symbol         string    `gorm:"size:20;index" json:"symbol"`
// 	EntryPrice     float64   `gorm:"type:numeric(18,8)" json:"entry_price"`
// 	TargetPrice    float64   `gorm:"type:numeric(18,8)" json:"target_price"`
// 	StopLoss       float64   `gorm:"type:numeric(18,8)" json:"stop_loss"`
// 	CurrentPrice   float64   `gorm:"type:numeric(18,8)" json:"current_price"`
// 	TradeStartDate time.Time `json:"trade_start_date"` // ✅ Add this

//		Strategy      string    `gorm:"size:50" json:"strategy"`
//		Risk          string    `gorm:"size:20" json:"risk"`
//		StartDate     time.Time `json:"start_date"`
//		EndDate       time.Time `json:"end_date"`
//		TotalDuration string    `gorm:"size:100" json:"total_duration"`
//		Status        string    `gorm:"size:50;default:'Pending'" json:"status"`
//	}
// type Signal struct {
// 	gorm.Model
// 	TraderID       uint      `gorm:"not null;index" json:"trader_id"`
// 	TraderName     string    `json:"traderName"`
// 	Symbol         string    `json:"symbol" gorm:"index"`
// 	Strategy       string    `gorm:"size:50" json:"strategy"`
// 	Risk           string    `gorm:"size:20" json:"risk"`
// 	TradeStartDate time.Time `json:"tradeStartDate"`
// 	TradeEndDate   time.Time `json:"tradeEndDate"`
// 	TotalDuration  string    `json:"totalDuration"`
// 	EntryPrice     float64   `json:"entryPrice"`
// 	TargetPrice    float64   `json:"targetPrice"`
// 	StopLoss       float64   `json:"stopLoss"`
// 	CurrentPrice   float64   `json:"currentPrice" gorm:"default:0"` // Ensure this field exists and is exported
// 	Status         string    `json:"status" gorm:"default:'Pending'"`
// 	// Add other fields as necessary
// }

type Signal struct {
	gorm.Model
	TraderID      uint   `gorm:"index;not null" json:"trader_id"` // Who created the signal
	TraderName    string `json:"traderName"`
	TotalDuration string `json:"totalDuration"`

	Symbol         string  `gorm:"size:20;not null" json:"symbol"`
	EntryPrice     float64 `gorm:"type:numeric(18,4);not null" json:"entry_price"`
	CurrentPrice   float64 `gorm:"type:numeric(18,4)" json:"current_price"` // Updated periodically
	TargetPrice    float64 `gorm:"type:numeric(18,4);not null" json:"target_price"`
	StopLoss       float64 `gorm:"type:numeric(18,4);not null" json:"stop_loss"`
	Strategy       string  `gorm:"type:text" json:"strategy"`
	Risk           string  `gorm:"size:20" json:"risk"`                     // e.g., "Low", "Medium", "High"
	Status         string  `gorm:"size:20;default:'Pending'" json:"status"` // e.g., "Pending", "Active", "Target Hit", "Stop Loss"
	PublishedAt    time.Time
	DeactivatedAt  *time.Time `json:"deactivated_at"`
	TradeStartDate time.Time  `json:"tradeStartDate"`
	TradeEndDate   time.Time  `json:"tradeEndDate"`
	// Add other signal-related fields
}
