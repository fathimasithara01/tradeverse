package customerrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrTraderSubscriptionNotFound = errors.New("trader subscription not found")
	ErrSubscriptionAlreadyActive  = errors.New("customer already has an active subscription with this trader")
)

type ITraderSubscriptionRepository interface {
	CreateTraderSubscription(ctx context.Context, sub *models.TraderSubscription) error
	GetActiveTraderSubscriptionForCustomer(ctx context.Context, customerID, traderID uint) (*models.TraderSubscription, error)
	GetTraderSubscriptionByID(ctx context.Context, subscriptionID uint) (*models.TraderSubscription, error)
	UpdateTraderSubscription(ctx context.Context, sub *models.TraderSubscription) error
	DeactivateExpiredTraderSubscriptions(ctx context.Context) error
	GetCustomerTraderSubscriptions(ctx context.Context, customerID uint) ([]models.TraderSubscription, error)
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
	GetSubscriptionPlanByID(ctx context.Context, planID uint) (*models.SubscriptionPlan, error)
}

type traderSubscriptionRepository struct {
	db *gorm.DB
}

func NewTraderSubscriptionRepository(db *gorm.DB) ITraderSubscriptionRepository {
	return &traderSubscriptionRepository{db: db}
}

func (r *traderSubscriptionRepository) CreateTraderSubscription(ctx context.Context, sub *models.TraderSubscription) error {
	return r.db.WithContext(ctx).Create(sub).Error
}

func (r *traderSubscriptionRepository) GetActiveTraderSubscriptionForCustomer(ctx context.Context, customerID, traderID uint) (*models.TraderSubscription, error) {
	var sub models.TraderSubscription
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND trader_id = ? AND is_active = ? AND end_date > ?", customerID, traderID, true, time.Now()).
		First(&sub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active trader subscription: %w", err)
	}
	return &sub, nil
}

func (r *traderSubscriptionRepository) GetTraderSubscriptionByID(ctx context.Context, subscriptionID uint) (*models.TraderSubscription, error) {
	var sub models.TraderSubscription
	err := r.db.WithContext(ctx).First(&sub, subscriptionID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTraderSubscriptionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get trader subscription by ID: %w", err)
	}
	return &sub, nil
}

func (r *traderSubscriptionRepository) UpdateTraderSubscription(ctx context.Context, sub *models.TraderSubscription) error {
	return r.db.WithContext(ctx).Save(sub).Error
}

func (r *traderSubscriptionRepository) DeactivateExpiredTraderSubscriptions(ctx context.Context) error {
	result := r.db.WithContext(ctx).
		Model(&models.TraderSubscription{}).
		Where("is_active = ? AND end_date <= ?", true, time.Now()).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to deactivate expired trader subscriptions: %w", result.Error)
	}
	if result.RowsAffected > 0 {
		fmt.Printf("Deactivated %d expired trader subscriptions.\n", result.RowsAffected)
	}
	return nil
}

func (r *traderSubscriptionRepository) GetCustomerTraderSubscriptions(ctx context.Context, customerID uint) ([]models.TraderSubscription, error) {
	var subs []models.TraderSubscription
	err := r.db.WithContext(ctx).
		Preload("Trader").
		Preload("TraderSubscriptionPlan").
		Where("user_id = ?", customerID).
		Order("end_date DESC").
		Find(&subs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get customer's trader subscriptions: %w", err)
	}
	return subs, nil
}

func (r *traderSubscriptionRepository) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (r *traderSubscriptionRepository) GetSubscriptionPlanByID(ctx context.Context, planID uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	err := r.db.WithContext(ctx).First(&plan, planID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("subscription plan not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription plan by ID: %w", err)
	}
	return &plan, nil
}
