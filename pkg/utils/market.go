// internal/utils/market.go
package utils

import (
	"errors"
	"math/rand"
	"time"
)

// Mock function to get current market price (replace with real API later)
func GetCurrentMarketPrice(symbol string) (float64, error) {
	if symbol == "" {
		return 0, errors.New("invalid symbol")
	}
	// For now, return random price
	rand.Seed(time.Now().UnixNano())
	price := 100 + rand.Float64()*50 // between 100 - 150
	return price, nil
}
