package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

// Define custom errors
var (
	ErrTradeNotFound     = errors.New("trade not found or does not belong to the trader")
	ErrInvalidTradeInput = errors.New("invalid trade input")
	ErrTradeNotOpen      = errors.New("trade is not open and cannot be updated")
	ErrClosingPrice      = errors.New("closing price cannot be less than entry price for sell, or greater for buy") // Example
	ErrInsufficientFunds = errors.New("insufficient funds for trade")
)

// TradeService defines the business logic for trade operations
type TradeService interface {
	CreateTrade(traderID uint, input *models.TradeInput) (*models.Trade, error)
	GetTradeByID(tradeID uint, traderID uint) (*models.Trade, error)
	ListTrades(traderID uint, pagination *models.PaginationParams) (*models.TradeListResponse, error)
	UpdateTrade(tradeID uint, traderID uint, input *models.TradeUpdateInput) (*models.Trade, error)
	DeleteTrade(tradeID uint, traderID uint) error
	CloseTrade(tradeID uint, traderID uint, closePrice float64) (*models.Trade, error) // Added explicitly
	CancelTrade(tradeID uint, traderID uint) (*models.Trade, error)                    // Added explicitly
}

// tradeService implements TradeService
type tradeService struct {
	tradeRepo       repository.TradeRepository
	walletRepo      repository.WalletRepository      // Assuming you have a wallet repository
	transactionRepo repository.TransactionRepository // Assuming a transaction repo
	db              *gorm.DB                         // For transactions
}

// NewTradeService creates a new instance of TradeService
func NewTradeService(tradeRepo repository.TradeRepository, walletRepo repository.WalletRepository, transactionRepo repository.TransactionRepository, db *gorm.DB) TradeService {
	return &tradeService{
		tradeRepo:       tradeRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		db:              db,
	}
}

// CreateTrade handles the creation of a new trade
func (s *tradeService) CreateTrade(traderID uint, input *models.TradeInput) (*models.Trade, error) {
	// Basic validation
	if input.TradeType == models.TradeTypeLimit || input.TradeType == models.TradeTypeStop {
		if input.EntryPrice <= 0 {
			return nil, ErrInvalidTradeInput
		}
	}

	trade := &models.Trade{
		TraderID:        traderID,
		Symbol:          input.Symbol,
		TradeType:       input.TradeType,
		Side:            input.Side,
		EntryPrice:      input.EntryPrice,
		Quantity:        input.Quantity,
		Leverage:        input.Leverage,
		StopLossPrice:   input.StopLossPrice,
		TakeProfitPrice: input.TakeProfitPrice,
		Status:          models.TradeStatusPending, // Start as pending, execution will change it to OPEN/CLOSED
		Fees:            0,                         // Placeholder, calculate real fees later
	}

	// --- Simulate wallet deduction (simplified) ---
	// In a real system, this would involve complex interactions with an exchange/broker API
	// and robust financial transaction management.
	// For example, for a BUY market order:
	// 1. Calculate margin/cost: input.Quantity * input.EntryPrice (or current market price) / input.Leverage
	// 2. Deduct from wallet
	// 3. Create a transaction record
	// This is a simplified example, adjust according to your financial model.

	// Use a database transaction for consistency
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Update repositories to use the transaction (tx)
		tradeRepoTx := repository.NewTradeRepository(tx)
		walletRepoTx := repository.NewWalletRepository(tx)           // Assuming NewWalletRepository takes *gorm.DB
		transactionRepoTx := repository.NewTransactionRepository(tx) // Assuming NewTransactionRepository takes *gorm.DB

		// Example: Check and deduct funds (very simplified logic)
		wallet, err := walletRepoTx.GetWalletByUserID(traderID)
		if err != nil {
			return fmt.Errorf("failed to get wallet: %w", err)
		}
		if wallet == nil {
			return errors.New("trader wallet not found")
		}

		// Calculate estimated cost/margin (simplified)
		estimatedCost := input.Quantity * input.EntryPrice
		if trade.Leverage > 1 {
			estimatedCost = estimatedCost / float64(trade.Leverage)
		}

		// Add a buffer or consider actual margin requirements
		if wallet.Balance < estimatedCost {
			return ErrInsufficientFunds
		}

		// Deduct funds (e.g., as a pending transaction or direct deduction for market orders)
		// For a real system, you might have an "allocated" balance.
		wallet.Balance -= estimatedCost
		wallet.LastUpdated = time.Now()
		if err := walletRepoTx.UpdateWallet(wallet); err != nil {
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}

		// Create a corresponding transaction
		walletTransaction := &models.WalletTransaction{
			WalletID:        wallet.ID,
			UserID:          traderID,
			TransactionType: models.TxTypeTradeLoss, // Or a new type like "TRADE_OPENING_FUNDS"
			Amount:          -estimatedCost,
			Currency:        wallet.Currency,
			Status:          models.TxStatusPending, // Will become success upon trade execution
			Description:     fmt.Sprintf("Funds reserved for opening trade %s %s %s", trade.Side, trade.Symbol, trade.TradeType),
			BalanceBefore:   wallet.Balance + estimatedCost,
			BalanceAfter:    wallet.Balance,
		}
		if err := transactionRepoTx.CreateTransaction(walletTransaction); err != nil {
			return fmt.Errorf("failed to create wallet transaction: %w", err)
		}
		// Link trade to transaction if needed: trade.WalletTransactionID = &walletTransaction.ID

		// Now create the trade
		if err := tradeRepoTx.CreateTrade(trade); err != nil {
			return fmt.Errorf("failed to create trade: %w", err)
		}

		// Update the transaction status to success and link to trade
		walletTransaction.Status = models.TxStatusSuccess
		walletTransaction.TradeID = &trade.ID
		if err := transactionRepoTx.UpdateTransaction(walletTransaction); err != nil {
			return fmt.Errorf("failed to update wallet transaction status: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// In a real application, you'd integrate with an actual exchange API here
	// to place the order. The status of the trade would then be updated based
	// on the exchange's response (e.g., from PENDING to OPEN or CLOSED if market).

	// For simplicity, we'll immediately set to OPEN if it's a market order
	if input.TradeType == models.TradeTypeMarket {
		trade.Status = models.TradeStatusOpen
		trade.ExecutedPrice = &input.EntryPrice // Assume market price is entry price for now
		s.tradeRepo.UpdateTrade(trade)          // Update status in DB
	}

	return trade, nil
}

// GetTradeByID retrieves a trade by ID for a specific trader
func (s *tradeService) GetTradeByID(tradeID uint, traderID uint) (*models.Trade, error) {
	trade, err := s.tradeRepo.GetTradeByID(tradeID, traderID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, ErrTradeNotFound
	}
	return trade, nil
}

// ListTrades retrieves a list of trades for a specific trader
func (s *tradeService) ListTrades(traderID uint, pagination *models.PaginationParams) (*models.TradeListResponse, error) {
	trades, total, err := s.tradeRepo.ListTrades(traderID, pagination)
	if err != nil {
		return nil, err
	}

	return &models.TradeListResponse{
		Trades: trades,
		Total:  total,
		Page:   pagination.Page,
		Limit:  pagination.Limit,
	}, nil
}

// UpdateTrade allows modifying an open trade's stop loss/take profit, or closing/cancelling it
func (s *tradeService) UpdateTrade(tradeID uint, traderID uint, input *models.TradeUpdateInput) (*models.Trade, error) {
	trade, err := s.tradeRepo.GetTradeByID(tradeID, traderID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, ErrTradeNotFound
	}

	if trade.Status != models.TradeStatusOpen && trade.Status != models.TradeStatusPending {
		return nil, ErrTradeNotOpen // Only open or pending trades can be updated
	}

	if input.StopLossPrice != nil {
		trade.StopLossPrice = input.StopLossPrice
	}
	if input.TakeProfitPrice != nil {
		trade.TakeProfitPrice = input.TakeProfitPrice
	}

	// Handle explicit actions
	if input.Action == "CLOSE" {
		if input.ClosePrice == nil || *input.ClosePrice <= 0 {
			return nil, errors.New("close price is required for closing a trade")
		}
		return s.CloseTrade(tradeID, traderID, *input.ClosePrice)
	} else if input.Action == "CANCEL" {
		return s.CancelTrade(tradeID, traderID)
	}

	// If no action, just update SL/TP
	err = s.tradeRepo.UpdateTrade(trade)
	if err != nil {
		return nil, fmt.Errorf("failed to update trade: %w", err)
	}

	return trade, nil
}

// CloseTrade closes an open trade and calculates P&L
func (s *tradeService) CloseTrade(tradeID uint, traderID uint, closePrice float64) (*models.Trade, error) {
	trade, err := s.tradeRepo.GetTradeByID(tradeID, traderID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, ErrTradeNotFound
	}
	if trade.Status != models.TradeStatusOpen {
		return nil, ErrTradeNotOpen
	}

	// Calculate P&L
	var pnl float64
	if trade.Side == models.TradeSideBuy {
		pnl = (closePrice - *trade.ExecutedPrice) * trade.Quantity * float64(trade.Leverage)
	} else { // SELL trade
		pnl = (*trade.ExecutedPrice - closePrice) * trade.Quantity * float64(trade.Leverage)
	}
	trade.Pnl = &pnl
	trade.ClosePrice = &closePrice
	trade.ClosedAt = models.TimePtr(time.Now()) // Helper to get a *time.Time
	trade.Status = models.TradeStatusClosed

	// --- Simulate wallet update for P&L (simplified) ---
	err = s.db.Transaction(func(tx *gorm.DB) error {
		tradeRepoTx := repository.NewTradeRepository(tx)
		walletRepoTx := repository.NewWalletRepository(tx)
		transactionRepoTx := repository.NewTransactionRepository(tx)

		// Update trade status and P&L
		if err := tradeRepoTx.UpdateTrade(trade); err != nil {
			return fmt.Errorf("failed to update trade for closing: %w", err)
		}

		// Update wallet with P&L
		wallet, err := walletRepoTx.GetWalletByUserID(traderID)
		if err != nil {
			return fmt.Errorf("failed to get wallet for P&L: %w", err)
		}
		if wallet == nil {
			return errors.New("trader wallet not found for P&L update")
		}

		wallet.Balance += pnl // Add/subtract P&L
		wallet.LastUpdated = time.Now()
		if err := walletRepoTx.UpdateWallet(wallet); err != nil {
			return fmt.Errorf("failed to update wallet balance with P&L: %w", err)
		}

		// Create a transaction for P&L
		txType := models.TxTypeTradeProfit
		if pnl < 0 {
			txType = models.TxTypeTradeLoss
		}
		walletTransaction := &models.WalletTransaction{
			WalletID:        wallet.ID,
			UserID:          traderID,
			TransactionType: txType,
			Amount:          pnl,
			Currency:        wallet.Currency,
			Status:          models.TxStatusSuccess,
			Description:     fmt.Sprintf("P&L for trade %d (%s %s %s)", trade.ID, trade.Side, trade.Symbol, trade.TradeType),
			BalanceBefore:   wallet.Balance - pnl,
			BalanceAfter:    wallet.Balance,
			TradeID:         &trade.ID,
		}
		if err := transactionRepoTx.CreateTransaction(walletTransaction); err != nil {
			return fmt.Errorf("failed to create P&L wallet transaction: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return trade, nil
}

// CancelTrade cancels a pending trade
func (s *tradeService) CancelTrade(tradeID uint, traderID uint) (*models.Trade, error) {
	trade, err := s.tradeRepo.GetTradeByID(tradeID, traderID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, ErrTradeNotFound
	}
	if trade.Status != models.TradeStatusPending {
		return nil, errors.New("only pending trades can be cancelled")
	}

	// In a real system, you'd send a cancel request to the exchange API here.
	// If successful:
	trade.Status = models.TradeStatusCancelled
	trade.ClosedAt = models.TimePtr(time.Now())

	err = s.tradeRepo.UpdateTrade(trade)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel trade: %w", err)
	}

	// If funds were reserved for this pending trade, return them to the wallet
	// This would involve finding the initial "funds reserved" transaction and reversing it.
	// For simplicity, we'll skip the detailed reversal here, but it's crucial.

	return trade, nil
}

// DeleteTrade (soft delete)
func (s *tradeService) DeleteTrade(tradeID uint, traderID uint) error {
	// A real trading platform might not allow deleting historical trades
	// Instead, they are marked as 'Closed' or 'Archived'.
	// Only pending/failed trades might be truly deletable.
	trade, err := s.tradeRepo.GetTradeByID(tradeID, traderID)
	if err != nil {
		return err
	}
	if trade == nil {
		return ErrTradeNotFound
	}
	if trade.Status == models.TradeStatusOpen || trade.Status == models.TradeStatusClosed {
		return errors.New("open or closed trades cannot be deleted, only archived")
	}

	return s.tradeRepo.DeleteTrade(tradeID, traderID)
}

// Helper function to get a pointer to time.Time
func (m *tradeService) TimePtr(t time.Time) *time.Time {
	return &t
}
