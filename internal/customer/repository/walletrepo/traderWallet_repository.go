// internal/customer/repository/walletrepo/wallet_repository.go
package walletrepo

import (
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrWalletNotFound            = errors.New("wallet not found")
	ErrInsufficientFunds         = errors.New("insufficient funds")
	ErrTransactionFailed         = errors.New("wallet transaction failed")
	ErrDepositRequestNotFound    = errors.New("deposit request not found")
	ErrWithdrawalRequestNotFound = errors.New("withdrawal request not found")
)

type WalletRepository interface {
	GetUserWallet(userID uint) (*models.Wallet, error)
	GetOrCreateWallet(userID uint) (*models.Wallet, error)
	UpdateWallet(wallet *models.Wallet) error                // Simple update, consider if needed outside transaction
	UpdateWalletTx(tx *gorm.DB, wallet *models.Wallet) error // Transaction-safe wallet update

	DebitWallet(tx *gorm.DB, walletID uint, amount float64, txType models.TransactionType, referenceID, description string) error
	CreditWallet(tx *gorm.DB, walletID uint, amount float64, txType models.TransactionType, referenceID, description string) (*models.WalletTransaction, error)
	CreateWalletTransaction(tx *gorm.DB, transaction *models.WalletTransaction) error
	GetWalletTransactions(userID uint, pagination models.PaginationParams) ([]models.WalletTransaction, int64, error)

	CreateDepositRequest(req *models.DepositRequest) error                      // <--- ADD THIS
	GetDepositRequestByID(id uint) (*models.DepositRequest, error)              // <--- ADD THIS
	UpdateDepositRequest(req *models.DepositRequest) error                      // <--- ADD THIS (for non-tx updates)
	UpdateDepositRequestTx(tx *gorm.DB, req *models.DepositRequest) error       // <--- ADD THIS (for tx updates)
	CreateWithdrawalRequest(req *models.WithdrawalRequest) error                // <--- ADD THIS
	GetWithdrawalRequestByID(id uint) (*models.WithdrawalRequest, error)        // <--- ADD THIS
	UpdateWithdrawalRequestTx(tx *gorm.DB, req *models.WithdrawalRequest) error // <--- ADD THIS

}

type gormWalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &gormWalletRepository{db: db}
}

func (r *gormWalletRepository) GetUserWallet(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to get wallet for user %d: %w", userID, err)
	}
	return &wallet, nil
}

func (r *gormWalletRepository) GetOrCreateWallet(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new wallet if not found
		wallet = models.Wallet{
			UserID:      userID,
			Balance:     0,
			Currency:    "INR", // Default currency, adjust as needed
			LastUpdated: time.Now(),
		}
		if createErr := r.db.Create(&wallet).Error; createErr != nil {
			return nil, fmt.Errorf("failed to create new wallet for user %d: %w", userID, createErr)
		}
		return &wallet, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get or create wallet for user %d: %w", userID, err)
	}
	return &wallet, nil
}

func (r *gormWalletRepository) UpdateWallet(wallet *models.Wallet) error {
	wallet.LastUpdated = time.Now()
	return r.db.Save(wallet).Error
}

func (r *gormWalletRepository) UpdateWalletTx(tx *gorm.DB, wallet *models.Wallet) error {
	wallet.LastUpdated = time.Now()
	return tx.Save(wallet).Error
}

func (r *gormWalletRepository) DebitWallet(tx *gorm.DB, walletID uint, amount float64, txType models.TransactionType, referenceID, description string) error {
	var wallet models.Wallet
	if err := tx.First(&wallet, walletID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrWalletNotFound
		}
		return fmt.Errorf("failed to find wallet %d for debit: %w", walletID, err)
	}

	if wallet.Balance < amount {
		return ErrInsufficientFunds
	}

	balanceBefore := wallet.Balance
	wallet.Balance -= amount
	wallet.LastUpdated = time.Now()
	if err := tx.Save(&wallet).Error; err != nil {
		return fmt.Errorf("failed to update wallet balance during debit: %w", err)
	}

	transaction := &models.WalletTransaction{
		WalletID:        walletID,
		UserID:          wallet.UserID,
		TransactionType: txType,
		Amount:          amount,
		Currency:        wallet.Currency,
		Status:          models.TxStatusSuccess, // Assuming success if debit goes through
		ReferenceID:     referenceID,
		Description:     description,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    wallet.Balance,
	}
	if err := tx.Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create debit transaction record: %w", err)
	}
	return nil
}

func (r *gormWalletRepository) CreditWallet(tx *gorm.DB, walletID uint, amount float64, txType models.TransactionType, referenceID, description string) (*models.WalletTransaction, error) {
	var wallet models.Wallet
	if err := tx.First(&wallet, walletID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to find wallet %d for credit: %w", walletID, err)
	}

	balanceBefore := wallet.Balance
	wallet.Balance += amount
	wallet.LastUpdated = time.Now()
	if err := tx.Save(&wallet).Error; err != nil {
		return nil, fmt.Errorf("failed to update wallet balance during credit: %w", err)
	}

	transaction := &models.WalletTransaction{
		WalletID:        walletID,
		UserID:          wallet.UserID,
		TransactionType: txType,
		Amount:          amount,
		Currency:        wallet.Currency,
		Status:          models.TxStatusSuccess, // Assuming success if credit goes through
		ReferenceID:     referenceID,
		Description:     description,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    wallet.Balance, // Wallet.Balance is now updated
	}
	if err := tx.Create(transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to create credit transaction record: %w", err)
	}
	return transaction, nil
}

func (r *gormWalletRepository) CreateWalletTransaction(tx *gorm.DB, transaction *models.WalletTransaction) error {
	return tx.Create(transaction).Error
}

func (r *gormWalletRepository) GetWalletTransactions(userID uint, pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	var transactions []models.WalletTransaction
	var totalCount int64

	query := r.db.Model(&models.WalletTransaction{}).Where("user_id = ?", userID)

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count wallet transactions: %w", err)
	}

	// Apply pagination
	if pagination.Limit == 0 {
		pagination.Limit = 10 // Default limit
	}
	if pagination.Page == 0 {
		pagination.Page = 1 // Default page
	}
	offset := (pagination.Page - 1) * pagination.Limit

	err := query.Order("created_at DESC").
		Limit(pagination.Limit).
		Offset(offset).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve wallet transactions: %w", err)
	}
	return transactions, totalCount, nil
}

func (r *gormWalletRepository) CreateDepositRequest(req *models.DepositRequest) error {
	return r.db.Create(req).Error
}

func (r *gormWalletRepository) GetDepositRequestByID(id uint) (*models.DepositRequest, error) {
	var req models.DepositRequest
	if err := r.db.First(&req, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDepositRequestNotFound
		}
		return nil, fmt.Errorf("failed to get deposit request by ID %d: %w", id, err)
	}
	return &req, nil
}

func (r *gormWalletRepository) UpdateDepositRequest(req *models.DepositRequest) error {
	return r.db.Save(req).Error
}

func (r *gormWalletRepository) UpdateDepositRequestTx(tx *gorm.DB, req *models.DepositRequest) error {
	return tx.Save(req).Error
}

func (r *gormWalletRepository) CreateWithdrawalRequest(req *models.WithdrawalRequest) error {
	return r.db.Create(req).Error
}

func (r *gormWalletRepository) GetWithdrawalRequestByID(id uint) (*models.WithdrawalRequest, error) {
	var req models.WithdrawalRequest
	if err := r.db.First(&req, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWithdrawalRequestNotFound
		}
		return nil, fmt.Errorf("failed to get withdrawal request by ID %d: %w", id, err)
	}
	return &req, nil
}

func (r *gormWalletRepository) UpdateWithdrawalRequestTx(tx *gorm.DB, req *models.WithdrawalRequest) error {
	return tx.Save(req).Error
}
