package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type SubscriptionRepository struct{}

func (r *SubscriptionRepository) GetAll() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := db.DB.Find(&subscriptions).Error
	return subscriptions, err
}
