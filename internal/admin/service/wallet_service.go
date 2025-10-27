package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type IAdminWalletService interface {
	GetAdminWalletSummary() (*models.WalletSummaryResponse, error)
	AdminInitiateDeposit(input models.DepositRequestInput) (*models.DepositResponse, error)
	AdminVerifyDeposit(depositID uint, input models.DepositVerifyInput) (*models.DepositResponse, error)
	AdminRequestWithdrawal(input models.WithdrawalRequestInput) (*models.WithdrawalResponse, error)
	AdminGetWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error)
	CreditAdminWallet(tx *gorm.DB, amount float64, currency, description string) error

	GetPendingWithdrawalRequests(pagination models.PaginationParams) ([]models.WithdrawRequest, int64, error)
	ApproveWithdrawalRequest(withdrawalID uint) error
	RejectWithdrawalRequest(withdrawalID uint) error

	GetAllWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) // All platform transactions
	GetAllCustomerTransactionsWithUserDetails(pagination models.PaginationParams) ([]models.AdminTransactionDisplayDTO, int64, error)
}

type AdminWalletService struct {
	Repo repository.IAdminWalletRepository
	DB   *gorm.DB
}

func NewAdminWalletService(repo repository.IAdminWalletRepository, db *gorm.DB) *AdminWalletService {
	return &AdminWalletService{
		Repo: repo,
		DB:   db,
	}
}

// GetAllCustomerTransactionsWithUserDetails retrieves all customer transactions with linked user data.
func (s *AdminWalletService) GetAllCustomerTransactionsWithUserDetails(pagination models.PaginationParams) ([]models.AdminTransactionDisplayDTO, int64, error) {
	return s.Repo.GetAllCustomerTransactions(pagination)
}

// GetAllWalletTransactions retrieves all wallet transactions across the platform.
func (s *AdminWalletService) GetAllWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	return s.Repo.GetAllWalletTransactions(pagination)
}

// GetAdminWalletSummary retrieves the current balance and details of the admin wallet.
func (s *AdminWalletService) GetAdminWalletSummary() (*models.WalletSummaryResponse, error) {
	wallet, err := s.Repo.GetAdminWallet()
	if err != nil {
		return nil, fmt.Errorf("error getting admin wallet: %w", err)
	}
	return &models.WalletSummaryResponse{
		UserID:      wallet.UserID,
		WalletID:    wallet.ID,
		Balance:     wallet.Balance,
		Currency:    wallet.Currency,
		LastUpdated: wallet.LastUpdated,
	}, nil
}

// AdminInitiateDeposit creates a pending deposit request for the admin wallet.
func (s *AdminWalletService) AdminInitiateDeposit(input models.DepositRequestInput) (*models.DepositResponse, error) {
	adminUser, err := s.Repo.FindAdminUser()
	if err != nil {
		return nil, fmt.Errorf("admin user not found: %w", err)
	}

	depositRequest := models.DepositRequest{
		UserID:         adminUser.ID,
		Amount:         input.Amount,
		Currency:       input.Currency,
		Status:         models.TxStatusPending,
		PaymentGateway: "AdminManual",
		PaymentMethod:  input.PaymentMethod, // Ensure payment method is passed
		RedirectURL:    "",
	}

	if err := s.Repo.CreateDepositRequest(&depositRequest); err != nil {
		return nil, fmt.Errorf("failed to create admin deposit request: %w", err)
	}

	return &models.DepositResponse{
		DepositID: depositRequest.ID,
		Amount:    depositRequest.Amount,
		Currency:  depositRequest.Currency,
		Status:    depositRequest.Status,
		Message:   "Admin deposit initiated. Awaiting verification (manual or simulated).",
	}, nil
}

// AdminVerifyDeposit processes the verification of an admin deposit, updating wallet balance and transaction records.
func (s *AdminWalletService) AdminVerifyDeposit(depositID uint, input models.DepositVerifyInput) (*models.DepositResponse, error) {
	depositRequest, err := s.Repo.GetDepositRequestByID(depositID)
	if err != nil {
		return nil, fmt.Errorf("deposit request with ID %d not found: %w", depositID, err)
	}

	if depositRequest.Status != models.TxStatusPending {
		return nil, errors.New("deposit request is not in pending state or already processed")
	}

	// The input.Status refers to the webhook/external status. We map it to our internal TxStatus.
	// Assume 'SUCCESS' from input maps to models.TxStatusSuccess.
	// If the input.Status is not 'SUCCESS', we consider it failed for admin manual deposit.
	if input.Status != string(models.TxStatusSuccess) {
		depositRequest.Status = models.TxStatusFailed
		if err := s.Repo.UpdateDepositRequest(depositRequest); err != nil {
			log.Printf("Warning: Failed to update deposit request %d status to FAILED: %v", depositID, err)
		}
		return nil, fmt.Errorf("deposit verification explicitly failed with status: %s", input.Status)
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		wallet, err := s.Repo.GetAdminWallet()
		if err != nil {
			return fmt.Errorf("failed to get admin wallet for deposit verification: %w", err)
		}

		balanceBefore := wallet.Balance
		wallet.Balance += input.Amount
		wallet.LastUpdated = time.Now()

		if err := s.Repo.UpdateWalletBalance(tx, wallet); err != nil {
			return fmt.Errorf("failed to update admin wallet balance during deposit: %w", err)
		}

		depositRequest.Status = models.TxStatusSuccess
		depositRequest.PaymentGatewayTxID = input.PaymentGatewayTxID
		depositRequest.CompletionTime = models.TimePtr(time.Now()) // Set completion time
		if err := s.Repo.UpdateDepositRequest(depositRequest); err != nil {
			return fmt.Errorf("failed to update admin deposit request status: %w", err)
		}

		walletTx := models.WalletTransaction{
			WalletID:           wallet.ID,
			UserID:             wallet.UserID,
			TransactionType:    models.TxTypeDeposit,
			Amount:             input.Amount,
			Currency:           wallet.Currency,
			Status:             models.TxStatusSuccess,
			ReferenceID:        fmt.Sprintf("DEPOSIT-%d", depositRequest.ID),
			PaymentGatewayTxID: input.PaymentGatewayTxID,
			Description:        "Admin deposit via manual verification",
			BalanceBefore:      balanceBefore,
			BalanceAfter:       wallet.Balance,
		}
		if err := s.Repo.CreateWalletTransaction(tx, &walletTx); err != nil {
			return fmt.Errorf("failed to create admin wallet transaction record for deposit: %w", err)
		}

		// Link the created wallet transaction to the deposit request
		depositRequest.WalletTransactionID = &walletTx.ID
		return s.Repo.UpdateDepositRequest(depositRequest)
	})

	if err != nil {
		return nil, err
	}

	return &models.DepositResponse{
		DepositID: depositRequest.ID,
		Amount:    depositRequest.Amount,
		Currency:  depositRequest.Currency,
		Status:    depositRequest.Status,
		Message:   "Admin deposit successfully verified and credited.",
	}, nil
}

// AdminRequestWithdrawal processes a withdrawal request from the admin wallet.
func (s *AdminWalletService) AdminRequestWithdrawal(input models.WithdrawalRequestInput) (*models.WithdrawalResponse, error) {
	adminUser, err := s.Repo.FindAdminUser()
	if err != nil {
		return nil, fmt.Errorf("admin user not found: %w", err)
	}
	wallet, err := s.Repo.GetAdminWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to get admin wallet: %w", err)
	}

	if wallet.Balance < input.Amount {
		return nil, errors.New("insufficient balance for withdrawal")
	}
	if wallet.Currency != input.Currency {
		return nil, errors.New("currency mismatch for withdrawal")
	}

	var finalWithdrawRequest models.WithdrawRequest // Declare here to be accessible after transaction

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		// Re-fetch wallet within transaction to ensure latest state
		currentWallet, err := s.Repo.GetAdminWallet()
		if err != nil {
			return fmt.Errorf("failed to get admin wallet in transaction: %w", err)
		}

		if currentWallet.Balance < input.Amount {
			return errors.New("insufficient balance for withdrawal (transaction re-check)")
		}
		if currentWallet.Currency != input.Currency {
			return errors.New("currency mismatch for withdrawal (transaction re-check)")
		}

		balanceBefore := currentWallet.Balance
		currentWallet.Balance -= input.Amount // Deduct balance immediately
		currentWallet.LastUpdated = time.Now()

		if err := s.Repo.UpdateWalletBalance(tx, currentWallet); err != nil {
			return fmt.Errorf("failed to update admin wallet balance during withdrawal: %w", err)
		}

		// Populate all fields for the new WithdrawRequest struct
		withdrawRequest := models.WithdrawRequest{
			UserID:             adminUser.ID,
			Amount:             input.Amount,
			Currency:           input.Currency,
			Status:             models.TxStatusPending, // Withdrawal is pending approval/processing
			BankAccountNumber:  input.BankAccountNumber,  // Use from input
			BankAccountHolder:  input.BankAccountHolder,  // Use from input
			IFSCCode:           input.IFSCCode,           // Use from input
			BeneficiaryAccount: input.BeneficiaryAccount, // This might be redundant if individual fields are stored
			PaymentGateway:     "AdminManualTransfer",
			PaymentGatewayTxID: "", // Will be filled upon approval
			RequestTime:        time.Now(),
		}
		if err := s.Repo.CreateWithdrawRequest(&withdrawRequest); err != nil {
			return fmt.Errorf("failed to create admin withdrawal request: %w", err)
		}
		finalWithdrawRequest = withdrawRequest // Assign to outer variable

		walletTx := models.WalletTransaction{
			WalletID:           currentWallet.ID,
			UserID:             currentWallet.UserID,
			TransactionType:    models.TxTypeWithdrawal,
			Amount:             input.Amount,
			Currency:           currentWallet.Currency,
			Status:             models.TxStatusPending, // Wallet transaction status is also pending
			ReferenceID:        fmt.Sprintf("WITHDRAW-%d", withdrawRequest.ID),
			PaymentGatewayTxID: "",
			Description:        fmt.Sprintf("Admin withdrawal request to %s (Account: %s)", input.BankAccountHolder, input.BankAccountNumber), // More descriptive
			BalanceBefore:      balanceBefore,
			BalanceAfter:       currentWallet.Balance,
		}
		if err := s.Repo.CreateWalletTransaction(tx, &walletTx); err != nil {
			return fmt.Errorf("failed to create admin withdrawal transaction record: %w", err)
		}
		// Link the created wallet transaction to the withdrawal request
		finalWithdrawRequest.WalletTransactionID = &walletTx.ID
		return s.Repo.UpdateWithdrawRequest(&finalWithdrawRequest)
	})

	if err != nil {
		return nil, err
	}

	return &models.WithdrawalResponse{
		WithdrawalID:       finalWithdrawRequest.ID,
		Amount:             input.Amount,
		Currency:           input.Currency,
		Status:             finalWithdrawRequest.Status, // Should be pending
		Message:            "Admin withdrawal request submitted for processing.",
		PaymentGatewayTxID: finalWithdrawRequest.PaymentGatewayTxID,
	}, nil
}

// AdminGetWalletTransactions retrieves transaction records specifically for the admin's wallet.
func (s *AdminWalletService) AdminGetWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	return s.Repo.AdminGetWalletTransactions(pagination)
}

// CreditAdminWallet adds funds to the admin wallet and records a transaction. This is typically used internally by the system.
func (s *AdminWalletService) CreditAdminWallet(tx *gorm.DB, amount float64, currency, description string) error {
	adminWallet, err := s.Repo.GetAdminWallet()
	if err != nil {
		return fmt.Errorf("failed to get admin wallet for credit: %w", err)
	}

	if adminWallet.Currency != currency {
		return fmt.Errorf("currency mismatch for admin wallet credit: expected %s, got %s", adminWallet.Currency, currency)
	}

	balanceBefore := adminWallet.Balance
	adminWallet.Balance += amount
	adminWallet.LastUpdated = time.Now()

	if err := s.Repo.UpdateWalletBalance(tx, adminWallet); err != nil {
		return fmt.Errorf("failed to update admin wallet balance for credit: %w", err)
	}

	walletTx := models.WalletTransaction{
		WalletID:        adminWallet.ID,
		UserID:          adminWallet.UserID,
		TransactionType: models.TxTypeDeposit, // Or a more specific type like TxTypeCredit if applicable
		Amount:          amount,
		Currency:        currency,
		Status:          models.TxStatusSuccess,
		Description:     description,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    adminWallet.Balance,
	}
	if err := s.Repo.CreateWalletTransaction(tx, &walletTx); err != nil {
		return fmt.Errorf("failed to create admin wallet transaction record for credit: %w", err)
	}
	return nil
}

// GetPendingWithdrawalRequests retrieves all withdrawal requests that are in a 'Pending' state.
func (s *AdminWalletService) GetPendingWithdrawalRequests(pagination models.PaginationParams) ([]models.WithdrawRequest, int64, error) {
	return s.Repo.GetPendingWithdrawalRequests(pagination)
}

// ApproveWithdrawalRequest marks a pending withdrawal request as 'Success' and updates the associated wallet transaction.
func (s *AdminWalletService) ApproveWithdrawalRequest(withdrawalID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		withdrawal, err := s.Repo.GetWithdrawRequestByID(withdrawalID)
		if err != nil {
			return fmt.Errorf("withdrawal request %d not found: %w", withdrawalID, err)
		}
		if withdrawal.Status != models.TxStatusPending {
			return errors.New("withdrawal request is not pending or already processed")
		}

		withdrawal.Status = models.TxStatusSuccess
		// Assign a mock payment gateway transaction ID upon approval
		withdrawal.PaymentGatewayTxID = fmt.Sprintf("PAYOUT-%d-%s", withdrawal.ID, time.Now().Format("20060102150405"))
		withdrawal.ProcessingTime = models.TimePtr(time.Now())
		withdrawal.CompletionTime = models.TimePtr(time.Now())

		if err := s.Repo.UpdateWithdrawRequest(withdrawal); err != nil {
			return fmt.Errorf("failed to update withdrawal request %d status to SUCCESS: %w", withdrawalID, err)
		}

		// If there's a linked wallet transaction, update its status too
		if withdrawal.WalletTransactionID != nil {
			var walletTx models.WalletTransaction
			// Ensure to use the transaction (tx) for fetching the wallet transaction
			if err := tx.First(&walletTx, *withdrawal.WalletTransactionID).Error; err != nil {
				// Log a warning but don't fail the transaction if the linked transaction is missing
				log.Printf("Warning: Wallet transaction for withdrawal %d (ID %d) not found: %v", withdrawal.ID, *withdrawal.WalletTransactionID, err)
			} else {
				walletTx.Status = models.TxStatusSuccess
				if err := tx.Save(&walletTx).Error; err != nil {
					return fmt.Errorf("failed to update linked wallet transaction %d status to SUCCESS: %w", *withdrawal.WalletTransactionID, err)
				}
			}
		}

		return nil
	})
}

// RejectWithdrawalRequest marks a pending withdrawal request as 'Rejected', reverses the funds to the customer,
// and updates the associated wallet transaction and creates a reversal transaction.
func (s *AdminWalletService) RejectWithdrawalRequest(withdrawalID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		withdrawal, err := s.Repo.GetWithdrawRequestByID(withdrawalID)
		if err != nil {
			return fmt.Errorf("withdrawal request %d not found: %w", withdrawalID, err)
		}
		if withdrawal.Status != models.TxStatusPending {
			return errors.New("withdrawal request is not pending or already processed")
		}

		// Get the customer's wallet to reverse funds
		customerWallet, err := s.Repo.GetCustomerWallet(withdrawal.UserID)
		if err != nil {
			return fmt.Errorf("customer wallet not found for user %d during rejection: %w", withdrawal.UserID, err)
		}

		// Revert funds to customer's wallet
		balanceBeforeReversal := customerWallet.Balance
		customerWallet.Balance += withdrawal.Amount
		customerWallet.LastUpdated = time.Now()

		if err := s.Repo.UpdateCustomerWalletBalance(tx, customerWallet); err != nil {
			return fmt.Errorf("failed to revert funds to customer wallet for withdrawal %d: %w", withdrawal.ID, err)
		}

		// Update withdrawal request status to Rejected
		withdrawal.Status = models.TxStatusRejected
		withdrawal.ProcessingTime = models.TimePtr(time.Now())
		withdrawal.CompletionTime = models.TimePtr(time.Now()) // Or a specific rejection time
		if err := s.Repo.UpdateWithdrawRequest(withdrawal); err != nil {
			return fmt.Errorf("failed to update withdrawal request %d status to REJECTED: %w", withdrawal.ID, err)
		}

		// Create a new wallet transaction for the reversal
		reversalTx := models.WalletTransaction{
			WalletID:        customerWallet.ID,
			UserID:          customerWallet.UserID,
			TransactionType: models.TxTypeReversal,
			Amount:          withdrawal.Amount,
			Currency:        withdrawal.Currency,
			Status:          models.TxStatusSuccess, // Reversal itself is a successful action
			ReferenceID:     fmt.Sprintf("REVERSAL-WITHDRAW-%d", withdrawal.ID),
			Description:     fmt.Sprintf("Withdrawal request %d rejected, funds returned to wallet.", withdrawal.ID),
			BalanceBefore:   balanceBeforeReversal,
			BalanceAfter:    customerWallet.Balance,
		}
		if err := s.Repo.CreateWalletTransaction(tx, &reversalTx); err != nil {
			return fmt.Errorf("failed to create reversal transaction for customer for withdrawal %d: %w", withdrawal.ID, err)
		}

		// Update the original wallet transaction linked to the withdrawal request
		if withdrawal.WalletTransactionID != nil {
			var originalWalletTx models.WalletTransaction
			if err := tx.First(&originalWalletTx, *withdrawal.WalletTransactionID).Error; err != nil {
				log.Printf("Warning: Original wallet transaction for withdrawal %d (ID %d) not found for status update: %v", withdrawal.ID, *withdrawal.WalletTransactionID, err)
			} else {
				originalWalletTx.Status = models.TxStatusReversed // Mark original withdrawal transaction as reversed
				originalWalletTx.Description = fmt.Sprintf("Original withdrawal reversed: %s", originalWalletTx.Description)
				if err := tx.Save(&originalWalletTx).Error; err != nil {
					return fmt.Errorf("failed to update original linked wallet transaction %d status to REVERSED: %w", *withdrawal.WalletTransactionID, err)
				}
			}
		}

		return nil
	})
}

// Helper function to get a pointer to a time.Time value
func (s *AdminWalletService) TimePtr(t time.Time) *time.Time {
    return &t
}