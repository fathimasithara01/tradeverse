package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ISubscriptionRepository interface {
	CreateSubscription(subscription *models.Subscription) error
	GetAllSubscriptions() ([]models.Subscription, error)
	GetSubscriptionByID(id uint) (*models.Subscription, error)
	GetSubscriptionsByUserID(userID uint) ([]models.Subscription, error)
	UpdateSubscription(subscription *models.Subscription) error
	DeleteSubscription(id uint) error
}

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) CreateSubscription(subscription *models.Subscription) error {
	return r.DB.Create(subscription).Error
}

func (r *SubscriptionRepository) GetAllSubscriptions() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.DB.Preload("User").Preload("SubscriptionPlan").Find(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) GetSubscriptionByID(id uint) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.DB.Preload("User").Preload("SubscriptionPlan").First(&subscription, id).Error
	return &subscription, err
}

func (r *SubscriptionRepository) GetSubscriptionsByUserID(userID uint) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.DB.Where("user_id = ?", userID).Preload("SubscriptionPlan").Find(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) UpdateSubscription(subscription *models.Subscription) error {
	return r.DB.Save(subscription).Error
}

func (r *SubscriptionRepository) DeleteSubscription(id uint) error {
	return r.DB.Unscoped().Delete(&models.Subscription{}, id).Error
}
