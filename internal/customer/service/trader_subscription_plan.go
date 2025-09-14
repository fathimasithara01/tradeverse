package service

import (
	"errors"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type AdminTraderSubscriptionPlanService interface {
	CreateTraderSubscriptionPlan(plan *models.TraderSubscriptionPlan) (*models.TraderSubscriptionPlan, error)
	GetTraderSubscriptionPlanByID(id uint) (*models.TraderSubscriptionPlan, error)
	ListTraderSubscriptionPlans() ([]models.TraderSubscriptionPlan, error)
	UpdateTraderSubscriptionPlan(id uint, updates map[string]interface{}) (*models.TraderSubscriptionPlan, error)
	DeleteTraderSubscriptionPlan(id uint) error
	ToggleTraderSubscriptionPlanStatus(id uint, isActive bool) (*models.TraderSubscriptionPlan, error)
}

type adminTraderSubscriptionPlanService struct {
	repo repository.AdminTraderSubscriptionPlanRepository
}

func NewAdminTraderSubscriptionPlanService(repo repository.AdminTraderSubscriptionPlanRepository) AdminTraderSubscriptionPlanService {
	return &adminTraderSubscriptionPlanService{repo: repo}
}

func (s *adminTraderSubscriptionPlanService) CreateTraderSubscriptionPlan(plan *models.TraderSubscriptionPlan) (*models.TraderSubscriptionPlan, error) {
	// Basic validation
	if plan.Name == "" || plan.Price <= 0 || plan.Duration <= 0 || plan.Interval == "" {
		return nil, errors.New("invalid plan data: name, price, duration, and interval are required")
	}
	err := s.repo.CreateTraderSubscriptionPlan(plan)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *adminTraderSubscriptionPlanService) GetTraderSubscriptionPlanByID(id uint) (*models.TraderSubscriptionPlan, error) {
	plan, err := s.repo.GetTraderSubscriptionPlanByID(id)
	if err != nil {
		return nil, errors.New("trader subscription plan not found")
	}
	return plan, nil
}

func (s *adminTraderSubscriptionPlanService) ListTraderSubscriptionPlans() ([]models.TraderSubscriptionPlan, error) {
	return s.repo.ListTraderSubscriptionPlans()
}

func (s *adminTraderSubscriptionPlanService) UpdateTraderSubscriptionPlan(id uint, updates map[string]interface{}) (*models.TraderSubscriptionPlan, error) {
	plan, err := s.repo.GetTraderSubscriptionPlanByID(id)
	if err != nil {
		return nil, errors.New("trader subscription plan not found")
	}

	// Apply updates. GORM's `Updates` method can handle this, but for explicit control
	// and validation, you might iterate through `updates` map.
	// For simplicity, directly map common fields if present in updates
	if name, ok := updates["name"].(string); ok && name != "" {
		plan.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		plan.Description = description
	}
	if price, ok := updates["price"].(float64); ok && price > 0 {
		plan.Price = price
	}
	if duration, ok := updates["duration"].(float64); ok && duration > 0 { // JSON numbers are often float64
		plan.Duration = int(duration)
	}
	if interval, ok := updates["interval"].(string); ok && interval != "" {
		plan.Interval = interval
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		plan.IsActive = isActive
	}
	if features, ok := updates["features"].(string); ok {
		plan.Features = features
	}
	if maxFollowers, ok := updates["max_followers"].(float64); ok && maxFollowers >= 0 {
		plan.MaxFollowers = int(maxFollowers)
	}
	if commissionRate, ok := updates["commission_rate"].(float64); ok && commissionRate >= 0 && commissionRate <= 1 {
		plan.CommissionRate = commissionRate
	}
	if analyticsAccess, ok := updates["analytics_access"].(string); ok {
		plan.AnalyticsAccess = analyticsAccess
	}

	err = s.repo.UpdateTraderSubscriptionPlan(plan)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *adminTraderSubscriptionPlanService) DeleteTraderSubscriptionPlan(id uint) error {
	_, err := s.repo.GetTraderSubscriptionPlanByID(id) // Check existence first
	if err != nil {
		return errors.New("trader subscription plan not found")
	}
	return s.repo.DeleteTraderSubscriptionPlan(id)
}

func (s *adminTraderSubscriptionPlanService) ToggleTraderSubscriptionPlanStatus(id uint, isActive bool) (*models.TraderSubscriptionPlan, error) {
	plan, err := s.repo.GetTraderSubscriptionPlanByID(id)
	if err != nil {
		return nil, errors.New("trader subscription plan not found")
	}
	plan.IsActive = isActive
	err = s.repo.UpdateTraderSubscriptionPlan(plan)
	if err != nil {
		return nil, err
	}
	return plan, nil
}
