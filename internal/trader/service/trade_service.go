package service

import (
	"context" // Ensure context is imported
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/fathimasithara01/tradeverse/pkg/utils/constants"
)

type TradeService interface {
	CreateTrade(ctx context.Context, traderID uint, input models.TradeInput) (*models.Trade, error)
	GetTraderTrades(ctx context.Context, traderID uint, page, limit int) (*models.TradeListResponse, error)
	UpdateTradeStatus(ctx context.Context, traderID, tradeID uint, input models.TradeUpdateInput) (*models.Trade, error)
	RemoveTrade(ctx context.Context, traderID, tradeID uint) error
}

type tradeService struct {
	tradeRepo repository.TradeRepository
}

func NewTradeService(tradeRepo repository.TradeRepository) TradeService {
	return &tradeService{
		tradeRepo: tradeRepo,
	}
}

func (s *tradeService) CreateTrade(ctx context.Context, traderID uint, input models.TradeInput) (*models.Trade, error) {
	// Basic validation for trade input
	if input.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	if input.TradeType == models.TradeTypeLimit || input.TradeType == models.TradeTypeStop {
		if input.EntryPrice <= 0 {
			return nil, errors.New("entry price is required for LIMIT and STOP orders")
		}
	} else if input.TradeType == models.TradeTypeMarket {
		input.EntryPrice = 0 // Or handle appropriately based on your market order logic
	}

	trade := &models.Trade{
		TraderID:        traderID,
		Symbol:          input.Symbol,
		TradeType:       input.TradeType,
		Side:            input.Side,
		EntryPrice:      input.EntryPrice,
		ExecutedPrice:   nil, // Market orders might not have this until execution
		Quantity:        input.Quantity,
		Leverage:        input.Leverage,
		StopLossPrice:   input.StopLossPrice,
		TakeProfitPrice: input.TakeProfitPrice,
		Status:          models.TradeStatusPending, // Initially pending
		OpenedAt:        models.TimePtr(time.Now()),
	}

	err := s.tradeRepo.CreateTrade(ctx, trade) // Pass ctx here
	if err != nil {
		return nil, fmt.Errorf("failed to create trade: %w", err)
	}

	return trade, nil
}

func (s *tradeService) GetTraderTrades(ctx context.Context, traderID uint, page, limit int) (*models.TradeListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // Default limit
	}
	offset := (page - 1) * limit

	// Fix: Call GetTradesByTraderID, passing ctx correctly
	trades, total, err := s.tradeRepo.GetTradesByTraderID(ctx, traderID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades for trader %d: %w", traderID, err)
	}

	return &models.TradeListResponse{
		Trades: trades,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

func (s *tradeService) UpdateTradeStatus(ctx context.Context, traderID, tradeID uint, input models.TradeUpdateInput) (*models.Trade, error) {
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve trade: %w", err)
	}
	if trade == nil {
		return nil, constants.ErrNotFound("Trade")
	}
	if trade.TraderID != traderID {
		return nil, constants.ErrForbidden
	}

	if input.Action == "CLOSE" {
		if input.ClosePrice == nil || *input.ClosePrice <= 0 {
			return nil, errors.New("close price is required to close a trade")
		}
		// Fix: Pass ctx to CloseTrade
		return s.tradeRepo.CloseTrade(ctx, tradeID, *input.ClosePrice)
	} else if input.Action == "CANCEL" {
		// Fix: Pass ctx to CancelTrade
		return s.tradeRepo.CancelTrade(ctx, tradeID)
	}

	// Handle general status updates (e.g., if an admin or automated system changes it)
	if input.Status != "" {
		if trade.Status == models.TradeStatusOpen || trade.Status == models.TradeStatusPending {
			if input.StopLossPrice != nil {
				trade.StopLossPrice = input.StopLossPrice
			}
			if input.TakeProfitPrice != nil {
				trade.TakeProfitPrice = input.TakeProfitPrice
			}
			if input.Status != "" {
				trade.Status = input.Status
			}
			// Fix: Pass ctx to UpdateTrade
			err = s.tradeRepo.UpdateTrade(ctx, trade)
			if err != nil {
				return nil, fmt.Errorf("failed to update trade: %w", err)
			}
			return trade, nil
		} else {
			return nil, errors.New("cannot modify a trade that is not open or pending")
		}
	}

	// If no action or status update, but SL/TP might have been updated
	if input.StopLossPrice != nil || input.TakeProfitPrice != nil {
		if trade.Status == models.TradeStatusOpen || trade.Status == models.TradeStatusPending {
			if input.StopLossPrice != nil {
				trade.StopLossPrice = input.StopLossPrice
			}
			if input.TakeProfitPrice != nil {
				trade.TakeProfitPrice = input.TakeProfitPrice
			}
			err = s.tradeRepo.UpdateTrade(ctx, trade)
			if err != nil {
				return nil, fmt.Errorf("failed to update trade: %w", err)
			}
			return trade, nil
		}
	}

	return nil, errors.New("no valid update action or status provided")
}

// RemoveTrade handles the deletion of a trade. This might be restricted to PENDING trades only.
func (s *tradeService) RemoveTrade(ctx context.Context, traderID, tradeID uint) error {
	// Fix: Pass ctx to GetTradeByID
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("failed to retrieve trade: %w", err)
	}
	if trade == nil {
		return constants.ErrNotFound("Trade")
	}
	if trade.TraderID != traderID {
		return constants.ErrForbidden
	}
	if trade.Status != models.TradeStatusPending {
		return errors.New("only pending trades can be removed")
	}

	// Fix: Pass ctx to DeleteTrade
	err = s.tradeRepo.DeleteTrade(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("failed to remove trade: %w", err)
	}
	return nil
}
