package customerrepo

import (
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrSubscriptionNotFound = errors.New("subscription not found") // Added a new error for clarity
)

// IAdminSubscriptionRepository defines the methods for managing admin subscriptions,
// specifically for 'upgrade to trader' functionalities.
type IAdminSubscriptionRepository interface {
	GetUserByID(userID uint) (*models.User, error)
	UpdateUserRole(userID uint, role models.UserRole) error
	CreateTraderProfile(profile *models.TraderProfile) error

	GetTraderSubscriptionPlans() ([]models.SubscriptionPlan, error) // Returns SubscriptionPlans that are IsUpgradeToTrader
	GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)

	// --- Methods for models.Subscription (the actual user subscriptions) ---
	CreateSubscription(sub *models.Subscription) error // New: For creating a general subscription record
	UpdateSubscription(sub *models.Subscription) error // New: For updating a general subscription record

	// New: Gets a single active 'upgrade to trader' subscription for a user
	GetUserActiveTraderSubscription(userID uint) (*models.Subscription, error)
	// New: Gets all active 'upgrade to trader' subscriptions for a user
	GetUserActiveTraderSubscriptions(userID uint) ([]models.Subscription, error)
	// New: Gets a specific subscription by ID and UserID
	GetSubscriptionByIDAndUserID(subscriptionID, userID uint) (*models.Subscription, error)
	// New: Gets expired active 'upgrade to trader' subscriptions for cron job
	GetExpiredActiveUpgradeToTraderSubscriptions() ([]models.Subscription, error)

	// --- Old methods that need to be removed or adapted ---
	// CreateTraderSubscription(sub *models.TraderSubscriptionPlan) error // No longer used for user subscriptions
	// GetUserTraderSubscription(userID uint) (*models.TraderSubscriptionPlan, error) // No longer used for user subscriptions
	// CancelTraderSubscription(userID uint, subscriptionID uint) error // No longer used, handled by new GetSubscriptionByIDAndUserID + UpdateSubscription
	// GetExpiredActiveTraderSubscriptions() ([]models.TraderSubscriptionPlan, error) // No longer used, replaced by GetExpiredActiveUpgradeToTraderSubscriptions
	// UpdateTraderSubscription(sub *models.TraderSubscriptionPlan) error // No longer used, replaced by UpdateSubscription
}

type adminSubscriptionRepository struct {
	db *gorm.DB
}

func NewIAdminSubscriptionRepository(db *gorm.DB) IAdminSubscriptionRepository {
	return &adminSubscriptionRepository{db: db}
}

// GetUserByID fetches a user by their ID.
func (r *adminSubscriptionRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserRole updates a user's role.
func (r *adminSubscriptionRepository) UpdateUserRole(userID uint, role models.UserRole) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

// CreateTraderProfile creates a new trader profile.
func (r *adminSubscriptionRepository) CreateTraderProfile(profile *models.TraderProfile) error {
	return r.db.Create(profile).Error
}

// GetTraderSubscriptionPlans fetches all active subscription plans marked as 'IsUpgradeToTrader'.
func (r *adminSubscriptionRepository) GetTraderSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	if err := r.db.Where("is_upgrade_to_trader = ? AND is_active = ?", true, true).
		Order("price ASC").
		Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

// GetSubscriptionPlanByID fetches a subscription plan by its ID.
func (r *adminSubscriptionRepository) GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	if err := r.db.First(&plan, id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

// CreateSubscription creates a new subscription record.
func (r *adminSubscriptionRepository) CreateSubscription(sub *models.Subscription) error {
	return r.db.Create(sub).Error
}

// UpdateSubscription updates an existing subscription record.
func (r *adminSubscriptionRepository) UpdateSubscription(sub *models.Subscription) error {
	return r.db.Save(sub).Error
}

// GetUserActiveTraderSubscription fetches a single active subscription for a user
// that is specifically for upgrading to a trader role.
func (r *adminSubscriptionRepository) GetUserActiveTraderSubscription(userID uint) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Preload("SubscriptionPlan").
		Joins("JOIN subscription_plans ON subscriptions.subscription_plan_id = subscription_plans.id").
		Where("subscriptions.user_id = ? AND subscriptions.is_active = ? AND subscription_plans.is_upgrade_to_trader = ? AND subscriptions.end_date > ?",
			userID, true, true, time.Now()).
		First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if not found, not an error
		}
		return nil, err
	}
	return &subscription, nil
}

// GetUserActiveTraderSubscriptions fetches all active subscriptions for a user
// that are specifically for upgrading to a trader role.
func (r *adminSubscriptionRepository) GetUserActiveTraderSubscriptions(userID uint) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Preload("SubscriptionPlan").
		Joins("JOIN subscription_plans ON subscriptions.subscription_plan_id = subscription_plans.id").
		Where("subscriptions.user_id = ? AND subscriptions.is_active = ? AND subscription_plans.is_upgrade_to_trader = ? AND subscriptions.end_date > ?",
			userID, true, true, time.Now()).
		Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}

// GetSubscriptionByIDAndUserID fetches a specific subscription by its ID and the UserID.
func (r *adminSubscriptionRepository) GetSubscriptionByIDAndUserID(subscriptionID, userID uint) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Preload("SubscriptionPlan").
		Where("id = ? AND user_id = ?", subscriptionID, userID).
		First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}
	return &subscription, nil
}

// GetExpiredActiveUpgradeToTraderSubscriptions fetches all active 'upgrade to trader' subscriptions
// whose end date is in the past.
func (r *adminSubscriptionRepository) GetExpiredActiveUpgradeToTraderSubscriptions() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Preload("SubscriptionPlan").
		Joins("JOIN subscription_plans ON subscriptions.subscription_plan_id = subscription_plans.id").
		Where("subscriptions.is_active = ? AND subscriptions.end_date < ? AND subscription_plans.is_upgrade_to_trader = ?",
			true, time.Now(), true).
		Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}

/*
   --- Removed/Adapted Old Methods ---
   These methods are commented out because they were for models.TraderSubscriptionPlan,
   which is no longer used for user subscriptions. The new models.Subscription
   and its associated methods handle this functionality.

func (r *adminSubscriptionRepository) CreateTraderSubscription(sub *models.TraderSubscriptionPlan) error {
	return r.db.Create(sub).Error
}

func (r *adminSubscriptionRepository) GetUserTraderSubscription(userID uint) (*models.TraderSubscriptionPlan, error) {
	var sub models.TraderSubscriptionPlan
	err := r.db.
		Where("user_id = ? AND is_active = ? AND end_date > ?", userID, true, time.Now()).
		Preload("TraderSubscriptionPlan").
		First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (r *adminSubscriptionRepository) CancelTraderSubscription(userID uint, subscriptionID uint) error {
	return r.db.Model(&models.TraderSubscriptionPlan{}).
		Where("id = ? AND user_id = ?", subscriptionID, userID).
		Updates(map[string]interface{}{"is_active": false, "end_date": time.Now()}).Error
}

func (r *adminSubscriptionRepository) GetExpiredActiveTraderSubscriptions() ([]models.TraderSubscriptionPlan, error) {
	var subs []models.TraderSubscriptionPlan
	err := r.db.
		Where("is_active = ? AND end_date < ?", true, time.Now()).
		Preload("TraderSubscriptionPlan"). // Preload plan details for logging
		Find(&subs).Error
	return subs, err
}

func (r *adminSubscriptionRepository) UpdateTraderSubscription(sub *models.TraderSubscriptionPlan) error {
	return r.db.Save(sub).Error
}
*/
