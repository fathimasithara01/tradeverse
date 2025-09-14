package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type AdminTraderSubscriptionPlanRepository interface {
	CreateTraderSubscriptionPlan(plan *models.TraderSubscriptionPlan) error
	GetTraderSubscriptionPlanByID(id uint) (*models.TraderSubscriptionPlan, error)
	ListTraderSubscriptionPlans() ([]models.TraderSubscriptionPlan, error)
	UpdateTraderSubscriptionPlan(plan *models.TraderSubscriptionPlan) error
	DeleteTraderSubscriptionPlan(id uint) error
}

type adminTraderSubscriptionPlanRepository struct {
	db *gorm.DB
}

func NewAdminTraderSubscriptionPlanRepository(db *gorm.DB) AdminTraderSubscriptionPlanRepository {
	return &adminTraderSubscriptionPlanRepository{db: db}
}

func (r *adminTraderSubscriptionPlanRepository) CreateTraderSubscriptionPlan(plan *models.TraderSubscriptionPlan) error {
	return r.db.Create(plan).Error
}

func (r *adminTraderSubscriptionPlanRepository) GetTraderSubscriptionPlanByID(id uint) (*models.TraderSubscriptionPlan, error) {
	var plan models.TraderSubscriptionPlan
	err := r.db.First(&plan, id).Error
	return &plan, err
}

func (r *adminTraderSubscriptionPlanRepository) ListTraderSubscriptionPlans() ([]models.TraderSubscriptionPlan, error) {
	var plans []models.TraderSubscriptionPlan
	err := r.db.Find(&plans).Error
	return plans, err
}

func (r *adminTraderSubscriptionPlanRepository) UpdateTraderSubscriptionPlan(plan *models.TraderSubscriptionPlan) error {
	return r.db.Save(plan).Error
}

func (r *adminTraderSubscriptionPlanRepository) DeleteTraderSubscriptionPlan(id uint) error {
	return r.db.Delete(&models.TraderSubscriptionPlan{}, id).Error
}
