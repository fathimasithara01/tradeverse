package service

import (
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
	ErrInsufficientFunds            = errors.New("insufficient funds in user's wallet")
	ErrUserIsAlreadyTrader          = errors.New("user is already a trader")
	ErrAdminWalletNotFound          = errors.New("admin wallet not found")
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
	SubscribeToTraderPlan(userID uint, planID uint) (*UserTraderSubscriptionResponse, error)
	GetCustomerTraderSubscription(userID uint) (*UserTraderSubscriptionResponse, error)
	CancelCustomerTraderSubscription(userID uint, subscriptionID uint) error
	DeactivateExpiredTraderSubscriptions() error
}

type adminSubscriptionService struct {
	repo       customerrepo.IAdminSubscriptionRepository
	walletSvc  WalletService
	walletRepo walletrepo.WalletRepository
	db         *gorm.DB
}

func NewCustomerService(repo customerrepo.IAdminSubscriptionRepository, walletSvc WalletService, walletRepo walletrepo.WalletRepository, db *gorm.DB) AdminSubscriptionService {
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
			Duration:        plan.Duration,
			Interval:        plan.Interval,
			Features:        plan.Features,
			MaxFollowers:    plan.MaxFollowers,
			CommissionRate:  plan.CommissionRate,
			AnalyticsAccess: plan.AnalyticsAccess,
		})
	}
	return responses, nil
}

const AdminUserID uint = 1

func (s *adminSubscriptionService) SubscribeToTraderPlan(userID uint, planID uint) (*UserTraderSubscriptionResponse, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user.IsTrader() {
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

	existingSub, err := s.repo.GetUserTraderSubscription(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing trader subscription: %w", err)
	}
	if existingSub != nil && existingSub.IsActive {
		return nil, ErrAlreadyHasTraderSubscription
	}

	customerWalletSummary, err := s.walletSvc.GetWalletSummary(userID)
	if err != nil {
		if errors.Is(err, ErrUserWalletNotFound) {
			return nil, fmt.Errorf("customer wallet not found, please create one or contact support: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve customer wallet summary: %w", err)
	}

	if customerWalletSummary.Balance < plan.Price {
		return nil, ErrInsufficientFunds
	}

	adminWallet, err := s.walletRepo.GetUserWallet(AdminUserID)
	if err != nil {
		if errors.Is(err, walletrepo.ErrWalletNotFound) {
			newAdminWallet := &models.Wallet{
				UserID:   AdminUserID,
				Balance:  0,
				Currency: "INR",
			}
			if createErr := s.walletRepo.CreateWallet(newAdminWallet); createErr != nil {
				return nil, fmt.Errorf("failed to auto-create admin wallet: %w", createErr)
			}
			adminWallet = newAdminWallet
		} else {
			return nil, fmt.Errorf("failed to retrieve admin wallet: %w", err)
		}
	}

	paymentReferenceID := fmt.Sprintf("TRADER_SUB_%d_USER_%d_%s", planID, userID, time.Now().Format("20060102150405"))
	paymentDescription := fmt.Sprintf("Payment for Trader Subscription Plan '%s' by User ID %d", plan.Name, userID)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		err = s.walletRepo.DebitWallet(tx, customerWalletSummary.WalletID, plan.Price, models.TxTypeSubscription, paymentReferenceID, "Trader subscription payment (debit)")
		if err != nil {
			if errors.Is(err, walletrepo.ErrInsufficientFunds) {
				return ErrInsufficientFunds
			}
			return fmt.Errorf("failed to debit customer wallet in transaction: %w", err)
		}

		err = s.walletRepo.CreditWallet(tx, adminWallet.ID, plan.Price, models.TxTypeSubscription, paymentReferenceID, paymentDescription)
		if err != nil {
			return fmt.Errorf("failed to credit admin wallet in transaction: %w", err)
		}

		now := time.Now()
		endDate := calculateEndDate(now, plan.Interval, plan.Duration)

		newSubscription := models.Subscription{
			UserID:             userID,
			SubscriptionPlanID: plan.ID,
			StartDate:          now,
			EndDate:            endDate,
			IsActive:           true,
			PaymentStatus:      string(models.TxStatusSuccess),
			AmountPaid:         plan.Price,
			TransactionID:      paymentReferenceID,
		}

		if err := tx.Create(&newSubscription).Error; err != nil {
			return fmt.Errorf("failed to create trader subscription record in transaction: %w", err)
		}

		user.Role = models.RoleTrader
		if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("role", models.RoleTrader).Error; err != nil {
			return fmt.Errorf("failed to update user role to trader in transaction: %w", err)
		}

		var traderProfile models.TraderProfile
		if err := tx.Where("user_id = ?", userID).First(&traderProfile).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				traderProfile = models.TraderProfile{
					UserID:     userID,
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
			traderProfile.Status = models.StatusPending
			if err := tx.Save(&traderProfile).Error; err != nil {
				return fmt.Errorf("failed to update trader profile in transaction: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, ErrInsufficientFunds) {
			return nil, ErrInsufficientFunds
		}
		return nil, err
	}

	finalSub, err := s.repo.GetUserTraderSubscription(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch newly created trader subscription: %w", err)
	}
	if finalSub == nil {
		return nil, fmt.Errorf("newly created trader subscription not found after successful transaction")
	}

	return &UserTraderSubscriptionResponse{
		ID:        finalSub.ID,
		PlanName:  plan.Name,
		Price:     finalSub.AmountPaid,
		StartDate: finalSub.StartDate,
		EndDate:   finalSub.EndDate,
		IsActive:  finalSub.IsActive,
		Status:    finalSub.PaymentStatus,
	}, nil
}

func (s *adminSubscriptionService) GetCustomerTraderSubscription(userID uint) (*UserTraderSubscriptionResponse, error) {
	sub, err := s.repo.GetUserTraderSubscription(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user's trader subscription: %w", err)
	}
	if sub == nil {
		return nil, nil
	}
	if sub.TraderSubscriptionPlan.ID == 0 {
		return nil, fmt.Errorf("trader subscription plan not loaded for subscription ID %d", sub.ID)
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

func (s *adminSubscriptionService) CancelCustomerTraderSubscription(userID uint, subscriptionID uint) error {
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
			if sub.TraderSubscriptionPlan != nil {
				planName = sub.TraderSubscriptionPlan.Name
			}
			log.Printf("Deactivated trader subscription ID %d for user %d (Plan: %s). EndDate: %v", sub.ID, sub.UserID, planName, sub.EndDate)

			var activeTraderSubsCount int64
			s.db.Model(&models.TraderSubscription{}).Where("user_id = ? AND is_active = ?", sub.UserID, true).Count(&activeTraderSubsCount)
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
