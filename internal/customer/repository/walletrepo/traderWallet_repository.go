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
}
type traderWalletRepo struct {
	DB *gorm.DB
}

func NewTraderWalletRepository(db *gorm.DB) WalletTraderRepository {
	return &traderWalletRepo{DB: db}
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

func (r *traderWalletRepo) UpdateWallet(wallet *models.Wallet) error {
	return r.DB.Save(wallet).Error
}

func (r *traderWalletRepo) CreateTransaction(tx *models.WalletTransaction) error {
	return r.DB.Create(tx).Error
}
