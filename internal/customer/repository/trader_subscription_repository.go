package repository

import (
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

// ITraderSubscriptionRepository defines the interface for customer-facing trader subscription data operations.
type ITraderSubscriptionRepository interface {
	CreateTraderSubscription(subscription *models.TraderSubscription) error
	GetTraderSubscriptionByID(id uint) (*models.TraderSubscription, error) // No user ID needed here, service handles authorization
	GetTraderSubscriptionsByUserID(userID uint) ([]models.TraderSubscription, error)
	UpdateTraderSubscription(subscription *models.TraderSubscription) error
	UpdateSubscriptionStatus(id uint, isActive, isPaused bool, endDate *time.Time) error

	PauseTraderSubscription(id uint) error
	ResumeTraderSubscription(id uint) error

	// Get available trader subscription plans
	GetActiveTraderSubscriptionPlans() ([]models.SubscriptionPlan, error)
	GetTraderSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)
	GetTraderProfileByID(id uint) (*models.TraderProfile, error)
}

// TraderSubscriptionRepository implements ITraderSubscriptionRepository
type TraderSubscriptionRepository struct {
	DB *gorm.DB
}

// NewTraderSubscriptionRepository creates a new TraderSubscriptionRepository
func NewTraderSubscriptionRepository(db *gorm.DB) *TraderSubscriptionRepository {
	return &TraderSubscriptionRepository{DB: db}
}

// CreateTraderSubscription creates a new trader subscription record in the database.
func (r *TraderSubscriptionRepository) CreateTraderSubscription(subscription *models.TraderSubscription) error {
	return r.DB.Create(subscription).Error
}

// GetTraderSubscriptionByID fetches a single trader subscription by its ID, including related User, Plan, and TraderProfile.
// This method takes only the subscription ID. Authorization logic (checking userID) should be in the service layer.
func (r *TraderSubscriptionRepository) GetTraderSubscriptionByID(id uint) (*models.TraderSubscription, error) {
	var subscription models.TraderSubscription
	err := r.DB.Preload("User").Preload("TraderSubscriptionPlan").Preload("TraderProfile").First(&subscription, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // Return nil if not found
	}
	return &subscription, err
}

// GetTraderSubscriptionsByUserID fetches all trader subscriptions for a given user ID, including related data.
func (r *TraderSubscriptionRepository) GetTraderSubscriptionsByUserID(userID uint) ([]models.TraderSubscription, error) {
	var subscriptions []models.TraderSubscription
	err := r.DB.Where("user_id = ?", userID).Preload("TraderSubscriptionPlan").Preload("TraderProfile").Find(&subscriptions).Error
	return subscriptions, err
}

// UpdateTraderSubscription updates an existing trader subscription record.
func (r *TraderSubscriptionRepository) UpdateTraderSubscription(subscription *models.TraderSubscription) error {
	return r.DB.Save(subscription).Error
}

// UpdateSubscriptionStatus updates the active and paused status of a subscription, and optionally its end date.
// This replaces the old DeleteTraderSubscription logic for cancellation.
func (r *TraderSubscriptionRepository) UpdateSubscriptionStatus(id uint, isActive, isPaused bool, endDate *time.Time) error {
	updates := map[string]interface{}{
		"is_active": isActive,
		"is_paused": isPaused,
	}
	if endDate != nil {
		updates["end_date"] = *endDate
	}
	return r.DB.Model(&models.TraderSubscription{}).Where("id = ?", id).Updates(updates).Error
}

// PauseTraderSubscription sets the IsPaused flag to true and updates last_pause_date.
func (r *TraderSubscriptionRepository) PauseTraderSubscription(id uint) error {
	return r.DB.Model(&models.TraderSubscription{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_paused":       true,
		"last_pause_date": time.Now(),
	}).Error
}

// ResumeTraderSubscription sets the IsPaused flag to false and updates last_resume_date.
func (r *TraderSubscriptionRepository) ResumeTraderSubscription(id uint) error {
	return r.DB.Model(&models.TraderSubscription{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_paused":        false,
		"last_resume_date": time.Now(),
	}).Error
}

// GetActiveTraderSubscriptionPlans fetches all active subscription plans marked as trader plans.
func (r *TraderSubscriptionRepository) GetActiveTraderSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	// Filter for plans that are active and specifically marked as trader plans
	err := r.DB.Where("is_active = ? AND is_trader_plan = ?", true, true).Find(&plans).Error
	return plans, err
}

// GetTraderSubscriptionPlanByID fetches a specific active trader subscription plan by ID.
func (r *TraderSubscriptionRepository) GetTraderSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	// This is the query causing "record not found" if the plan doesn't meet conditions
	err := r.DB.Where("id = ? AND is_active = ? AND is_trader_plan = ?", id, true, true).First(&plan).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &plan, err
}

// GetTraderProfileByID fetches a trader profile by ID.
func (r *TraderSubscriptionRepository) GetTraderProfileByID(id uint) (*models.TraderProfile, error) {
	var traderProfile models.TraderProfile
	err := r.DB.Preload("User").First(&traderProfile, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &traderProfile, err
}
