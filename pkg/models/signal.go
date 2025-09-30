// pkg/models/signal.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Signal struct {
	gorm.Model
	TraderName   string  `gorm:"not null" json:"traderName"`
	StopLoss     float64 `gorm:"type:numeric(18,8)" json:"stopLoss"`
	EntryPrice   float64 `gorm:"type:numeric(18,8)" json:"entryPrice"`
	TargetPrice  float64 `gorm:"type:numeric(18,8)" json:"targetPrice"`
	CurrentPrice float64 `gorm:"type:numeric(20,8)" json:"currentPrice"` // keep more for live data

	Strategy       string    `gorm:"size:50" json:"strategy"`
	Risk           string    `gorm:"size:20" json:"risk"` // Low, Medium, High
	TradeStartDate time.Time `json:"startDate"`
	TradeEndDate   time.Time `json:"endDate"`
	TotalDuration  string    `gorm:"size:100" json:"totalDuration"`           // e.g., "4 weeks 2 days"
	Status         string    `gorm:"size:50;default:'Pending'" json:"status"` // e.g., Pending, Target Hit, Stop Loss
	Symbol         string    `gorm:"size:20;index" json:"symbol"`             // e.g., "BTCUSDT", "ETHUSDT"
}
