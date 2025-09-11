package paymentgateway

import (
	"errors"
	"fmt"
	"time"
)

type SimulatedPaymentClient struct {
	// In a real client, you'd have API keys, base URLs, etc.
}

func NewSimulatedPaymentClient() *SimulatedPaymentClient {
	return &SimulatedPaymentClient{}
}

func (c *SimulatedPaymentClient) CreateDepositInitiation(amount float64, currency, customerID string) (string, string, error) {
	if amount <= 0 {
		return "", "", errors.New("deposit amount must be positive")
	}
	paymentID := fmt.Sprintf("pg_dep_%d", time.Now().UnixNano())
	redirectURL := fmt.Sprintf("https://mock-payment-gateway.com/pay?id=%s&amount=%.2f&currency=%s", paymentID, amount, currency)
	fmt.Printf("Simulated Payment Gateway: Deposit initiated for customer %s, amount %.2f, paymentID %s\n", customerID, amount, paymentID)
	return paymentID, redirectURL, nil
}

func (c *SimulatedPaymentClient) ProcessWithdrawal(amount float64, currency, beneficiaryAccount string) (string, error) {
	if amount <= 0 {
		return "", errors.New("withdrawal amount must be positive")
	}
	transactionID := fmt.Sprintf("pg_wd_%d", time.Now().UnixNano())
	fmt.Printf("Simulated Payment Gateway: Withdrawal processed for amount %.2f to account %s, transactionID %s\n", amount, beneficiaryAccount, transactionID)
	return transactionID, nil
}

func (c *SimulatedPaymentClient) VerifyDeposit(paymentID string) (bool, error) {

	if paymentID == "" {
		return false, errors.New("payment ID cannot be empty")
	}
	fmt.Printf("Simulated Payment Gateway: Verifying payment ID %s... assuming success.\n", paymentID)
	return true, nil
}
