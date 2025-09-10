package service

import (
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/fathimasithara01/tradeverse/pkg/repository"
)

type ISubscriptionService interface {
	CreateSubscription(userID, planID uint, amount float64, transactionID string) (*models.Subscription, error)
	GetAllSubscriptions() ([]models.Subscription, error)
	GetSubscriptionByID(id uint) (*models.Subscription, error)
	GetSubscriptionsByUserID(userID uint) ([]models.Subscription, error)
	UpdateSubscription(subscription *models.Subscription) error
	DeleteSubscription(id uint) error
	GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)
}

type SubscriptionService struct {
	subscriptionRepo repository.ISubscriptionRepository
	planRepo         repository.ISubscriptionPlanRepository
}

func NewSubscriptionService(subRepo repository.ISubscriptionRepository, planRepo repository.ISubscriptionPlanRepository) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subRepo,
		planRepo:         planRepo,
	}
}

func (s *SubscriptionService) CreateSubscription(userID, planID uint, amount float64, transactionID string) (*models.Subscription, error) {
	plan, err := s.planRepo.GetSubscriptionPlanByID(planID)
	if err != nil {
		return nil, err
	}

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
		endDate = startDate.AddDate(0, 1, 0)
	}

	subscription := &models.Subscription{
		UserID:             userID,
		SubscriptionPlanID: planID,
		StartDate:          startDate,
		EndDate:            endDate,
		IsActive:           true,
		PaymentStatus:      "paid",
		AmountPaid:         amount,
		TransactionID:      transactionID,
	}

	err = s.subscriptionRepo.CreateSubscription(subscription)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *SubscriptionService) GetAllSubscriptions() ([]models.Subscription, error) {
	return s.subscriptionRepo.GetAllSubscriptions()
}

func (s *SubscriptionService) GetSubscriptionByID(id uint) (*models.Subscription, error) {
	return s.subscriptionRepo.GetSubscriptionByID(id)
}

func (s *SubscriptionService) GetSubscriptionsByUserID(userID uint) ([]models.Subscription, error) {
	return s.subscriptionRepo.GetSubscriptionsByUserID(userID)
}

func (s *SubscriptionService) UpdateSubscription(subscription *models.Subscription) error {
	return s.subscriptionRepo.UpdateSubscription(subscription)
}

func (s *SubscriptionService) DeleteSubscription(id uint) error {
	return s.subscriptionRepo.DeleteSubscription(id)
}

func (s *SubscriptionService) GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error) {
	return s.planRepo.GetSubscriptionPlanByID(id)
}
