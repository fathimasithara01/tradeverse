package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ISubscriptionPlanRepository interface {
	CreateSubscriptionPlan(plan *models.SubscriptionPlan) error
	GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error)
	GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)
	UpdateSubscriptionPlan(plan *models.SubscriptionPlan) error
	DeleteSubscriptionPlan(id uint) error
}

type SubscriptionPlanRepository struct {
	DB *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) *SubscriptionPlanRepository {
	return &SubscriptionPlanRepository{DB: db}
}

func (r *SubscriptionPlanRepository) CreateSubscriptionPlan(plan *models.SubscriptionPlan) error {
	return r.DB.Create(plan).Error
}

func (r *SubscriptionPlanRepository) GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	err := r.DB.Find(&plans).Error
	return plans, err
}

func (r *SubscriptionPlanRepository) GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	err := r.DB.First(&plan, id).Error
	return &plan, err
}

func (r *SubscriptionPlanRepository) UpdateSubscriptionPlan(plan *models.SubscriptionPlan) error {
	// GORM's Save method works for both creation and update.
	// When a primary key is set (ID > 0), it updates; otherwise, it creates.
	return r.DB.Save(plan).Error
}

func (r *SubscriptionPlanRepository) DeleteSubscriptionPlan(id uint) error {
	// Unscoped() is used for hard delete. Without it, GORM performs a soft delete
	// if the model has gorm.DeletedAt field.
	return r.DB.Unscoped().Delete(&models.SubscriptionPlan{}, id).Error
}
