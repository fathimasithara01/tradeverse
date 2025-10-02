package models

import (
	"time"

	"gorm.io/gorm"
)

// MarketData represents the real-time or near real-time market information for a cryptocurrency.
type MarketData struct {
	gorm.Model
	Symbol         string    `gorm:"uniqueIndex;not null" json:"symbol"` // e.g., "BTCUSDT"
	Name           string    `json:"name"`                               // e.g., "Bitcoin"
	CurrentPrice   float64   `gorm:"type:numeric(20,8);not null" json:"current_price"`
	PriceChange24H float64   `gorm:"type:numeric(10,4)" json:"price_change_24h"` // Percentage change
	Volume24H      float64   `gorm:"type:numeric(25,8)" json:"volume_24h"`
	MarketCap      float64   `gorm:"type:numeric(30,8)" json:"market_cap"`
	LogoURL        string    `json:"logo_url"` // URL to the coin's logo
	LastUpdated    time.Time `gorm:"index" json:"last_updated"`
}

// MarketDataAPIResponse structure for data sent to the frontend
type MarketDataAPIResponse struct {
	Symbol         string  `json:"symbol"`
	Name           string  `json:"name"`
	CurrentPrice   float64 `json:"current_price"`
	PriceChange24H float64 `json:"price_change_24h"`
	LogoURL        string  `json:"logo_url"`
}
