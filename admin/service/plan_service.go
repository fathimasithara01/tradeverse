package service

import (
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type PlanService struct {
	Repo repository.PlanRepository
}

func (s *PlanService) GetPendingPlans() ([]models.Plan, error) {
	return s.Repo.GetAllPendingPlans()
}

func (s *PlanService) ApprovePlan(id uint) error {
	return s.Repo.UpdatePlanStatus(id, "approved")
}

func (s *PlanService) RejectPlan(id uint) error {
	return s.Repo.UpdatePlanStatus(id, "rejected")
}

func (s *PlanService) GetAllPlans() ([]models.Plan, error) {
	return s.Repo.GetAll()
}

func (s *PlanService) CreatePlan(plan models.Plan) (models.Plan, error) {
	return s.Repo.Create(plan)
}

func (s *PlanService) UpdatePlan(id uint, plan models.Plan) error {
	return s.Repo.Update(id, plan)
}

func (s *PlanService) DeletePlan(id uint) error {
	return s.Repo.Delete(id)
}
