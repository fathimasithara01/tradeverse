package repository

// import (
// 	"github.com/fathimasithara01/tradeverse/pkg/models"
// 	"gorm.io/gorm"
// )

// type SubscriptionRepository interface {
// 	ListTraderSubscribers(traderID uint) ([]models.TraderSubscription, error)
// 	GetTraderSubscriberDetails(traderID, subscriptionID uint) (*models.TraderSubscription, error)
// }

// type subscriptionRepository struct {
// 	db *gorm.DB
// }

// func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
// 	return &subscriptionRepository{db: db}
// }

// func (r *subscriptionRepository) ListTraderSubscribers(traderID uint) ([]models.TraderSubscription, error) {
// 	var subscriptions []models.TraderSubscription
// 	err := r.db.
// 		Preload("User").
// 		Preload("TraderSubscriptionPlan").
// 		Joins("JOIN trader_profiles ON trader_profiles.user_id = trader_subscriptions.trader_subscription_plan_id"). // Assuming plan is tied to trader's profile
// 		Where("trader_profiles.user_id = ?", traderID).
// 		Where("trader_subscriptions.is_active = ?", true).
// 		Find(&subscriptions).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return subscriptions, nil
// }

// func (r *subscriptionRepository) GetTraderSubscriberDetails(traderID, subscriptionID uint) (*models.TraderSubscription, error) {
// 	var subscription models.TraderSubscription
// 	err := r.db.
// 		Preload("User").
// 		Preload("TraderSubscriptionPlan").
// 		Joins("JOIN trader_profiles ON trader_profiles.user_id = trader_subscriptions.trader_subscription_plan_id"). // Assuming plan is tied to trader's profile
// 		Where("trader_profiles.user_id = ?", traderID).
// 		Where("trader_subscriptions.id = ?", subscriptionID).
// 		First(&subscription).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &subscription, nil
// }
