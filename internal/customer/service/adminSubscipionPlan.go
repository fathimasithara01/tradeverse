package service

import (
	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type ICustomerSubscriptionPlanService interface {
	GetAllSubscriptionPlans() ([]models.AdminTraderSubscriptionPlan, error)
	GetSubscriptionPlanByID(id uint) (*models.AdminTraderSubscriptionPlan, error)
}

type CustomerSubscriptionPlanService struct {
	repo customerrepo.ICustomerSubscriptionPlanRepository
}

func NewCustomerSubscriptionPlanService(repo customerrepo.ICustomerSubscriptionPlanRepository) *CustomerSubscriptionPlanService {
	return &CustomerSubscriptionPlanService{repo: repo}
}

func (s *CustomerSubscriptionPlanService) GetAllSubscriptionPlans() ([]models.AdminTraderSubscriptionPlan, error) {
	return s.repo.GetAllSubscriptionPlans()
}

func (s *CustomerSubscriptionPlanService) GetSubscriptionPlanByID(id uint) (*models.AdminTraderSubscriptionPlan, error) {
	return s.repo.GetSubscriptionPlanByID(id)
}
