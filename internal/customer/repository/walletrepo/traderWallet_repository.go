package walletrepo

import (
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type CustomerWalletRepository interface {
	GetWalletByUserID(userID uint) (*models.Wallet, error)
	UpdateWalletBalance(userID uint, amount float64) error
	CreateTransaction(tx *models.WalletTransaction) error
	CreateDepositRequest(dr *models.DepositRequest) error
	UpdateDepositRequest(dr *models.DepositRequest) error
	CreateWithdrawRequest(wr *models.WithdrawRequest) error
}

type customerWalletRepository struct {
	db *gorm.DB
}

func NewCustomerWalletRepository(db *gorm.DB) CustomerWalletRepository {
	return &customerWalletRepository{db: db}
}

func (r *customerWalletRepository) GetWalletByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *customerWalletRepository) UpdateWalletBalance(userID uint, amount float64) error {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return err
	}
	wallet.Balance += amount
	wallet.LastUpdated = time.Now()
	return r.db.Save(&wallet).Error
}

func (r *customerWalletRepository) CreateTransaction(tx *models.WalletTransaction) error {
	return r.db.Create(tx).Error
}

func (r *customerWalletRepository) CreateDepositRequest(dr *models.DepositRequest) error {
	return r.db.Create(dr).Error
}

func (r *customerWalletRepository) UpdateDepositRequest(dr *models.DepositRequest) error {
	return r.db.Save(dr).Error
}

func (r *customerWalletRepository) CreateWithdrawRequest(wr *models.WithdrawRequest) error {
	return r.db.Create(wr).Error
}
