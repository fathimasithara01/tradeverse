package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	customerrepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo" // Changed import path
	walletrepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"     // Changed import path
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrPlanNotFound                 = errors.New("subscription plan not found")
	ErrNotTraderPlan                = errors.New("this is not a trader subscription plan")
	ErrAlreadyHasTraderSubscription = errors.New("user already has an active trader subscription")
	ErrNoActiveTraderSubscription   = errors.New("active trader subscription not found for this user and ID")
	ErrInsufficientFunds            = errors.New("insufficient funds in user's wallet") // This is distinct from repository.ErrInsufficientFunds
	ErrUserIsAlreadyTrader          = errors.New("user is already a trader")
	ErrAdminWalletNotFound          = errors.New("admin wallet not found") // This is a specific scenario, will handle in code
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

type CustomerService interface {
	ListTraderSubscriptionPlans() ([]TraderSubscriptionPlanResponse, error)
	SubscribeToTraderPlan(userID uint, planID uint) (*UserTraderSubscriptionResponse, error)
	GetCustomerTraderSubscription(userID uint) (*UserTraderSubscriptionResponse, error)
	CancelCustomerTraderSubscription(userID uint, subscriptionID uint) error
}

type customerService struct {
	repo       customerrepo.CustomerRepository // Changed type
	walletSvc  WalletService
	walletRepo walletrepo.WalletRepository // Changed type
	db         *gorm.DB
}

func NewCustomerService(repo customerrepo.CustomerRepository, walletSvc WalletService, walletRepo walletrepo.WalletRepository, db *gorm.DB) CustomerService { // Changed type
	return &customerService{
		repo:       repo,
		walletSvc:  walletSvc,
		walletRepo: walletRepo,
		db:         db,
	}
}

func (s *customerService) ListTraderSubscriptionPlans() ([]TraderSubscriptionPlanResponse, error) {
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

const AdminUserID uint = 1 // Assuming admin user ID is 1 (Can be configured)

func (s *customerService) SubscribeToTraderPlan(userID uint, planID uint) (*UserTraderSubscriptionResponse, error) {
	// 1. Fetch User and Plan
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

	// 2. Check for existing active trader subscription
	existingSub, err := s.repo.GetUserTraderSubscription(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing trader subscription: %w", err)
	}
	if existingSub != nil {
		return nil, ErrAlreadyHasTraderSubscription
	}

	// 3. Get Wallets (Customer's and Admin's) - via walletSvc and walletRepo
	customerWalletSummary, err := s.walletSvc.GetWalletSummary(userID)
	if err != nil {
		if errors.Is(err, ErrUserWalletNotFound) { // ErrUserWalletNotFound is in THIS service package (wallet_service.go)
			return nil, fmt.Errorf("customer wallet not found, please create one or contact support: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve customer wallet summary: %w", err)
	}

	if customerWalletSummary.Balance < plan.Price {
		return nil, ErrInsufficientFunds
	}

	// Get admin's wallet (directly from walletRepo as it's an internal transfer, not external payment)
	adminWallet, err := s.walletRepo.GetUserWallet(AdminUserID)
	if err != nil {
		if errors.Is(err, walletrepo.ErrWalletNotFound) { // Correct: walletrepo.ErrWalletNotFound
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
			if errors.Is(err, walletrepo.ErrInsufficientFunds) { // Correct: walletrepo.ErrInsufficientFunds
				return ErrInsufficientFunds // Return your service's specific ErrInsufficientFunds
			}
			return fmt.Errorf("failed to debit customer wallet in transaction: %w", err)
		}

		err = s.walletRepo.CreditWallet(tx, adminWallet.ID, plan.Price, models.TxTypeSubscription, paymentReferenceID, paymentDescription)
		if err != nil {
			return fmt.Errorf("failed to credit admin wallet in transaction: %w", err)
		}

		now := time.Now()
		endDate := calculateEndDate(now, plan.Interval, plan.Duration)

		newSubscription := models.TraderSubscription{
			UserID:                   userID,
			TraderSubscriptionPlanID: plan.ID,
			StartDate:                now,
			EndDate:                  endDate,
			IsActive:                 true,
			PaymentStatus:            string(models.TxStatusSuccess),
			AmountPaid:               plan.Price,
			TransactionID:            paymentReferenceID,
		}

		if err := tx.Create(&newSubscription).Error; err != nil {
			return fmt.Errorf("failed to create trader subscription record in transaction: %w", err)
		}

		if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("role", models.RoleTrader).Error; err != nil {
			return fmt.Errorf("failed to update user role to trader in transaction: %w", err)
		}

		newTraderProfile := models.TraderProfile{
			UserID:     userID,
			Status:     models.StatusPending,
			IsVerified: false,
		}
		if err := tx.Create(&newTraderProfile).Error; err != nil {
			return fmt.Errorf("failed to create trader profile in transaction: %w", err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, ErrInsufficientFunds) { // This error is local to customer service
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

func (s *customerService) GetCustomerTraderSubscription(userID uint) (*UserTraderSubscriptionResponse, error) {
	sub, err := s.repo.GetUserTraderSubscription(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user's trader subscription: %w", err)
	}
	if sub == nil {
		return nil, nil
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

func (s *customerService) CancelCustomerTraderSubscription(userID uint, subscriptionID uint) error {
	existingSub, err := s.repo.GetUserTraderSubscription(userID)
	if err != nil {
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
	case "day", "days", "d":
		return start.AddDate(0, 0, duration)
	case "month", "months", "m", "monthly":
		return start.AddDate(0, duration, 0)
	case "year", "years", "y", "yearly":
		return start.AddDate(duration, 0, 0)
	default:
		return start.AddDate(0, duration, 0)
	}
}
