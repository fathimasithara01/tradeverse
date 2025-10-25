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

	GetAllWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error)

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

func (s *AdminWalletService) GetAllCustomerTransactionsWithUserDetails(pagination models.PaginationParams) ([]models.AdminTransactionDisplayDTO, int64, error) {
	return s.Repo.GetAllCustomerTransactions(pagination)
}

func (s *AdminWalletService) GetAllWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	return s.Repo.GetAllWalletTransactions(pagination) // This needs to be implemented in repository
}

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
		PaymentGateway: "AdminManual",
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

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		wallet, err := s.Repo.GetAdminWallet()
		if err != nil {
			return err
		}

		balanceBefore := wallet.Balance
		wallet.Balance += input.Amount
		wallet.LastUpdated = time.Now()

		if err := s.Repo.UpdateWalletBalance(tx, wallet); err != nil {
			return fmt.Errorf("failed to update admin wallet balance: %w", err)
		}

		depositRequest.Status = models.TxStatusSuccess
		depositRequest.PaymentGatewayTxID = input.PaymentGatewayTxID
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
			return fmt.Errorf("failed to create admin wallet transaction record: %w", err)
		}

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
		wallet.Balance -= input.Amount
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

		walletTx := models.WalletTransaction{
			WalletID:           wallet.ID,
			UserID:             wallet.UserID,
			TransactionType:    models.TxTypeWithdrawal,
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
		withdrawRequest.WalletTransactionID = &walletTx.ID
		return s.Repo.UpdateWithdrawRequest(&withdrawRequest)
	})

	if err != nil {
		return nil, err
	}

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
		wallet.Balance -= input.Amount
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
		finalWithdrawRequest = withdrawRequest

		walletTx := models.WalletTransaction{
			WalletID:           wallet.ID,
			UserID:             wallet.UserID,
			TransactionType:    models.TxTypeWithdrawal,
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
		WithdrawalID:       finalWithdrawRequest.ID,
		Amount:             input.Amount,
		Currency:           input.Currency,
		Status:             models.TxStatusPending,
		Message:            "Admin withdrawal request submitted for processing.",
		PaymentGatewayTxID: "",
	}, nil
}

func (s *AdminWalletService) AdminGetWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	return s.Repo.AdminGetWalletTransactions(pagination)
}

func (s *AdminWalletService) CreditAdminWallet(tx *gorm.DB, amount float64, currency, description string) error {
	adminWallet, err := s.Repo.GetAdminWallet()
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

	walletTx := models.WalletTransaction{
		WalletID:        adminWallet.ID,
		UserID:          adminWallet.UserID,
		TransactionType: models.TxTypeDeposit,
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

func (s *AdminWalletService) GetPendingWithdrawalRequests(pagination models.PaginationParams) ([]models.WithdrawRequest, int64, error) {
	return s.Repo.GetPendingWithdrawalRequests(pagination)
}

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

		withdrawal.PaymentGatewayTxID = fmt.Sprintf("PAYOUT-%d-%s", withdrawal.ID, time.Now().Format("20060102")) // Mock ID

		if err := s.Repo.UpdateWithdrawRequest(withdrawal); err != nil { // Corrected: Removed `tx`
			return fmt.Errorf("failed to update withdrawal request status: %w", err)
		}

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
