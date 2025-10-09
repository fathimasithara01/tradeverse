package models

import (
	"time"

	"gorm.io/gorm"
)

type Signal struct {
	gorm.Model
	TraderID       uint      `gorm:"not null;index" json:"trader_id"`
	TraderName     string    `gorm:"not null" json:"trader_name"`
	Symbol         string    `gorm:"size:20;index" json:"symbol"`
	EntryPrice     float64   `gorm:"type:numeric(18,8)" json:"entry_price"`
	TargetPrice    float64   `gorm:"type:numeric(18,8)" json:"target_price"`
	StopLoss       float64   `gorm:"type:numeric(18,8)" json:"stop_loss"`
	CurrentPrice   float64   `gorm:"type:numeric(18,8)" json:"current_price"`
	TradeStartDate time.Time `json:"trade_start_date"` // ✅ Add this

	Strategy      string    `gorm:"size:50" json:"strategy"`
	Risk          string    `gorm:"size:20" json:"risk"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	TotalDuration string    `gorm:"size:100" json:"total_duration"`
	Status        string    `gorm:"size:50;default:'Pending'" json:"status"`
}

// package models

// import (
// 	"time"

// 	"gorm.io/gorm"
// )

// type Signal struct {
// 	gorm.Model

// 	UserID   uint  `gorm:"not null;default:1"`
// 	TraderID *uint `gorm:"index" json:"trader_id,omitempty"` // Link to the Trader (User.ID) if IsTraderPlan is true

// 	TraderName   string  `gorm:"not null" json:"traderName"`
// 	StopLoss     float64 `gorm:"type:numeric(18,8)" json:"stopLoss"`
// 	EntryPrice   float64 `gorm:"type:numeric(18,8)" json:"entryPrice"`
// 	TargetPrice  float64 `gorm:"type:numeric(18,8)" json:"targetPrice"`
// 	CurrentPrice float64 `gorm:"type:numeric(20,8)" json:"currentPrice"` // keep more for live data

// 	Strategy       string    `gorm:"size:50" json:"strategy"`
// 	Risk           string    `gorm:"size:20" json:"risk"` // Low, Medium, High
// 	TradeStartDate time.Time `json:"startDate"`
// 	TradeEndDate   time.Time `json:"endDate"`
// 	TotalDuration  string    `gorm:"size:100" json:"totalDuration"`           // e.g., "4 weeks 2 days"
// 	Status         string    `gorm:"size:50;default:'Pending'" json:"status"` // e.g., Pending, Target Hit, Stop Loss
// 	Symbol         string    `gorm:"size:20;index" json:"symbol"`             // e.g., "BTCUSDT", "ETHUSDT"
// }
