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

type IAdminSubscriptionRepository interface {
	GetTraderSubscriptionPlans() ([]models.SubscriptionPlan, error)
	GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)
	CreateTraderSubscription(sub *models.TraderSubscriptionPlan) error
	GetUserTraderSubscription(userID uint) (*models.TraderSubscriptionPlan, error)
	CancelTraderSubscription(userID uint, subscriptionID uint) error

	GetUserByID(userID uint) (*models.User, error)
	UpdateUserRole(userID uint, role models.UserRole) error
	CreateTraderProfile(profile *models.TraderProfile) error

	GetExpiredActiveTraderSubscriptions() ([]models.TraderSubscriptionPlan, error)
	UpdateTraderSubscription(sub *models.TraderSubscriptionPlan) error
}

type adminSubscriptionRepository struct {
	db *gorm.DB
}

func NewIAdminSubscriptionRepository(db *gorm.DB) IAdminSubscriptionRepository { // Changed return type
	return &adminSubscriptionRepository{db: db}
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
func (r *adminSubscriptionRepository) GetTraderSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	if err := r.db.Where("is_trader_plan = ? AND is_active = ?", true, true).
		Order("price ASC").
		Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *adminSubscriptionRepository) GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	if err := r.db.First(&plan, id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

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

func (r *adminSubscriptionRepository) UpdateUserRole(userID uint, role models.UserRole) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

func (r *adminSubscriptionRepository) CreateTraderProfile(profile *models.TraderProfile) error {
	return r.db.Create(profile).Error
}
