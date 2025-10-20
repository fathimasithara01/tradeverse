package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type ILiveSignalService interface {
	CreateSignal(ctx context.Context, signal *models.Signal) (*models.Signal, error)
	GetAllSignals(ctx context.Context) ([]models.Signal, error)
	UpdateAllSignalsCurrentPrices(ctx context.Context) error
	CheckAndSetSignalStatuses(ctx context.Context) error
}

type liveSignalService struct {
	signalRepo repository.ISignalRepository
}

func NewLiveSignalService(signalRepo repository.ISignalRepository) ILiveSignalService {
	return &liveSignalService{signalRepo: signalRepo}
}

func (s *liveSignalService) CreateSignal(ctx context.Context, signal *models.Signal) (*models.Signal, error) {
	if signal.TradeStartDate.After(time.Now()) {
		signal.Status = "Pending"
	} else {
		signal.Status = "Active"
	}
	log.Printf("Creating signal: Symbol=%s, Trader=%s, Entry=%.4f, Target=%.4f, SL=%.4f, InitialStatus=%s",
		signal.Symbol, signal.TraderName, signal.EntryPrice, signal.TargetPrice, signal.StopLoss, signal.Status)
	return s.signalRepo.CreateSignal(ctx, signal)
}

func (s *liveSignalService) GetAllSignals(ctx context.Context) ([]models.Signal, error) {
	signals, err := s.signalRepo.GetAllSignals(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to get all signals in GetAllSignals service: %v", err)
		return nil, fmt.Errorf("failed to get all signals: %w", err)
	}
	log.Printf("Retrieved %d signals from DB for GetAllSignals.", len(signals))
	return signals, nil
}

func (s *liveSignalService) UpdateAllSignalsCurrentPrices(ctx context.Context) error {
	log.Println("Starting to update current prices for all signals...")

	signals, err := s.signalRepo.GetAllSignals(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all signals: %w", err)
	}

	for _, signal := range signals {
		if signal.Status == "Target Hit" || signal.Status == "Stop Loss" {
			continue
		}

		if signal.Symbol == "" {
			log.Printf("Signal ID %d has empty symbol, skipping.", signal.ID)
			continue
		}

		marketData, err := s.signalRepo.GetMarketDataBySymbol(ctx, signal.Symbol)
		if err != nil {
			log.Printf("Error fetching market data for %s: %v", signal.Symbol, err)
			continue
		}
		if marketData == nil {
			log.Printf("No market data found for symbol %s (signal ID %d), skipping.", signal.Symbol, signal.ID)
			continue
		}

		if signal.CurrentPrice != marketData.CurrentPrice {
			err := s.signalRepo.UpdateSignalCurrentPrice(ctx, signal.ID, marketData.CurrentPrice)
			if err != nil {
				log.Printf("Error updating current price for signal ID %d: %v", signal.ID, err)
			} else {
				log.Printf("Updated signal ID %d current price: %.4f -> %.4f", signal.ID, signal.CurrentPrice, marketData.CurrentPrice)
			}
		}
	}

	log.Println("Finished updating current prices for all signals.")
	return nil
}

func (s *liveSignalService) CheckAndSetSignalStatuses(ctx context.Context) error {
	log.Println("Starting signal status check (SL/Target/Activation)...")

	signals, err := s.signalRepo.GetActiveAndPendingSignals(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to get active/pending signals for status check in service: %v", err)
		return fmt.Errorf("failed to get active/pending signals for status check: %w", err)
	}
	log.Printf("Found %d active/pending signals to check status.", len(signals))

	for _, signal := range signals {
		log.Printf("Checking signal ID %d (Symbol: %s, Status: %s, Current: %.4f, Entry: %.4f, Target: %.4f, SL: %.4f)",
			signal.ID, signal.Symbol, signal.Status, signal.CurrentPrice, signal.EntryPrice, signal.TargetPrice, signal.StopLoss)

		// 1. Check for Pending to Active transition
		if signal.Status == "Pending" && !signal.TradeStartDate.After(time.Now()) {
			log.Printf("Signal ID %d (by %s, %s) is now Active.", signal.ID, signal.TraderName, signal.Symbol)
			err := s.signalRepo.UpdateSignalStatus(ctx, signal.ID, "Active")
			if err != nil {
				log.Printf("Error setting signal ID %d to Active: %v", signal.ID, err)
			}
			continue
		}

		if signal.Status != "Active" {
			log.Printf("Signal ID %d (Symbol: %s) is not Active, skipping SL/Target check. Current Status: %s", signal.ID, signal.Symbol, signal.Status)
			continue
		}

		if signal.CurrentPrice == 0 {
			log.Printf("Warning: Signal ID %d (%s) has zero current price, skipping SL/Target check.", signal.ID, signal.Symbol)
			continue
		}
		if signal.TargetPrice == 0 && signal.StopLoss == 0 {
			log.Printf("Warning: Signal ID %d (%s) has zero TargetPrice and StopLoss, skipping SL/Target check.", signal.ID, signal.Symbol)
			continue
		}

		// 2. Check for Stop Loss Hit
		if signal.StopLoss != 0 && signal.CurrentPrice <= signal.StopLoss {
			log.Printf("Signal ID %d (by %s, %s) hit Stop Loss at %.4f (SL: %.4f). Updating status.", signal.ID, signal.TraderName, signal.Symbol, signal.CurrentPrice, signal.StopLoss)
			err := s.signalRepo.UpdateSignalStatus(ctx, signal.ID, "Stop Loss")
			if err != nil {
				log.Printf("Error setting signal ID %d to Stop Loss: %v", signal.ID, err)
			}
			continue
		}

		// 3. Check for Target Hit
		if signal.TargetPrice != 0 && signal.CurrentPrice >= signal.TargetPrice {
			log.Printf("Signal ID %d (by %s, %s) hit Target at %.4f (Target: %.4f). Updating status.", signal.ID, signal.TraderName, signal.Symbol, signal.CurrentPrice, signal.TargetPrice)
			err := s.signalRepo.UpdateSignalStatus(ctx, signal.ID, "Target Hit")
			if err != nil {
				log.Printf("Error setting signal ID %d to Target Hit: %v", signal.ID, err)
			}
			continue
		}
	}
	log.Println("Finished signal status check.")
	return nil
}
