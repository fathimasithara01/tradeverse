package repository

import (
	"context"
	"errors"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetWalletByUserID(ctx context.Context, userID uint) (*models.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *models.Wallet) error
	CreateTransaction(ctx context.Context, tx *models.WalletTransaction) error
	GetTransactionsByWalletID(ctx context.Context, walletID uint) ([]models.WalletTransaction, error)
}

type gormWalletRepository struct {
	db *gorm.DB
}

func NewGormWalletRepository(db *gorm.DB) WalletRepository {
	return &gormWalletRepository{db: db}
}

func (r *gormWalletRepository) GetWalletByUserID(ctx context.Context, userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *gormWalletRepository) UpdateWallet(ctx context.Context, wallet *models.Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}

func (r *gormWalletRepository) CreateTransaction(ctx context.Context, tx *models.WalletTransaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *gormWalletRepository) GetTransactionsByWalletID(ctx context.Context, walletID uint) ([]models.WalletTransaction, error) {
	var txs []models.WalletTransaction
	err := r.db.WithContext(ctx).Where("wallet_id = ?", walletID).Order("created_at desc").Find(&txs).Error
	return txs, err
}
