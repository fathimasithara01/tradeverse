package utils

import (
	"errors"
	"math/rand"
	"time"
)

func GetCurrentMarketPrice(symbol string) (float64, error) {
	if symbol == "" {
		return 0, errors.New("invalid symbol")
	}
	rand.Seed(time.Now().UnixNano())
	price := 100 + rand.Float64()*50 
	return price, nil
}
