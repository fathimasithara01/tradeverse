package models

import (
	"time"

	"gorm.io/gorm"
)

type TraderPerformance struct {
	gorm.Model
	TraderID        uint      `gorm:"uniqueIndex;not null" json:"trader_id"`
	Trader          User      `gorm:"foreignKey:TraderID" json:"trader"`
	TotalROI        float64   `gorm:"type:numeric(10,4);default:0.00" json:"total_roi"` // Total Return on Investment
	DailyROI        float64   `gorm:"type:numeric(10,4);default:0.00" json:"daily_roi"`
	WeeklyROI       float64   `gorm:"type:numeric(10,4);default:0.00" json:"weekly_roi"`
	MonthlyROI      float64   `gorm:"type:numeric(10,4);default:0.00" json:"monthly_roi"`
	MaxDrawdown     float64   `gorm:"type:numeric(10,4);default:0.00" json:"max_drawdown"` // Max percentage loss from a peak
	WinRate         float64   `gorm:"type:numeric(5,2);default:0.00" json:"win_rate"`
	LossRate        float64   `gorm:"type:numeric(5,2);default:0.00" json:"loss_rate"`
	AverageProfit   float64   `gorm:"type:numeric(18,8);default:0.00" json:"average_profit"`
	AverageLoss     float64   `gorm:"type:numeric(18,8);default:0.00" json:"average_loss"`
	TotalTrades     uint      `gorm:"default:0" json:"total_trades"`
	WinningTrades   uint      `gorm:"default:0" json:"winning_trades"`
	LosingTrades    uint      `gorm:"default:0" json:"losing_trades"`
	ActiveCopiers   uint      `gorm:"default:0" json:"active_copiers"` // Number of customers currently copying
	LastUpdated     time.Time `json:"last_updated"`
	TradingStyle    string    `gorm:"size:255" json:"trading_style,omitempty"` // e.g., "Scalping", "Swing Trading"
	Bio             string    `gorm:"type:text" json:"bio,omitempty"`          // Trader's self-description
	IsPublicProfile bool      `gorm:"default:true" json:"is_public_profile"`   // If trader wants to be listed on leaderboard
}
