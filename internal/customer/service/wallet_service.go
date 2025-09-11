package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	paymentgateway "github.com/fathimasithara01/tradeverse/pkg/payment_gateway.go" // Assuming this path is correct
	// Import gorm for transaction error checking
)

type WalletServicer interface {
	GetWalletSummary(userID uint) (*models.WalletSummaryResponse, error)
	InitiateDeposit(userID uint, amount float64, currency string) (*models.DepositResponse, error)
	VerifyDeposit(pgTxID string, amount float64, status string) error // Webhook/Callback
	InitiateWithdrawal(userID uint, amount float64, currency, beneficiaryAccount string) (*models.WithdrawalResponse, error)
	ListTransactions(userID uint, page, limit int) ([]models.WalletTransaction, int64, error)
	// Additional admin methods might be needed, e.g., for deposit/withdrawal approval/rejection
}

type walletService struct {
	walletRepo repository.WalletRepository
	pgClient   *paymentgateway.SimulatedPaymentClient // Simulated Payment Gateway Client
}

func NewWalletService(walletRepo repository.WalletRepository, pgClient *paymentgateway.SimulatedPaymentClient) WalletServicer {
	return &walletService{
		walletRepo: walletRepo,
		pgClient:   pgClient,
	}
}

func (s *walletService) GetWalletSummary(userID uint) (*models.WalletSummaryResponse, error) {
	wallet, err := s.walletRepo.GetWalletByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve wallet: %w", err)
	}
	if wallet == nil {
		newWallet := &models.Wallet{
			UserID:      userID,
			Balance:     0.0,
			Currency:    "INR", // Default currency
			LastUpdated: time.Now(),
		}
		if err := s.walletRepo.CreateWallet(newWallet); err != nil {
			return nil, fmt.Errorf("failed to create default wallet: %w", err)
		}
		wallet = newWallet
	}

	return &models.WalletSummaryResponse{
		UserID:      wallet.UserID,
		Balance:     wallet.Balance,
		Currency:    wallet.Currency,
		LastUpdated: wallet.LastUpdated,
	}, nil
}

func (s *walletService) InitiateDeposit(userID uint, amount float64, currency string) (*models.DepositResponse, error) {
	// 1. Initiate with Payment Gateway
	pgTxID, redirectURL, err := s.pgClient.CreateDepositInitiation(amount, currency, fmt.Sprintf("user_%d", userID))
	if err != nil {
		return nil, fmt.Errorf("payment gateway initiation failed: %w", err)
	}

	// 2. Record the deposit request in our database
	depositReq := &models.DepositRequest{
		UserID:             userID,
		Amount:             amount,
		Currency:           currency,
		Status:             models.TxStatusPending,
		PaymentGateway:     "SimulatedPG", // Or actual PG name
		PaymentGatewayTxID: pgTxID,
		RedirectURL:        redirectURL,
	}
	if err := s.walletRepo.CreateDepositRequest(depositReq); err != nil {
		return nil, fmt.Errorf("failed to save deposit request: %w", err)
	}

	return &models.DepositResponse{
		DepositID:          depositReq.ID,
		Amount:             depositReq.Amount,
		Currency:           depositReq.Currency,
		Status:             depositReq.Status,
		RedirectURL:        depositReq.RedirectURL,
		PaymentGatewayTxID: pgTxID,
		Message:            "Deposit initiated. Please complete payment via redirect URL.",
	}, nil
}

func (s *walletService) VerifyDeposit(pgTxID string, amount float64, status string) error {
	// 1. Find the corresponding deposit request
	depositReq, err := s.walletRepo.FindDepositRequestByPGTxID(pgTxID)
	if err != nil {
		return fmt.Errorf("failed to find deposit request by PG TxID %s: %w", pgTxID, err)
	}
	if depositReq == nil {
		return fmt.Errorf("deposit request not found for payment gateway transaction ID: %s", pgTxID)
	}

	// Idempotency check: If already successfully processed, do nothing.
	if depositReq.Status == models.TxStatusSuccess {
		fmt.Printf("Deposit %d for PG_TxID %s already SUCCESS. Idempotent call ignored.\n", depositReq.ID, pgTxID)
		return nil // Already processed, no error
	}

	// Important: Check if the amount matches to prevent tampering
	if depositReq.Amount != amount {
		return fmt.Errorf("amount mismatch for deposit %d. Expected %.2f, got %.2f", depositReq.ID, depositReq.Amount, amount)
	}

	originalDepositStatus := depositReq.Status

	// Update deposit request status based on callback status
	switch status {
	case "SUCCESS":
		depositReq.Status = models.TxStatusSuccess
	case "FAILED":
		depositReq.Status = models.TxStatusFailed
	case "PENDING":
		// If it's already pending and we get another pending, just update the timestamp or similar.
		// No balance change, but we proceed to update the request record.
		depositReq.Status = models.TxStatusPending // Explicitly set to PENDING
	case "CANCELLED": // Add handling for cancelled status
		depositReq.Status = models.TxStatusCancelled
	default:
		return fmt.Errorf("unsupported payment gateway status: %s for PG_TxID %s", status, pgTxID)
	}

	// Only process balance update if the status changed to SUCCESS and it wasn't already success
	if depositReq.Status == models.TxStatusSuccess && originalDepositStatus != models.TxStatusSuccess {
		// 2. Perform the actual balance update and transaction logging in a database transaction
		tx := s.walletRepo.BeginTransaction()
		if tx.Error != nil {
			return fmt.Errorf("failed to begin database transaction: %w", tx.Error)
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				fmt.Printf("Recovered from panic during deposit verification: %v\n", r)
			}
		}()

		wallet, err := s.walletRepo.GetWalletByUserID(depositReq.UserID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to retrieve user wallet: %w", err)
		}
		if wallet == nil {
			// This shouldn't happen if user exists, but handle defensively
			newWallet := &models.Wallet{
				UserID:      depositReq.UserID,
				Balance:     0.0,
				Currency:    depositReq.Currency,
				LastUpdated: time.Now(),
			}
			if err := s.walletRepo.CreateWallet(newWallet); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create default wallet for user: %w", err)
			}
			wallet = newWallet
		}

		balanceBefore := wallet.Balance
		if err := s.walletRepo.UpdateWalletBalance(tx, wallet, depositReq.Amount); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}

		walletTx := &models.WalletTransaction{
			WalletID:           wallet.ID,
			UserID:             depositReq.UserID,
			TransactionType:    models.TxTypeDeposit,
			Amount:             depositReq.Amount,
			Currency:           depositReq.Currency,
			Status:             models.TxStatusSuccess, // Wallet transaction is only created on success
			ReferenceID:        fmt.Sprintf("DEPOSIT_%d", depositReq.ID),
			PaymentGatewayTxID: depositReq.PaymentGatewayTxID,
			Description:        "Deposit via " + depositReq.PaymentGateway,
			BalanceBefore:      balanceBefore,
			BalanceAfter:       wallet.Balance,
		}
		if err := s.walletRepo.CreateWalletTransaction(tx, walletTx); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create wallet transaction record: %w", err)
		}

		// Link transaction to deposit request
		depositReq.WalletTransactionID = &walletTx.ID

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit deposit transaction: %w", err)
		}
	}

	// 3. Update the deposit request in the database (even if status is PENDING/FAILED/CANCELLED)
	// This ensures the record reflects the latest status from the PG.
	if err := s.walletRepo.UpdateDepositRequest(depositReq); err != nil {
		return fmt.Errorf("failed to update deposit request status for PG_TxID %s: %w", pgTxID, err)
	}

	return nil
}

// InitiateWithdrawal processes a user's request to withdraw funds.
func (s *walletService) InitiateWithdrawal(userID uint, amount float64, currency, beneficiaryAccount string) (*models.WithdrawalResponse, error) {
	// 1. Check wallet balance and ensure sufficient funds
	wallet, err := s.walletRepo.GetWalletByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve wallet: %w", err)
	}
	if wallet == nil {
		return nil, errors.New("user wallet not found")
	}
	if wallet.Balance < amount {
		return nil, fmt.Errorf("insufficient funds. Available: %.2f %s, Requested: %.2f %s", wallet.Balance, wallet.Currency, amount, currency)
	}
	if wallet.Currency != currency {
		return nil, fmt.Errorf("currency mismatch. Wallet currency: %s, Requested currency: %s", wallet.Currency, currency)
	}

	// 2. Create a withdrawal request and mark as PENDING
	withdrawReq := &models.WithdrawRequest{
		UserID:             userID,
		Amount:             amount,
		Currency:           currency,
		Status:             models.TxStatusPending, // Initially pending
		BeneficiaryAccount: beneficiaryAccount,
		PaymentGateway:     "SimulatedPG",
	}
	if err := s.walletRepo.CreateWithdrawRequest(withdrawReq); err != nil {
		return nil, fmt.Errorf("failed to create withdrawal request: %w", err)
	}

	// 3. Deduct funds and log transaction in a database transaction
	tx := s.walletRepo.BeginTransaction()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin database transaction for withdrawal: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Printf("Recovered from panic during withdrawal initiation: %v\n", r)
		}
	}()

	balanceBefore := wallet.Balance
	if err := s.walletRepo.UpdateWalletBalance(tx, wallet, -amount); err != nil { // Deduct amount
		tx.Rollback()
		return nil, fmt.Errorf("failed to update wallet balance during withdrawal: %w", err)
	}

	walletTx := &models.WalletTransaction{
		WalletID:        wallet.ID,
		UserID:          userID,
		TransactionType: models.TxTypeWithdraw,
		Amount:          amount,
		Currency:        currency,
		Status:          models.TxStatusPending, // Still pending until PG confirms
		ReferenceID:     fmt.Sprintf("WITHDRAW_%d", withdrawReq.ID),
		Description:     "Withdrawal to " + beneficiaryAccount,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    wallet.Balance,
	}
	if err := s.walletRepo.CreateWalletTransaction(tx, walletTx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create wallet transaction record for withdrawal: %w", err)
	}

	// Link transaction to withdrawal request
	withdrawReq.WalletTransactionID = &walletTx.ID
	if err := s.walletRepo.UpdateWithdrawRequest(withdrawReq); err != nil {
		// Even if linking fails, we ideally want to commit the balance deduction and transaction record
		// but since this is part of the same logical unit, we'll roll back for consistency.
		tx.Rollback()
		return nil, fmt.Errorf("failed to link transaction to withdrawal request: %w", err)
	}

	// 4. Call Payment Gateway to process withdrawal (this might be asynchronous)
	pgWithdrawTxID, err := s.pgClient.ProcessWithdrawal(amount, currency, beneficiaryAccount)
	if err != nil {
		tx.Rollback() // If PG call fails, roll back everything
		return nil, fmt.Errorf("payment gateway withdrawal failed: %w", err)
	}

	// Update withdrawal request with PG transaction ID and potentially status (if synchronous)
	withdrawReq.PaymentGatewayTxID = pgWithdrawTxID

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit withdrawal transaction: %w", err)
	}

	// After a successful commit, if updating withdrawReq fails, it's a warning, not a rollback reason.
	// The core financial transaction is already committed.
	if err := s.walletRepo.UpdateWithdrawRequest(withdrawReq); err != nil {
		fmt.Printf("Warning: Failed to update withdrawal request %d with PG transaction ID %s after successful commit: %v\n", withdrawReq.ID, pgWithdrawTxID, err)
	}

	return &models.WithdrawalResponse{
		WithdrawalID:       withdrawReq.ID,
		Amount:             withdrawReq.Amount,
		Currency:           withdrawReq.Currency,
		Status:             withdrawReq.Status, // Will be PENDING
		PaymentGatewayTxID: pgWithdrawTxID,
		Message:            "Withdrawal initiated and funds deducted. Processing with payment gateway.",
	}, nil
}

func (s *walletService) ListTransactions(userID uint, page, limit int) ([]models.WalletTransaction, int64, error) {
	transactions, total, err := s.walletRepo.FindWalletTransactions(userID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve transactions: %w", err)
	}
	return transactions, total, nil
}
