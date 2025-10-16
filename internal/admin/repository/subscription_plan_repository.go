package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ISubscriptionPlanRepository interface {
	CreateSubscriptionPlan(plan *models.AdminSubscriptionPlan) error
	GetAllSubscriptionPlans() ([]models.AdminSubscriptionPlan, error)
	GetSubscriptionPlanByID(id uint) (*models.AdminSubscriptionPlan, error)
	UpdateSubscriptionPlan(plan *models.AdminSubscriptionPlan) error
	DeleteSubscriptionPlan(id uint) error
}

type SubscriptionPlanRepository struct {
	DB *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) *SubscriptionPlanRepository {
	return &SubscriptionPlanRepository{DB: db}
}

func (r *SubscriptionPlanRepository) CreateSubscriptionPlan(plan *models.AdminSubscriptionPlan) error {
	return r.DB.Create(plan).Error
}

func (r *SubscriptionPlanRepository) GetAllSubscriptionPlans() ([]models.AdminSubscriptionPlan, error) {
	var plans []models.AdminSubscriptionPlan
	err := r.DB.Find(&plans).Error
	return plans, err
}

func (r *SubscriptionPlanRepository) GetSubscriptionPlanByID(id uint) (*models.AdminSubscriptionPlan, error) {
	var plan models.AdminSubscriptionPlan
	err := r.DB.First(&plan, id).Error
	return &plan, err
}

func (r *SubscriptionPlanRepository) UpdateSubscriptionPlan(plan *models.AdminSubscriptionPlan) error {

	return r.DB.Save(plan).Error
}

func (r *SubscriptionPlanRepository) DeleteSubscriptionPlan(id uint) error {
	return r.DB.Delete(&models.AdminSubscriptionPlan{}, id).Error

	// return r.DB.Unscoped().Delete(&models.SubscriptionPlan{}, id).Error
}
