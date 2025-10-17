package customerrepo

import (
	"errors"
	"fmt"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ICustomerSubscriptionPlanRepository interface {
	GetAllSubscriptionPlans() ([]models.AdminTraderSubscriptionPlan, error)
	GetSubscriptionPlanByID(id uint) (*models.AdminTraderSubscriptionPlan, error)
}

type CustomerSubscriptionPlanRepository struct {
	DB *gorm.DB
}

func NewCustomerSubscriptionPlanRepository(db *gorm.DB) *CustomerSubscriptionPlanRepository {
	return &CustomerSubscriptionPlanRepository{DB: db}
}

func (r *CustomerSubscriptionPlanRepository) GetAllSubscriptionPlans() ([]models.AdminTraderSubscriptionPlan, error) {
	var plans []models.AdminTraderSubscriptionPlan
	err := r.DB.Where("is_active = ?", true).Find(&plans).Error
	return plans, err
}

func (r *CustomerSubscriptionPlanRepository) GetSubscriptionPlanByID(id uint) (*models.AdminTraderSubscriptionPlan, error) {
	var plan models.AdminTraderSubscriptionPlan
	err := r.DB.Where("is_active = ?", true).First(&plan, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("subscription plan not found or not active")
	}
	return &plan, err
}
