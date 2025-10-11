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
	// UpdateWalletBalance takes a GORM transaction object for atomic updates
	UpdateWalletBalance(tx *gorm.DB, adminWallet *models.Wallet) error
	// You might also have methods for admin-specific transactions
	CreateAdminWallet(adminWallet *models.Wallet) error // To initialize if not exists
}

type adminWalletRepository struct {
	db *gorm.DB
}

func NewAdminWalletRepository(db *gorm.DB) IAdminWalletRepository {
	return &adminWalletRepository{db: db}
}

func (r *adminWalletRepository) GetAdminWallet() (*models.Wallet, error) {
	var adminWallet models.Wallet
	// Assuming there's only one admin wallet or a specific way to identify it (e.g., ID 1)
	err := r.db.First(&adminWallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAdminWalletNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin wallet: %w", err)
	}
	return &adminWallet, nil
}

// UpdateWalletBalance updates the admin wallet within a transaction
func (r *adminWalletRepository) UpdateWalletBalance(tx *gorm.DB, adminWallet *models.Wallet) error {
	// Use the provided transaction object for saving
	return tx.Save(adminWallet).Error
}

func (r *adminWalletRepository) CreateAdminWallet(adminWallet *models.Wallet) error {
	return r.db.Create(adminWallet).Error
}
