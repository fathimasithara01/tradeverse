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
	// When creating a signal, ensure its initial status is Pending or Active based on tradeStartDate
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

	signals, err := s.signalRepo.GetAllSignals(ctx) // Get all signals, not just active/pending, to update current price for all
	if err != nil {
		log.Printf("ERROR: Failed to get all signals for current price update in service: %v", err)
		return fmt.Errorf("failed to get all signals for current price update: %w", err)
	}
	log.Printf("Found %d signals to potentially update current prices.", len(signals))

	for _, signal := range signals {
		if signal.Status == "Target Hit" || signal.Status == "Stop Loss" {
			log.Printf("Skipping current price update for signal ID %d (by %s, %s) with status '%s'", signal.ID, signal.TraderName, signal.Symbol, signal.Status)
			continue // No need to update current price for finished signals
		}
		if signal.EntryPrice == 0 && signal.TargetPrice == 0 && signal.StopLoss == 0 {
			log.Printf("Skipping current price update for signal ID %d (by %s, %s) because all prices are zero. Signal not properly configured?", signal.ID, signal.TraderName, signal.Symbol)
			continue
		}

		marketData, err := s.signalRepo.GetMarketDataBySymbol(ctx, signal.Symbol)
		if err != nil {
			log.Printf("Warning: Could not fetch market data for symbol %s (from signal ID %d, trader %s): %v", signal.Symbol, signal.ID, signal.TraderName, err)
			continue
		}
		if marketData == nil {
			log.Printf("Info: No market data found for symbol %s (from signal ID %d, trader %s) in DB, skipping current price update for this signal.", signal.Symbol, signal.ID, signal.TraderName)
			continue
		}

		// Only update if the price has actually changed to avoid unnecessary DB writes
		if signal.CurrentPrice != marketData.CurrentPrice {
			log.Printf("Updating signal ID %d (by %s, %s) current price from %.4f to %.4f", signal.ID, signal.TraderName, signal.Symbol, signal.CurrentPrice, marketData.CurrentPrice)
			err := s.signalRepo.UpdateSignalCurrentPrice(ctx, signal.ID, marketData.CurrentPrice)
			if err != nil {
				log.Printf("Error updating current price for signal ID %d: %v", signal.ID, err)
			}
		} else {
			log.Printf("Signal ID %d (by %s, %s) current price %.4f is already up-to-date.", signal.ID, signal.TraderName, signal.Symbol, signal.CurrentPrice)
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
