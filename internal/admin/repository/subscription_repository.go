package repository

import (
	"fmt" // Import fmt for structured errors
	"log"
	"time"

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
	GetExpiredActiveSubscriptions() ([]models.Subscription, error) // New method

}

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) GetExpiredActiveSubscriptions() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.DB.
		Where("is_active = ? AND end_date < ?", true, time.Now()).
		Preload("User").
		Find(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) CreateSubscription(subscription *models.Subscription) error {
	return r.DB.Create(subscription).Error
}

func (r *SubscriptionRepository) GetAllSubscriptions() ([]models.Subscription, error) {
	log.Println("DEBUG: SubscriptionRepository.GetAllSubscriptions was called.") // Add this
	var subscriptions []models.Subscription
	err := r.DB.
		Preload("User").
		Preload("User.TraderProfile").
		Preload("SubscriptionPlan").
		Find(&subscriptions).Error

	if err != nil {
		log.Printf("ERROR: failed to fetch all subscriptions from DB: %v", err) // Add this
		return nil, fmt.Errorf("failed to fetch all subscriptions from DB: %w", err)
	}
	log.Printf("DEBUG: SubscriptionRepository found %d subscriptions.", len(subscriptions)) // Add this
	return subscriptions, nil
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
