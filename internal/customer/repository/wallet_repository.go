package repository

import (
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

// WalletRepository defines the interface for wallet data operations.
type WalletRepository interface {
	GetWalletByUserID(userID uint) (*models.Wallet, error)
	CreateWallet(wallet *models.Wallet) error
	UpdateWalletBalance(tx *gorm.DB, wallet *models.Wallet, amount float64) error // Use transaction for update
	CreateWalletTransaction(tx *gorm.DB, transaction *models.WalletTransaction) error
	FindWalletTransactions(userID uint, page, limit int) ([]models.WalletTransaction, int64, error)

	CreateDepositRequest(req *models.DepositRequest) error
	FindDepositRequestByPGTxID(pgTxID string) (*models.DepositRequest, error)
	UpdateDepositRequest(req *models.DepositRequest) error

	CreateWithdrawRequest(req *models.WithdrawRequest) error
	FindWithdrawRequestByID(withdrawID uint) (*models.WithdrawRequest, error)
	UpdateWithdrawRequest(req *models.WithdrawRequest) error

	// New transaction method to begin a database transaction
	BeginTransaction() *gorm.DB
}

// walletRepository implements WalletRepository with GORM.
type walletRepository struct {
	db *gorm.DB
}

// NewWalletRepository creates a new WalletRepository instance.
func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

// BeginTransaction starts a new database transaction.
func (r *walletRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

// GetWalletByUserID retrieves a user's wallet.
func (r *walletRepository) GetWalletByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Wallet not found, not an error
		}
		return nil, err
	}
	return &wallet, nil
}

// CreateWallet creates a new wallet for a user.
func (r *walletRepository) CreateWallet(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

// UpdateWalletBalance updates the wallet balance within a transaction.
func (r *walletRepository) UpdateWalletBalance(tx *gorm.DB, wallet *models.Wallet, amount float64) error {
	// IMPORTANT: In a real-world scenario, you would likely use a database-level
	// atomic update like "UPDATE wallets SET balance = balance + ? WHERE id = ? AND version = ?"
	// to prevent race conditions without explicit locking.
	// For simplicity, this example updates the balance directly on the model and saves it.
	// GORM's .Save() will update all fields.
	wallet.Balance += amount
	wallet.LastUpdated = time.Now()
	return tx.Save(wallet).Error
}

// CreateWalletTransaction records a new transaction.
func (r *walletRepository) CreateWalletTransaction(tx *gorm.DB, transaction *models.WalletTransaction) error {
	return tx.Create(transaction).Error
}

// FindWalletTransactions lists a user's wallet transactions with pagination.
func (r *walletRepository) FindWalletTransactions(userID uint, page, limit int) ([]models.WalletTransaction, int64, error) {
	var transactions []models.WalletTransaction
	var total int64

	offset := (page - 1) * limit

	// Count total transactions for pagination metadata
	err := r.db.Model(&models.WalletTransaction{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated transactions
	err = r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// CreateDepositRequest saves a new deposit request.
func (r *walletRepository) CreateDepositRequest(req *models.DepositRequest) error {
	return r.db.Create(req).Error
}

// FindDepositRequestByPGTxID finds a deposit request by the payment gateway's transaction ID.
func (r *walletRepository) FindDepositRequestByPGTxID(pgTxID string) (*models.DepositRequest, error) {
	var req models.DepositRequest
	err := r.db.Where("payment_gateway_tx_id = ?", pgTxID).First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

// UpdateDepositRequest updates an existing deposit request.
func (r *walletRepository) UpdateDepositRequest(req *models.DepositRequest) error {
	return r.db.Save(req).Error // Save updates all fields
}

// CreateWithdrawRequest saves a new withdrawal request.
func (r *walletRepository) CreateWithdrawRequest(req *models.WithdrawRequest) error {
	return r.db.Create(req).Error
}

func (r *walletRepository) FindWithdrawRequestByID(withdrawID uint) (*models.WithdrawRequest, error) {
	var req models.WithdrawRequest
	err := r.db.First(&req, withdrawID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

func (r *walletRepository) UpdateWithdrawRequest(req *models.WithdrawRequest) error {
	return r.db.Save(req).Error
}
