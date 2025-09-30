package service

import (
	"context"
	"fmt"
	"log"
	"time" // Added for time comparisons

	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type ILiveSignalService interface {
	CreateSignal(ctx context.Context, signal *models.Signal) (*models.Signal, error)
	GetAllSignals(ctx context.Context) ([]models.Signal, error)
	// Renamed for clarity:
	UpdateAllSignalsCurrentPrices(ctx context.Context) error
	// NEW: Function to check and set signal statuses
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
	return s.signalRepo.CreateSignal(ctx, signal)
}

func (s *liveSignalService) GetAllSignals(ctx context.Context) ([]models.Signal, error) {
	return s.signalRepo.GetAllSignals(ctx)
}

// UpdateAllSignalsCurrentPrices fetches market data and updates the `CurrentPrice` for all signals.
func (s *liveSignalService) UpdateAllSignalsCurrentPrices(ctx context.Context) error {
	log.Println("Starting to update current prices for all signals...")

	signals, err := s.signalRepo.GetAllSignals(ctx) // Get all signals, not just active/pending, to update current price for all
	if err != nil {
		return fmt.Errorf("failed to get all signals for current price update: %w", err)
	}

	for _, signal := range signals {
		// Only update current price if the signal is not yet "Target Hit" or "Stop Loss"
		// If you want to update CurrentPrice even for finished signals, remove this check.
		// For display purposes, it might be fine to keep updating for a while.
		// For status checks, we'll only look at "Active" or "Pending".
		if signal.Status == "Target Hit" || signal.Status == "Stop Loss" {
			continue // No need to update current price for finished signals
		}

		marketData, err := s.signalRepo.GetMarketDataBySymbol(ctx, signal.Symbol)
		if err != nil {
			log.Printf("Warning: Could not fetch market data for symbol %s (from signal ID %d, trader %s): %v", signal.Symbol, signal.ID, signal.TraderName, err)
			continue
		}
		if marketData == nil {
			log.Printf("Info: No market data found for symbol %s (from signal ID %d, trader %s), skipping current price update.", signal.Symbol, signal.ID, signal.TraderName)
			continue
		}

		// Only update if the price has actually changed to avoid unnecessary DB writes
		if signal.CurrentPrice != marketData.CurrentPrice {
			log.Printf("Updating signal ID %d (by %s, %s) current price from %.4f to %.4f", signal.ID, signal.TraderName, signal.Symbol, signal.CurrentPrice, marketData.CurrentPrice)
			err := s.signalRepo.UpdateSignalCurrentPrice(ctx, signal.ID, marketData.CurrentPrice)
			if err != nil {
				log.Printf("Error updating current price for signal ID %d: %v", signal.ID, err)
			}
		}
	}
	log.Println("Finished updating current prices for all signals.")
	return nil
}

// NEW: CheckAndSetSignalStatuses reviews active/pending signals and updates their status if SL/Target is hit, or if it becomes active.
func (s *liveSignalService) CheckAndSetSignalStatuses(ctx context.Context) error {
	log.Println("Starting signal status check (SL/Target/Activation)...")

	signals, err := s.signalRepo.GetActiveAndPendingSignals(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active/pending signals for status check: %w", err)
	}

	for _, signal := range signals {
		// 1. Check for Pending to Active transition
		if signal.Status == "Pending" && !signal.TradeStartDate.After(time.Now()) {
			log.Printf("Signal ID %d (by %s, %s) is now Active.", signal.ID, signal.TraderName, signal.Symbol)
			err := s.signalRepo.UpdateSignalStatus(ctx, signal.ID, "Active")
			if err != nil {
				log.Printf("Error setting signal ID %d to Active: %v", signal.ID, err)
			}
			// Continue to next signal, or re-evaluate with "Active" status in the same loop if needed
			// For simplicity, we'll let the next cron run pick it up as active for SL/Target checks
			continue
		}

		// Only process 'Active' signals for SL/Target hits
		if signal.Status != "Active" {
			continue
		}

		// Ensure CurrentPrice is available (it should be updated by the other cron job)
		if signal.CurrentPrice == 0 {
			log.Printf("Warning: Signal ID %d (%s) has zero current price, skipping SL/Target check.", signal.ID, signal.Symbol)
			continue
		}

		// 2. Check for Stop Loss
		// Assuming for a BUY signal, SL is below entry, Target is above entry.
		// For a SELL signal, SL is above entry, Target is below entry.
		// Your current model doesn't explicitly state buy/sell, so we'll assume a "long" position where CurrentPrice > EntryPrice is profit.
		// You might need to add a `TradeType` field (e.g., "Long", "Short") to the Signal model for more precise logic.
		// For now, let's assume a "long" bias (buy low, sell high).

		if signal.CurrentPrice <= signal.StopLoss {
			log.Printf("Signal ID %d (by %s, %s) hit Stop Loss at %.4f (SL: %.4f).", signal.ID, signal.TraderName, signal.Symbol, signal.CurrentPrice, signal.StopLoss)
			err := s.signalRepo.UpdateSignalStatus(ctx, signal.ID, "Stop Loss")
			if err != nil {
				log.Printf("Error setting signal ID %d to Stop Loss: %v", signal.ID, err)
			}
			// Add logic here to trigger further actions, e.g., send notifications, close copy trades, etc.
			continue // This signal is done, move to the next
		}

		// 3. Check for Target Hit
		if signal.CurrentPrice >= signal.TargetPrice {
			log.Printf("Signal ID %d (by %s, %s) hit Target at %.4f (Target: %.4f).", signal.ID, signal.TraderName, signal.Symbol, signal.CurrentPrice, signal.TargetPrice)
			err := s.signalRepo.UpdateSignalStatus(ctx, signal.ID, "Target Hit")
			if err != nil {
				log.Printf("Error setting signal ID %d to Target Hit: %v", signal.ID, err)
			}
			// Add logic here to trigger further actions, e.g., send notifications, close copy trades, etc.
			continue // This signal is done, move to the next
		}
	}
	log.Println("Finished signal status check.")
	return nil
}
