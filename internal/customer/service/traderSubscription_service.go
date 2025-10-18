// internal/customer/service/trader_subscription.go - **NEW FILE**
package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ICustomerTraderSignalSubscriptionService interface {
	GetAvailableTradersWithPlans(ctx context.Context) ([]models.User, error)
	SubscribeToTrader(ctx context.Context, customerID uint, input models.SubscribeToTraderInput) error
	GetSubscribedTradersSignals(ctx context.Context, customerID uint) ([]models.Signal, error)
	GetActiveSubscriptions(ctx context.Context, customerID uint) ([]models.CustomerTraderSignalSubscription, error)
	IsCustomerSubscribedToTrader(ctx context.Context, customerID, traderID uint) (bool, error)
}

type CustomerTraderSignalSubscriptionService struct {
	repo customerrepo.ICustomerTraderSignalSubscriptionRepository
	db   *gorm.DB // Pass the DB instance for transactions
}

func NewCustomerTraderSignalSubscriptionService(repo customerrepo.ICustomerTraderSignalSubscriptionRepository, db *gorm.DB) ICustomerTraderSignalSubscriptionService {
	return &CustomerTraderSignalSubscriptionService{repo: repo, db: db}
}

func (s *CustomerTraderSignalSubscriptionService) GetAvailableTradersWithPlans(ctx context.Context) ([]models.User, error) {
	return s.repo.GetTradersWithPlans(ctx)
}

func (s *CustomerTraderSignalSubscriptionService) SubscribeToTrader(ctx context.Context, customerID uint, input models.SubscribeToTraderInput) error {
	log.Printf("Attempting to subscribe customer %d to plan ID %d", customerID, input.TraderSubscriptionPlanID)

	plan, err := s.repo.GetTraderSubscriptionPlanByID(ctx, input.TraderSubscriptionPlanID)
	if err != nil {
		log.Printf("Error getting trader subscription plan ID %d: %v", input.TraderSubscriptionPlanID, err)
		if errors.Is(err, errors.New("trader subscription plan not found")) { // Match the error string from repo
			return errors.New("invalid trader subscription plan: trader subscription plan not found")
		}
		return fmt.Errorf("failed to get trader subscription plan: %w", err)
	}
	log.Printf("Found plan: %+v", plan)

	if !plan.IsActive {
		return fmt.Errorf("trader subscription plan is not active")
	}

	// Check if customer is already subscribed to this plan
	isSubscribed, err := s.repo.IsCustomerSubscribedToPlan(ctx, customerID, plan.ID)
	if err != nil {
		return fmt.Errorf("failed to check existing subscription: %w", err)
	}
	if isSubscribed {
		return fmt.Errorf("you are already subscribed to this plan")
	}

	// Start a database transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Deduct money from customer's wallet
	customerWallet, err := s.repo.GetTraderWallet(ctx, customerID) // Assuming GetTraderWallet can get any user's wallet
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get customer wallet: %w", err)
	}

	if customerWallet.Balance < plan.Price {
		tx.Rollback()
		return fmt.Errorf("insufficient funds in wallet. current balance: %.2f, required: %.2f", customerWallet.Balance, plan.Price)
	}

	if err := s.repo.UpdateWalletBalance(ctx, customerID, -plan.Price, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deduct funds from customer wallet: %w", err)
	}
	customerBalanceAfter := customerWallet.Balance - plan.Price

	// 2. Distribute funds: Admin commission + Trader revenue
	adminCommissionAmount := plan.Price * (plan.AdminCommission / 100.0)
	traderRevenueAmount := plan.Price - adminCommissionAmount
	// Get Admin wallet and update balance
	adminWallet, err := s.repo.GetAdminWallet(ctx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get admin wallet: %w", err)
	}
	if err := s.repo.UpdateWalletBalance(ctx, adminWallet.UserID, adminCommissionAmount, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to credit admin commission: %w", err)
	}
	adminBalanceAfter := adminWallet.Balance + adminCommissionAmount

	// Get Trader wallet and update balance
	traderWallet, err := s.repo.GetTraderWallet(ctx, plan.TraderID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get trader wallet: %w", err)
	}
	if err := s.repo.UpdateWalletBalance(ctx, plan.TraderID, traderRevenueAmount, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to credit trader revenue: %w", err)
	}
	traderBalanceAfter := traderWallet.Balance + traderRevenueAmount

	customerTx := models.WalletTransaction{
		WalletID:        customerWallet.ID,
		UserID:          customerID,
		Type:            models.TxTypeSubscription,
		TransactionType: models.TxTypeDebit,
		Name:            "Trader Subscription Debit",
		Amount:          plan.Price,
		Currency:        plan.Currency,
		Status:          models.TxStatusSuccess,
		Description:     fmt.Sprintf("Subscription to trader %d's plan '%s'", plan.TraderID, plan.Name),
		BalanceBefore:   customerWallet.Balance,
		BalanceAfter:    customerBalanceAfter,
		ReferenceID:     fmt.Sprintf("TRADER_SUB_%d_PLAN_%d", customerID, plan.ID),
		TransactionID:   fmt.Sprintf("TRADER_SUB_%d_%d_%d", customerID, plan.TraderID, time.Now().UnixNano()), // unique
	}

	if err := s.repo.CreateWalletTransaction(ctx, &customerTx, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record customer transaction: %w", err)
	}

	adminTx := models.WalletTransaction{
		WalletID:        adminWallet.ID,
		UserID:          adminWallet.UserID,
		Type:            models.TxTypeAdminCommission,
		TransactionType: models.TxTypeCredit,
		Name:            "Credit from Trader Subscription Commission",
		Amount:          adminCommissionAmount,
		Currency:        plan.Currency,
		Status:          models.TxStatusSuccess,
		Description:     fmt.Sprintf("Commission from customer %d subscribing to trader %d's plan '%s'", customerID, plan.TraderID, plan.Name),
		BalanceBefore:   adminWallet.Balance,
		BalanceAfter:    adminBalanceAfter,
		ReferenceID:     fmt.Sprintf("TRADER_SUB_COMMISSION_%d_PLAN_%d", customerID, plan.ID),
		TransactionID:   fmt.Sprintf("ADMIN_COMM_%d_%d_%d", customerID, plan.TraderID, time.Now().UnixNano()),
	}
	if err := s.repo.CreateWalletTransaction(ctx, &adminTx, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record admin transaction: %w", err)
	}

	// Trader transaction (credit)
	traderTx := models.WalletTransaction{
		WalletID:        traderWallet.ID,
		UserID:          plan.TraderID,
		Type:            models.TxTypeTraderRevenue,
		TransactionType: models.TxTypeCredit,
		Name:            "Credit from Customer Subscription",
		Amount:          traderRevenueAmount,
		Currency:        plan.Currency,
		Status:          models.TxStatusSuccess,
		Description:     fmt.Sprintf("Revenue from customer %d subscribing to plan '%s'", customerID, plan.Name),
		BalanceBefore:   traderWallet.Balance,
		BalanceAfter:    traderBalanceAfter,
		ReferenceID:     fmt.Sprintf("TRADER_REVENUE_%d_PLAN_%d", customerID, plan.ID),
		TransactionID:   fmt.Sprintf("TRADER_REV_%d_%d_%d", customerID, plan.TraderID, time.Now().UnixNano()),
	}
	if err := s.repo.CreateWalletTransaction(ctx, &traderTx, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record trader transaction: %w", err)
	}

	// 4. Create CustomerTraderSubscription record
	startDate := time.Now()
	endDate := startDate.Add(time.Duration(plan.DurationDays) * 24 * time.Hour)

	newSubscription := &models.CustomerTraderSignalSubscription{
		CustomerID:               customerID,
		TraderID:                 plan.TraderID,
		TraderSubscriptionPlanID: plan.ID,
		StartDate:                startDate,
		EndDate:                  endDate,
		IsActive:                 true,
		WalletTransactionID:      &customerTx.ID, // Link to the customer's payment transaction
	}

	if _, err := s.repo.CreateCustomerTraderSubscription(ctx, newSubscription); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create customer-trader subscription record: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Customer %d successfully subscribed to trader %d's plan %d. Admin got %.2f, Trader got %.2f",
		customerID, plan.TraderID, plan.ID, adminCommissionAmount, traderRevenueAmount)

	return nil
}

func (s *CustomerTraderSignalSubscriptionService) GetSubscribedTradersSignals(ctx context.Context, customerID uint) ([]models.Signal, error) {
	signals, err := s.repo.GetAllSignalsFromSubscribedTraders(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve signals from subscribed traders: %w", err)
	}
	return signals, nil
}

func (s *CustomerTraderSignalSubscriptionService) GetActiveSubscriptions(ctx context.Context, customerID uint) ([]models.CustomerTraderSignalSubscription, error) {
	subscriptions, err := s.repo.GetActiveTraderSubscriptionsForCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscriptions: %w", err)
	}
	return subscriptions, nil
}

func (s *CustomerTraderSignalSubscriptionService) IsCustomerSubscribedToTrader(ctx context.Context, customerID, traderID uint) (bool, error) {
	return s.repo.IsCustomerSubscribedToTrader(ctx, customerID, traderID)
}
