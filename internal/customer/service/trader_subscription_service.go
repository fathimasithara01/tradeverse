package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

// ITraderSubscriptionService defines the interface for customer-facing trader subscription business logic.
type ITraderSubscriptionService interface {
	GetTraderSubscriptionPlanByID(planID uint) (*models.SubscriptionPlan, error)
	GetAllActiveTraderSubscriptionPlans() ([]models.SubscriptionPlan, error)
	GetTraderProfile(traderProfileID uint) (*models.TraderProfile, error)

	SubscribeToTrader(userID, planID, traderProfileID uint, amountPaid float64, transactionID string) (*models.TraderSubscription, error)
	GetMyTraderSubscription(subscriptionID, userID uint) (*models.TraderSubscription, error)
	ListMyTraderSubscriptions(userID uint) ([]models.TraderSubscription, error)
	UpdateTraderSubscriptionSettings(subscriptionID, userID uint, allocation, riskMultiplier float64) (*models.TraderSubscription, error)
	PauseTraderCopyTrading(subscriptionID, userID uint) error
	ResumeTraderCopyTrading(subscriptionID, userID uint) error
	CancelTraderSubscription(subscriptionID, userID uint) error
	SimulateTraderSubscription(planID uint, initialCapital float64) (interface{}, error) // Placeholder for simulation
}

type TraderSubscriptionService struct {
	repo repository.ITraderSubscriptionRepository
}

func NewTraderSubscriptionService(repo repository.ITraderSubscriptionRepository) *TraderSubscriptionService {
	return &TraderSubscriptionService{repo: repo}
}

func (s *TraderSubscriptionService) GetTraderSubscriptionPlanByID(planID uint) (*models.SubscriptionPlan, error) {
	plan, err := s.repo.GetTraderSubscriptionPlanByID(planID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trader subscription plan: %w", err)
	}
	if plan == nil {
		return nil, errors.New("trader subscription plan not found or not active")
	}
	return plan, nil
}

// GetAllActiveTraderSubscriptionPlans retrieves all active trader subscription plans.
func (s *TraderSubscriptionService) GetAllActiveTraderSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	plans, err := s.repo.GetActiveTraderSubscriptionPlans()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch active trader subscription plans: %w", err)
	}
	return plans, nil
}

// GetTraderProfile retrieves a trader's profile.
func (s *TraderSubscriptionService) GetTraderProfile(traderProfileID uint) (*models.TraderProfile, error) {
	profile, err := s.repo.GetTraderProfileByID(traderProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trader profile: %w", err)
	}
	if profile == nil {
		return nil, errors.New("trader profile not found")
	}
	// Optionally, ensure the trader is approved before allowing subscription
	// if !profile.IsApproved {
	// 	return nil, errors.New("trader not yet approved")
	// }
	return profile, nil
}

// SubscribeToTrader creates a new trader subscription for a customer.
func (s *TraderSubscriptionService) SubscribeToTrader(userID, planID, traderProfileID uint, amountPaid float64, transactionID string) (*models.TraderSubscription, error) {
	plan, err := s.GetTraderSubscriptionPlanByID(planID)
	if err != nil {
		return nil, err
	}
	if plan.Price > amountPaid {
		return nil, errors.New("amount paid is less than plan price")
	}

	// Validate TraderProfile exists and is associated with a trader
	traderProfile, err := s.GetTraderProfile(traderProfileID)
	if err != nil {
		return nil, err
	}
	if traderProfile == nil {
		return nil, errors.New("trader profile not found")
	}
	// TODO: Add logic to check if UserID is a customer and TraderProfileID is a trader
	// This would typically involve user service calls.

	// Calculate subscription end date
	startDate := time.Now()
	var endDate time.Time
	switch plan.Interval {
	case "days":
		endDate = startDate.AddDate(0, 0, plan.Duration)
	case "monthly":
		endDate = startDate.AddDate(0, plan.Duration, 0)
	case "yearly":
		endDate = startDate.AddDate(plan.Duration, 0, 0)
	default:
		return nil, errors.New("invalid subscription plan interval")
	}

	subscription := &models.TraderSubscription{
		UserID:                   userID,
		TraderSubscriptionPlanID: planID,
		TraderProfileID:          &traderProfileID, // Link to the specific trader
		StartDate:                startDate,
		EndDate:                  endDate,
		IsActive:                 true,
		PaymentStatus:            "paid", // Assuming successful payment
		AmountPaid:               amountPaid,
		TransactionID:            transactionID,
		Allocation:               1.0, // Default allocation
		RiskMultiplier:           1.0, // Default risk multiplier
	}

	if err := s.repo.CreateTraderSubscription(subscription); err != nil {
		return nil, fmt.Errorf("failed to create trader subscription: %w", err)
	}

	// TODO: Integrate with a wallet service here to deposit `amountPaid` into the admin's wallet.
	// Example: s.walletService.DepositToAdmin(amountPaid, transactionID, userID, "TraderSubscription")

	return subscription, nil
}

// GetMyTraderSubscription fetches a single trader subscription for the authenticated user.
func (s *TraderSubscriptionService) GetMyTraderSubscription(subscriptionID, userID uint) (*models.TraderSubscription, error) {
	subscription, err := s.repo.GetTraderSubscriptionByID(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve subscription: %w", err)
	}
	if subscription == nil || subscription.UserID != userID { // Authorization check in service
		return nil, errors.New("subscription not found or not authorized")
	}
	return subscription, nil
}

// ListMyTraderSubscriptions fetches all trader subscriptions for the authenticated user.
func (s *TraderSubscriptionService) ListMyTraderSubscriptions(userID uint) ([]models.TraderSubscription, error) {
	subscriptions, err := s.repo.GetTraderSubscriptionsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	return subscriptions, nil
}

// UpdateTraderSubscriptionSettings allows a customer to update allocation and risk multiplier.
func (s *TraderSubscriptionService) UpdateTraderSubscriptionSettings(subscriptionID, userID uint, allocation, riskMultiplier float64) (*models.TraderSubscription, error) {
	subscription, err := s.repo.GetTraderSubscriptionByID(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve subscription: %w", err)
	}
	if subscription == nil || subscription.UserID != userID { // Authorization check in service
		return nil, errors.New("subscription not found or not authorized")
	}
	if !subscription.IsActive {
		return nil, errors.New("cannot update an inactive subscription")
	}

	// Basic validation for allocation and risk (e.g., must be positive)
	if allocation <= 0 || riskMultiplier <= 0 {
		return nil, errors.New("allocation and risk multiplier must be positive")
	}
	// Add more complex validation as needed (e.g., max limits)

	subscription.Allocation = allocation
	subscription.RiskMultiplier = riskMultiplier

	if err := s.repo.UpdateTraderSubscription(subscription); err != nil {
		return nil, fmt.Errorf("failed to update subscription settings: %w", err)
	}
	return subscription, nil
}

// PauseTraderCopyTrading pauses a customer's copy trading for a specific subscription.
func (s *TraderSubscriptionService) PauseTraderCopyTrading(subscriptionID, userID uint) error {
	subscription, err := s.repo.GetTraderSubscriptionByID(subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve subscription: %w", err)
	}
	if subscription == nil || subscription.UserID != userID { // Authorization check in service
		return errors.New("subscription not found or not authorized")
	}
	if !subscription.IsActive {
		return errors.New("cannot pause an inactive subscription")
	}
	if subscription.IsPaused {
		return errors.New("subscription is already paused")
	}

	return s.repo.PauseTraderSubscription(subscriptionID)
}

func (s *TraderSubscriptionService) ResumeTraderCopyTrading(subscriptionID, userID uint) error {
	subscription, err := s.repo.GetTraderSubscriptionByID(subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve subscription: %w", err)
	}
	if subscription == nil || subscription.UserID != userID { // Authorization check in service
		return errors.New("subscription not found or not authorized")
	}
	if !subscription.IsActive {
		return errors.New("cannot resume an inactive subscription")
	}
	if !subscription.IsPaused {
		return errors.New("subscription is not paused")
	}

	return s.repo.ResumeTraderSubscription(subscriptionID)
}

func (s *TraderSubscriptionService) CancelTraderSubscription(subscriptionID, userID uint) error {
	subscription, err := s.repo.GetTraderSubscriptionByID(subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve subscription: %w", err)
	}
	if subscription == nil || subscription.UserID != userID { // Authorization check in service
		return errors.New("subscription not found or not authorized")
	}
	if !subscription.IsActive {
		return errors.New("subscription is already inactive or cancelled")
	}

	now := time.Now()
	if err := s.repo.UpdateSubscriptionStatus(subscriptionID, false, true, &now); err != nil { // Set isActive=false, isPaused=true
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}
	return nil
}

func (s *TraderSubscriptionService) SimulateTraderSubscription(planID uint, initialCapital float64) (interface{}, error) {
	plan, err := s.GetTraderSubscriptionPlanByID(planID)
	if err != nil {
		return nil, err
	}
	if plan == nil { // This check is technically redundant if GetTraderSubscriptionPlanByID returns an error for nil.
		return nil, errors.New("trader subscription plan not found for simulation")
	}

	simulationResult := map[string]interface{}{
		"plan_name":        plan.Name,
		"initial_capital":  initialCapital,
		"simulated_return": fmt.Sprintf("%.2f", initialCapital*1.15), // Dummy 15% return
		"max_drawdown":     fmt.Sprintf("%.2f", initialCapital*0.05), // Dummy 5% drawdown
		"duration":         fmt.Sprintf("%d %s", plan.Duration, plan.Interval),
		"message":          "This is a placeholder simulation. Implement detailed logic here.",
	}

	return simulationResult, nil
}
