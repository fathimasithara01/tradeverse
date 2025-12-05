package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ISubscriptionPlanRepository interface {
	CreateSubscriptionPlan(plan *models.AdminTraderSubscriptionPlan) error
	GetAllSubscriptionPlans() ([]models.AdminTraderSubscriptionPlan, error)
	GetSubscriptionPlanByID(id uint) (*models.AdminTraderSubscriptionPlan, error)
	UpdateSubscriptionPlan(plan *models.AdminTraderSubscriptionPlan) error
	DeleteSubscriptionPlan(id uint) error
}

type SubscriptionPlanRepository struct {
	DB *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) *SubscriptionPlanRepository {
	return &SubscriptionPlanRepository{DB: db}
}

func (r *SubscriptionPlanRepository) CreateSubscriptionPlan(plan *models.AdminTraderSubscriptionPlan) error {
	return r.DB.Create(plan).Error
}

func (r *SubscriptionPlanRepository) GetAllSubscriptionPlans() ([]models.AdminTraderSubscriptionPlan, error) {
	var plans []models.AdminTraderSubscriptionPlan
	err := r.DB.Find(&plans).Error
	return plans, err
}

func (r *SubscriptionPlanRepository) GetSubscriptionPlanByID(id uint) (*models.AdminTraderSubscriptionPlan, error) {
	var plan models.AdminTraderSubscriptionPlan
	err := r.DB.First(&plan, id).Error
	return &plan, err
}

func (r *SubscriptionPlanRepository) UpdateSubscriptionPlan(plan *models.AdminTraderSubscriptionPlan) error {

	return r.DB.Save(plan).Error
}

func (r *SubscriptionPlanRepository) DeleteSubscriptionPlan(id uint) error {
	return r.DB.Delete(&models.AdminTraderSubscriptionPlan{}, id).Error

}
