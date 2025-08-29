package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/fathimasithara01/tradeverse/pkg/repository"
)

type ILiveSignalService interface {
	getCurrentPrice(symbol string) (float64, error)
	GetLiveSignals() ([]TraderSignal, error)
}

type TraderSignal struct {
	Trader       models.User  `json:"trader"`
	LatestTrade  models.Trade `json:"latest_trade"`
	CurrentPrice float64      `json:"current_price"`
	LivePNL      float64      `json:"live_pnl"`
	LivePNLPct   float64      `json:"live_pnl_pct"`
}

type PolygonResponse struct {
	Results []struct {
		ClosePrice float64 `json:"c"`
	} `json:"results"`
	Status string `json:"status"`
}

type LiveSignalService struct {
	UserRepo repository.IUserRepository
}

func NewLiveSignalService(userRepo repository.IUserRepository) ILiveSignalService {
	return &LiveSignalService{UserRepo: userRepo}
}

func (s *LiveSignalService) getCurrentPrice(symbol string) (float64, error) {
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

func (s *LiveSignalService) GetLiveSignals() ([]TraderSignal, error) {
	approvedTraders, err := s.UserRepo.FindTradersByStatus(models.StatusApproved)
	if err != nil {
		return nil, err
	}
	signals := make([]TraderSignal, 0)

	for _, trader := range approvedTraders {
		latestTrade, err := s.UserRepo.GetLatestOpenTradeForUser(trader.ID)
		if err != nil {
			continue
		}

		currentPrice, err := s.getCurrentPrice(latestTrade.Symbol)
		if err != nil {
			fmt.Printf("Warning: could not fetch price for %s: %v\n", latestTrade.Symbol, err)
			continue
		}

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
