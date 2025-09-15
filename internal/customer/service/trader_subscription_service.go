package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

// DTOs for request/response
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
	// plans, err := s.repo.GetTraderSubscriptionPlans()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to fetch trader subscription plans: %w", err)
	// }

	// var responses []TraderSubscriptionPlanResponse
	// for _, plan := range plans {
	// 	responses = append(responses, TraderSubscriptionPlanResponse{
	// 		ID:              plan.ID,
	// 		Name:            plan.Name,
	// 		Description:     plan.Description,
	// 		Price:           plan.Price,
	// 		Duration:        plan.Duration,
	// 		Interval:        plan.Interval,
	// 		Features:        plan.Features,
	// 		MaxFollowers:    plan.MaxFollowers,
	// 		CommissionRate:  plan.CommissionRate,
	// 		AnalyticsAccess: plan.AnalyticsAccess,
	// 	})
	// }
	return s.repo.GetTraderSubscriptionPlans()

	// return responses, nil
}

func (s *customerService) SubscribeToTraderPlan(userID uint, planID uint) (*UserTraderSubscriptionResponse, error) {
	// 1. Get the plan details
	plan, err := s.repo.GetSubscriptionPlanByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription plan not found")
		}
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}

	if !plan.IsTraderPlan {
		return nil, errors.New("this is not a trader subscription plan")
	}

	// 2. Check if user already has an active trader subscription
	existingSub, err := s.repo.GetUserTraderSubscription(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing trader subscription: %w", err)
	}
	if existingSub != nil {
		return nil, errors.New("user already has an active trader subscription")
	}

	// 3. Simulate payment (in a real app, this integrates with a payment gateway)
	// For simplicity, we assume payment is successful and funds go to admin.
	adminWallet, err := s.repo.GetAdminWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin wallet: %w", err)
	}

	paymentReferenceID := fmt.Sprintf("SUB_%d_USER_%d_%s", planID, userID, time.Now().Format("20060102150405"))
	paymentDescription := fmt.Sprintf("Payment for Trader Subscription Plan '%s' by User ID %d", plan.Name, userID)

	// Credit admin's wallet
	err = s.repo.CreditWallet(adminWallet.ID, plan.Price, models.TxTypeDeposit, paymentReferenceID, paymentDescription)
	if err != nil {
		return nil, fmt.Errorf("failed to credit admin wallet: %w", err)
	}

	// 4. Create TraderSubscription record
	now := time.Now()
	endDate := now.AddDate(0, 0, plan.Duration) // Assuming Duration is in days. Adjust if 'Interval' is used.
	if plan.Interval == "monthly" {
		endDate = now.AddDate(0, plan.Duration, 0)
	} else if plan.Interval == "yearly" {
		endDate = now.AddDate(plan.Duration, 0, 0)
	}

	newSubscription := models.TraderSubscription{
		UserID:                   userID,
		TraderSubscriptionPlanID: plan.ID,
		StartDate:                now,
		EndDate:                  endDate,
		IsActive:                 true,
		// PaymentStatus:            models.TxStatusSuccess, // Assuming success
		AmountPaid:    plan.Price,
		TransactionID: paymentReferenceID, // Store our internal reference
	}

	if err := s.repo.CreateTraderSubscription(&newSubscription); err != nil {
		// IMPORTANT: In a real scenario, if subscription creation fails *after* payment,
		// you need to either refund the user or flag it for manual review.
		return nil, fmt.Errorf("failed to create trader subscription record: %w", err)
	}

	// 5. Update user role to Trader
	if err := s.repo.UpdateUserRole(userID, models.RoleTrader); err != nil {
		// Similar to above, consider rollbacks or flags
		return nil, fmt.Errorf("failed to update user role to trader: %w", err)
	}

	return &UserTraderSubscriptionResponse{
		ID:        newSubscription.ID,
		PlanName:  plan.Name,
		Price:     newSubscription.AmountPaid,
		StartDate: newSubscription.StartDate,
		EndDate:   newSubscription.EndDate,
		IsActive:  newSubscription.IsActive,
		Status:    string(newSubscription.PaymentStatus),
	}, nil
}

func (s *customerService) GetCustomerTraderSubscription(userID uint) (*UserTraderSubscriptionResponse, error) {
	sub, err := s.repo.GetUserTraderSubscription(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user's trader subscription: %w", err)
	}
	if sub == nil {
		return nil, nil // No active subscription
	}

	return &UserTraderSubscriptionResponse{
		ID:        sub.ID,
		PlanName:  sub.TraderSubscriptionPlan.Name,
		Price:     sub.AmountPaid,
		StartDate: sub.StartDate,
		EndDate:   sub.EndDate,
		IsActive:  sub.IsActive,
		Status:    string(sub.PaymentStatus),
	}, nil
}

func (s *customerService) CancelCustomerTraderSubscription(userID uint, subscriptionID uint) error {
	// Optional: You might want to add logic here to check if a refund is due
	// This example simply marks it inactive.
	existingSub, err := s.repo.GetUserTraderSubscription(userID)
	if err != nil {
		return fmt.Errorf("failed to check existing trader subscription: %w", err)
	}
	if existingSub == nil || existingSub.ID != subscriptionID || !existingSub.IsActive {
		return errors.New("active trader subscription not found for this user and ID")
	}

	err = s.repo.CancelTraderSubscription(userID, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to cancel trader subscription: %w", err)
	}

	// Optional: Revert user role if no other active trader subscriptions (complex, depends on business logic)
	// For simplicity, we won't revert the role automatically here, as they might have other trader activities.
	// Reverting role might be an admin-only action or part of a more complex lifecycle.

	return nil
}
