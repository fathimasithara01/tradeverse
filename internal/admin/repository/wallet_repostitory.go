package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type IAdminWalletRepository interface {
	GetAdminWallet() (*models.Wallet, error)
	CreateWalletTransaction(tx *gorm.DB, transaction *models.WalletTransaction) error
	UpdateWalletBalance(tx *gorm.DB, wallet *models.Wallet) error
	CreateDepositRequest(deposit *models.DepositRequest) error
	GetDepositRequestByID(depositID uint) (*models.DepositRequest, error)
	UpdateDepositRequest(deposit *models.DepositRequest) error
	CreateWithdrawRequest(withdraw *models.WithdrawRequest) error
	GetWithdrawRequestByID(withdrawID uint) (*models.WithdrawRequest, error)
	UpdateWithdrawRequest(withdraw *models.WithdrawRequest) error
	GetAdminWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error)
	FindAdminUser() (*models.User, error)

	GetPendingWithdrawalRequests(pagination models.PaginationParams) ([]models.WithdrawRequest, int64, error)
	UpdateCustomerWalletBalance(tx *gorm.DB, wallet *models.Wallet) error
	GetCustomerWallet(userID uint) (*models.Wallet, error)
}

type AdminWalletRepository struct {
	DB *gorm.DB
}

func NewAdminWalletRepository(db *gorm.DB) *AdminWalletRepository {
	return &AdminWalletRepository{DB: db}
}

func (r *AdminWalletRepository) FindAdminUser() (*models.User, error) {
	var adminUser models.User
	err := r.DB.Where("role = ?", models.RoleAdmin).First(&adminUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin user not found")
		}
		return nil, fmt.Errorf("failed to find admin user: %w", err)
	}
	return &adminUser, nil
}

func (r *AdminWalletRepository) GetAdminWallet() (*models.Wallet, error) {
	adminUser, err := r.FindAdminUser()
	if err != nil {
		return nil, err
	}

	var wallet models.Wallet
	err = r.DB.Where("user_id = ?", adminUser.ID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin wallet not found")
		}
		return nil, fmt.Errorf("failed to retrieve admin wallet: %w", err)
	}
	return &wallet, nil
}

func (r *AdminWalletRepository) CreateWalletTransaction(tx *gorm.DB, transaction *models.WalletTransaction) error {
	return tx.Create(transaction).Error
}

func (r *AdminWalletRepository) UpdateWalletBalance(tx *gorm.DB, wallet *models.Wallet) error {
	wallet.LastUpdated = time.Now()
	return tx.Save(wallet).Error
}

func (r *AdminWalletRepository) CreateDepositRequest(deposit *models.DepositRequest) error {
	return r.DB.Create(deposit).Error
}

func (r *AdminWalletRepository) GetDepositRequestByID(depositID uint) (*models.DepositRequest, error) {
	var deposit models.DepositRequest
	err := r.DB.First(&deposit, depositID).Error
	return &deposit, err
}

func (r *AdminWalletRepository) UpdateDepositRequest(deposit *models.DepositRequest) error {
	return r.DB.Save(deposit).Error
}

func (r *AdminWalletRepository) CreateWithdrawRequest(withdraw *models.WithdrawRequest) error {
	return r.DB.Create(withdraw).Error
}

func (r *AdminWalletRepository) GetWithdrawRequestByID(withdrawID uint) (*models.WithdrawRequest, error) {
	var withdraw models.WithdrawRequest
	err := r.DB.First(&withdraw, withdrawID).Error
	return &withdraw, err
}

func (r *AdminWalletRepository) UpdateWithdrawRequest(withdraw *models.WithdrawRequest) error {
	return r.DB.Save(withdraw).Error
}

func (r *AdminWalletRepository) GetAdminWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	adminUser, err := r.FindAdminUser()
	if err != nil {
		return nil, 0, err
	}

	var transactions []models.WalletTransaction
	var total int64

	query := r.DB.Where("user_id = ?", adminUser.ID).Order("created_at DESC")

	err = query.Model(&models.WalletTransaction{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count admin wallet transactions: %w", err)
	}

	err = query.Offset((pagination.Page - 1) * pagination.Limit).Limit(pagination.Limit).Find(&transactions).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve admin wallet transactions: %w", err)
	}

	return transactions, total, nil
}

func (r *AdminWalletRepository) GetPendingWithdrawalRequests(pagination models.PaginationParams) ([]models.WithdrawRequest, int64, error) {
	var withdrawals []models.WithdrawRequest
	var total int64

	query := r.DB.Where("status = ?", models.TxStatusPending).Order("created_at DESC")

	err := query.Model(&models.WithdrawRequest{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pending withdrawal requests: %w", err)
	}

	err = query.Offset((pagination.Page - 1) * pagination.Limit).Limit(pagination.Limit).Find(&withdrawals).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve pending withdrawal requests: %w", err)
	}

	return withdrawals, total, nil
}

func (r *AdminWalletRepository) UpdateCustomerWalletBalance(tx *gorm.DB, wallet *models.Wallet) error {
	wallet.LastUpdated = time.Now()
	return tx.Save(wallet).Error
}

func (r *AdminWalletRepository) GetCustomerWallet(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.DB.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer wallet not found")
		}
		return nil, fmt.Errorf("failed to retrieve customer wallet: %w", err)
	}
	return &wallet, nil
}
