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

// IAdminWalletService defines the interface for admin wallet services.
type IAdminWalletService interface {
	GetAdminWalletSummary() (*models.WalletSummaryResponse, error)
	AdminInitiateDeposit(input models.DepositRequestInput) (*models.DepositResponse, error)
	AdminVerifyDeposit(depositID uint, input models.DepositVerifyInput) (*models.DepositResponse, error)
	AdminRequestWithdrawal(input models.WithdrawalRequestInput) (*models.WithdrawalResponse, error)
	AdminGetWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error)
	CreditAdminWallet(tx *gorm.DB, amount float64, currency, description string) error // For subscription payments

	// For approving/rejecting customer withdrawals
	GetPendingWithdrawalRequests(pagination models.PaginationParams) ([]models.WithdrawRequest, int64, error)
	ApproveWithdrawalRequest(withdrawalID uint) error
	RejectWithdrawalRequest(withdrawalID uint) error
}

type AdminWalletService struct {
	Repo repository.IAdminWalletRepository
	DB   *gorm.DB // Inject DB for transaction management
}

func NewAdminWalletService(repo repository.IAdminWalletRepository, db *gorm.DB) *AdminWalletService {
	return &AdminWalletService{
		Repo: repo,
		DB:   db,
	}
}

// GetAdminWalletSummary retrieves the summary of the admin's wallet.
func (s *AdminWalletService) GetAdminWalletSummary() (*models.WalletSummaryResponse, error) {
	wallet, err := s.Repo.GetAdminWallet()
	if err != nil {
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

// AdminInitiateDeposit simulates initiating a deposit for the admin.
func (s *AdminWalletService) AdminInitiateDeposit(input models.DepositRequestInput) (*models.DepositResponse, error) {
	adminUser, err := s.Repo.FindAdminUser()
	if err != nil {
		return nil, err
	}

	depositRequest := models.DepositRequest{
		UserID:         adminUser.ID,
		Amount:         input.Amount,
		Currency:       input.Currency,
		Status:         models.TxStatusPending,
		PaymentGateway: "AdminManual", // Or a specific gateway for admin deposits
		RedirectURL:    "",            // Not applicable for internal admin deposits
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

// AdminVerifyDeposit simulates verifying a deposit for the admin.
func (s *AdminWalletService) AdminVerifyDeposit(depositID uint, input models.DepositVerifyInput) (*models.DepositResponse, error) {
	depositRequest, err := s.Repo.GetDepositRequestByID(depositID)
	if err != nil {
		return nil, fmt.Errorf("deposit request not found: %w", err)
	}

	if depositRequest.Status != models.TxStatusPending {
		return nil, errors.New("deposit request is not in pending state")
	}

	if input.Status != string(models.TxStatusSuccess) {
		depositRequest.Status = models.TxStatusFailed
		s.Repo.UpdateDepositRequest(depositRequest)
		return nil, errors.New("deposit verification failed as per input status")
	}

	// Start a database transaction
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		wallet, err := s.Repo.GetAdminWallet() // Get admin wallet within the transaction
		if err != nil {
			return err
		}

		balanceBefore := wallet.Balance
		wallet.Balance += input.Amount // Update balance
		wallet.LastUpdated = time.Now()

		if err := s.Repo.UpdateWalletBalance(tx, wallet); err != nil {
			return fmt.Errorf("failed to update admin wallet balance: %w", err)
		}

		// Update deposit request status
		depositRequest.Status = models.TxStatusSuccess
		depositRequest.PaymentGatewayTxID = input.PaymentGatewayTxID // Assuming webhook provides this
		if err := s.Repo.UpdateDepositRequest(depositRequest); err != nil {
			return fmt.Errorf("failed to update admin deposit request status: %w", err)
		}

		// Create a wallet transaction record
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
			return fmt.Errorf("failed to create admin wallet transaction record: %w", err)
		}

		depositRequest.WalletTransactionID = &walletTx.ID // Link transaction to deposit request
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

func (s *AdminWalletService) AdminRequestWithdrawal(input models.WithdrawalRequestInput) (*models.WithdrawalResponse, error) {
	adminUser, err := s.Repo.FindAdminUser()
	if err != nil {
		return nil, err
	}
	wallet, err := s.Repo.GetAdminWallet()
	if err != nil {
		return nil, err
	}

	if wallet.Balance < input.Amount {
		return nil, errors.New("insufficient balance for withdrawal")
	}
	if wallet.Currency != input.Currency {
		return nil, errors.New("currency mismatch for withdrawal")
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		balanceBefore := wallet.Balance
		wallet.Balance -= input.Amount // Deduct balance immediately
		wallet.LastUpdated = time.Now()

		if err := s.Repo.UpdateWalletBalance(tx, wallet); err != nil {
			return fmt.Errorf("failed to update admin wallet balance during withdrawal: %w", err)
		}

		withdrawRequest := models.WithdrawRequest{
			UserID:             adminUser.ID,
			Amount:             input.Amount,
			Currency:           input.Currency,
			Status:             models.TxStatusPending, // Withdrawal requests usually go into pending for review/processing
			BeneficiaryAccount: input.BeneficiaryAccount,
			PaymentGateway:     "AdminManualTransfer", // Or actual payment gateway
			PaymentGatewayTxID: "",
		}
		if err := s.Repo.CreateWithdrawRequest(&withdrawRequest); err != nil {
			return fmt.Errorf("failed to create admin withdrawal request: %w", err)
		}

		// Create a pending wallet transaction for the withdrawal
		walletTx := models.WalletTransaction{
			WalletID:           wallet.ID,
			UserID:             wallet.UserID,
			TransactionType:    models.TxTypeWithdraw,
			Amount:             input.Amount,
			Currency:           wallet.Currency,
			Status:             models.TxStatusPending,
			ReferenceID:        fmt.Sprintf("WITHDRAW-%d", withdrawRequest.ID),
			PaymentGatewayTxID: "", // Will be updated upon actual processing
			Description:        fmt.Sprintf("Admin withdrawal request to %s", input.BeneficiaryAccount),
			BalanceBefore:      balanceBefore,
			BalanceAfter:       wallet.Balance, // Balance is deducted
		}
		if err := s.Repo.CreateWalletTransaction(tx, &walletTx); err != nil {
			return fmt.Errorf("failed to create admin withdrawal transaction record: %w", err)
		}
		withdrawRequest.WalletTransactionID = &walletTx.ID    // Link transaction to withdrawal request
		return s.Repo.UpdateWithdrawRequest(&withdrawRequest) // Corrected: Pass pointer
	})

	if err != nil {
		return nil, err
	}

	// The `WithdrawalID` in the response needs to come from the `withdrawRequest` object after it's created and saved
	// The ID is generated by GORM after `CreateWithdrawRequest`
	// So, we need to access withdrawRequest.ID after the transaction
	// However, `withdrawRequest` is defined inside the transaction closure.
	// We need to declare it outside if we want to use its ID in the return.
	// For simplicity, I'll assume the repository returns the created request with ID
	// or we can refactor `AdminRequestWithdrawal` to get the ID.
	// For now, let's assume `withdrawRequest` (from inside the transaction) is available and populated with ID.
	// A more robust solution would be to return `withdrawRequest` from the transaction or retrieve it.
	// Given the original code, the `withdrawRequest` defined in this scope would not have the ID.
	// I'll add a placeholder and note this for a potential follow-up if it's not correctly propagated.
	// To actually get the ID, you'd need to fetch `withdrawRequest.ID` after the transaction completes.
	// For the compiler, let's use a dummy value or ensure `withdrawRequest` is accessible if its ID is needed.
	// A simpler fix for now is to move the declaration outside.

	// Let's modify AdminRequestWithdrawal to properly return the ID.
	// Re-reading your code, it seems `withdrawRequest` is correctly populated with ID *within* the transaction
	// but the `WithdrawalResponse` is created *outside* the transaction with a dummy `0`.
	// We need to move the `WithdrawalResponse` creation inside or store the ID from `withdrawRequest`.

	// I will declare withdrawRequest outside the transaction to capture its ID
	var finalWithdrawRequest models.WithdrawRequest

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		wallet, err := s.Repo.GetAdminWallet()
		if err != nil {
			return err
		}

		if wallet.Balance < input.Amount {
			return errors.New("insufficient balance for withdrawal (re-check in transaction)")
		}
		if wallet.Currency != input.Currency {
			return errors.New("currency mismatch for withdrawal (re-check in transaction)")
		}

		balanceBefore := wallet.Balance
		wallet.Balance -= input.Amount // Deduct balance immediately
		wallet.LastUpdated = time.Now()

		if err := s.Repo.UpdateWalletBalance(tx, wallet); err != nil {
			return fmt.Errorf("failed to update admin wallet balance during withdrawal: %w", err)
		}

		withdrawRequest := models.WithdrawRequest{
			UserID:             adminUser.ID,
			Amount:             input.Amount,
			Currency:           input.Currency,
			Status:             models.TxStatusPending,
			BeneficiaryAccount: input.BeneficiaryAccount,
			PaymentGateway:     "AdminManualTransfer",
			PaymentGatewayTxID: "",
		}
		if err := s.Repo.CreateWithdrawRequest(&withdrawRequest); err != nil {
			return fmt.Errorf("failed to create admin withdrawal request: %w", err)
		}
		finalWithdrawRequest = withdrawRequest // Assign to outer scope variable

		walletTx := models.WalletTransaction{
			WalletID:           wallet.ID,
			UserID:             wallet.UserID,
			TransactionType:    models.TxTypeWithdraw,
			Amount:             input.Amount,
			Currency:           wallet.Currency,
			Status:             models.TxStatusPending,
			ReferenceID:        fmt.Sprintf("WITHDRAW-%d", withdrawRequest.ID),
			PaymentGatewayTxID: "",
			Description:        fmt.Sprintf("Admin withdrawal request to %s", input.BeneficiaryAccount),
			BalanceBefore:      balanceBefore,
			BalanceAfter:       wallet.Balance,
		}
		if err := s.Repo.CreateWalletTransaction(tx, &walletTx); err != nil {
			return fmt.Errorf("failed to create admin withdrawal transaction record: %w", err)
		}
		finalWithdrawRequest.WalletTransactionID = &walletTx.ID
		return s.Repo.UpdateWithdrawRequest(&finalWithdrawRequest)
	})

	if err != nil {
		return nil, err
	}

	return &models.WithdrawalResponse{
		WithdrawalID:       finalWithdrawRequest.ID, // Corrected: use actual ID
		Amount:             input.Amount,
		Currency:           input.Currency,
		Status:             models.TxStatusPending,
		Message:            "Admin withdrawal request submitted for processing.",
		PaymentGatewayTxID: "",
	}, nil
}

// AdminGetWalletTransactions retrieves all transactions for the admin's wallet.
func (s *AdminWalletService) AdminGetWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	return s.Repo.GetAdminWalletTransactions(pagination)
}

// CreditAdminWallet is a helper function to credit the admin's wallet, specifically for subscription payments.
func (s *AdminWalletService) CreditAdminWallet(tx *gorm.DB, amount float64, currency, description string) error {
	adminWallet, err := s.Repo.GetAdminWallet() // Get admin wallet within the transaction
	if err != nil {
		return fmt.Errorf("failed to get admin wallet: %w", err)
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

	// Create a wallet transaction record for the credit
	walletTx := models.WalletTransaction{
		WalletID:        adminWallet.ID,
		UserID:          adminWallet.UserID,
		TransactionType: models.TxTypeDeposit, // Treat as a deposit to admin wallet
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

// GetPendingWithdrawalRequests retrieves all pending withdrawal requests from customers.
func (s *AdminWalletService) GetPendingWithdrawalRequests(pagination models.PaginationParams) ([]models.WithdrawRequest, int64, error) {
	return s.Repo.GetPendingWithdrawalRequests(pagination)
}

// ApproveWithdrawalRequest approves a customer's withdrawal request.
func (s *AdminWalletService) ApproveWithdrawalRequest(withdrawalID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		withdrawal, err := s.Repo.GetWithdrawRequestByID(withdrawalID)
		if err != nil {
			return fmt.Errorf("withdrawal request not found: %w", err)
		}
		if withdrawal.Status != models.TxStatusPending {
			return errors.New("withdrawal request is not pending")
		}

		withdrawal.Status = models.TxStatusSuccess
		// Simulate integration with a payment gateway here to actually disburse funds.
		// For now, we'll just update the status.
		withdrawal.PaymentGatewayTxID = fmt.Sprintf("PAYOUT-%d-%s", withdrawal.ID, time.Now().Format("20060102")) // Mock ID

		if err := s.Repo.UpdateWithdrawRequest(withdrawal); err != nil { // Corrected: Removed `tx`
			return fmt.Errorf("failed to update withdrawal request status: %w", err)
		}

		// Update the corresponding WalletTransaction status to SUCCESS
		if withdrawal.WalletTransactionID != nil {
			var walletTx models.WalletTransaction
			// Ensure walletTx is fetched within the current transaction `tx`
			if err := tx.First(&walletTx, *withdrawal.WalletTransactionID).Error; err != nil {
				log.Printf("Warning: Wallet transaction for withdrawal %d not found: %v", withdrawal.ID, err)
			} else {
				walletTx.Status = models.TxStatusSuccess
				if err := tx.Save(&walletTx).Error; err != nil {
					return fmt.Errorf("failed to update linked wallet transaction status: %w", err)
				}
			}
		}

		return nil
	})
}

func (s *AdminWalletService) RejectWithdrawalRequest(withdrawalID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		withdrawal, err := s.Repo.GetWithdrawRequestByID(withdrawalID)
		if err != nil {
			return fmt.Errorf("withdrawal request not found: %w", err)
		}
		if withdrawal.Status != models.TxStatusPending {
			return errors.New("withdrawal request is not pending")
		}

		customerWallet, err := s.Repo.GetCustomerWallet(withdrawal.UserID)
		if err != nil {

			return fmt.Errorf("customer wallet not found for user %d: %w", withdrawal.UserID, err)
		}

		balanceBeforeReversal := customerWallet.Balance
		customerWallet.Balance += withdrawal.Amount
		customerWallet.LastUpdated = time.Now()

		if err := s.Repo.UpdateCustomerWalletBalance(tx, customerWallet); err != nil {
			return fmt.Errorf("failed to revert funds to customer wallet: %w", err)
		}

		withdrawal.Status = models.TxStatusRejected
		if err := s.Repo.UpdateWithdrawRequest(withdrawal); err != nil {
			return fmt.Errorf("failed to update withdrawal request status to rejected: %w", err)
		}

		reversalTx := models.WalletTransaction{
			WalletID:        customerWallet.ID,
			UserID:          customerWallet.UserID,
			TransactionType: models.TxTypeReversal,
			Amount:          withdrawal.Amount,
			Currency:        withdrawal.Currency,
			Status:          models.TxStatusSuccess,
			ReferenceID:     fmt.Sprintf("REVERSAL-WITHDRAW-%d", withdrawal.ID),
			Description:     fmt.Sprintf("Withdrawal request %d rejected, funds returned.", withdrawal.ID),
			BalanceBefore:   balanceBeforeReversal,
			BalanceAfter:    customerWallet.Balance,
		}
		if err := s.Repo.CreateWalletTransaction(tx, &reversalTx); err != nil {
			return fmt.Errorf("failed to create reversal transaction for customer: %w", err)
		}

		if withdrawal.WalletTransactionID != nil {
			var originalWalletTx models.WalletTransaction
			if err := tx.First(&originalWalletTx, *withdrawal.WalletTransactionID).Error; err != nil {
				log.Printf("Warning: Original wallet transaction for withdrawal %d not found: %v", withdrawal.ID, err)
			} else {
				originalWalletTx.Status = models.TxStatusReversed
				if err := tx.Save(&originalWalletTx).Error; err != nil {
					return fmt.Errorf("failed to update original linked wallet transaction status: %w", err)
				}
			}
		}

		return nil
	})
}
