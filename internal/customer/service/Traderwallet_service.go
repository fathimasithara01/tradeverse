package service

import (
	"errors"
	"fmt"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	paymentgateway "github.com/fathimasithara01/tradeverse/pkg/payment_gateway.go"
)

type CustomerWalletService interface {
	GetBalance(userID uint) (*models.WalletSummaryResponse, error)
	Deposit(userID uint, input models.DepositRequestInput) (*models.DepositResponse, error)
	VerifyDeposit(pgTxID string, userID uint) error
	Withdraw(userID uint, input models.WithdrawalRequestInput) (*models.WithdrawalResponse, error)
}

type customerWalletService struct {
	repo walletrepo.CustomerWalletRepository
	pg   paymentgateway.SimulatedPaymentClient
}

func NewCustomerWalletService(repo walletrepo.CustomerWalletRepository, pg paymentgateway.SimulatedPaymentClient) CustomerWalletService {
	return &customerWalletService{repo: repo, pg: pg}
}

func (s *customerWalletService) GetBalance(userID uint) (*models.WalletSummaryResponse, error) {
	wallet, err := s.repo.GetWalletByUserID(userID)
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

func (s *customerWalletService) Deposit(userID uint, input models.DepositRequestInput) (*models.DepositResponse, error) {
	pgTxID, redirectURL, err := s.pg.CreateDepositInitiation(input.Amount, input.Currency, fmt.Sprint(userID))
	if err != nil {
		return nil, err
	}

	dr := &models.DepositRequest{
		UserID:             userID,
		Amount:             input.Amount,
		Currency:           input.Currency,
		Status:             models.TxStatusPending,
		PaymentGateway:     "SimulatedPG",
		PaymentGatewayTxID: pgTxID,
		RedirectURL:        redirectURL,
	}

	if err := s.repo.CreateDepositRequest(dr); err != nil {
		return nil, err
	}

	return &models.DepositResponse{
		DepositID:          dr.ID,
		Amount:             dr.Amount,
		Currency:           dr.Currency,
		Status:             dr.Status,
		RedirectURL:        dr.RedirectURL,
		PaymentGatewayTxID: pgTxID,
		Message:            "Deposit initiated",
	}, nil
}

func (s *customerWalletService) VerifyDeposit(pgTxID string, userID uint) error {
	verified, err := s.pg.VerifyDeposit(pgTxID)
	if err != nil || !verified {
		return errors.New("deposit verification failed")
	}

	// Credit wallet
	if err := s.repo.UpdateWalletBalance(userID, 100); err != nil { // Amount should be fetched from DepositRequest
		return err
	}

	// Record transaction
	tx := &models.WalletTransaction{
		WalletID:        userID,
		UserID:          userID,
		TransactionType: models.TxTypeDeposit,
		Amount:          100, // Example
		Currency:        "USD",
		Status:          models.TxStatusSuccess,
		Description:     "Deposit verified",
	}
	return s.repo.CreateTransaction(tx)
}

func (s *customerWalletService) Withdraw(userID uint, input models.WithdrawalRequestInput) (*models.WithdrawalResponse, error) {
	wallet, err := s.repo.GetWalletByUserID(userID)
	if err != nil {
		return nil, err
	}
	if wallet.Balance < input.Amount {
		return nil, errors.New("insufficient funds")
	}

	pgTxID, err := s.pg.ProcessWithdrawal(input.Amount, input.Currency, input.BeneficiaryAccount)
	if err != nil {
		return nil, err
	}

	// Deduct balance
	if err := s.repo.UpdateWalletBalance(userID, -input.Amount); err != nil {
		return nil, err
	}

	// Record withdrawal
	wr := &models.WithdrawRequest{
		UserID:             userID,
		Amount:             input.Amount,
		Currency:           input.Currency,
		Status:             models.TxStatusSuccess,
		BeneficiaryAccount: input.BeneficiaryAccount,
		PaymentGateway:     "SimulatedPG",
		PaymentGatewayTxID: pgTxID,
	}
	if err := s.repo.CreateWithdrawRequest(wr); err != nil {
		return nil, err
	}

	return &models.WithdrawalResponse{
		WithdrawalID:       wr.ID,
		Amount:             wr.Amount,
		Currency:           wr.Currency,
		Status:             wr.Status,
		PaymentGatewayTxID: wr.PaymentGatewayTxID,
		Message:            "Withdrawal successful",
	}, nil
}
