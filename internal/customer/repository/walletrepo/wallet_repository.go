package walletrepo // Changed package name

import (
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrWalletNotFound          = errors.New("wallet not found for user")
	ErrInsufficientFunds       = errors.New("insufficient funds")
	ErrDepositRequestNotFound  = errors.New("deposit request not found")
	ErrWithdrawRequestNotFound = errors.New("withdrawal request not found")
)

type WalletRepository interface {
	GetUserWallet(userID uint) (*models.Wallet, error)
	CreateWallet(wallet *models.Wallet) error
	CreditWallet(tx *gorm.DB, walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error
	DebitWallet(tx *gorm.DB, walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error
	CreateWalletTransaction(tx *gorm.DB, walletTx *models.WalletTransaction) error
	GetWalletTransactions(userID uint, pagination models.PaginationParams) ([]models.WalletTransaction, int64, error)
	GetWalletTransactionByID(txID uint) (*models.WalletTransaction, error)

	CreateDepositRequest(req *models.DepositRequest) error
	GetDepositRequestByID(reqID uint) (*models.DepositRequest, error)
	UpdateDepositRequest(req *models.DepositRequest) error

	CreateWithdrawRequest(req *models.WithdrawRequest) error
	GetWithdrawRequestByID(reqID uint) (*models.WithdrawRequest, error)
	UpdateWithdrawRequest(req *models.WithdrawRequest) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

// ... (Rest of your walletRepository methods remain the same) ...
func (r *walletRepository) GetUserWallet(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) CreateWallet(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepository) CreditWallet(tx *gorm.DB, walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error {
	var wallet models.Wallet
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, walletID).Error; err != nil {
		return err
	}

	balanceBefore := wallet.Balance
	wallet.Balance += amount
	wallet.LastUpdated = time.Now()
	if err := tx.Save(&wallet).Error; err != nil {
		return err
	}

	walletTx := &models.WalletTransaction{
		WalletID:        walletID,
		UserID:          wallet.UserID,
		TransactionType: transactionType,
		Amount:          amount,
		Currency:        wallet.Currency,
		Status:          models.TxStatusSuccess,
		ReferenceID:     referenceID,
		Description:     description,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    wallet.Balance,
	}
	return r.CreateWalletTransaction(tx, walletTx)
}

func (r *walletRepository) DebitWallet(tx *gorm.DB, walletID uint, amount float64, transactionType models.TransactionType, referenceID string, description string) error {
	var wallet models.Wallet
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, walletID).Error; err != nil {
		return err
	}

	if wallet.Balance < amount {
		return ErrInsufficientFunds
	}

	balanceBefore := wallet.Balance
	wallet.Balance -= amount
	wallet.LastUpdated = time.Now()
	if err := tx.Save(&wallet).Error; err != nil {
		return err
	}

	walletTx := &models.WalletTransaction{
		WalletID:        walletID,
		UserID:          wallet.UserID,
		TransactionType: transactionType,
		Amount:          amount,
		Currency:        wallet.Currency,
		Status:          models.TxStatusSuccess,
		ReferenceID:     referenceID,
		Description:     description,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    wallet.Balance,
	}
	return r.CreateWalletTransaction(tx, walletTx)
}

func (r *walletRepository) CreateWalletTransaction(tx *gorm.DB, walletTx *models.WalletTransaction) error {
	return tx.Create(walletTx).Error
}

func (r *walletRepository) GetWalletTransactions(userID uint, pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	var transactions []models.WalletTransaction
	var total int64

	query := r.db.Where("user_id = ?", userID)

	if err := query.Model(&models.WalletTransaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(pagination.Limit).Offset((pagination.Page - 1) * pagination.Limit).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *walletRepository) GetWalletTransactionByID(txID uint) (*models.WalletTransaction, error) {
	var transaction models.WalletTransaction
	if err := r.db.First(&transaction, txID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *walletRepository) CreateDepositRequest(req *models.DepositRequest) error {
	return r.db.Create(req).Error
}

func (r *walletRepository) GetDepositRequestByID(reqID uint) (*models.DepositRequest, error) {
	var req models.DepositRequest
	if err := r.db.First(&req, reqID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDepositRequestNotFound
		}
		return nil, err
	}
	return &req, nil
}

func (r *walletRepository) UpdateDepositRequest(req *models.DepositRequest) error {
	return r.db.Save(req).Error
}

func (r *walletRepository) CreateWithdrawRequest(req *models.WithdrawRequest) error {
	return r.db.Create(req).Error
}

func (r *walletRepository) GetWithdrawRequestByID(reqID uint) (*models.WithdrawRequest, error) {
	var req models.WithdrawRequest
	if err := r.db.First(&req, reqID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWithdrawRequestNotFound
		}
		return nil, err
	}
	return &req, nil
}

func (r *walletRepository) UpdateWithdrawRequest(req *models.WithdrawRequest) error {
	return r.db.Save(req).Error
}
