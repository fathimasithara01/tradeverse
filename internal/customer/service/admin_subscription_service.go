package service // This file is in the 'service' package

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
	ErrNotTraderPlan                = errors.New("this is not a trader subscription plan")
	ErrAlreadyHasTraderSubscription = errors.New("user already has an active trader subscription")
	ErrNoActiveTraderSubscription   = errors.New("active trader subscription not found for this user and ID")
	ErrUserIsAlreadyTrader          = errors.New("user is already a trader")

	// Make sure ErrInsufficientFunds and ErrAdminWalletNotFound are defined if used.
	// You commented them out. If you need them, uncomment and ensure they are sourced correctly.
	// For now, let's use the ones from walletrepo package where applicable.
	ErrInsufficientFunds      = errors.New("insufficient funds in user's wallet") // Re-adding for this service's specific error.
	ErrAdminWalletNotFound    = errors.New("admin wallet not found")
	ErrCustomerWalletNotFound = errors.New("customer wallet not found, please create one or contact support")
)

type TraderSubscriptionPlanResponse struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Duration        int     `json:"duration"`
	Interval        string  `json:"interval"`
	Features        string  `json:"features"`
	MaxFollowers    int     `json:"max_followers,omitempty"`
	CommissionRate  float64 `json:"commission_rate,omitempty"`
	AnalyticsAccess string  `json:"analytics_access,omitempty"`
}

type UserTraderSubscriptionResponse struct {
	ID        uint      `json:"id"`
	PlanName  string    `json:"plan_name"`
	Price     float64   `json:"price"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	Status    string    `json:"payment_status"`
}

type AdminSubscriptionService interface {
	ListTraderSubscriptionPlans() ([]TraderSubscriptionPlanResponse, error)
	SubscribeToTraderPlan(customerID, planID uint) (*models.TraderSubscriptionResponse, error)
	GetCustomerTraderSubscription(userID uint) (*UserTraderSubscriptionResponse, error)
	// CancelCustomerTraderSubscription(userID uint, subscriptionID uint) error // Original commented out
	DeactivateExpiredTraderSubscriptions() error
	CancelCustomerTraderSubscription(ctx context.Context, userID, subscriptionID uint) error // Current signature
}

type adminSubscriptionService struct {
	repo       customerrepo.IAdminSubscriptionRepository
	walletSvc  IWalletService // <--- Use the interface directly, no package prefix needed as it's in the same 'service' package
	walletRepo walletrepo.WalletRepository
	db         *gorm.DB
}

// NewCustomerService is likely meant to be NewAdminSubscriptionService based on the return type
func NewAdminSubscriptionService(repo customerrepo.IAdminSubscriptionRepository, walletSvc IWalletService, walletRepo walletrepo.WalletRepository, db *gorm.DB) AdminSubscriptionService {
	return &adminSubscriptionService{
		repo:       repo,
		walletSvc:  walletSvc,
		walletRepo: walletRepo,
		db:         db,
	}
}

func (s *adminSubscriptionService) ListTraderSubscriptionPlans() ([]TraderSubscriptionPlanResponse, error) {
	plans, err := s.repo.GetTraderSubscriptionPlans()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trader subscription plans: %w", err)
	}

	var responses []TraderSubscriptionPlanResponse
	for _, plan := range plans {
		responses = append(responses, TraderSubscriptionPlanResponse{
			ID:              plan.ID,
			Name:            plan.Name,
			Description:     plan.Description,
			Price:           plan.Price,
			Duration:        int(plan.Duration),
			Interval:        plan.Interval,
			Features:        plan.Features,
			MaxFollowers:    plan.MaxFollowers,
			CommissionRate:  plan.CommissionRate,
			AnalyticsAccess: plan.AnalyticsAccess,
		})
	}
	return responses, nil
}

const AdminUserID uint = 1 // <--- Moved AdminUserID here so it's defined in this file

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
	if !plan.IsTraderPlan {
		return nil, ErrNotTraderPlan
	}

	existingSub, _ := s.repo.GetUserTraderSubscription(customerID) // Error is ignored, consider handling it for robustness.
	if existingSub != nil && existingSub.IsActive {
		return nil, ErrAlreadyHasTraderSubscription
	}

	// Use walletSvc to get wallet summary, which is a higher-level abstraction
	customerWalletSummary, err := s.walletSvc.GetUserWallet(context.Background(), customerID) // <--- Corrected method call and added context
	if err != nil {
		// Use specific wallet service errors or wrap them
		if errors.Is(err, ErrWalletServiceNotFound) { // Assuming ErrWalletServiceNotFound is defined in wallet_service.go
			return nil, fmt.Errorf("%w, please create one or contact support: %v", ErrCustomerWalletNotFound, err) // Re-using your defined error.
		}
		return nil, fmt.Errorf("failed to get customer wallet: %w", err)
	}

	if customerWalletSummary.Balance < plan.Price {
		return nil, ErrInsufficientFunds // <--- Using your service's ErrInsufficientFunds
	}

	adminWallet, err := s.walletRepo.GetUserWallet(AdminUserID) // Using the repository for admin wallet directly
	if err != nil {
		if errors.Is(err, walletrepo.ErrWalletNotFound) { // Check for walletrepo's ErrWalletNotFound
			// If admin wallet doesn't exist, create it (this logic can be in admin repo or here)
			newAdminWallet := &models.Wallet{
				UserID:   AdminUserID,
				Balance:  0,
				Currency: "INR", // Default currency
			}
			if createErr := s.walletRepo.UpdateWalletTx(s.db, newAdminWallet); createErr != nil { // Use UpdateWalletTx with current DB for creation
				return nil, fmt.Errorf("failed to auto-create admin wallet: %w", createErr)
			}
			adminWallet = newAdminWallet
		} else {
			return nil, fmt.Errorf("failed to get admin wallet: %w", err) // <--- Using your service's ErrAdminWalletNotFound
		}
	}

	paymentRef := fmt.Sprintf("TRADER_UPGRADE_%d_USER_%d_%s", planID, customerID, time.Now().Format("20060102150405"))

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Debit customer wallet using the wallet repository
		if err := s.walletRepo.DebitWallet(tx, customerWalletSummary.ID, plan.Price, models.TxTypeSubscription, paymentRef, "Trader upgrade"); err != nil {
			if errors.Is(err, walletrepo.ErrInsufficientFunds) {
				return ErrInsufficientFunds // Translate repo error to service error
			}
			return fmt.Errorf("failed to debit customer wallet: %w", err)
		}

		// Credit admin wallet using the wallet repository
		_, err := s.walletRepo.CreditWallet(tx, adminWallet.ID, plan.Price, models.TxTypeSubscription, paymentRef, "Trader upgrade")
		if err != nil {
			return fmt.Errorf("failed to credit admin wallet: %w", err)
		}

		now := time.Now()
		// Use calculateEndDate helper for consistency
		endDate := calculateEndDate(now, plan.Interval, int(plan.Duration))

		sub := models.TraderSubscriptionPlan{ // This should be models.TraderSubscription if you intend to upgrade to trader
			UserID: customerID,
			// SubscriptionPlanID: plan.ID,
			StartDate:     now,
			EndDate:       endDate,
			IsActive:      true,
			PaymentStatus: string(models.TxStatusSuccess),
			AmountPaid:    plan.Price,
			// TransactionID:      paymentRef,
			// AdminCommission:    plan.Price, // This value is implicitly handled by admin wallet credit
		}

		if err := tx.Create(&sub).Error; err != nil {
			return fmt.Errorf("failed to create subscription: %w", err)
		}

		// Upgrade user to Trader
		if err := tx.Model(&models.User{}).Where("id = ?", customerID).Update("role", models.RoleTrader).Error; err != nil {
			return fmt.Errorf("failed to upgrade user role: %w", err)
		}

		// Assuming you might need to create or update a TraderProfile as well, similar to your commented code
		var traderProfile models.TraderProfile
		if err := tx.Where("user_id = ?", customerID).First(&traderProfile).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				traderProfile = models.TraderProfile{
					UserID:     customerID,
					Status:     models.StatusPending, // Or models.StatusApproved if direct approval
					IsVerified: false,
				}
				if err := tx.Create(&traderProfile).Error; err != nil {
					return fmt.Errorf("failed to create trader profile in transaction: %w", err)
				}
			} else {
				return fmt.Errorf("failed to check existing trader profile: %w", err)
			}
		} else {
			traderProfile.Status = models.StatusPending // Update status if already exists
			if err := tx.Save(&traderProfile).Error; err != nil {
				return fmt.Errorf("failed to update trader profile in transaction: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &models.TraderSubscriptionResponse{
		PlanName:        plan.Name,
		AmountPaid:      plan.Price,
		AdminCommission: plan.Price, // This is technically the whole plan price going to admin, adjust if there's a commission rate.
		PaymentStatus:   string(models.TxStatusSuccess),
		TransactionID:   paymentRef,
		StartDate:       time.Now().Format(time.RFC3339),
		EndDate:         time.Now().Add(time.Duration(plan.Duration) * 24 * time.Hour).Format(time.RFC3339), // Re-calculate or use `endDate` variable.
		IsActive:        true,
		Message:         "Successfully upgraded to trader plan", // Add a message for response.
		Status:          string(models.TxStatusSuccess),         // Align with payment status.
	}, nil
}

// ... (rest of the functions remain the same)
func (s *adminSubscriptionService) GetCustomerTraderSubscription(userID uint) (*UserTraderSubscriptionResponse, error) {
	sub, err := s.repo.GetUserTraderSubscription(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user's trader subscription: %w", err)
	}
	if sub == nil {
		return nil, nil // No subscription found is not an error here, but nil response
	}
	// Ensure TraderSubscriptionPlan is preloaded by the repository, or fetch it here.
	// For now, assuming it's preloaded.
	if sub.TraderSubscriptionPlan == nil || sub.TraderSubscriptionPlan.ID == 0 { // Check for nil before accessing ID
		// Try to fetch plan if not preloaded
		plan, planErr := s.repo.GetSubscriptionPlanByID(sub.TraderSubscriptionPlanID)
		if planErr != nil {
			log.Printf("Warning: Failed to load TraderSubscriptionPlan for subscription ID %d: %v", sub.ID, planErr)
			return nil, fmt.Errorf("trader subscription plan not loaded for subscription ID %d: %w", sub.ID, planErr)
		}
		sub.TraderSubscriptionPlan = plan
	}

	return &UserTraderSubscriptionResponse{
		ID:        sub.ID,
		PlanName:  sub.TraderSubscriptionPlan.Name,
		Price:     sub.AmountPaid,
		StartDate: sub.StartDate,
		EndDate:   sub.EndDate,
		IsActive:  sub.IsActive,
		Status:    sub.PaymentStatus,
	}, nil
}

func (s *adminSubscriptionService) CancelCustomerTraderSubscription(ctx context.Context, userID, subscriptionID uint) error {
	existingSub, err := s.repo.GetUserTraderSubscription(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoActiveTraderSubscription
		}
		return fmt.Errorf("failed to check existing trader subscription: %w", err)
	}
	if existingSub == nil || existingSub.ID != subscriptionID || !existingSub.IsActive {
		return ErrNoActiveTraderSubscription
	}

	err = s.repo.CancelTraderSubscription(userID, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to cancel trader subscription: %w", err)
	}

	return nil
}

func calculateEndDate(start time.Time, interval string, duration int) time.Time {
	switch strings.ToLower(strings.TrimSpace(interval)) {
	case "day", "days", "d", "daily":
		return start.AddDate(0, 0, duration)
	case "week", "weeks", "w", "weekly":
		return start.AddDate(0, 0, duration*7)
	case "month", "months", "m", "monthly":
		return start.AddDate(0, duration, 0)
	case "year", "years", "y", "yearly":
		return start.AddDate(duration, 0, 0)
	default:
		fmt.Printf("Warning: Unknown interval '%s'. Defaulting to 1 month duration.\n", interval)
		return start.AddDate(0, 1, 0)
	}
}

func (s *adminSubscriptionService) DeactivateExpiredTraderSubscriptions() error {
	log.Println("Running cron job: Deactivating expired trader subscriptions...")
	expiredTraderSubs, err := s.repo.GetExpiredActiveTraderSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to get expired active trader subscriptions: %w", err)
	}

	if len(expiredTraderSubs) == 0 {
		log.Println("No expired trader subscriptions found to deactivate.")
		return nil
	}

	for _, sub := range expiredTraderSubs {
		sub.IsActive = false
		sub.PaymentStatus = "expired"
		if err := s.repo.UpdateTraderSubscription(&sub); err != nil {
			log.Printf("Error deactivating trader subscription ID %d for user %d: %v", sub.ID, sub.UserID, err)
		} else {
			planName := "Unknown Plan"
			// Check for nil before accessing plan fields
			if sub.TraderSubscriptionPlan != nil {
				planName = sub.TraderSubscriptionPlan.Name
			}
			log.Printf("Deactivated trader subscription ID %d for user %d (Plan: %s). EndDate: %v", sub.ID, sub.UserID, planName, sub.EndDate)

			var activeTraderSubsCount int64
			// Ensure you're querying the correct table (models.Subscription or models.TraderSubscription)
			s.db.Model(&models.TraderSubscriptionPlan{}).Where("user_id = ? AND is_active = ?", sub.UserID, true).Count(&activeTraderSubsCount)
			if activeTraderSubsCount == 0 {
				if err := s.db.Model(&models.User{}).Where("id = ?", sub.UserID).Update("role", models.RoleCustomer).Error; err != nil {
					log.Printf("Warning: Failed to demote user %d to customer after their last trader subscription expired: %v", sub.UserID, err)
				} else {
					log.Printf("User %d demoted to customer role after last trader subscription expired.", sub.UserID)
				}
			}
		}
	}
	log.Printf("Cron job finished: Deactivated %d trader subscriptions.", len(expiredTraderSubs))
	return nil
}
