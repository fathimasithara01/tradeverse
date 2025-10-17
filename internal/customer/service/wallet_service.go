package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	walletrepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	paymentgateway "github.com/fathimasithara01/tradeverse/pkg/payment_gateway.go"
	"gorm.io/gorm"
)

// Define specific EXPORTED errors for the wallet service
var (
	ErrWalletServiceNotFound          = errors.New("wallet not found for user")
	ErrWalletServiceInsufficientFunds = errors.New("insufficient funds in wallet")
	ErrWalletServiceTransactionFailed = errors.New("wallet transaction failed")
	ErrDepositAlreadyProcessed        = errors.New("deposit already processed")
	ErrInvalidDepositStatus           = errors.New("invalid deposit status")
	ErrUserWalletNotFound             = errors.New("user wallet not found")
	ErrWithdrawalRequestNotFound      = errors.New("withdrawal request not found")
)

type IWalletService interface {
	GetUserWallet(ctx context.Context, userID uint) (*models.Wallet, error)
	DepositFunds(ctx context.Context, userID uint, amount float64, referenceID, description string) (*models.WalletTransaction, error)
	WithdrawFunds(ctx context.Context, userID uint, amount float64, referenceID, description string) (*models.WalletTransaction, error)
	GetWalletTransactions(ctx context.Context, userID uint, pagination models.PaginationParams) ([]models.WalletTransaction, int64, error)
	GetWalletSummary(userID uint) (*models.WalletSummaryResponse, error)
	InitiateDeposit(userID uint, input models.DepositRequestInput) (*models.DepositResponse, error)
	VerifyDeposit(depositID uint, input models.DepositVerifyInput) (*models.DepositVerifyResponse, error)
	RequestWithdrawal(userID uint, input models.WithdrawalRequestInput) (*models.WithdrawalResponse, error)
	GetTransactions(userID uint, pagination models.PaginationParams) ([]models.WalletTransaction, int64, error)
	DebitUserWallet(userID uint, amount float64, currency, description, transactionID string) error
}

type walletService struct {
	db             *gorm.DB
	walletRepo     walletrepo.WalletRepository
	paymentGateway paymentgateway.SimulatedPaymentClient // Correct: Field name is `paymentGateway`

}

func NewWalletService(db *gorm.DB, repo walletrepo.WalletRepository, pgClient paymentgateway.SimulatedPaymentClient) IWalletService {
	return &walletService{
		db:             db,
		walletRepo:     repo,
		paymentGateway: pgClient,
	}
}
func (s *walletService) DebitUserWallet(userID uint, amount float64, currency, description, transactionID string) error {
	var wallet models.Wallet

	// ✅ Find wallet
	err := s.db.Where("user_id = ?", userID).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// ✅ Create wallet automatically if missing
		wallet = models.Wallet{
			UserID:   userID,
			Balance:  0,
			Currency: currency,
		}
		if err := s.db.Create(&wallet).Error; err != nil {
			return fmt.Errorf("failed to create wallet: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to fetch wallet: %w", err)
	}

	// ✅ Check balance
	if wallet.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	// ✅ Calculate balances
	balanceBefore := wallet.Balance
	balanceAfter := wallet.Balance - amount

	// ✅ Create transaction record
	tx := models.WalletTransaction{
		WalletID:      wallet.ID,
		UserID:        userID,
		Type:          models.TxTypeDebit,
		Amount:        amount,
		Currency:      currency,
		Status:        models.TxStatusSuccess,
		Description:   description,
		TransactionID: transactionID,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
	}

	// ✅ Update wallet & save transaction
	err = s.db.Transaction(func(txn *gorm.DB) error {
		if err := txn.Create(&tx).Error; err != nil {
			return fmt.Errorf("failed to record transaction: %w", err)
		}

		if err := txn.Model(&wallet).Update("balance", balanceAfter).Error; err != nil {
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to record debit transaction in transaction: %w", err)
	}

	return nil
}

func (s *walletService) GetUserWallet(ctx context.Context, userID uint) (*models.Wallet, error) {
	wallet, err := s.walletRepo.GetUserWallet(userID)
	if err != nil {
		if errors.Is(err, walletrepo.ErrWalletNotFound) {
			return nil, ErrWalletServiceNotFound
		}
		return nil, fmt.Errorf("failed to get user wallet: %w", err)
	}
	return wallet, nil
}

func (s *walletService) GetWalletSummary(userID uint) (*models.WalletSummaryResponse, error) {
	wallet, err := s.GetUserWallet(context.Background(), userID)
	if err != nil {
		if errors.Is(err, ErrWalletServiceNotFound) {
			return nil, ErrUserWalletNotFound
		}
		return nil, err
	}
	return &models.WalletSummaryResponse{
		UserID:      wallet.UserID,
		WalletID:    wallet.ID,
		Balance:     wallet.Balance,
		Currency:    wallet.Currency,
		LastUpdated: wallet.LastUpdated,
	}, nil
}

func (s *walletService) InitiateDeposit(userID uint, input models.DepositRequestInput) (*models.DepositResponse, error) {
	log.Printf("Initiating deposit for user %d, amount %.2f, method %s", userID, input.Amount, input.PaymentMethod)

	pgTxID, redirectURL, err := s.paymentGateway.CreateDepositInitiation(input.Amount, input.Currency, fmt.Sprint(userID))
	if err != nil {
		return nil, fmt.Errorf("payment gateway initiation failed: %w", err)
	}

	depositRequest := &models.DepositRequest{
		UserID:             userID,
		Amount:             input.Amount,
		PaymentMethod:      input.PaymentMethod,
		Currency:           input.Currency,
		Status:             models.TxStatusPending,
		RequestTime:        time.Now(),
		PaymentGatewayTxID: pgTxID, // Store the PG transaction ID
	}
	if err := s.walletRepo.CreateDepositRequest(depositRequest); err != nil {
		return nil, fmt.Errorf("failed to create deposit request: %w", err)
	}

	return &models.DepositResponse{
		DepositID:   depositRequest.ID,
		Status:      models.TxStatusPending,
		Message:     "Deposit initiated. Awaiting payment confirmation.",
		RedirectURL: redirectURL, // Use the redirect URL from PG
		Amount:      depositRequest.Amount,
		Currency:    depositRequest.Currency,
	}, nil
}

func (s *walletService) VerifyDeposit(depositID uint, input models.DepositVerifyInput) (*models.DepositVerifyResponse, error) {
	depositRequest, err := s.walletRepo.GetDepositRequestByID(depositID)
	if err != nil {
		if errors.Is(err, walletrepo.ErrDepositRequestNotFound) {
			return nil, walletrepo.ErrDepositRequestNotFound
		}
		return nil, fmt.Errorf("failed to get deposit request: %w", err)
	}

	if depositRequest.Status == models.TxStatusSuccess || depositRequest.Status == models.TxStatusFailed {
		return nil, ErrDepositAlreadyProcessed
	}

	// Step 1: Verify with Payment Gateway
	isVerified, pgErr := s.paymentGateway.VerifyDeposit(depositRequest.PaymentGatewayTxID)
	if pgErr != nil || !isVerified {
		// Payment Gateway verification failed
		depositRequest.Status = models.TxStatusFailed
		if err := s.walletRepo.UpdateDepositRequest(depositRequest); err != nil {
			log.Printf("Failed to update deposit request status to failed after PG verification failure for ID %d: %v", depositID, err)
		}
		return nil, fmt.Errorf("payment gateway verification failed or not successful: %w", ErrInvalidDepositStatus)
	}

	var createdTransaction *models.WalletTransaction // Will store the transaction created by CreditWallet
	err = s.db.Transaction(func(tx *gorm.DB) error {
		wallet, err := s.walletRepo.GetUserWallet(depositRequest.UserID)
		if err != nil {
			return fmt.Errorf("%w for user %d: %v", ErrUserWalletNotFound, depositRequest.UserID, err) // Use service-level error
		}

		// Correctly capture the transaction created by CreditWallet
		createdTransaction, err = s.walletRepo.CreditWallet(tx, wallet.ID, depositRequest.Amount, models.TxTypeDeposit,
			fmt.Sprintf("DEPOSIT_%d", depositRequest.ID), "Funds added via deposit verification")
		if err != nil {
			return fmt.Errorf("failed to credit user wallet during deposit verification: %w", err)
		}

		depositRequest.Status = models.TxStatusSuccess
		now := time.Now()
		depositRequest.CompletionTime = &now
		// Use the PG Transaction ID from the initial request.
		// If `input.PaymentGatewayTxID` is important and different, handle carefully.
		// For now, we assume `depositRequest.PaymentGatewayTxID` stored from initiation is sufficient.
		// depositRequest.PaymentGatewayTxID = input.PaymentGatewayTxID // Uncomment if you want to update PG TxID from input
		if err := s.walletRepo.UpdateDepositRequestTx(tx, depositRequest); err != nil {
			return fmt.Errorf("failed to update deposit request status: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWalletServiceTransactionFailed, err)
	}

	// Ensure createdTransaction is not nil before accessing its fields
	transactionID := ""
	if createdTransaction != nil {
		transactionID = createdTransaction.ReferenceID
	}

	return &models.DepositVerifyResponse{
		DepositID:     depositID,
		Status:        models.TxStatusSuccess,
		TransactionID: transactionID, // Use the ReferenceID from the actual created transaction
		Message:       "Deposit verified and funds added.",
	}, nil
}

func (s *walletService) RequestWithdrawal(userID uint, input models.WithdrawalRequestInput) (*models.WithdrawalResponse, error) {
	log.Printf("Requesting withdrawal for user %d, amount %.2f, bank account %s", userID, input.Amount, input.BankAccountNumber)

	if input.Amount <= 0 {
		return nil, errors.New("withdrawal amount must be positive")
	}

	var withdrawalRequest *models.WithdrawalRequest
	err := s.db.Transaction(func(tx *gorm.DB) error {
		wallet, err := s.walletRepo.GetUserWallet(userID)
		if err != nil {
			return fmt.Errorf("%w for user %d: %v", ErrWalletServiceNotFound, userID, err)
		}
		if wallet.Balance < input.Amount {
			return ErrWalletServiceInsufficientFunds
		}

		// --- Initiate with Payment Gateway FIRST ---
		pgTxID, pgErr := s.paymentGateway.ProcessWithdrawal(input.Amount, input.Currency, input.BankAccountNumber)
		if pgErr != nil {
			return fmt.Errorf("payment gateway withdrawal failed: %w", pgErr)
		}

		// Debit wallet ONLY AFTER PG confirms initiation/processing
		debitRefID := fmt.Sprintf("WITHDRAW_REQ_%s", time.Now().Format("20060102150405"))
		err = s.walletRepo.DebitWallet(tx, wallet.ID, input.Amount, models.TxTypeWithdrawal,
			debitRefID, "Withdrawal request debit (pending processing by PG)")
		if err != nil {
			return fmt.Errorf("failed to debit wallet for withdrawal request: %w", err)
		}

		withdrawalRequest = &models.WithdrawalRequest{
			UserID:             userID,
			Amount:             input.Amount,
			Currency:           input.Currency,
			BankAccountNumber:  input.BankAccountNumber,
			BankAccountHolder:  input.BankAccountHolder,
			IFSCCode:           input.IFSCCode,
			Status:             models.TxStatusSuccess, // Assume success if PG processed and debit succeeded
			RequestTime:        time.Now(),
			PaymentGatewayTxID: pgTxID,
		}
		if err := s.walletRepo.CreateWithdrawalRequest(withdrawalRequest); err != nil {
			return fmt.Errorf("failed to create withdrawal request: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWalletServiceTransactionFailed, err)
	}

	return &models.WithdrawalResponse{
		WithdrawalID: withdrawalRequest.ID,
		Amount:       withdrawalRequest.Amount,
		Currency:     withdrawalRequest.Currency,
		Status:       models.TxStatusSuccess,
		Message:      "Withdrawal request submitted and processed.",
	}, nil
}

func (s *walletService) GetTransactions(userID uint, pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	return s.GetWalletTransactions(context.Background(), userID, pagination)
}

func (s *walletService) DepositFunds(ctx context.Context, userID uint, amount float64, referenceID, description string) (*models.WalletTransaction, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be positive")
	}

	var transaction *models.WalletTransaction
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		currentWallet, err := s.walletRepo.GetOrCreateWallet(userID)
		if err != nil {
			return fmt.Errorf("failed to get or create wallet: %w", err)
		}

		balanceBefore := currentWallet.Balance
		currentWallet.Balance += amount
		currentWallet.LastUpdated = time.Now()
		if err := s.walletRepo.UpdateWalletTx(tx, currentWallet); err != nil {
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}

		newTx := &models.WalletTransaction{
			WalletID:        currentWallet.ID,
			UserID:          userID,
			TransactionType: models.TxTypeDeposit,
			Amount:          amount,
			Currency:        currentWallet.Currency,
			Status:          models.TxStatusSuccess,
			ReferenceID:     referenceID,
			Description:     description,
			BalanceBefore:   balanceBefore,
			BalanceAfter:    currentWallet.Balance,
		}
		if err := s.walletRepo.CreateWalletTransaction(tx, newTx); err != nil {
			return fmt.Errorf("failed to create wallet transaction: %w", err)
		}
		transaction = newTx
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWalletServiceTransactionFailed, err)
	}
	return transaction, nil
}

func (s *walletService) WithdrawFunds(ctx context.Context, userID uint, amount float64, referenceID, description string) (*models.WalletTransaction, error) {
	if amount <= 0 {
		return nil, errors.New("withdrawal amount must be positive")
	}

	var transaction *models.WalletTransaction
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		currentWallet, err := s.walletRepo.GetUserWallet(userID)
		if err != nil {
			if errors.Is(err, walletrepo.ErrWalletNotFound) {
				return ErrWalletServiceNotFound
			}
			return fmt.Errorf("failed to get user wallet: %w", err)
		}

		if currentWallet.Balance < amount {
			return ErrWalletServiceInsufficientFunds
		}

		balanceBefore := currentWallet.Balance
		currentWallet.Balance -= amount
		currentWallet.LastUpdated = time.Now()
		if err := s.walletRepo.UpdateWalletTx(tx, currentWallet); err != nil {
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}

		newTx := &models.WalletTransaction{
			WalletID:        currentWallet.ID,
			UserID:          userID,
			TransactionType: models.TxTypeWithdrawal,
			Amount:          amount,
			Currency:        currentWallet.Currency,
			Status:          models.TxStatusSuccess,
			ReferenceID:     referenceID,
			Description:     description,
			BalanceBefore:   balanceBefore,
			BalanceAfter:    currentWallet.Balance,
		}
		if err := s.walletRepo.CreateWalletTransaction(tx, newTx); err != nil {
			return fmt.Errorf("failed to create wallet transaction: %w", err)
		}
		transaction = newTx
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWalletServiceTransactionFailed, err)
	}
	return transaction, nil
}

func (s *walletService) GetWalletTransactions(ctx context.Context, userID uint, pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	return s.walletRepo.GetWalletTransactions(userID, pagination)
}
