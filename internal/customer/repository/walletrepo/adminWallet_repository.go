package walletrepo

import (
	"errors"
	"fmt"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrAdminWalletNotFound = errors.New("admin wallet not found")
)

type IAdminWalletRepository interface {
	GetAdminWallet() (*models.Wallet, error)
	UpdateWalletBalance(tx *gorm.DB, adminWallet *models.Wallet) error
	CreateAdminWallet(adminWallet *models.Wallet) error 
}

type adminWalletRepository struct {
	db *gorm.DB
}

func NewAdminWalletRepository(db *gorm.DB) IAdminWalletRepository {
	return &adminWalletRepository{db: db}
}

func (r *adminWalletRepository) GetAdminWallet() (*models.Wallet, error) {
	var adminWallet models.Wallet
	err := r.db.First(&adminWallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAdminWalletNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin wallet: %w", err)
	}
	return &adminWallet, nil
}

func (r *adminWalletRepository) UpdateWalletBalance(tx *gorm.DB, adminWallet *models.Wallet) error {
	return tx.Save(adminWallet).Error
}

func (r *adminWalletRepository) CreateAdminWallet(adminWallet *models.Wallet) error {
	return r.db.Create(adminWallet).Error
}
