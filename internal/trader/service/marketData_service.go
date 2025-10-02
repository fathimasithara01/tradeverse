package service

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
)

type MarketDataService interface {
	CreateMarketData(symbol string, price float64) (*models.MarketData, error)
}

type marketDataService struct {
	repo repository.MarketDataRepository
}

func NewMarketDataService(repo repository.MarketDataRepository) MarketDataService {
	return &marketDataService{repo: repo}
}

func (s *marketDataService) CreateMarketData(symbol string, price float64) (*models.MarketData, error) {
	data := &models.MarketData{
		Symbol: symbol,
		CurrentPrice:  price,
	}
	if err := s.repo.Create(data); err != nil {
		return nil, err
	}
	return data, nil
}
