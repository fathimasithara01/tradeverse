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
	// Add GetMarketDataBySymbol to the interface as it's used by the controller
	GetMarketDataBySymbol(ctx context.Context, symbol string) (*models.MarketData, error)
}

type liveSignalService struct {
	signalRepo repository.ISignalRepository
}

func NewLiveSignalService(signalRepo repository.ISignalRepository) ILiveSignalService {
	return &liveSignalService{signalRepo: signalRepo}
}

func (s *liveSignalService) CreateSignal(ctx context.Context, signal *models.Signal) (*models.Signal, error) {
	// The service layer is the authoritative place to set the initial status
	if signal.TradeStartDate.After(time.Now()) {
		signal.Status = "Pending"
	} else {
		signal.Status = "Active"
	}
	log.Printf("Attempting to create signal: Symbol=%s, Trader=%s, Entry=%.4f, Target=%.4f, SL=%.4f, InitialStatus=%s, InitialCurrentPrice=%.4f",
		signal.Symbol, signal.TraderName, signal.EntryPrice, signal.TargetPrice, signal.StopLoss, signal.Status, signal.CurrentPrice)
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

// GetMarketDataBySymbol implementation for the service layer
func (s *liveSignalService) GetMarketDataBySymbol(ctx context.Context, symbol string) (*models.MarketData, error) {
	return s.signalRepo.GetMarketDataBySymbol(ctx, symbol)
}

func (s *liveSignalService) UpdateAllSignalsCurrentPrices(ctx context.Context) error {
	log.Println("Starting to update current prices for all signals...")

	signals, err := s.signalRepo.GetAllSignals(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all signals for price update: %w", err)
	}

	for _, signal := range signals {
		// Skip signals that have reached their target or stop loss
		if signal.Status == "Target Hit" || signal.Status == "Stop Loss" {
			continue
		}

		if signal.Symbol == "" {
			log.Printf("Warning: Signal ID %d has empty symbol, skipping price update.", signal.ID)
			continue
		}

		marketData, err := s.signalRepo.GetMarketDataBySymbol(ctx, signal.Symbol)
		if err != nil {
			log.Printf("Error fetching market data for %s (Signal ID %d): %v", signal.Symbol, signal.ID, err)
			continue // Continue to next signal on error
		}
		if marketData == nil {
			log.Printf("No market data found for symbol %s (Signal ID %d), skipping price update.", signal.Symbol, signal.ID)
			continue
		}

		// Only update if the current price has actually changed to avoid unnecessary DB writes
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

	// Only fetch signals that are 'Active' or 'Pending' to optimize checks
	signals, err := s.signalRepo.GetActiveAndPendingSignals(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to get active/pending signals for status check in service: %v", err)
		return fmt.Errorf("failed to get active/pending signals for status check: %w", err)
	}
	log.Printf("Found %d active/pending signals to check status.", len(signals))

	for _, signal := range signals {
		log.Printf("Checking signal ID %d (Symbol: %s, Current Status: %s, Current Price: %.4f, Entry: %.4f, Target: %.4f, SL: %.4f)",
			signal.ID, signal.Symbol, signal.Status, signal.CurrentPrice, signal.EntryPrice, signal.TargetPrice, signal.StopLoss)

		// 1. Check for Pending to Active transition
		if signal.Status == "Pending" && !signal.TradeStartDate.After(time.Now()) {
			log.Printf("Signal ID %d (by %s, %s) TradeStartDate has passed. Transitioning to Active.", signal.ID, signal.TraderName, signal.Symbol)
			err := s.signalRepo.UpdateSignalStatus(ctx, signal.ID, "Active")
			if err != nil {
				log.Printf("Error setting signal ID %d to Active: %v", signal.ID, err)
			}
			// Important: Continue to next signal, as this one just became active and its price/SL/Target might
			// be checked in the *next* cycle after current price update, or immediately if currentPrice is already valid.
			// For simplicity, we let the next cycle handle SL/Target if it was Pending.
			continue
		}

		// 2. Only check Stop Loss/Target for *Active* signals
		if signal.Status != "Active" {
			log.Printf("Signal ID %d (Symbol: %s) is not Active (Status: %s), skipping SL/Target check.", signal.ID, signal.Symbol, signal.Status)
			continue
		}

		// Basic sanity checks before comparing prices
		if signal.CurrentPrice == 0 {
			log.Printf("Warning: Signal ID %d (%s) has zero current price, skipping SL/Target check. Ensure market data is updating.", signal.ID, signal.Symbol)
			continue
		}
		if signal.TargetPrice == 0 && signal.StopLoss == 0 {
			log.Printf("Warning: Signal ID %d (%s) has zero TargetPrice and StopLoss, skipping SL/Target check. Please review signal data.", signal.ID, signal.Symbol)
			continue
		}

		// 3. Check for Stop Loss hit
		// Only check if StopLoss is defined (not zero)
		if signal.StopLoss != 0 && signal.CurrentPrice <= signal.StopLoss {
			log.Printf("Signal ID %d (by %s, %s) hit Stop Loss at %.4f (SL: %.4f). Updating status.", signal.ID, signal.TraderName, signal.Symbol, signal.CurrentPrice, signal.StopLoss)
			err := s.signalRepo.UpdateSignalStatus(ctx, signal.ID, "Stop Loss")
			if err != nil {
				log.Printf("Error setting signal ID %d to Stop Loss: %v", signal.ID, err)
			}
			continue // Once hit, move to the next signal
		}

		// 4. Check for Target Hit
		// Only check if TargetPrice is defined (not zero)
		if signal.TargetPrice != 0 && signal.CurrentPrice >= signal.TargetPrice {
			log.Printf("Signal ID %d (by %s, %s) hit Target at %.4f (Target: %.4f). Updating status.", signal.ID, signal.TraderName, signal.Symbol, signal.CurrentPrice, signal.TargetPrice)
			err := s.signalRepo.UpdateSignalStatus(ctx, signal.ID, "Target Hit")
			if err != nil {
				log.Printf("Error setting signal ID %d to Target Hit: %v", signal.ID, err)
			}
			continue // Once hit, move to the next signal
		}
	}
	log.Println("Finished signal status check.")
	return nil
}
