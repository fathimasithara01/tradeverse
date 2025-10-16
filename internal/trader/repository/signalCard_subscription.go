package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ITraderSubscriptionRepository interface {
	CreateTraderSubscriptionPlan(ctx context.Context, plan *models.TraderSignalSubscriptionPlan) (*models.TraderSignalSubscriptionPlan, error)
	GetTraderSubscriptionPlanByID(ctx context.Context, planID uint) (*models.TraderSignalSubscriptionPlan, error)
	GetTraderSubscriptionPlansByTraderID(ctx context.Context, traderID uint) ([]models.TraderSignalSubscriptionPlan, error)
	UpdateTraderSubscriptionPlan(ctx context.Context, plan *models.TraderSignalSubscriptionPlan) error
	DeleteTraderSubscriptionPlan(ctx context.Context, planID, traderID uint) error

	CheckIfUserIsActiveTrader(ctx context.Context, userID uint) (bool, error)
	GetUserActiveUpgradeSubscription(ctx context.Context, userID uint, planID uint) (*models.UserSubscription, error) // Specific check

	CreateCustomerTraderSubscription(ctx context.Context, sub *models.CustomerTraderSignalSubscription) error
	CheckIfCustomerIsSubscribedToTraderPlan(ctx context.Context, customerID uint, traderPlanID uint) (bool, error)

	SetUserRole(ctx context.Context, userID uint, role models.UserRole, tx *gorm.DB) error // Update user role during transaction
	GetAllTraderUpgradePlans(ctx context.Context) ([]models.AdminTraderSubscriptionPlan, error)
	GetTraderUpgradePlanByID(ctx context.Context, planID uint) (*models.AdminTraderSubscriptionPlan, error)
	CreateUserSubscription(ctx context.Context, sub *models.UserSubscription) error

	// Wallet and Transaction Management (shared)
	GetUserWallet(ctx context.Context, userID uint) (*models.Wallet, error)
	UpdateWalletBalance(ctx context.Context, userID uint, amount float64, tx *gorm.DB) error
	CreateWalletTransaction(ctx context.Context, transaction *models.WalletTransaction, tx *gorm.DB) error
	GetAdminWallet(ctx context.Context) (*models.Wallet, error)
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
}

type TraderSubscriptionRepository struct {
	db *gorm.DB
}

func NewTraderSubscriptionRepository(db *gorm.DB) ITraderSubscriptionRepository {
	return &TraderSubscriptionRepository{db: db}
}
func (r *TraderSubscriptionRepository) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *TraderSubscriptionRepository) CreateTraderSubscriptionPlan(ctx context.Context, plan *models.TraderSignalSubscriptionPlan) (*models.TraderSignalSubscriptionPlan, error) {

	if err := r.db.WithContext(ctx).Create(plan).Error; err != nil {
		return nil, fmt.Errorf("failed to create trader subscription plan: %w", err)
	}
	return plan, nil
}

// func (r *TraderSubscriptionRepository) GetTraderSubscriptionPlanByID(ctx context.Context, planID uint) (*models.TraderSubscriptionPlan, error) {
// 	var plan models.TraderSubscriptionPlan
// 	if err := r.db.WithContext(ctx).First(&plan, planID).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, fmt.Errorf("trader subscription plan not found")
// 		}
// 		return nil, fmt.Errorf("failed to get trader subscription plan: %w", err)
// 	}
// 	return &plan, nil
// }

func (r *TraderSubscriptionRepository) GetTraderSubscriptionPlanByID(ctx context.Context, planID uint) (*models.TraderSignalSubscriptionPlan, error) {
	var plan models.TraderSignalSubscriptionPlan
	if err := r.db.WithContext(ctx).First(&plan, planID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("trader subscription plan not found") // Or define a specific repo error, like customerrepo.ErrPlanNotFound
		}
		return nil, fmt.Errorf("failed to get trader subscription plan by ID %d: %w", planID, err)
	}
	return &plan, nil
}
func (r *TraderSubscriptionRepository) GetTraderSubscriptionPlansByTraderID(ctx context.Context, traderID uint) ([]models.TraderSignalSubscriptionPlan, error) {
	var plans []models.TraderSignalSubscriptionPlan
	err := r.db.WithContext(ctx).Where("trader_id = ?", traderID).Find(&plans).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find trader subscription plans for trader ID %d: %w", traderID, err)
	}
	return plans, nil
}

// func (r *TraderSubscriptionRepository) GetTraderSubscriptionPlansByTraderID(ctx context.Context, traderID uint) ([]models.TraderSubscriptionPlan, error) {
// 	var plans []models.TraderSubscriptionPlan
// 	if err := r.db.WithContext(ctx).Where("trader_id = ?", traderID).Find(&plans).Error; err != nil {
// 		return nil, fmt.Errorf("failed to get trader subscription plans for trader %d: %w", traderID, err)
// 	}
// 	return plans, nil
// }

func (r *TraderSubscriptionRepository) UpdateTraderSubscriptionPlan(ctx context.Context, plan *models.TraderSignalSubscriptionPlan) error {
	if err := r.db.WithContext(ctx).Save(plan).Error; err != nil {
		return fmt.Errorf("failed to update trader subscription plan: %w", err)
	}
	return nil
}

func (r *TraderSubscriptionRepository) DeleteTraderSubscriptionPlan(ctx context.Context, planID, traderID uint) error {
	// Ensure only the owner trader can delete their plan
	result := r.db.WithContext(ctx).Where("id = ? AND trader_id = ?", planID, traderID).Delete(&models.TraderSignalSubscriptionPlan{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete trader subscription plan: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("trader subscription plan not found or not owned by this trader")
	}
	return nil
}

func (r *TraderSubscriptionRepository) CheckIfUserIsActiveTrader(ctx context.Context, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.UserSubscription{}).
		Joins("JOIN subscription_plans ON user_subscriptions.subscription_plan_id = subscription_plans.id").
		Where("user_subscriptions.user_id = ? AND user_subscriptions.is_active = ? AND user_subscriptions.end_date > ? AND subscription_plans.is_upgrade_to_trader = ?",
			userID, true, time.Now(), true).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check active trader upgrade subscription status for user %d: %w", userID, err)
	}
	return count > 0, nil
}

func (r *TraderSubscriptionRepository) GetUserActiveUpgradeSubscription(ctx context.Context, userID uint, planID uint) (*models.UserSubscription, error) {
	var sub models.UserSubscription
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND subscription_plan_id = ? AND is_active = ? AND end_date > ?", userID, planID, true, time.Now()).
		First(&sub).Error
	if err != nil {
		return nil, err // Returns gorm.ErrRecordNotFound if not found
	}
	return &sub, nil
}

func (r *TraderSubscriptionRepository) SetUserRole(ctx context.Context, userID uint, role models.UserRole, tx *gorm.DB) error {
	return tx.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

// --- Customer Subscription to Trader's plan ---
func (r *TraderSubscriptionRepository) CreateCustomerTraderSubscription(ctx context.Context, sub *models.CustomerTraderSignalSubscription) error {
	if err := r.db.WithContext(ctx).Create(sub).Error; err != nil {
		return fmt.Errorf("failed to create customer trader subscription: %w", err)
	}
	return nil
}

func (r *TraderSubscriptionRepository) CheckIfCustomerIsSubscribedToTraderPlan(ctx context.Context, customerID uint, traderPlanID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.CustomerTraderSignalSubscription{}).
		Where("customer_id = ? AND trader_subscription_plan_id = ? AND is_active = ? AND end_date > ?",
			customerID, traderPlanID, true, time.Now()).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check customer subscription status for trader plan: %w", err)
	}
	return count > 0, nil
}

// --- Admin-defined Subscription Plan (for upgrading to trader) ---

func (r *TraderSubscriptionRepository) GetAllTraderUpgradePlans(ctx context.Context) ([]models.AdminTraderSubscriptionPlan, error) {
	var plans []models.AdminTraderSubscriptionPlan
	// Only fetch plans that are specifically for upgrading to a trader role
	if err := r.db.WithContext(ctx).Where("is_upgrade_to_trader = ?", true).Find(&plans).Error; err != nil {
		return nil, fmt.Errorf("failed to get admin subscription plans for trader upgrade: %w", err)
	}
	return plans, nil
}

func (r *TraderSubscriptionRepository) GetTraderUpgradePlanByID(ctx context.Context, planID uint) (*models.AdminTraderSubscriptionPlan, error) {
	var plan models.AdminTraderSubscriptionPlan
	// Ensure it's an upgrade plan
	if err := r.db.WithContext(ctx).Where("id = ? AND is_upgrade_to_trader = ?", planID, true).First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("trader upgrade subscription plan not found")
		}
		return nil, fmt.Errorf("failed to get trader upgrade subscription plan: %w", err)
	}
	return &plan, nil
}

func (r *TraderSubscriptionRepository) CreateUserSubscription(ctx context.Context, sub *models.UserSubscription) error {
	if err := r.db.WithContext(ctx).Create(sub).Error; err != nil {
		return fmt.Errorf("failed to create user subscription: %w", err)
	}
	return nil
}

// --- Shared Wallet & Transaction Functions ---

func (r *TraderSubscriptionRepository) GetUserWallet(ctx context.Context, userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, fmt.Errorf("wallet not found for user %d: %w", userID, err)
	}
	return &wallet, nil
}

func (r *TraderSubscriptionRepository) UpdateWalletBalance(ctx context.Context, userID uint, amount float64, tx *gorm.DB) error {
	return tx.WithContext(ctx).Model(&models.Wallet{}).Where("user_id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *TraderSubscriptionRepository) CreateWalletTransaction(ctx context.Context, transaction *models.WalletTransaction, tx *gorm.DB) error {
	if err := tx.WithContext(ctx).Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create wallet transaction: %w", err)
	}
	return nil
}

func (r *TraderSubscriptionRepository) GetAdminWallet(ctx context.Context) (*models.Wallet, error) {
	var adminUser models.User
	// Assuming there's an admin user with Role = RoleAdmin
	if err := r.db.WithContext(ctx).Where("role = ?", models.RoleAdmin).First(&adminUser).Error; err != nil {
		return nil, fmt.Errorf("admin user not found: %w", err)
	}
	var adminWallet models.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", adminUser.ID).First(&adminWallet).Error; err != nil {
		return nil, fmt.Errorf("admin wallet not found for admin user %d: %w", adminUser.ID, err)
	}
	return &adminWallet, nil
}
