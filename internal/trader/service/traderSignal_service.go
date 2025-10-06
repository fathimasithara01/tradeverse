package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type ISignalService interface {
	CreateSignal(ctx context.Context, signal *models.Signal) (*models.Signal, error)
	GetAllSignals(ctx context.Context) ([]models.Signal, error)
	UpdatePendingSignalsCurrentPrice(ctx context.Context) error
	UpdateActiveSignalStatuses(ctx context.Context) error
	GetSignalByID(ctx context.Context, id uint) (*models.Signal, error)
	UpdateSignal(ctx context.Context, updated *models.Signal) (*models.Signal, error)
	DeleteSignal(ctx context.Context, id uint) error
}

type SignalService struct {
	repo repository.ISignalRepository
}

func NewSignalService(repo repository.ISignalRepository) ISignalService {
	return &SignalService{repo: repo}
}

func (s *SignalService) CreateSignal(ctx context.Context, signal *models.Signal) (*models.Signal, error) {
	if len(signal.Symbol) > 0 && !strings.HasSuffix(signal.Symbol, "USDT") {
		signal.Symbol = strings.ToUpper(signal.Symbol) + "USDT"
	} else {
		signal.Symbol = strings.ToUpper(signal.Symbol)
	}

	signal.Status = "Pending"

	return s.repo.CreateSignal(ctx, signal)
}

func (s *SignalService) GetAllSignals(ctx context.Context) ([]models.Signal, error) {
	return s.repo.GetAllSignals(ctx)
}

func (s *SignalService) UpdatePendingSignalsCurrentPrice(ctx context.Context) error {
	pendingSignals, err := s.repo.GetPendingSignals(ctx)
	if err != nil {
		return err
	}

	for _, signal := range pendingSignals {
		md, err := s.repo.GetMarketDataBySymbol(ctx, signal.Symbol)
		if err != nil || md == nil {
			log.Printf("No market data for %s", signal.Symbol)
			continue
		}

		if md.CurrentPrice >= signal.EntryPrice {
			err := s.repo.UpdateSignalStatus(ctx, signal.ID, "Active")
			if err != nil {
				log.Printf("Failed to activate signal %d: %v", signal.ID, err)
			} else {
				log.Printf("Signal %d is now Active", signal.ID)
			}
		}

		err = s.repo.UpdateSignalCurrentPrice(ctx, signal.ID, md.CurrentPrice)
		if err != nil {
			log.Printf("Failed to update current price for signal %d: %v", signal.ID, err)
		}
	}

	return nil
}

func (s *SignalService) UpdateActiveSignalStatuses(ctx context.Context) error {
	activeSignals, err := s.repo.GetActiveSignals(ctx)
	if err != nil {
		return err
	}

	for _, signal := range activeSignals {
		md, err := s.repo.GetMarketDataBySymbol(ctx, signal.Symbol)
		if err != nil || md == nil {
			continue
		}

		_ = s.repo.UpdateSignalCurrentPrice(ctx, signal.ID, md.CurrentPrice)

		if md.CurrentPrice <= signal.StopLoss {
			_ = s.repo.UpdateSignalStatus(ctx, signal.ID, "Stop Loss")
			log.Printf("Signal %d hit Stop Loss", signal.ID)
			continue
		}

		if md.CurrentPrice >= signal.TargetPrice {
			_ = s.repo.UpdateSignalStatus(ctx, signal.ID, "Target Hit")
			log.Printf("Signal %d hit Target Price", signal.ID)
			continue
		}
	}

	return nil
}

func (s *SignalService) GetSignalByID(ctx context.Context, id uint) (*models.Signal, error) {
	return s.repo.GetSignalByID(ctx, id)
}

func (s *SignalService) UpdateSignal(ctx context.Context, updated *models.Signal) (*models.Signal, error) {
	existing, err := s.repo.GetSignalByID(ctx, updated.ID)
	if err != nil {
		return nil, fmt.Errorf("signal not found: %w", err)
	}

	existing.EntryPrice = updated.EntryPrice
	existing.StopLoss = updated.StopLoss
	existing.TargetPrice = updated.TargetPrice
	existing.Strategy = updated.Strategy
	existing.Risk = updated.Risk
	existing.Symbol = strings.ToUpper(updated.Symbol)

	if err := s.repo.UpdateSignal(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update signal: %w", err)
	}
	return existing, nil
}

func (s *SignalService) DeleteSignal(ctx context.Context, id uint) error {
	return s.repo.DeleteSignal(ctx, id)
}
