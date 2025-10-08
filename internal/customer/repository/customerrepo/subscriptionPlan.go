package customerrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

// ISubscriptionPlanRepository defines methods for interacting with SubscriptionPlan data.
type ISubscriptionPlanRepository interface {
	GetSubscriptionPlanByID(ctx context.Context, planID uint) (*models.SubscriptionPlan, error)
	// Add other subscription plan-related methods as needed (e.g., GetTraderSubscriptionPlans, GetAllActivePlans)
}

type subscriptionPlanRepository struct {
	db *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) ISubscriptionPlanRepository {
	return &subscriptionPlanRepository{db: db}
}

func (r *subscriptionPlanRepository) GetSubscriptionPlanByID(ctx context.Context, planID uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	err := r.db.WithContext(ctx).First(&plan, planID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription plan not found")
		}
		return nil, fmt.Errorf("failed to get subscription plan by ID: %w", err)
	}
	return &plan, nil
}
