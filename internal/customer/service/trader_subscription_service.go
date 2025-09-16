// internal/customer/service/customer_service.go
package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrPlanNotFound                 = errors.New("subscription plan not found")
	ErrNotTraderPlan                = errors.New("this is not a trader subscription plan")
	ErrAlreadyHasTraderSubscription = errors.New("user already has an active trader subscription")
	ErrNoActiveTraderSubscription   = errors.New("active trader subscription not found for this user and ID")
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
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) ListTraderSubscriptionPlans() ([]TraderSubscriptionPlanResponse, error) {
	plans, err := s.repo.GetTraderSubscriptionPlans() // Calls the correct repository method
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

func (s *customerService) SubscribeToTraderPlan(userID uint, planID uint) (*UserTraderSubscriptionResponse, error) {
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
	if err != nil {
		return nil, fmt.Errorf("failed to check existing trader subscription: %w", err)
	}
	if existingSub != nil {
		return nil, ErrAlreadyHasTraderSubscription
	}

	adminWallet, err := s.repo.GetAdminWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin wallet: %w", err)
	}

	paymentReferenceID := fmt.Sprintf("SUB_%d_USER_%d_%s", planID, userID, time.Now().Format("20060102150405"))
	paymentDescription := fmt.Sprintf("Payment for Trader Subscription Plan '%s' by User ID %d", plan.Name, userID)

	err = s.repo.CreditWallet(adminWallet.ID, plan.Price, models.TxTypeDeposit, paymentReferenceID, paymentDescription)
	if err != nil {
		return nil, fmt.Errorf("failed to credit admin wallet: %w", err)
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

	if err := s.repo.CreateTraderSubscription(&newSubscription); err != nil {
		return nil, fmt.Errorf("failed to create trader subscription record: %w", err)
	}

	if err := s.repo.UpdateUserRole(userID, models.RoleTrader); err != nil {
		return nil, fmt.Errorf("failed to update user role to trader: %w", err)
	}

	return &UserTraderSubscriptionResponse{
		ID:        newSubscription.ID,
		PlanName:  plan.Name,
		Price:     newSubscription.AmountPaid,
		StartDate: newSubscription.StartDate,
		EndDate:   newSubscription.EndDate,
		IsActive:  newSubscription.IsActive,
		Status:    newSubscription.PaymentStatus,
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
		return start.AddDate(0, 0, duration)
	}
}
