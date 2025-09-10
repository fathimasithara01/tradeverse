package service

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/fathimasithara01/tradeverse/pkg/repository"
)

type ISubscriptionPlanService interface {
	CreateSubscriptionPlan(plan *models.SubscriptionPlan) error
	GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error)
	GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)
	UpdateSubscriptionPlan(plan *models.SubscriptionPlan) error
	DeleteSubscriptionPlan(id uint) error
}

type SubscriptionPlanService struct {
	repo repository.ISubscriptionPlanRepository
}

func NewSubscriptionPlanService(repo repository.ISubscriptionPlanRepository) *SubscriptionPlanService {
	return &SubscriptionPlanService{repo: repo}
}

func (s *SubscriptionPlanService) CreateSubscriptionPlan(plan *models.SubscriptionPlan) error {
	return s.repo.CreateSubscriptionPlan(plan)
}

func (s *SubscriptionPlanService) GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	return s.repo.GetAllSubscriptionPlans()
}

func (s *SubscriptionPlanService) GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error) {
	return s.repo.GetSubscriptionPlanByID(id)
}

func (s *SubscriptionPlanService) UpdateSubscriptionPlan(plan *models.SubscriptionPlan) error {
	return s.repo.UpdateSubscriptionPlan(plan)
}

func (s *SubscriptionPlanService) DeleteSubscriptionPlan(id uint) error {
	return s.repo.DeleteSubscriptionPlan(id)
}
