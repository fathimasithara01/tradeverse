package service

import (
	"errors"
	"fmt"

	walletrepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo" // Changed import path
	"github.com/fathimasithara01/tradeverse/pkg/models"
	paymentgateway "github.com/fathimasithara01/tradeverse/pkg/payment_gateway.go"
	"gorm.io/gorm"
)

var (
	ErrUserWalletNotFound         = errors.New("user wallet not found")
	ErrInvalidDepositStatus       = errors.New("invalid deposit status")
	ErrDepositAlreadyProcessed    = errors.New("deposit already processed")
	ErrWithdrawalAlreadyProcessed = errors.New("withdrawal already processed")
	ErrInvalidWithdrawalStatus    = errors.New("invalid withdrawal status")
)

type WalletService interface {
	GetWalletSummary(userID uint) (*models.WalletSummaryResponse, error)
	InitiateDeposit(userID uint, input models.DepositRequestInput) (*models.DepositResponse, error)
	VerifyDeposit(depositID uint, input models.DepositVerifyInput) (*models.DepositResponse, error)
	RequestWithdrawal(userID uint, input models.WithdrawalRequestInput) (*models.WithdrawalResponse, error)
	GetTransactions(userID uint, pagination models.PaginationParams) (*models.TransactionListResponse, error)
}

type walletService struct {
	repo          walletrepo.WalletRepository
	paymentClient paymentgateway.SimulatedPaymentClient
	db            *gorm.DB
}

func NewWalletService(repo walletrepo.WalletRepository, paymentClient paymentgateway.SimulatedPaymentClient, db *gorm.DB) WalletService { // Changed type
	return &walletService{
		repo:          repo,
		paymentClient: paymentClient,
		db:            db,
	}
}

func (s *walletService) GetWalletSummary(userID uint) (*models.WalletSummaryResponse, error) {
	wallet, err := s.repo.GetUserWallet(userID)
	if err != nil {
		if errors.Is(err, walletrepo.ErrWalletNotFound) { 
			newWallet := &models.Wallet{
				UserID:   userID,
				Balance:  0,
				Currency: "INR",
			}
			if createErr := s.repo.CreateWallet(newWallet); createErr != nil {
				return nil, fmt.Errorf("failed to create wallet for user (from wallet service fallback): %w", createErr)
			}
			wallet = newWallet
		} else {
			return nil, fmt.Errorf("failed to get user wallet: %w", err)
		}
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
	_, err := s.repo.GetUserWallet(userID)
	if err != nil {
		if errors.Is(err, walletrepo.ErrWalletNotFound) { 
			return nil, ErrUserWalletNotFound
		}
		return nil, fmt.Errorf("failed to get user wallet: %w", err)
	}

	depositReq := &models.DepositRequest{
		UserID:         userID,
		Amount:         input.Amount,
		Currency:       input.Currency,
		Status:         models.TxStatusPending,
		PaymentGateway: "SimulatedPG",
	}
	if err := s.repo.CreateDepositRequest(depositReq); err != nil {
		return nil, fmt.Errorf("failed to create deposit request: %w", err)
	}

	pgTxID, redirectURL, err := s.paymentClient.CreateDepositInitiation(input.Amount, input.Currency, fmt.Sprintf("%d", userID))
	if err != nil {
		depositReq.Status = models.TxStatusFailed
		s.repo.UpdateDepositRequest(depositReq)
		return nil, fmt.Errorf("payment gateway initiation failed: %w", err)
	}

	depositReq.PaymentGatewayTxID = pgTxID
	depositReq.RedirectURL = redirectURL
	if err := s.repo.UpdateDepositRequest(depositReq); err != nil {
		return nil, fmt.Errorf("failed to update deposit request with PG info: %w", err)
	}

	return &models.DepositResponse{
		DepositID:          depositReq.ID,
		Amount:             depositReq.Amount,
		Currency:           depositReq.Currency,
		Status:             depositReq.Status,
		RedirectURL:        depositReq.RedirectURL,
		PaymentGatewayTxID: depositReq.PaymentGatewayTxID,
		Message:            "Deposit initiated. Please complete payment via redirect URL.",
	}, nil
}

func (s *walletService) VerifyDeposit(depositID uint, input models.DepositVerifyInput) (*models.DepositResponse, error) {
	depositReq, err := s.repo.GetDepositRequestByID(depositID)
	if err != nil {
		return nil, err
	}

	if depositReq.Status == models.TxStatusSuccess || depositReq.Status == models.TxStatusFailed {
		return nil, ErrDepositAlreadyProcessed
	}
	if depositReq.Status != models.TxStatusPending {
		return nil, ErrInvalidDepositStatus
	}

	isVerified, err := s.paymentClient.VerifyDeposit(input.PaymentGatewayTxID)
	if err != nil || !isVerified || input.Status != string(models.TxStatusSuccess) {
		depositReq.Status = models.TxStatusFailed
		s.repo.UpdateDepositRequest(depositReq)
		return nil, fmt.Errorf("deposit verification failed: %w", err)
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		depositReq.Status = models.TxStatusSuccess
		depositReq.PaymentGatewayTxID = input.PaymentGatewayTxID
		if err := tx.Save(depositReq).Error; err != nil {
			return fmt.Errorf("failed to update deposit request status: %w", err)
		}

		wallet, err := s.repo.GetUserWallet(depositReq.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user wallet during deposit verification: %w", err)
		}

		err = s.repo.CreditWallet(tx, wallet.ID, depositReq.Amount, models.TxTypeDeposit,
			fmt.Sprintf("DEPOSIT_%d_PGTX_%s", depositReq.ID, depositReq.PaymentGatewayTxID),
			fmt.Sprintf("Deposit from Payment Gateway (Request ID: %d)", depositReq.ID),
		)
		if err != nil {
			return fmt.Errorf("failed to credit wallet during deposit verification: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &models.DepositResponse{
		DepositID:          depositReq.ID,
		Amount:             depositReq.Amount,
		Currency:           depositReq.Currency,
		Status:             depositReq.Status,
		PaymentGatewayTxID: depositReq.PaymentGatewayTxID,
		Message:            "Deposit successful.",
	}, nil
}

func (s *walletService) RequestWithdrawal(userID uint, input models.WithdrawalRequestInput) (*models.WithdrawalResponse, error) {
	wallet, err := s.repo.GetUserWallet(userID)
	if err != nil {
		if errors.Is(err, walletrepo.ErrWalletNotFound) {
			return nil, ErrUserWalletNotFound
		}
		return nil, fmt.Errorf("failed to get user wallet: %w", err)
	}

	if wallet.Balance < input.Amount {
		return nil, walletrepo.ErrInsufficientFunds 
	}

	withdrawReq := &models.WithdrawRequest{
		UserID:             userID,
		Amount:             input.Amount,
		Currency:           input.Currency,
		BeneficiaryAccount: input.BeneficiaryAccount,
		Status:             models.TxStatusPending,
		PaymentGateway:     "SimulatedPG",
	}
	if err := s.repo.CreateWithdrawRequest(withdrawReq); err != nil {
		return nil, fmt.Errorf("failed to create withdrawal request: %w", err)
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		err := s.repo.DebitWallet(tx, wallet.ID, input.Amount, models.TxTypeWithdraw,
			fmt.Sprintf("WITHDRAW_REQ_%d", withdrawReq.ID),
			fmt.Sprintf("Withdrawal request to %s (Request ID: %d)", input.BeneficiaryAccount, withdrawReq.ID),
		)
		if err != nil {
			withdrawReq.Status = models.TxStatusFailed
			tx.Save(withdrawReq)
			return err
		}

		pgTxID, pgErr := s.paymentClient.ProcessWithdrawal(input.Amount, input.Currency, input.BeneficiaryAccount)
		if pgErr != nil {
		
			withdrawReq.Status = models.TxStatusFailed
			tx.Save(withdrawReq) 
			return fmt.Errorf("payment gateway withdrawal processing failed: %w", pgErr)
		}
		withdrawReq.PaymentGatewayTxID = pgTxID
		withdrawReq.Status = models.TxStatusSuccess

		if err := tx.Save(withdrawReq).Error; err != nil {
			return fmt.Errorf("failed to update withdrawal request status: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &models.WithdrawalResponse{
		WithdrawalID:       withdrawReq.ID,
		Amount:             withdrawReq.Amount,
		Currency:           withdrawReq.Currency,
		Status:             withdrawReq.Status,
		PaymentGatewayTxID: withdrawReq.PaymentGatewayTxID,
		Message:            "Withdrawal request processed successfully.",
	}, nil
}

func (s *walletService) GetTransactions(userID uint, pagination models.PaginationParams) (*models.TransactionListResponse, error) {
	transactions, total, err := s.repo.GetWalletTransactions(userID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve wallet transactions: %w", err)
	}

	return &models.TransactionListResponse{
		Transactions: transactions,
		Total:        total,
		Page:         pagination.Page,
		Limit:        pagination.Limit,
	}, nil
}
