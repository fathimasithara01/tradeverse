package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrAlreadySubscribed  = errors.New("customer is already actively subscribed to this trader")
	ErrSelfSubscription   = errors.New("traders cannot subscribe to their own plans")
	ErrPlanNotForTrader   = errors.New("the selected plan is not a valid trader subscription plan")
	ErrTraderNotMatchPlan = errors.New("the selected plan does not belong to the specified trader")
)

type ITraderSubscriptionService interface {
	SubscribeToTrader(ctx context.Context, customerID, traderID, planID uint) (*models.TraderSubscription, error)
	IsCustomerSubscribedToTrader(ctx context.Context, customerID, traderID uint) (bool, error)
	GetCustomerTraderSubscriptions(ctx context.Context, customerID uint) ([]models.TraderSubscription, error)
	GetTraderPlans(ctx context.Context, traderID uint) ([]models.SubscriptionPlan, error)
	// Additional methods if needed for managing subscriptions (e.g., cancelling)
}

type traderSubscriptionService struct {
	repo customerrepo.ITraderSubscriptionRepository
	db   *gorm.DB
}

func NewTraderSubscriptionService(repo customerrepo.ITraderSubscriptionRepository, db *gorm.DB) ITraderSubscriptionService {
	return &traderSubscriptionService{repo: repo, db: db}
}

// SubscribeToTrader handles the entire subscription process including payment and commission split.
func (s *traderSubscriptionService) SubscribeToTrader(ctx context.Context, customerID, traderID, planID uint) (*models.TraderSubscription, error) {
	if customerID == traderID {
		return nil, ErrSelfSubscription
	}

	// 1. Check if customer already has an active subscription to this trader
	existingSub, err := s.repo.GetActiveTraderSubscription(customerID, traderID)
	if err != nil && !errors.Is(err, customerrepo.ErrSubscriptionNotFound) {
		return nil, fmt.Errorf("failed to check existing subscription: %w", err)
	}
	if existingSub != nil {
		return nil, ErrAlreadySubscribed
	}

	// 2. Get the subscription plan
	plan, err := s.repo.GetTraderSubscriptionPlan(planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}
	if plan.TraderID == nil || *plan.TraderID != traderID {
		return nil, ErrTraderNotMatchPlan
	}

	// 3. Get customer's wallet
	customerWallet, err := s.repo.GetUserWallet(customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer wallet: %w", err)
	}

	if customerWallet.Balance < plan.Price {
		return nil, ErrInsufficientFunds
	}

	// 4. Get trader's wallet
	traderWallet, err := s.repo.GetUserWallet(traderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trader wallet: %w", err)
	}

	// 5. Get admin's user and wallet for commission
	adminUser, err := s.repo.GetAdminUser()
	if err != nil {
		return nil, fmt.Errorf("failed to find admin user for commission: %w", err)
	}
	adminWallet, err := s.repo.GetUserWallet(adminUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin wallet for commission: %w", err)
	}

	var newSubscription *models.TraderSubscription
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Debit customer's wallet for the full price
		debitDesc := fmt.Sprintf("Subscription to Trader %d (Plan: %s)", traderID, plan.Name)
		if err := s.repo.DebitWallet(tx, customerWallet.ID, plan.Price, models.TxTypeSubscription, fmt.Sprintf("SUB_CUS_%d_T_%d_P_%d", customerID, traderID, planID), debitDesc); err != nil {
			return fmt.Errorf("failed to debit customer wallet: %w", err)
		}

		// Calculate commission
		commissionAmount := plan.Price * plan.CommissionRate
		traderShare := plan.Price - commissionAmount

		// Credit admin's wallet with commission
		adminCreditDesc := fmt.Sprintf("Commission from Customer %d's subscription to Trader %d (Plan: %s)", customerID, traderID, plan.Name)
		if err := s.repo.CreditWallet(tx, adminWallet.ID, commissionAmount, models.TxTypeFee, fmt.Sprintf("COMM_CUS_%d_T_%d_P_%d", customerID, traderID, planID), adminCreditDesc); err != nil {
			return fmt.Errorf("failed to credit admin wallet with commission: %w", err)
		}

		// Credit trader's wallet with remaining amount
		traderCreditDesc := fmt.Sprintf("Earnings from Customer %d's subscription (Plan: %s)", customerID, plan.Name)
		if err := s.repo.CreditWallet(tx, traderWallet.ID, traderShare, models.TxTypeSubscription, fmt.Sprintf("EARN_CUS_%d_T_%d_P_%d", customerID, traderID, planID), traderCreditDesc); err != nil {
			return fmt.Errorf("failed to credit trader wallet: %w", err)
		}

		// Create the TraderSubscription record
		now := time.Now()
		endDate := now.Add(time.Duration(plan.Duration) * s.getDurationUnit(plan.Interval))

		newSubscription = &models.TraderSubscription{
			UserID:                   customerID,
			TraderID:                 traderID,
			TraderSubscriptionPlanID: planID,
			StartDate:                now,
			EndDate:                  endDate,
			IsActive:                 true,
			PaymentStatus:            "paid", // Assuming immediate payment from wallet
			AmountPaid:               plan.Price,
			TransactionID:            fmt.Sprintf("SUB_TX_%d_%d", customerID, now.Unix()), // Generate a unique transaction ID
		}
		if err := s.repo.CreateTraderSubscription(tx, newSubscription); err != nil {
			return fmt.Errorf("failed to create trader subscription record: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return newSubscription, nil
}

// Helper function to convert interval string to time.Duration
func (s *traderSubscriptionService) getDurationUnit(interval string) time.Duration {
	switch interval {
	case "day":
		return 24 * time.Hour
	case "month":
		return 30 * 24 * time.Hour // Approximation
	case "year":
		return 365 * 24 * time.Hour // Approximation
	default:
		return 30 * 24 * time.Hour // Default to month if unknown
	}
}

// IsCustomerSubscribedToTrader checks if a customer has an active subscription to a specific trader.
func (s *traderSubscriptionService) IsCustomerSubscribedToTrader(ctx context.Context, customerID, traderID uint) (bool, error) {
	sub, err := s.repo.GetActiveTraderSubscription(customerID, traderID)
	if err != nil {
		if errors.Is(err, customerrepo.ErrSubscriptionNotFound) {
			return false, nil // No active subscription found
		}
		return false, fmt.Errorf("failed to check subscription status: %w", err)
	}
	// Also ensure EndDate is in the future
	return sub.IsActive && sub.EndDate.After(time.Now()), nil
}

// GetCustomerTraderSubscriptions retrieves all subscriptions a customer has made to traders.
func (s *traderSubscriptionService) GetCustomerTraderSubscriptions(ctx context.Context, customerID uint) ([]models.TraderSubscription, error) {
	var subs []models.TraderSubscription
	err := s.db.WithContext(ctx).Preload("Trader").Preload("SubscriptionPlan").
		Where("user_id = ?", customerID).Find(&subs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get customer trader subscriptions: %w", err)
	}
	return subs, nil
}

// GetTraderPlans retrieves all active subscription plans offered by a specific trader.
func (s *traderSubscriptionService) GetTraderPlans(ctx context.Context, traderID uint) ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	err := s.db.WithContext(ctx).Where("trader_id = ? AND is_trader_plan = ? AND is_active = ?", traderID, true, true).Find(&plans).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get trader's subscription plans: %w", err)
	}
	return plans, nil
}

// Add other service methods as needed for subscription management (e.g., extend, cancel, get details)
