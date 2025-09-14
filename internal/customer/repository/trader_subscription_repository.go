package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type TraderSubscriptionRepository interface {
	GetTraderSubscriptionPlanByID(id uint) (*models.TraderSubscriptionPlan, error)
	ListActiveTraderSubscriptionPlans() ([]models.TraderSubscriptionPlan, error)
	CreateTraderSubscription(subscription *models.TraderSubscription) error
	GetTraderSubscriptionByUserID(userID uint) (*models.TraderSubscription, error) // Get active trader sub by user
	CreateTraderProfile(profile *models.TraderProfile) error
	UpdateUserRole(userID uint, role models.UserRole) error
	GetRoleByName(roleName models.UserRole) (*models.Role, error) // Assuming you have a Role model/table
	// Add other methods as needed, e.g., update TraderSubscription
}

type traderSubscriptionRepository struct {
	db *gorm.DB
}

func NewTraderSubscriptionRepository(db *gorm.DB) TraderSubscriptionRepository {
	return &traderSubscriptionRepository{db: db}
}

func (r *traderSubscriptionRepository) GetTraderSubscriptionPlanByID(id uint) (*models.TraderSubscriptionPlan, error) {
	var plan models.TraderSubscriptionPlan
	err := r.db.First(&plan, id).Error
	return &plan, err
}

func (r *traderSubscriptionRepository) ListActiveTraderSubscriptionPlans() ([]models.TraderSubscriptionPlan, error) {
	var plans []models.TraderSubscriptionPlan
	err := r.db.Where("is_active = ?", true).Find(&plans).Error
	return plans, err
}

func (r *traderSubscriptionRepository) CreateTraderSubscription(subscription *models.TraderSubscription) error {
	return r.db.Create(subscription).Error
}

func (r *traderSubscriptionRepository) GetTraderSubscriptionByUserID(userID uint) (*models.TraderSubscription, error) {
	var traderSub models.TraderSubscription
	err := r.db.Preload("TraderSubscriptionPlan").Where("user_id = ? AND is_active = ?", userID, true).First(&traderSub).Error
	return &traderSub, err
}

func (r *traderSubscriptionRepository) CreateTraderProfile(profile *models.TraderProfile) error {
	return r.db.Create(profile).Error
}

func (r *traderSubscriptionRepository) UpdateUserRole(userID uint, role models.UserRole) error {
	// Find the RoleID for the given role name
	roleModel, err := r.GetRoleByName(role)
	if err != nil {
		return err
	}

	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"role":    role,
		"role_id": roleModel.ID, // Update role_id if you are using a separate roles table
	}).Error
}

// Assuming you have a models.Role struct and a roles table
func (r *traderSubscriptionRepository) GetRoleByName(roleName models.UserRole) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", roleName).First(&role).Error
	return &role, err
}

// Add a basic Role model if you don't have one, or adjust GetRoleByName
// For simplicity, if you only rely on the string 'Role' field in User, you can remove RoleID related logic
// and simplify UpdateUserRole:
// return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
