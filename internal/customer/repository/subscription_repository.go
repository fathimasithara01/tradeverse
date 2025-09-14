package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	CreateSubscription(subscription *models.Subscription) error
	GetSubscriptionByID(id uint, userID uint) (*models.Subscription, error)
	ListSubscriptionsByUserID(userID uint) ([]models.Subscription, error)
	UpdateSubscription(subscription *models.Subscription) error
	DeleteSubscription(id uint, userID uint) error
	// Add methods for fetching subscription plans
	GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)
	ListActiveSubscriptionPlans() ([]models.SubscriptionPlan, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) CreateSubscription(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *subscriptionRepository) GetSubscriptionByID(id uint, userID uint) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Preload("SubscriptionPlan").Preload("User").Where("id = ? AND user_id = ?", id, userID).First(&subscription).Error
	return &subscription, err
}

func (r *subscriptionRepository) ListSubscriptionsByUserID(userID uint) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Preload("SubscriptionPlan").Preload("User").Where("user_id = ?", userID).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepository) UpdateSubscription(subscription *models.Subscription) error {
	return r.db.Save(subscription).Error
}

func (r *subscriptionRepository) DeleteSubscription(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Subscription{}).Error
}

func (r *subscriptionRepository) GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	err := r.db.First(&plan, id).Error
	return &plan, err
}

func (r *subscriptionRepository) ListActiveSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	err := r.db.Where("is_active = ?", true).Find(&plans).Error
	return plans, err
}
