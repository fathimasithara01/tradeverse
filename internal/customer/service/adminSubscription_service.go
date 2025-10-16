package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	customerrepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	walletrepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrPlanNotFound                 = errors.New("subscription plan not found")
	ErrNotUpgradeToTraderPlan       = errors.New("this is not an upgrade to trader subscription plan") // Renamed for clarity
	ErrAlreadyHasTraderSubscription = errors.New("user already has an active trader subscription")
	ErrNoActiveTraderSubscription   = errors.New("active trader subscription not found for this user and ID")
	ErrUserIsAlreadyTrader          = errors.New("user is already a trader")

	ErrInsufficientFunds      = errors.New("insufficient funds in user's wallet")
	ErrAdminWalletNotFound    = errors.New("admin wallet not found")
	ErrCustomerWalletNotFound = errors.New("customer wallet not found, please create one or contact support")
	// ErrWalletServiceNotFound should be defined in wallet_service.go, assuming it is.
	ErrNotTraderPlan = errors.New("error")
)

type TraderSubscriptionPlanResponse struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Price             float64 `json:"price"`
	Duration          int     `json:"duration"` // Now duration in days for response struct
	Interval          string  `json:"interval"`
	Features          string  `json:"features"`
	MaxFollowers      int     `json:"max_followers,omitempty"`
	CommissionRate    float64 `json:"commission_rate,omitempty"`
	AnalyticsAccess   string  `json:"analytics_access,omitempty"`
	IsUpgradeToTrader bool    `json:"is_upgrade_to_trader"`
}

type UserTraderSubscriptionResponse struct {
	ID             uint      `json:"id"`
	PlanName       string    `json:"plan_name"`
	Price          float64   `json:"price"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	IsActive       bool      `json:"is_active"`
	PaymentStatus  string    `json:"payment_status"` // Renamed from 'Status' for clarity
	TransactionID  string    `json:"transaction_id"`
	CommissionRate float64   `json:"commission_rate,omitempty"`
}

type AdminSubscriptionService interface {
	ListTraderSubscriptionPlans() ([]TraderSubscriptionPlanResponse, error)
	SubscribeToTraderPlan(customerID, planID uint) (*models.TraderSubscriptionResponse, error)
	GetCustomerTraderSubscription(userID uint) (*UserTraderSubscriptionResponse, error)
	DeactivateExpiredTraderSubscriptions() error
	CancelCustomerTraderSubscription(ctx context.Context, userID, subscriptionID uint) error
}

type adminSubscriptionService struct {
	repo       customerrepo.IAdminSubscriptionRepository
	walletSvc  IWalletService
	walletRepo walletrepo.WalletRepository
	db         *gorm.DB
}

func NewAdminSubscriptionService(repo customerrepo.IAdminSubscriptionRepository, walletSvc IWalletService, walletRepo walletrepo.WalletRepository, db *gorm.DB) AdminSubscriptionService {
	return &adminSubscriptionService{
		repo:       repo,
		walletSvc:  walletSvc,
		walletRepo: walletRepo,
		db:         db,
	}
}

func (s *adminSubscriptionService) ListTraderSubscriptionPlans() ([]TraderSubscriptionPlanResponse, error) {

	plans, err := s.repo.GetTraderSubscriptionPlans() // Assuming this method now returns SubscriptionPlan models
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trader subscription plans: %w", err)
	}

	var responses []TraderSubscriptionPlanResponse
	for _, plan := range plans {
		if !plan.IsUpgradeToTrader { // Filter only plans meant for upgrading to a trader
			continue
		}
		responses = append(responses, TraderSubscriptionPlanResponse{
			ID:                plan.ID,
			Name:              plan.Name,
			Description:       plan.Description,
			Price:             plan.Price,
			Duration:          int(plan.Duration / (24 * time.Hour)), // Convert duration to days for response
			Interval:          plan.Interval,
			Features:          plan.Features,
			MaxFollowers:      plan.MaxFollowers,
			CommissionRate:    plan.CommissionRate,
			AnalyticsAccess:   plan.AnalyticsAccess,
			IsUpgradeToTrader: plan.IsUpgradeToTrader,
		})
	}
	return responses, nil
}

const AdminUserID uint = 1

func (s *adminSubscriptionService) SubscribeToTraderPlan(customerID, planID uint) (*models.TraderSubscriptionResponse, error) {
	customer, err := s.repo.GetUserByID(customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	if customer.IsTrader() {
		return nil, ErrUserIsAlreadyTrader
	}

	plan, err := s.repo.GetSubscriptionPlanByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlanNotFound
		}
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}
	if !plan.IsUpgradeToTrader {
		return nil, ErrNotUpgradeToTraderPlan // Use the renamed error
	}

	// Check for active trader subscription for this user
	activeTraderSubs, err := s.repo.GetUserActiveTraderSubscriptions(customerID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // Only return error if it's not simply "not found"
		return nil, fmt.Errorf("failed to check for existing trader subscriptions: %w", err)
	}
	if len(activeTraderSubs) > 0 {
		return nil, ErrAlreadyHasTraderSubscription
	}

	customerWalletSummary, err := s.walletSvc.GetUserWallet(context.Background(), customerID)
	if err != nil {
		if errors.Is(err, ErrWalletServiceNotFound) {
			return nil, fmt.Errorf("%w, please create one or contact support: %v", ErrCustomerWalletNotFound, err)
		}
		return nil, fmt.Errorf("failed to get customer wallet: %w", err)
	}

	if customerWalletSummary.Balance < plan.Price {
		return nil, ErrInsufficientFunds
	}

	// Ensure Admin wallet exists. If not, create it.
	adminWallet, err := s.walletRepo.GetUserWallet(AdminUserID)
	if err != nil {
		if errors.Is(err, walletrepo.ErrWalletNotFound) {
			newAdminWallet := &models.Wallet{
				UserID:   AdminUserID,
				Balance:  0,
				Currency: "INR", // Default currency
			}
			if createErr := s.walletRepo.UpdateWalletTx(s.db, newAdminWallet); createErr != nil {
				return nil, fmt.Errorf("failed to auto-create admin wallet: %w", createErr)
			}
			adminWallet = newAdminWallet
		} else {
			return nil, fmt.Errorf("failed to get admin wallet: %w", err)
		}
	}

	paymentRef := fmt.Sprintf("TRADER_UPGRADE_PLAN_%d_USER_%d_%s", planID, customerID, time.Now().Format("20060102150405"))

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Debit customer wallet
		if err := s.walletRepo.DebitWallet(tx, customerWalletSummary.ID, plan.Price, models.TxTypeSubscription, paymentRef, "Trader upgrade subscription"); err != nil {
			if errors.Is(err, walletrepo.ErrInsufficientFunds) {
				return ErrInsufficientFunds
			}
			return fmt.Errorf("failed to debit customer wallet: %w", err)
		}

		// Credit admin wallet (since the full price goes to admin for trader upgrade plans)
		_, err := s.walletRepo.CreditWallet(tx, adminWallet.ID, plan.Price, models.TxTypeSubscription, paymentRef, "Trader upgrade revenue")
		if err != nil {
			return fmt.Errorf("failed to credit admin wallet: %w", err)
		}

		now := time.Now()
		// Convert plan.Duration (which is time.Duration) to days/months/years based on plan.Interval
		endDate := calculateEndDate(now, plan.Interval, int(plan.Duration/time.Hour)) // Duration in hours, calculateEndDate will handle interval logic

		// Create a new models.Subscription record
		subscription := models.CustomerToTraderSub{
			UserID:             customerID,
			SubscriptionPlanID: plan.ID,
			TraderID:           &customerID, // Set TraderID to customerID as they are becoming a trader
			StartDate:          now,
			EndDate:            endDate,
			IsActive:           true,
			PaymentStatus:      string(models.TxStatusSuccess),
			AmountPaid:         plan.Price,
			TransactionID:      paymentRef,
		}

		if err := tx.Create(&subscription).Error; err != nil {
			return fmt.Errorf("failed to create subscription: %w", err)
		}

		// Upgrade user to Trader role
		if err := tx.Model(&models.User{}).Where("id = ?", customerID).Update("role", models.RoleTrader).Error; err != nil {
			return fmt.Errorf("failed to upgrade user role to trader: %w", err)
		}

		// Create or update TraderProfile
		var traderProfile models.TraderProfile
		if err := tx.Where("user_id = ?", customerID).First(&traderProfile).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				traderProfile = models.TraderProfile{
					UserID:     customerID,
					Status:     models.StatusApproved, // Automatically approved if paying for a plan
					IsVerified: true,                  // Or based on other criteria
				}
				if err := tx.Create(&traderProfile).Error; err != nil {
					return fmt.Errorf("failed to create trader profile: %w", err)
				}
			} else {
				return fmt.Errorf("failed to check existing trader profile: %w", err)
			}
		} else {
			traderProfile.Status = models.StatusApproved // Update status if already exists
			traderProfile.IsVerified = true
			if err := tx.Save(&traderProfile).Error; err != nil {
				return fmt.Errorf("failed to update trader profile: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// For the response, we might still want a models.TraderSubscriptionResponse format
	return &models.TraderSubscriptionResponse{
		PlanName:        plan.Name,
		AmountPaid:      plan.Price,
		AdminCommission: plan.Price, // Full price as commission for upgrade plans
		PaymentStatus:   string(models.TxStatusSuccess),
		TransactionID:   paymentRef,
		StartDate:       time.Now().Format(time.RFC3339),
		EndDate:         calculateEndDate(time.Now(), plan.Interval, int(plan.Duration/time.Hour)).Format(time.RFC3339),
		IsActive:        true,
		Message:         "Successfully upgraded to trader plan",
		Status:          string(models.TxStatusSuccess),
	}, nil
}

func (s *adminSubscriptionService) GetCustomerTraderSubscription(userID uint) (*UserTraderSubscriptionResponse, error) {
	// Fetch active subscriptions where TraderID is set to the user and associated plan is an upgrade plan
	sub, err := s.repo.GetUserActiveTraderSubscription(userID) // This method needs to be implemented in repository to fetch models.Subscription
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active trader subscription found for this user
		}
		return nil, fmt.Errorf("failed to fetch user's trader subscription: %w", err)
	}
	if sub == nil {
		return nil, nil
	}

	// Ensure SubscriptionPlan is preloaded by the repository, or fetch it here.
	if sub.SubscriptionPlan.ID == 0 { // Check if plan is loaded
		plan, planErr := s.repo.GetSubscriptionPlanByID(sub.SubscriptionPlanID)
		if planErr != nil {
			log.Printf("Warning: Failed to load SubscriptionPlan for subscription ID %d: %v", sub.ID, planErr)
			return nil, fmt.Errorf("subscription plan not loaded for subscription ID %d: %w", sub.ID, planErr)
		}
		sub.SubscriptionPlan = *plan // <--- CORRECTED LINE: Dereference the pointer here
	}

	// Double-check if the plan is actually an upgrade to trader plan
	if !sub.SubscriptionPlan.IsUpgradeToTrader {
		return nil, fmt.Errorf("the active subscription found (ID: %d) is not an upgrade to trader plan", sub.ID)
	}

	return &UserTraderSubscriptionResponse{
		ID:             sub.ID,
		PlanName:       sub.SubscriptionPlan.Name,
		Price:          sub.AmountPaid,
		StartDate:      sub.StartDate,
		EndDate:        sub.EndDate,
		IsActive:       sub.IsActive,
		PaymentStatus:  sub.PaymentStatus,
		TransactionID:  sub.TransactionID,
		CommissionRate: sub.SubscriptionPlan.CommissionRate, // Assuming you want to display this
	}, nil
}

func (s *adminSubscriptionService) CancelCustomerTraderSubscription(ctx context.Context, userID, subscriptionID uint) error {
	// This should now cancel a models.Subscription record
	existingSub, err := s.repo.GetSubscriptionByIDAndUserID(subscriptionID, userID) // Need a repo method to get a specific subscription
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoActiveTraderSubscription
		}
		return fmt.Errorf("failed to find subscription: %w", err)
	}
	if existingSub == nil || !existingSub.IsActive || existingSub.TraderID == nil || *existingSub.TraderID != userID {
		return ErrNoActiveTraderSubscription
	}

	// Perform the cancellation logic (set IsActive to false, update status)
	err = s.db.Transaction(func(tx *gorm.DB) error {
		existingSub.IsActive = false
		existingSub.PaymentStatus = string(models.TxStatusCancelled)
		if err := tx.Save(existingSub).Error; err != nil {
			return fmt.Errorf("failed to update subscription status to cancelled: %w", err)
		}

		// Check if the user has any other active 'IsUpgradeToTrader' subscriptions
		var count int64
		err := tx.Model(&models.CustomerToTraderSub{}).
			Joins("JOIN subscription_plans ON subscriptions.subscription_plan_id = subscription_plans.id").
			Where("subscriptions.user_id = ? AND subscriptions.is_active = ? AND subscription_plans.is_upgrade_to_trader = ?", userID, true, true).
			Count(&count).Error
		if err != nil {
			return fmt.Errorf("failed to count other active trader subscriptions: %w", err)
		}

		if count == 0 {
			// If no other active trader subscriptions, demote the user
			if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("role", models.RoleCustomer).Error; err != nil {
				return fmt.Errorf("failed to demote user role to customer: %w", err)
			}
			log.Printf("User %d demoted to customer role after canceling last trader subscription.", userID)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to cancel trader subscription in transaction: %w", err)
	}
	log.Printf("Trader subscription ID %d for user %d cancelled successfully.", subscriptionID, userID)
	return nil
}

// calculateEndDate now takes duration in hours (as plan.Duration is time.Duration)
func calculateEndDate(start time.Time, interval string, durationHours int) time.Time {
	// Convert durationHours back to time.Duration for Add methods
	totalDuration := time.Duration(durationHours) * time.Hour

	switch strings.ToLower(strings.TrimSpace(interval)) {
	case "day", "days", "d", "daily":
		return start.AddDate(0, 0, int(totalDuration.Hours()/24))
	case "week", "weeks", "w", "weekly":
		return start.AddDate(0, 0, int(totalDuration.Hours()/(24*7)))
	case "month", "months", "m", "monthly":
		return start.AddDate(0, int(totalDuration.Hours()/(24*30)), 0) // Approximation for month
	case "year", "years", "y", "yearly":
		return start.AddDate(int(totalDuration.Hours()/(24*365)), 0, 0) // Approximation for year
	default:
		fmt.Printf("Warning: Unknown interval '%s'. Defaulting to 1 month duration.\n", interval)
		return start.AddDate(0, 1, 0) // Default to 1 month
	}
}

func (s *adminSubscriptionService) DeactivateExpiredTraderSubscriptions() error {
	log.Println("Running cron job: Deactivating expired trader subscriptions...")
	// Fetch all active subscriptions that are 'IsUpgradeToTrader' and whose EndDate is in the past
	expiredSubs, err := s.repo.GetExpiredActiveUpgradeToTraderSubscriptions() // This needs to be implemented in repository
	if err != nil {
		return fmt.Errorf("failed to get expired active upgrade to trader subscriptions: %w", err)
	}

	if len(expiredSubs) == 0 {
		log.Println("No expired upgrade to trader subscriptions found to deactivate.")
		return nil
	}

	for _, sub := range expiredSubs {
		err := s.db.Transaction(func(tx *gorm.DB) error {
			sub.IsActive = false
			sub.PaymentStatus = "expired"
			if err := tx.Save(&sub).Error; err != nil {
				return fmt.Errorf("failed to update subscription ID %d: %w", sub.ID, err)
			}

			// Check if this user has any other *active* 'IsUpgradeToTrader' subscriptions
			var activeTraderSubsCount int64
			err := tx.Model(&models.CustomerToTraderSub{}).
				Joins("JOIN subscription_plans ON subscriptions.subscription_plan_id = subscription_plans.id").
				Where("subscriptions.user_id = ? AND subscriptions.is_active = ? AND subscription_plans.is_upgrade_to_trader = ?", sub.UserID, true, true).
				Count(&activeTraderSubsCount).Error
			if err != nil {
				return fmt.Errorf("failed to count active trader subscriptions for user %d: %w", sub.UserID, err)
			}

			if activeTraderSubsCount == 0 {
				// If no other active 'upgrade to trader' subscriptions, demote the user
				if err := tx.Model(&models.User{}).Where("id = ?", sub.UserID).Update("role", models.RoleCustomer).Error; err != nil {
					return fmt.Errorf("failed to demote user %d to customer: %w", sub.UserID, err)
				}
				log.Printf("User %d demoted to customer role after their last trader upgrade subscription expired.", sub.UserID)
			}
			return nil
		})

		if err != nil {
			log.Printf("Error processing expired trader subscription ID %d for user %d: %v", sub.ID, sub.UserID, err)
		} else {
			planName := "Unknown Plan"
			if sub.SubscriptionPlan.ID != 0 { // Check if plan was loaded
				planName = sub.SubscriptionPlan.Name
			}
			log.Printf("Deactivated upgrade to trader subscription ID %d for user %d (Plan: %s). EndDate: %v", sub.ID, sub.UserID, planName, sub.EndDate)
		}
	}
	log.Printf("Cron job finished: Deactivated %d upgrade to trader subscriptions.", len(expiredSubs))
	return nil
}
