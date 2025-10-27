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
	FindAdminUser() (*models.User, error)

	// Specific to Admin Wallet Transactions
	AdminGetWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error)
	// Get all transactions across the entire platform, not just admin's
	GetAllWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error)

	GetPendingWithdrawalRequests(pagination models.PaginationParams) ([]models.WithdrawRequest, int64, error)
	UpdateCustomerWalletBalance(tx *gorm.DB, wallet *models.Wallet) error
	GetCustomerWallet(userID uint) (*models.Wallet, error)

	// For customer-specific transactions, including user details
	GetAllCustomerTransactions(pagination models.PaginationParams) ([]models.AdminTransactionDisplayDTO, int64, error)
	GetCustomerByUserID(userID uint) (*models.User, error)
}

type AdminWalletRepository struct {
	DB *gorm.DB
}

func NewAdminWalletRepository(db *gorm.DB) *AdminWalletRepository {
	return &AdminWalletRepository{DB: db}
}

// GetAllCustomerTransactions retrieves all wallet transactions made by customers, including their user details.
func (r *AdminWalletRepository) GetAllCustomerTransactions(pagination models.PaginationParams) ([]models.AdminTransactionDisplayDTO, int64, error) {
	var transactions []models.WalletTransaction
	var totalCount int64
	var displayTransactions []models.AdminTransactionDisplayDTO

	// Query for customer transactions, preloading User to get details
	// Make sure the `User` struct in `models.WalletTransaction` is correctly defined
	query := r.DB.
		Preload("User").
		Joins("INNER JOIN users ON wallet_transactions.user_id = users.id").
		Where("users.role = ?", models.RoleCustomer).
		Order("wallet_transactions.created_at DESC")

	// Handle Search Query
	if pagination.Search != "" {
		searchLike := "%" + pagination.Search + "%"
		query = query.Where(`
			wallet_transactions.transaction_type ILIKE ? OR
			wallet_transactions.description ILIKE ? OR
			wallet_transactions.reference_id ILIKE ? OR
			users.email ILIKE ? OR
			users.name ILIKE ? OR
			users.phone ILIKE ?`,
			searchLike, searchLike, searchLike, searchLike, searchLike, searchLike,
		)
	}

	// Count total before pagination
	if err := query.Model(&models.WalletTransaction{}).Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count all customer transactions: %w", err)
	}

	// Pagination logic
	offset := (pagination.Page - 1) * pagination.Limit
	if err := query.Offset(offset).Limit(pagination.Limit).Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve all customer transactions: %w", err)
	}

	// Map transactions to DTOs
	for _, tx := range transactions {
		userName := "N/A"
		userEmail := "N/A"
		userPhone := "N/A"

		// Check if User relationship was successfully loaded
		if tx.UserID != 0 && tx.User.ID != 0 {
			userName = tx.User.Name
			userEmail = tx.User.Email
			userPhone = tx.User.Phone
		}

		displayTransactions = append(displayTransactions, models.AdminTransactionDisplayDTO{
			ID:              tx.ID,
			UserID:          tx.UserID,
			UserName:        userName,
			UserEmail:       userEmail,
			UserPhone:       userPhone,
			TransactionType: tx.TransactionType,
			Amount:          tx.Amount,
			Currency:        tx.Currency,
			Status:          tx.Status,
			ReferenceID:     tx.ReferenceID,
			Description:     tx.Description,
			CreatedAt:       tx.CreatedAt,
			BalanceBefore:   tx.BalanceBefore,
			BalanceAfter:    tx.BalanceAfter,
		})
	}

	return displayTransactions, totalCount, nil
}

// GetCustomerByUserID retrieves a customer's user details by their ID.
func (r *AdminWalletRepository) GetCustomerByUserID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("id = ? AND role = ?", userID, models.RoleCustomer).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer user not found")
		}
		return nil, fmt.Errorf("failed to retrieve customer user: %w", err)
	}
	return &user, nil
}

// AdminGetWalletTransactions retrieves wallet transactions specifically for the admin user.
func (r *AdminWalletRepository) AdminGetWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	adminUser, err := r.FindAdminUser()
	if err != nil {
		return nil, 0, err
	}

	var transactions []models.WalletTransaction
	var total int64

	// Filter by admin user ID
	query := r.DB.Where("user_id = ?", adminUser.ID).Order("created_at DESC")

	// Add search filter if SearchQuery is provided
	if pagination.Search != "" {
		searchLike := "%" + pagination.Search + "%"
		query = query.Where(
			"transaction_type ILIKE ? OR description ILIKE ? OR reference_id ILIKE ? OR payment_gateway_tx_id ILIKE ?",
			searchLike, searchLike, searchLike, searchLike,
		)
	}

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

// GetAllWalletTransactions retrieves all wallet transactions across the entire platform (all users, including admin).
func (r *AdminWalletRepository) GetAllWalletTransactions(pagination models.PaginationParams) ([]models.WalletTransaction, int64, error) {
	var transactions []models.WalletTransaction
	var total int64

	// No user_id filter here, fetches all transactions
	query := r.DB.Order("created_at DESC")

	// Add search filter if SearchQuery is provided (for all transactions)
	if pagination.Search != "" {
		searchLike := "%" + pagination.Search + "%"
		query = query.Where(
			"transaction_type ILIKE ? OR description ILIKE ? OR reference_id ILIKE ? OR payment_gateway_tx_id ILIKE ?",
			searchLike, searchLike, searchLike, searchLike,
		)
	}

	err := query.Model(&models.WalletTransaction{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count all wallet transactions: %w", err)
	}

	err = query.Offset((pagination.Page - 1) * pagination.Limit).Limit(pagination.Limit).Find(&transactions).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve all wallet transactions: %w", err)
	}

	return transactions, total, nil
}

// FindAdminUser finds the first user with the RoleAdmin.
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

// GetAdminWallet retrieves the wallet associated with the admin user.
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

// CreateWalletTransaction creates a new wallet transaction record within a GORM transaction.
func (r *AdminWalletRepository) CreateWalletTransaction(tx *gorm.DB, transaction *models.WalletTransaction) error {
	return tx.Create(transaction).Error
}

// UpdateWalletBalance updates the balance and last updated timestamp of a wallet within a GORM transaction.
func (r *AdminWalletRepository) UpdateWalletBalance(tx *gorm.DB, wallet *models.Wallet) error {
	wallet.LastUpdated = time.Now()
	return tx.Save(wallet).Error
}

// CreateDepositRequest creates a new deposit request.
func (r *AdminWalletRepository) CreateDepositRequest(deposit *models.DepositRequest) error {
	return r.DB.Create(deposit).Error
}

// GetDepositRequestByID retrieves a deposit request by its ID.
func (r *AdminWalletRepository) GetDepositRequestByID(depositID uint) (*models.DepositRequest, error) {
	var deposit models.DepositRequest
	err := r.DB.First(&deposit, depositID).Error
	return &deposit, err
}

// UpdateDepositRequest updates an existing deposit request.
func (r *AdminWalletRepository) UpdateDepositRequest(deposit *models.DepositRequest) error {
	return r.DB.Save(deposit).Error
}

// CreateWithdrawRequest creates a new withdrawal request.
func (r *AdminWalletRepository) CreateWithdrawRequest(withdraw *models.WithdrawRequest) error {
	return r.DB.Create(withdraw).Error
}

// GetWithdrawRequestByID retrieves a withdrawal request by its ID.
func (r *AdminWalletRepository) GetWithdrawRequestByID(withdrawID uint) (*models.WithdrawRequest, error) {
	var withdraw models.WithdrawRequest
	// Ensure User is preloaded for later use (e.g., in service for customer wallet)
	// This Preload should now work correctly due to the updated models.WithdrawRequest struct
	err := r.DB.Preload("User").First(&withdraw, withdrawID).Error
	return &withdraw, err
}

// UpdateWithdrawRequest updates an existing withdrawal request.
func (r *AdminWalletRepository) UpdateWithdrawRequest(withdraw *models.WithdrawRequest) error {
	return r.DB.Save(withdraw).Error
}

// GetPendingWithdrawalRequests retrieves all withdrawal requests with a 'Pending' status.
func (r *AdminWalletRepository) GetPendingWithdrawalRequests(pagination models.PaginationParams) ([]models.WithdrawRequest, int64, error) {
	var withdrawals []models.WithdrawRequest
	var total int64

	// Preload User to get customer details
	// This Preload should now work correctly due to the updated models.WithdrawRequest struct
	query := r.DB.Preload("User").Where("status = ?", models.TxStatusPending).Order("created_at DESC")

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

// UpdateCustomerWalletBalance updates the balance and last updated timestamp of a customer's wallet within a GORM transaction.
func (r *AdminWalletRepository) UpdateCustomerWalletBalance(tx *gorm.DB, wallet *models.Wallet) error {
	wallet.LastUpdated = time.Now()
	return tx.Save(wallet).Error
}

// GetCustomerWallet retrieves a customer's wallet by their user ID.
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