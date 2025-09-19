// Example of how WalletRepository and TransactionRepository might look:
// internal/trader/repository/wallet_repository.go
package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetWalletByUserID(userID uint) (*models.Wallet, error)
	UpdateWallet(wallet *models.Wallet) error
	// ... other wallet methods
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetWalletByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &wallet, err
}

func (r *walletRepository) UpdateWallet(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}
