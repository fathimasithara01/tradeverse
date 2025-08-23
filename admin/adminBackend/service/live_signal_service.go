package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/repository"
)

// This struct defines the rich data object we'll send to our frontend.
type TraderSignal struct {
	Trader       models.User  `json:"trader"`
	LatestTrade  models.Trade `json:"latest_trade"`
	CurrentPrice float64      `json:"current_price"`
	LivePNL      float64      `json:"live_pnl"`
	LivePNLPct   float64      `json:"live_pnl_pct"`
}

// This struct matches the JSON response from the Polygon.io "Previous Close" API.
type PolygonResponse struct {
	Results []struct {
		ClosePrice float64 `json:"c"`
	} `json:"results"`
	Status string `json:"status"`
}

type LiveSignalService struct{ UserRepo *repository.UserRepository }

func NewLiveSignalService(userRepo *repository.UserRepository) *LiveSignalService {
	return &LiveSignalService{UserRepo: userRepo}
}

// getCurrentPrice fetches the latest price for a given trading symbol from Polygon.io.
func (s *LiveSignalService) getCurrentPrice(symbol string) (float64, error) {
	// Example Symbol for Polygon: Forex -> "C:EURUSD", Crypto -> "X:BTCUSD", Stock -> "AAPL"
	// We might need to format our stored symbol to match this. For now, let's assume it matches.
	apiKey := config.AppConfig.PolygonApiKey
	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/prev?adjusted=true&apiKey=%s", symbol, apiKey)

	var polygonData PolygonResponse
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&polygonData); err != nil {
		return 0, err
	}

	if polygonData.Status == "OK" && len(polygonData.Results) > 0 {
		return polygonData.Results[0].ClosePrice, nil
	}

	return 0, fmt.Errorf("could not get price for symbol %s", symbol)
}

// GetLiveSignals is the main function that builds the data for the page.
func (s *LiveSignalService) GetLiveSignals() ([]TraderSignal, error) {
	approvedTraders, err := s.UserRepo.FindTradersByStatus(models.StatusApproved)
	if err != nil {
		return nil, err
	}
	signals := make([]TraderSignal, 0)

	for _, trader := range approvedTraders {
		latestTrade, err := s.UserRepo.GetLatestOpenTradeForUser(trader.ID)
		if err != nil {
			// This is normal; it just means the trader has no open trades. We skip them.
			continue
		}

		// 3. Get the REAL, LIVE price for that trade's specific symbol.
		currentPrice, err := s.getCurrentPrice(latestTrade.Symbol)
		if err != nil {
			// If the API fails for one symbol, we log it and continue to the next trader.
			fmt.Printf("Warning: could not fetch price for %s: %v\n", latestTrade.Symbol, err)
			continue
		}

		// 4. Calculate P&L and build the final signal object.
		pnl := currentPrice - latestTrade.EntryPrice
		pnlPct := (pnl / latestTrade.EntryPrice) * 100

		signals = append(signals, TraderSignal{
			Trader:       trader,
			LatestTrade:  latestTrade,
			CurrentPrice: currentPrice,
			LivePNL:      pnl,
			LivePNLPct:   pnlPct,
		})
	}

	return signals, nil
}
