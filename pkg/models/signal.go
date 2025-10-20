package models

import (
	"time"

	"gorm.io/gorm"
)

type Signal struct {
	gorm.Model
	TraderID      uint   `gorm:"index;not null" json:"trader_id"`
	TraderName    string `json:"traderName"`
	TotalDuration string `json:"totalDuration"`

	Symbol         string  `gorm:"size:20;not null" json:"symbol"`
	EntryPrice     float64 `gorm:"type:numeric(18,4);not null" json:"entry_price"`
	CurrentPrice   float64 `gorm:"type:numeric(18,4)" json:"current_price"`
	TargetPrice    float64 `gorm:"type:numeric(18,4);not null" json:"target_price"`
	StopLoss       float64 `gorm:"type:numeric(18,4);not null" json:"stop_loss"`
	Strategy       string  `gorm:"type:text" json:"strategy"`
	Risk           string  `gorm:"size:20" json:"risk"`
	Status         string  `gorm:"size:20;default:'Pending'" json:"status"`
	PublishedAt    time.Time
	DeactivatedAt  *time.Time `json:"deactivated_at"`
	TradeStartDate time.Time  `json:"tradeStartDate"`
	TradeEndDate   time.Time  `json:"tradeEndDate"`
	CreatedBy      string     `json:"createdBy"`
	CreatorID      uint       `json:"creatorId"`
}
