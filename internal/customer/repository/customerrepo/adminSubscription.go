package customerrepo

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ICustomerSubscriptionRepository interface {
	CreateSubscription(subscription *models.CustomerToTraderSub) error
	GetSubscriptionByID(id uint) (*models.CustomerToTraderSub, error)
	GetSubscriptionsByUserID(userID uint) ([]models.CustomerToTraderSub, error)
	UpdateSubscription(subscription *models.CustomerToTraderSub) error

	GetActiveTraderSubscriptions() ([]models.CustomerToTraderSub, error)
	UpdateTraderSubscription(subscription *models.CustomerToTraderSub) error
}

type CustomerSubscriptionRepository struct {
	DB *gorm.DB
}

func NewCustomerSubscriptionRepository(db *gorm.DB) *CustomerSubscriptionRepository {
	return &CustomerSubscriptionRepository{DB: db}
}

func (r *CustomerSubscriptionRepository) CreateSubscription(subscription *models.CustomerToTraderSub) error {
	return r.DB.Create(subscription).Error
}
func (r *CustomerSubscriptionRepository) GetActiveTraderSubscriptions() ([]models.CustomerToTraderSub, error) {
	var subscriptions []models.CustomerToTraderSub
	err := r.DB.Where("is_active = ?", true).Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (r *CustomerSubscriptionRepository) UpdateTraderSubscription(subscription *models.CustomerToTraderSub) error {
	return r.DB.Save(subscription).Error
}

func (r *CustomerSubscriptionRepository) GetSubscriptionByID(id uint) (*models.CustomerToTraderSub, error) {
	var subscription models.CustomerToTraderSub
	err := r.DB.Preload("SubscriptionPlan").First(&subscription, id).Error
	return &subscription, err
}

func (r *CustomerSubscriptionRepository) GetSubscriptionsByUserID(userID uint) ([]models.CustomerToTraderSub, error) {
	var subscriptions []models.CustomerToTraderSub
	err := r.DB.Where("user_id = ?", userID).Preload("SubscriptionPlan").Find(&subscriptions).Error
	return subscriptions, err
}

func (r *CustomerSubscriptionRepository) UpdateSubscription(subscription *models.CustomerToTraderSub) error {
	return r.DB.Save(subscription).Error
}
