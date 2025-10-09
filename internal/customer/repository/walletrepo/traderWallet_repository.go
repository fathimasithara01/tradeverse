package walletrepo

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type WalletTraderRepository interface {
	Create(sub *models.TraderSubscription) error
	GetPlanByID(planID uint) (*models.SubscriptionPlan, error)
	GetActiveByCustomerID(customerID uint) ([]models.TraderSubscription, error)
	GetWalletByUserID(userID uint) (*models.Wallet, error)
	UpdateWallet(wallet *models.Wallet) error
	CreateTransaction(tx *models.WalletTransaction) error
	GetUserWallet(userID uint) (*models.Wallet, error)
	UpdateWalletTx(tx *gorm.DB, wallet *models.Wallet) error // ✅ transaction-safe

}
type traderWalletRepo struct {
	DB *gorm.DB
}

func NewTraderWalletRepository(db *gorm.DB) WalletTraderRepository {
	return &traderWalletRepo{DB: db}
}

func (r *traderWalletRepo) UpdateWallet(wallet *models.Wallet) error {
	return r.DB.Save(wallet).Error
}

// ✅ transaction-safe update
func (r *traderWalletRepo) UpdateWalletTx(tx *gorm.DB, wallet *models.Wallet) error {
	return tx.Save(wallet).Error
}
func (r *traderWalletRepo) Create(sub *models.TraderSubscription) error {
	return r.DB.Create(sub).Error
}

func (r *traderWalletRepo) GetActiveByCustomerID(customerID uint) ([]models.TraderSubscription, error) {
	var subs []models.TraderSubscription
	err := r.DB.Preload("Trader").Preload("TraderSubscriptionPlan").
		Where("user_id = ? AND is_active = ?", customerID, true).Find(&subs).Error
	return subs, err
}

func (r *traderWalletRepo) GetPlanByID(planID uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	if err := r.DB.First(&plan, planID).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *traderWalletRepo) GetWalletByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *traderWalletRepo) CreateTransaction(tx *models.WalletTransaction) error {
	return r.DB.Create(tx).Error
}

// Implement GetUserWallet to satisfy the interface
func (r *traderWalletRepo) GetUserWallet(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}
