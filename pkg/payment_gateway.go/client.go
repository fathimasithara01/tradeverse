package paymentgateway

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrPaymentGatewayFailed = errors.New("payment gateway operation failed")
	ErrInvalidAmount        = errors.New("invalid amount for payment gateway")
	ErrTransactionNotFound  = errors.New("payment gateway transaction not found")
)

type SimulatedPaymentClient interface {
	CreateDepositInitiation(amount float64, currency, userID string) (pgTxID, redirectURL string, err error)
	VerifyDeposit(pgTxID string) (isVerified bool, err error)
	ProcessWithdrawal(amount float64, currency, beneficiaryAccount string) (pgTxID string, err error)
}

type simulatedPaymentClient struct {
}

func NewSimulatedPaymentClient() SimulatedPaymentClient {
	return &simulatedPaymentClient{}
}

func (s *simulatedPaymentClient) CreateDepositInitiation(amount float64, currency, userID string) (string, string, error) {
	if amount <= 0 {
		return "", "", ErrInvalidAmount
	}
	pgTxID := fmt.Sprintf("PG_DEPOSIT_%s_%d", userID, time.Now().UnixNano())
	redirectURL := fmt.Sprintf("https://simulated-pg.com/pay?tx=%s", pgTxID)
	fmt.Printf("Simulated PG: Deposit initiated for User %s, Amount %.2f %s. PG Transaction ID: %s\n", userID, amount, currency, pgTxID)
	return pgTxID, redirectURL, nil
}

func (s *simulatedPaymentClient) VerifyDeposit(pgTxID string) (bool, error) {
	if pgTxID == "" {
		return false, ErrTransactionNotFound
	}
	fmt.Printf("Simulated PG: Deposit verification requested for PG Transaction ID: %s. Assuming successful.\n", pgTxID)
	return true, nil
}

func (s *simulatedPaymentClient) ProcessWithdrawal(amount float64, currency, beneficiaryAccount string) (string, error) {
	if amount <= 0 {
		return "", ErrInvalidAmount
	}
	pgTxID := fmt.Sprintf("PG_WITHDRAW_%s_%d", beneficiaryAccount, time.Now().UnixNano())
	fmt.Printf("Simulated PG: Withdrawal processed for Account %s, Amount %.2f %s. PG Transaction ID: %s\n", beneficiaryAccount, amount, currency, pgTxID)
	return pgTxID, nil
}
