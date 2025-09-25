package customerrepo

import (
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type ITraderSubscriptionRepository interface {
	GetTraderSubscriptionPlans() ([]models.SubscriptionPlan, error)
	GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)
	CreateTraderSubscription(sub *models.TraderSubscription) error
	GetUserTraderSubscription(userID uint) (*models.TraderSubscription, error)
	CancelTraderSubscription(userID uint, subscriptionID uint) error

	GetUserByID(userID uint) (*models.User, error)
	UpdateUserRole(userID uint, role models.UserRole) error
	CreateTraderProfile(profile *models.TraderProfile) error

	GetExpiredActiveTraderSubscriptions() ([]models.TraderSubscription, error)
	UpdateTraderSubscription(sub *models.TraderSubscription) error
}



type traderSubscriptionRepository struct {
	db *gorm.DB
}

func NewTraderSubscriptionRepository(db *gorm.DB) ITraderSubscriptionRepository { // Changed return type
	return &traderSubscriptionRepository{db: db}
}

func (r *traderSubscriptionRepository) GetExpiredActiveTraderSubscriptions() ([]models.TraderSubscription, error) {
	var subs []models.TraderSubscription
	err := r.db.
		Where("is_active = ? AND end_date < ?", true, time.Now()).
		Preload("TraderSubscriptionPlan"). // Preload plan details for logging
		Find(&subs).Error
	return subs, err
}

func (r *traderSubscriptionRepository) UpdateTraderSubscription(sub *models.TraderSubscription) error {
	return r.db.Save(sub).Error
}
func (r *traderSubscriptionRepository) GetTraderSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	if err := r.db.Where("is_trader_plan = ? AND is_active = ?", true, true).
		Order("price ASC").
		Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *traderSubscriptionRepository) GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	if err := r.db.First(&plan, id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *traderSubscriptionRepository) CreateTraderSubscription(sub *models.TraderSubscription) error {
	return r.db.Create(sub).Error
}

func (r *traderSubscriptionRepository) GetUserTraderSubscription(userID uint) (*models.TraderSubscription, error) {
	var sub models.TraderSubscription
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

func (r *traderSubscriptionRepository) CancelTraderSubscription(userID uint, subscriptionID uint) error {
	return r.db.Model(&models.TraderSubscription{}).
		Where("id = ? AND user_id = ?", subscriptionID, userID).
		Updates(map[string]interface{}{"is_active": false, "end_date": time.Now()}).Error
}

func (r *traderSubscriptionRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *traderSubscriptionRepository) UpdateUserRole(userID uint, role models.UserRole) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

func (r *traderSubscriptionRepository) CreateTraderProfile(profile *models.TraderProfile) error {
	return r.db.Create(profile).Error
}
