package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type PlanRepository struct{}

func (r *PlanRepository) GetAllPendingPlans() ([]models.Plan, error) {
	var plans []models.Plan
	err := db.DB.Where("status = ?", "pending").Find(&plans).Error
	return plans, err
}

func (r *PlanRepository) UpdatePlanStatus(id uint, status string) error {
	return db.DB.Model(&models.Plan{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *PlanRepository) GetAll() ([]models.Plan, error) {
	var plans []models.Plan
	err := db.DB.Find(&plans).Error
	return plans, err
}

func (r *PlanRepository) Create(plan models.Plan) (models.Plan, error) {
	err := db.DB.Create(&plan).Error
	return plan, err
}

func (r *PlanRepository) Update(id uint, data models.Plan) error {
	return db.DB.Model(&models.Plan{}).Where("id = ?", id).Updates(data).Error
}

func (r *PlanRepository) Delete(id uint) error {
	return db.DB.Delete(&models.Plan{}, id).Error
}
