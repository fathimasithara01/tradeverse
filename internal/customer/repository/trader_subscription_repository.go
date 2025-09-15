package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
	clause "gorm.io/gorm/clause"
)

type modelsPackage struct {
	SubscriptionPlan   models.SubscriptionPlan
	TraderSubscription models.TraderSubscription
	User               models.User
	Wallet             models.Wallet
	WalletTransaction  models.WalletTransaction
	TransactionType    models.TransactionType
	TransactionStatus  models.TransactionStatus
	UserRole           models.UserRole
}

type CustomerRepository interface {
	GetTraderSubscriptionPlans() ([]models.SubscriptionPlan, error)
	GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)
	CreateTraderSubscription(sub *models.TraderSubscription) error
	GetUserTraderSubscription(userID uint) (*models.TraderSubscription, error)
	CancelTraderSubscription(userID uint, subscriptionID uint) error
	GetUserByID(userID uint) (*models.User, error)
	UpdateUserRole(userID uint, role models.UserRole) error
	GetAdminWallet() (*models.Wallet, error)
	CreditWallet(walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error
	DebitWallet(walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error
	CreateWalletTransaction(tx *models.WalletTransaction) error
	// For customer's general subscriptions (if they exist)
	GetActiveSubscriptionByUserID(userID uint) (*models.Subscription, error)
}

type customerRepository struct {
	db *gorm.DB
	// Assuming models are defined in pkg/models based on your folder structure
	// You might need to adjust this if your models are nested differently
	models *modelsPackage
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{
		db: db,
		models: &modelsPackage{
			SubscriptionPlan:   models.SubscriptionPlan{},
			TraderSubscription: models.TraderSubscription{},
			User:               models.User{},
			Wallet:             models.Wallet{},
			WalletTransaction:  models.WalletTransaction{},
			TransactionType:    "",
			TransactionStatus:  "",
			UserRole:           "",
		},
	}
}

func (r *customerRepository) GetTraderSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	// Fetch only plans designated for traders
	if err := r.db.Where("is_trader_plan = ?", true).Find(&plans).Error; err != nil {
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
	// Fetch active trader subscriptions for the user
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Preload("TraderSubscriptionPlan").First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active subscription found
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
		return nil, err
	}
	return &user, nil
}

func (r *customerRepository) UpdateUserRole(userID uint, role models.UserRole) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

// GetAdminWallet fetches the wallet of the admin user.
// In a real application, AdminUserID should be configured or dynamically found.
// For this example, I'll assume an admin user exists with ID 1 and their role is 'admin'.
const AdminUserID uint = 1 // THIS IS A PLACEHOLDER. ADJUST AS PER YOUR ADMIN USER ID.

func (r *customerRepository) GetAdminWallet() (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", AdminUserID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("admin wallet not found for user ID %d", AdminUserID)
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *customerRepository) CreditWallet(walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		// Ensure 'clause' is accessible here
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, walletID).Error; err != nil {
			return err
		}

		balanceBefore := wallet.Balance
		wallet.Balance += amount
		wallet.LastUpdated = time.Now()
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		walletTx := models.WalletTransaction{
			WalletID:        walletID,
			UserID:          wallet.UserID,
			TransactionType: transactionType,
			Amount:          amount,
			Currency:        wallet.Currency,
			Status:          models.TxStatusSuccess,
			ReferenceID:     referenceID,
			Description:     description,
			BalanceBefore:   balanceBefore,
			BalanceAfter:    wallet.Balance,
		}
		if err := tx.Create(&walletTx).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *customerRepository) DebitWallet(walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, walletID).Error; err != nil {
			return err
		}

		if wallet.Balance < amount {
			return errors.New("insufficient funds")
		}

		balanceBefore := wallet.Balance
		wallet.Balance -= amount
		wallet.LastUpdated = time.Now()
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		walletTx := models.WalletTransaction{
			WalletID:        walletID,
			UserID:          wallet.UserID,
			TransactionType: transactionType,
			Amount:          amount,
			Currency:        wallet.Currency,
			Status:          models.TxStatusSuccess,
			ReferenceID:     referenceID,
			Description:     description,
			BalanceBefore:   balanceBefore,
			BalanceAfter:    wallet.Balance,
		}
		if err := tx.Create(&walletTx).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *customerRepository) CreateWalletTransaction(tx *models.WalletTransaction) error {
	return r.db.Create(tx).Error
}

func (r *customerRepository) GetActiveSubscriptionByUserID(userID uint) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.Where("user_id = ? AND is_active = ? AND end_date > ?", userID, true, time.Now()).Preload("SubscriptionPlan").First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}
