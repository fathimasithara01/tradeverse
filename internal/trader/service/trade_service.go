package service

import (
	"context"
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type TradeService interface {
	ListTrades(traderID uint, page, limit int) (*models.TradeListResponse, error)
	GetTrade(id uint, traderID uint) (*models.Trade, error)
	CreateTrade(ctx context.Context, req models.TradeRequest) (models.Trade, error)
	UpdateTrade(id uint, traderID uint, input models.TradeUpdateInput) (*models.Trade, error)
	DeleteTrade(id uint, traderID uint) error
}

type tradeService struct {
	repo repository.TradeRepository
}

func NewTradeService(repo repository.TradeRepository) TradeService {
	return &tradeService{repo: repo}
}

func (s *tradeService) ListTrades(traderID uint, page, limit int) (*models.TradeListResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	trades, total, err := s.repo.GetAllTrades(traderID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.TradeListResponse{
		Trades: trades,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

func (s *tradeService) GetTrade(id uint, traderID uint) (*models.Trade, error) {
	return s.repo.GetTradeByID(id, traderID)
}

func (s *tradeService) CreateTrade(ctx context.Context, req models.TradeRequest) (models.Trade, error) {
	if req.TraderID == 0 {
		return models.Trade{}, errors.New("invalid trader ID")
	}
	return s.repo.CreateTrade(ctx, req)
}

func (s *tradeService) UpdateTrade(id uint, traderID uint, input models.TradeUpdateInput) (*models.Trade, error) {
	trade, err := s.repo.GetTradeByID(id, traderID)
	if err != nil || trade == nil {
		return nil, errors.New("trade not found")
	}

	if input.StopLossPrice != nil {
		trade.StopLossPrice = input.StopLossPrice
	}
	if input.TakeProfitPrice != nil {
		trade.TakeProfitPrice = input.TakeProfitPrice
	}
	if input.Action == "CLOSE" {
		if input.ClosePrice != nil {
			trade.ClosePrice = input.ClosePrice
		}
		now := time.Now()
		trade.Status = models.TradeStatusClosed
		trade.ClosedAt = models.TimePtr(now)
	}
	if input.Action == "CANCEL" {
		trade.Status = models.TradeStatusCancelled
	}

	if input.Status != "" {
		trade.Status = input.Status
	}

	if err := s.repo.UpdateTrade(trade); err != nil {
		return nil, err
	}
	return trade, nil
}

func (s *tradeService) DeleteTrade(id uint, traderID uint) error {
	return s.repo.DeleteTrade(id, traderID)
}
