package service

import (
	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type ISubscriptionPlanService interface {
	CreateSubscriptionPlan(plan *models.AdminSubscriptionPlan) error
	GetAllSubscriptionPlans() ([]models.AdminSubscriptionPlan, error)
	GetSubscriptionPlanByID(id uint) (*models.AdminSubscriptionPlan, error)
	UpdateSubscriptionPlan(plan *models.AdminSubscriptionPlan) error
	DeleteSubscriptionPlan(id uint) error
}

type SubscriptionPlanService struct {
	repo repository.ISubscriptionPlanRepository
}

func NewSubscriptionPlanService(repo repository.ISubscriptionPlanRepository) *SubscriptionPlanService {
	return &SubscriptionPlanService{repo: repo}
}

func (s *SubscriptionPlanService) CreateSubscriptionPlan(plan *models.AdminSubscriptionPlan) error {
	return s.repo.CreateSubscriptionPlan(plan)
}

func (s *SubscriptionPlanService) GetAllSubscriptionPlans() ([]models.AdminSubscriptionPlan, error) {
	return s.repo.GetAllSubscriptionPlans()
}

func (s *SubscriptionPlanService) GetSubscriptionPlanByID(id uint) (*models.AdminSubscriptionPlan, error) {
	return s.repo.GetSubscriptionPlanByID(id)
}

func (s *SubscriptionPlanService) UpdateSubscriptionPlan(plan *models.AdminSubscriptionPlan) error {
	return s.repo.UpdateSubscriptionPlan(plan)
}

func (s *SubscriptionPlanService) DeleteSubscriptionPlan(id uint) error {
	return s.repo.DeleteSubscriptionPlan(id)
}
