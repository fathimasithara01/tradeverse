package walletrepo

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type TraderWalletRepository interface {
	GetByUserID(userID uint) (*models.Wallet, error)
	UpdateBalance(walletID uint, amount float64) error
	AddTransaction(tx *models.WalletTransaction) error
}
	
type traderwalletRepository struct {
	db *gorm.DB
}

func NewTraderWalletRepository(db *gorm.DB) TraderWalletRepository {
	return &traderwalletRepository{db}
}

func (r *traderwalletRepository) GetByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *traderwalletRepository) UpdateBalance(walletID uint, amount float64) error {
	return r.db.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *traderwalletRepository) AddTransaction(txn *models.WalletTransaction) error {
	return r.db.Create(txn).Error
}
