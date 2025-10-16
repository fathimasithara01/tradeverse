package customerrepo

// import (
// 	"errors"
// 	"time"

// 	"github.com/fathimasithara01/tradeverse/pkg/models"
// 	"gorm.io/gorm"
// )

// var (
// 	ErrUserNotFound         = errors.New("user not found")
// 	ErrSubscriptionNotFound = errors.New("subscription not found") // Added a new error for clarity
// )

// // IAdminSubscriptionRepository defines the methods for managing admin subscriptions,
// // specifically for 'upgrade to trader' functionalities.
// type IAdminSubscriptionRepository interface {
// 	GetUserByID(userID uint) (*models.User, error)
// 	UpdateUserRole(userID uint, role models.UserRole) error
// 	CreateTraderProfile(profile *models.TraderProfile) error

// 	GetTraderSubscriptionPlans() ([]models.AdminTraderSubscriptionPlan, error) // Returns SubscriptionPlans that are IsUpgradeToTrader
// 	GetSubscriptionPlanByID(id uint) (*models.AdminTraderSubscriptionPlan, error)

// 	CreateSubscription(sub *models.CustomerToTraderSub) error
// 	UpdateSubscription(sub *models.CustomerToTraderSub) error
// 	GetUserActiveTraderSubscription(userID uint) (*models.CustomerToTraderSub, error)
// 	GetUserActiveTraderSubscriptions(userID uint) ([]models.CustomerToTraderSub, error)
// 	GetSubscriptionByIDAndUserID(subscriptionID, userID uint) (*models.CustomerToTraderSub, error)
// 	GetExpiredActiveUpgradeToTraderSubscriptions() ([]models.CustomerToTraderSub, error)
// }

// type adminSubscriptionRepository struct {
// 	db *gorm.DB
// }

// func NewIAdminSubscriptionRepository(db *gorm.DB) IAdminSubscriptionRepository {
// 	return &adminSubscriptionRepository{db: db}
// }

// func (r *adminSubscriptionRepository) GetUserByID(userID uint) (*models.User, error) {
// 	var user models.User
// 	if err := r.db.First(&user, userID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrUserNotFound
// 		}
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func (r *adminSubscriptionRepository) UpdateUserRole(userID uint, role models.UserRole) error {
// 	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
// }

// func (r *adminSubscriptionRepository) CreateTraderProfile(profile *models.TraderProfile) error {
// 	return r.db.Create(profile).Error
// }

// func (r *adminSubscriptionRepository) GetTraderSubscriptionPlans() ([]models.AdminTraderSubscriptionPlan, error) {
// 	var plans []models.AdminTraderSubscriptionPlan
// 	if err := r.db.Where("is_upgrade_to_trader = ? AND is_active = ?", true, true).
// 		Order("price ASC").
// 		Find(&plans).Error; err != nil {
// 		return nil, err
// 	}
// 	return plans, nil
// }

// func (r *adminSubscriptionRepository) GetSubscriptionPlanByID(id uint) (*models.AdminTraderSubscriptionPlan, error) {
// 	var plan models.AdminTraderSubscriptionPlan
// 	if err := r.db.First(&plan, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &plan, nil
// }

// func (r *adminSubscriptionRepository) CreateSubscription(sub *models.CustomerToTraderSub) error {
// 	return r.db.Create(sub).Error
// }

// func (r *adminSubscriptionRepository) UpdateSubscription(sub *models.CustomerToTraderSub) error {
// 	return r.db.Save(sub).Error
// }

// func (r *adminSubscriptionRepository) GetUserActiveTraderSubscription(userID uint) (*models.CustomerToTraderSub, error) {
// 	var subscription models.CustomerToTraderSub
// 	err := r.db.Preload("SubscriptionPlan").
// 		Joins("JOIN subscription_plans ON subscriptions.subscription_plan_id = subscription_plans.id").
// 		Where("subscriptions.user_id = ? AND subscriptions.is_active = ? AND subscription_plans.is_upgrade_to_trader = ? AND subscriptions.end_date > ?",
// 			userID, true, true, time.Now()).
// 		First(&subscription).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &subscription, nil
// }

// func (r *adminSubscriptionRepository) GetUserActiveTraderSubscriptions(userID uint) ([]models.CustomerToTraderSub, error) {
// 	var subscriptions []models.CustomerToTraderSub
// 	err := r.db.Preload("SubscriptionPlan").
// 		Joins("JOIN subscription_plans ON subscriptions.subscription_plan_id = subscription_plans.id").
// 		Where("subscriptions.user_id = ? AND subscriptions.is_active = ? AND subscription_plans.is_upgrade_to_trader = ? AND subscriptions.end_date > ?",
// 			userID, true, true, time.Now()).
// 		Find(&subscriptions).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return subscriptions, nil
// }

// func (r *adminSubscriptionRepository) GetSubscriptionByIDAndUserID(subscriptionID, userID uint) (*models.CustomerToTraderSub, error) {
// 	var subscription models.CustomerToTraderSub
// 	err := r.db.Preload("SubscriptionPlan").
// 		Where("id = ? AND user_id = ?", subscriptionID, userID).
// 		First(&subscription).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrSubscriptionNotFound
// 		}
// 		return nil, err
// 	}
// 	return &subscription, nil
// }

// func (r *adminSubscriptionRepository) GetExpiredActiveUpgradeToTraderSubscriptions() ([]models.CustomerToTraderSub, error) {
// 	var subscriptions []models.CustomerToTraderSub
// 	err := r.db.Preload("SubscriptionPlan").
// 		Joins("JOIN subscription_plans ON subscriptions.subscription_plan_id = subscription_plans.id").
// 		Where("subscriptions.is_active = ? AND subscriptions.end_date < ? AND subscription_plans.is_upgrade_to_trader = ?",
// 			true, time.Now(), true).
// 		Find(&subscriptions).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return subscriptions, nil
// }
