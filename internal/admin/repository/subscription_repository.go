package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ISubscriptionRepository interface {
	CreateSubscription(subscription *models.CustomerToTraderSub) error
	GetAllSubscriptions() ([]models.CustomerToTraderSub, error)
	GetSubscriptionByID(id uint) (*models.CustomerToTraderSub, error)
	GetSubscriptionsByUserID(userID uint) ([]models.CustomerToTraderSub, error)
	UpdateSubscription(subscription *models.CustomerToTraderSub) error
	DeleteSubscription(id uint) error
	GetExpiredActiveSubscriptions() ([]models.CustomerToTraderSub, error)
}

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) GetExpiredActiveSubscriptions() ([]models.CustomerToTraderSub, error) {
	var subscriptions []models.CustomerToTraderSub
	err := r.DB.
		Where("is_active = ? AND end_date < ?", true, time.Now()).
		Preload("User").
		Find(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) CreateSubscription(subscription *models.CustomerToTraderSub) error {
	return r.DB.Create(subscription).Error
}

func (r *SubscriptionRepository) GetAllSubscriptions() ([]models.CustomerToTraderSub, error) {
	log.Println("DEBUG: SubscriptionRepository.GetAllSubscriptions was called.")
	var subscriptions []models.CustomerToTraderSub
	err := r.DB.
		Preload("User").
		Preload("User.TraderProfile").
		Preload("SubscriptionPlan").
		Find(&subscriptions).Error

	if err != nil {
		log.Printf("ERROR: failed to fetch all subscriptions from DB: %v", err)
		return nil, fmt.Errorf("failed to fetch all subscriptions from DB: %w", err)
	}
	log.Printf("DEBUG: SubscriptionRepository found %d subscriptions.", len(subscriptions))
	return subscriptions, nil
}

func (r *SubscriptionRepository) GetSubscriptionByID(id uint) (*models.CustomerToTraderSub, error) {
	var subscription models.CustomerToTraderSub
	err := r.DB.Preload("User").Preload("SubscriptionPlan").First(&subscription, id).Error
	return &subscription, err
}

func (r *SubscriptionRepository) GetSubscriptionsByUserID(userID uint) ([]models.CustomerToTraderSub, error) {
	var subscriptions []models.CustomerToTraderSub
	err := r.DB.Where("user_id = ?", userID).Preload("SubscriptionPlan").Find(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) UpdateSubscription(subscription *models.CustomerToTraderSub) error {
	return r.DB.Save(subscription).Error
}

func (r *SubscriptionRepository) DeleteSubscription(id uint) error {
	return r.DB.Unscoped().Delete(&models.CustomerToTraderSub{}, id).Error
}
