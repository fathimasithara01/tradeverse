package service

import (
	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type LiveTradeService interface {
	PublishLiveTrade(trade *models.LiveTrade) error
	GetActiveTrades(traderID uint) ([]models.LiveTrade, error)
}

type liveTradeService struct {
	repo repository.LiveTradeRepository
}

func NewLiveTradeService(repo repository.LiveTradeRepository) LiveTradeService {
	return &liveTradeService{repo: repo}
}

func (s *liveTradeService) PublishLiveTrade(trade *models.LiveTrade) error {
	return s.repo.Create(trade)
}

func (s *liveTradeService) GetActiveTrades(traderID uint) ([]models.LiveTrade, error) {
	return s.repo.GetActive(traderID)
}
