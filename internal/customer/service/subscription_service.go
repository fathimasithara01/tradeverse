package service

import (
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type SubscriptionService interface {
	SubscribeToTrader(userID uint, planID uint, amountPaid float64, transactionID string) (*models.Subscription, error)
	ListMySubscriptions(userID uint) ([]models.Subscription, error)
	GetSubscriptionDetails(id uint, userID uint) (*models.Subscription, error)
	UpdateSubscription(id uint, userID uint, updates map[string]interface{}) (*models.Subscription, error)
	PauseCopyTrading(id uint, userID uint) (*models.Subscription, error)
	ResumeCopyTrading(id uint, userID uint) (*models.Subscription, error)
	CancelSubscription(id uint, userID uint) error
	RunSimulation(id uint, userID uint) (string, error) // Placeholder for simulation logic
	ListAvailableSubscriptionPlans() ([]models.SubscriptionPlan, error)
}

type subscriptionService struct {
	subscriptionRepo repository.SubscriptionRepository
}

func NewSubscriptionService(subscriptionRepo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{subscriptionRepo: subscriptionRepo}
}

func (s *subscriptionService) SubscribeToTrader(userID uint, planID uint, amountPaid float64, transactionID string) (*models.Subscription, error) {
	plan, err := s.subscriptionRepo.GetSubscriptionPlanByID(planID)
	if err != nil {
		return nil, errors.New("subscription plan not found")
	}

	if !plan.IsActive {
		return nil, errors.New("subscription plan is not active")
	}

	// In a real application, you'd integrate with a payment gateway here
	// For now, assume payment is successful if amountPaid is positive
	if amountPaid <= 0 {
		return nil, errors.New("payment amount must be positive")
	}

	startDate := time.Now()
	endDate := calculateEndDate(startDate, plan.Duration, plan.Interval)

	subscription := &models.Subscription{
		UserID:             userID,
		SubscriptionPlanID: planID,
		StartDate:          startDate,
		EndDate:            endDate,
		IsActive:           true,
		PaymentStatus:      "paid", // Assume paid for now
		AmountPaid:         amountPaid,
		TransactionID:      transactionID,
	}

	err = s.subscriptionRepo.CreateSubscription(subscription)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *subscriptionService) ListMySubscriptions(userID uint) ([]models.Subscription, error) {
	return s.subscriptionRepo.ListSubscriptionsByUserID(userID)
}

func (s *subscriptionService) GetSubscriptionDetails(id uint, userID uint) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetSubscriptionByID(id, userID)
	if err != nil {
		return nil, errors.New("subscription not found or unauthorized")
	}
	return subscription, nil
}

func (s *subscriptionService) UpdateSubscription(id uint, userID uint, updates map[string]interface{}) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetSubscriptionByID(id, userID)
	if err != nil {
		return nil, errors.New("subscription not found or unauthorized")
	}

	// Apply updates (e.g., allocation/risk - these fields would be in your Subscription model)
	// Example:
	// if allocation, ok := updates["allocation"].(float64); ok {
	// 	subscription.Allocation = allocation
	// }
	// if riskLevel, ok := updates["risk_level"].(string); ok {
	// 	subscription.RiskLevel = riskLevel
	// }

	err = s.subscriptionRepo.UpdateSubscription(subscription)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *subscriptionService) PauseCopyTrading(id uint, userID uint) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetSubscriptionByID(id, userID)
	if err != nil {
		return nil, errors.New("subscription not found or unauthorized")
	}
	subscription.IsActive = false // Or add a specific 'IsPaused' field to the model
	err = s.subscriptionRepo.UpdateSubscription(subscription)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *subscriptionService) ResumeCopyTrading(id uint, userID uint) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetSubscriptionByID(id, userID)
	if err != nil {
		return nil, errors.New("subscription not found or unauthorized")
	}
	subscription.IsActive = true // Or set 'IsPaused' to false
	err = s.subscriptionRepo.UpdateSubscription(subscription)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *subscriptionService) CancelSubscription(id uint, userID uint) error {
	// In a real scenario, you might mark it as cancelled rather than deleting
	// and handle refunds if applicable.
	subscription, err := s.subscriptionRepo.GetSubscriptionByID(id, userID)
	if err != nil {
		return errors.New("subscription not found or unauthorized")
	}
	subscription.IsActive = false
	subscription.EndDate = time.Now() // End it immediately
	err = s.subscriptionRepo.UpdateSubscription(subscription)
	if err != nil {
		return err
	}
	return nil // Return s.subscriptionRepo.DeleteSubscription(id, userID) if you want to hard delete
}

func (s *subscriptionService) RunSimulation(id uint, userID uint) (string, error) {
	subscription, err := s.subscriptionRepo.GetSubscriptionByID(id, userID)
	if err != nil {
		return "", errors.New("subscription not found or unauthorized")
	}

	simulationResult := "Running simulation for subscription " + subscription.SubscriptionPlan.Name + " for user " + subscription.User.Name + "..."
	simulationResult += "\nExpected ROI: 15% (hypothetical)"
	simulationResult += "\nMaximum Drawdown: 5% (hypothetical)"
	simulationResult += "\nHistorical performance analysis shows potential for growth."

	return simulationResult, nil
}

func (s *subscriptionService) ListAvailableSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	return s.subscriptionRepo.ListActiveSubscriptionPlans()
}

func calculateEndDate(startDate time.Time, duration int, interval string) time.Time {
	switch interval {
	case "days":
		return startDate.Add(time.Duration(duration) * 24 * time.Hour)
	case "months":
		return startDate.AddDate(0, duration, 0)
	case "years":
		return startDate.AddDate(duration, 0, 0)
	default:
		return startDate.AddDate(0, 1, 0) // Default to 1 month
	}
}
