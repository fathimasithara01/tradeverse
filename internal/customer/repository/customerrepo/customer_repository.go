package customerrepo // Changed package name

import (
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type CustomerRepository interface {
	GetTraderSubscriptionPlans() ([]models.SubscriptionPlan, error)
	GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)
	CreateTraderSubscription(sub *models.TraderSubscription) error
	GetUserTraderSubscription(userID uint) (*models.TraderSubscription, error)
	CancelTraderSubscription(userID uint, subscriptionID uint) error

	GetUserByID(userID uint) (*models.User, error)
	UpdateUserRole(userID uint, role models.UserRole) error
	CreateTraderProfile(profile *models.TraderProfile) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

func (r *customerRepository) GetTraderSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	if err := r.db.Where("is_trader_plan = ? AND is_active = ?", true, true).
		Order("price ASC").
		Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *customerRepository) GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	if err := r.db.First(&plan, id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *customerRepository) CreateTraderSubscription(sub *models.TraderSubscription) error {
	return r.db.Create(sub).Error
}

func (r *customerRepository) GetUserTraderSubscription(userID uint) (*models.TraderSubscription, error) {
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

func (r *customerRepository) CancelTraderSubscription(userID uint, subscriptionID uint) error {
	return r.db.Model(&models.TraderSubscription{}).
		Where("id = ? AND user_id = ?", subscriptionID, userID).
		Updates(map[string]interface{}{"is_active": false, "end_date": time.Now()}).Error
}

func (r *customerRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *customerRepository) UpdateUserRole(userID uint, role models.UserRole) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

func (r *customerRepository) CreateTraderProfile(profile *models.TraderProfile) error {
	return r.db.Create(profile).Error
}
