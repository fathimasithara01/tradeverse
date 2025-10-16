package customerrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ICustomerTraderSubscriptionRepository interface {
	GetTradersWithPlans(ctx context.Context) ([]models.User, error)
	GetTraderSubscriptionPlanByID(ctx context.Context, planID uint) (*models.TraderSignalSubscriptionPlan, error)
	CreateCustomerTraderSubscription(ctx context.Context, sub *models.CustomerTraderSignalSubscription) (*models.CustomerTraderSignalSubscription, error)
	IsCustomerSubscribedToTrader(ctx context.Context, customerID, traderID uint) (bool, error)
	GetActiveTraderSubscriptionsForCustomer(ctx context.Context, customerID uint) ([]models.CustomerTraderSignalSubscription, error)
	GetAllSignalsFromSubscribedTraders(ctx context.Context, customerID uint) ([]models.Signal, error) 
	GetTraderByID(ctx context.Context, traderID uint) (*models.User, error)
	UpdateWalletBalance(ctx context.Context, userID uint, amount float64, tx *gorm.DB) error
	CreateWalletTransaction(ctx context.Context, transaction *models.WalletTransaction, tx *gorm.DB) error
	GetAdminWallet(ctx context.Context) (*models.Wallet, error)
	GetTraderWallet(ctx context.Context, traderID uint) (*models.Wallet, error)
	IsCustomerSubscribedToPlan(ctx context.Context, customerID, planID uint) (bool, error) // New
}

type CustomerTraderSubscriptionRepository struct {
	db *gorm.DB
}

func NewCustomerTraderSubscriptionRepository(db *gorm.DB) ICustomerTraderSubscriptionRepository {
	return &CustomerTraderSubscriptionRepository{db: db}
}

func (r *CustomerTraderSubscriptionRepository) GetTraderByID(ctx context.Context, traderID uint) (*models.User, error) {
	var trader models.User
	// --- FIX HERE ---
	// Changed "is_trader = ?" to "role = ?" and true to models.RoleTrader
	if err := r.db.WithContext(ctx).Where("id = ? AND role = ?", traderID, models.RoleTrader).First(&trader).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("trader not found")
		}
		return nil, fmt.Errorf("failed to get trader: %w", err)
	}
	return &trader, nil
}

func (r *CustomerTraderSubscriptionRepository) GetTradersWithPlans(ctx context.Context) ([]models.User, error) {
	var traders []models.User
	// Fetch users who are marked as traders and have at least one active TraderSubscriptionPlan
	err := r.db.WithContext(ctx).
		Preload("TraderSubscriptionPlans", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true)
		}).
		// --- FIX HERE ---
		// Changed "is_trader = ?" to "role = ?" and true to models.RoleTrader
		Where("role = ?", models.RoleTrader).
		Find(&traders).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get traders with plans: %w", err)
	}

	// Filter out traders who genuinely have no active plans after preloading
	var activeTraders []models.User
	for _, trader := range traders {
		if len(trader.TraderSubscriptionPlans) > 0 {
			activeTraders = append(activeTraders, trader)
		}
	}

	return activeTraders, nil
}

func (r *CustomerTraderSubscriptionRepository) GetTraderSubscriptionPlanByID(ctx context.Context, planID uint) (*models.TraderSignalSubscriptionPlan, error) {
	var plan models.TraderSignalSubscriptionPlan
	if err := r.db.WithContext(ctx).First(&plan, planID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("trader subscription plan not found")
		}
		return nil, fmt.Errorf("failed to get trader subscription plan: %w", err)
	}
	return &plan, nil
}

func (r *CustomerTraderSubscriptionRepository) CreateCustomerTraderSubscription(ctx context.Context, sub *models.CustomerTraderSignalSubscription) (*models.CustomerTraderSignalSubscription, error) {
	if err := r.db.WithContext(ctx).Create(sub).Error; err != nil {
		return nil, fmt.Errorf("failed to create customer-trader subscription: %w", err)
	}
	return sub, nil
}

func (r *CustomerTraderSubscriptionRepository) IsCustomerSubscribedToTrader(ctx context.Context, customerID, traderID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.CustomerTraderSignalSubscription{}).
		Where("customer_id = ? AND trader_id = ? AND end_date > ? AND is_active = ?", customerID, traderID, time.Now(), true).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check customer subscription status: %w", err)
	}
	return count > 0, nil
}

func (r *CustomerTraderSubscriptionRepository) GetActiveTraderSubscriptionsForCustomer(ctx context.Context, customerID uint) ([]models.CustomerTraderSignalSubscription, error) {
	var subscriptions []models.CustomerTraderSignalSubscription
	err := r.db.WithContext(ctx).
		Preload("Trader").Preload("Plan").
		Where("customer_id = ? AND end_date > ? AND is_active = ?", customerID, time.Now(), true).
		Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get active trader subscriptions for customer: %w", err)
	}
	return subscriptions, nil
}

func (r *CustomerTraderSubscriptionRepository) GetAllSignalsFromSubscribedTraders(ctx context.Context, customerID uint) ([]models.Signal, error) {
	var signals []models.Signal

	// First, get all trader IDs the customer is subscribed to
	var subscribedTraderIDs []uint
	err := r.db.WithContext(ctx).
		Model(&models.CustomerTraderSignalSubscription{}).
		Select("DISTINCT trader_id").
		Where("customer_id = ? AND end_date > ? AND is_active = ?", customerID, time.Now(), true).
		Pluck("trader_id", &subscribedTraderIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get subscribed trader IDs: %w", err)
	}

	if len(subscribedTraderIDs) == 0 {
		return []models.Signal{}, nil // No subscriptions, no signals
	}

	// Then, fetch all signals from those traders
	err = r.db.WithContext(ctx).
		Where("trader_id IN (?)", subscribedTraderIDs).
		Order("published_at DESC"). // Order by publication date, newest first
		Find(&signals).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get signals from subscribed traders: %w", err)
	}
	return signals, nil
}

func (r *CustomerTraderSubscriptionRepository) UpdateWalletBalance(ctx context.Context, userID uint, amount float64, tx *gorm.DB) error {
	return tx.WithContext(ctx).Model(&models.Wallet{}).Where("user_id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *CustomerTraderSubscriptionRepository) CreateWalletTransaction(ctx context.Context, transaction *models.WalletTransaction, tx *gorm.DB) error {
	return tx.WithContext(ctx).Create(transaction).Error
}

func (r *CustomerTraderSubscriptionRepository) GetAdminWallet(ctx context.Context) (*models.Wallet, error) {
	var adminUser models.User
	// Assuming there's an admin user with RoleAdmin
	// --- FIX HERE ---
	// Changed "is_admin = ?" to "role = ?" and true to models.RoleAdmin
	if err := r.db.WithContext(ctx).Where("role = ?", models.RoleAdmin).First(&adminUser).Error; err != nil {
		return nil, fmt.Errorf("admin user not found: %w", err)
	}
	var adminWallet models.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", adminUser.ID).First(&adminWallet).Error; err != nil {
		return nil, fmt.Errorf("admin wallet not found: %w", err)
	}
	return &adminWallet, nil
}

func (r *CustomerTraderSubscriptionRepository) GetTraderWallet(ctx context.Context, traderID uint) (*models.Wallet, error) {
	var traderWallet models.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", traderID).First(&traderWallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("trader wallet not found for trader ID %d", traderID)
		}
		return nil, fmt.Errorf("failed to get trader wallet: %w", err)
	}
	return &traderWallet, nil
}

func (r *CustomerTraderSubscriptionRepository) IsCustomerSubscribedToPlan(ctx context.Context, customerID, planID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.CustomerTraderSignalSubscription{}).
		Where("customer_id = ? AND trader_subscription_plan_id = ? AND end_date > ? AND is_active = ?", customerID, planID, time.Now(), true).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check customer subscription to plan status: %w", err)
	}
	return count > 0, nil
}
