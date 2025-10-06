package customerrepo

import (
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrSubscriptionNotFound = errors.New("trader subscription not found")
)

type ITraderSubscriptionRepository interface {
	CreateTraderSubscription(tx *gorm.DB, subscription *models.TraderSubscription) error
	GetActiveTraderSubscription(customerID, traderID uint) (*models.TraderSubscription, error)
	GetTraderSubscriptionByID(subscriptionID uint) (*models.TraderSubscription, error)
	GetTraderSubscriptionPlan(planID uint) (*models.SubscriptionPlan, error)
	GetTraderByID(traderID uint) (*models.User, error)
	GetUserWallet(userID uint) (*models.Wallet, error)
	UpdateTraderSubscription(tx *gorm.DB, subscription *models.TraderSubscription) error
	// Helper to debit/credit wallet within the repository
	CreditWallet(tx *gorm.DB, walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error
	DebitWallet(tx *gorm.DB, walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error
	CreateWalletTransaction(tx *gorm.DB, walletTx *models.WalletTransaction) error
	GetAdminUser() (*models.User, error) // To find the admin for commission
}

type traderSubscriptionRepository struct {
	db *gorm.DB
}

func NewTraderSubscriptionRepository(db *gorm.DB) ITraderSubscriptionRepository {
	return &traderSubscriptionRepository{db: db}
}

func (r *traderSubscriptionRepository) CreateTraderSubscription(tx *gorm.DB, subscription *models.TraderSubscription) error {
	return tx.Create(subscription).Error
}

func (r *traderSubscriptionRepository) GetActiveTraderSubscription(customerID, traderID uint) (*models.TraderSubscription, error) {
	var sub models.TraderSubscription
	err := r.db.
		Preload("TraderSubscriptionPlan").
		Preload("Trader").
		Where("user_id = ? AND trader_id = ? AND is_active = ? AND end_date > ?", customerID, traderID, true, time.Now()).
		First(&sub).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}
	return &sub, nil
}

func (r *traderSubscriptionRepository) GetTraderSubscriptionByID(subscriptionID uint) (*models.TraderSubscription, error) {
	var sub models.TraderSubscription
	err := r.db.
		Preload("TraderSubscriptionPlan").
		Preload("User").
		Preload("Trader").
		First(&sub, subscriptionID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}
	return &sub, nil
}

func (r *traderSubscriptionRepository) GetTraderSubscriptionPlan(planID uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	err := r.db.First(&plan, planID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription plan not found")
		}
		return nil, err
	}
	if !plan.IsTraderPlan {
		return nil, errors.New("plan is not a trader subscription plan")
	}
	return &plan, nil
}

func (r *traderSubscriptionRepository) GetTraderByID(traderID uint) (*models.User, error) {
	var trader models.User
	err := r.db.Where("id = ? AND role = ?", traderID, models.RoleTrader).First(&trader).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("trader not found")
		}
		return nil, err
	}
	return &trader, nil
}

func (r *traderSubscriptionRepository) GetUserWallet(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found for user")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *traderSubscriptionRepository) UpdateTraderSubscription(tx *gorm.DB, subscription *models.TraderSubscription) error {
	return tx.Save(subscription).Error
}

func (r *traderSubscriptionRepository) CreditWallet(tx *gorm.DB, walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error {
	var wallet models.Wallet
	if err := tx.First(&wallet, walletID).Error; err != nil {
		return err
	}
	balanceBefore := wallet.Balance
	wallet.Balance += amount
	wallet.LastUpdated = time.Now()
	if err := tx.Save(&wallet).Error; err != nil {
		return err
	}
	walletTx := &models.WalletTransaction{
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
	return r.CreateWalletTransaction(tx, walletTx)
}

func (r *traderSubscriptionRepository) DebitWallet(tx *gorm.DB, walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error {
	var wallet models.Wallet
	if err := tx.First(&wallet, walletID).Error; err != nil {
		return err
	}
	if wallet.Balance < amount {
		return errors.New("insufficient funds for debit")
	}
	balanceBefore := wallet.Balance
	wallet.Balance -= amount
	wallet.LastUpdated = time.Now()
	if err := tx.Save(&wallet).Error; err != nil {
		return err
	}
	walletTx := &models.WalletTransaction{
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
	return r.CreateWalletTransaction(tx, walletTx)
}

func (r *traderSubscriptionRepository) CreateWalletTransaction(tx *gorm.DB, walletTx *models.WalletTransaction) error {
	return tx.Create(walletTx).Error
}

func (r *traderSubscriptionRepository) GetAdminUser() (*models.User, error) {
	var adminUser models.User
	err := r.db.Where("role = ?", models.RoleAdmin).First(&adminUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin user not found")
		}
		return nil, err
	}
	return &adminUser, nil
}
