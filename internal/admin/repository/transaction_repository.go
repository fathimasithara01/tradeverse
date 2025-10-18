package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetAllTransactions(page, limit int, search string, year, month, day int) ([]models.WalletTransaction, int64, error)
	GetAvailableYears() ([]int, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) GetAllTransactions(page, limit int, search string, year, month, day int) ([]models.WalletTransaction, int64, error) {
	var transactions []models.WalletTransaction
	var total int64

	query := r.db.Model(&models.WalletTransaction{}).
		Preload("User")

	if search != "" {
		query = query.Joins("JOIN users ON users.id = wallet_transactions.user_id").
			Where("users.name ILIKE ? OR users.email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply date filters
	if year > 0 {
		query = query.Where("EXTRACT(YEAR FROM wallet_transactions.created_at) = ?", year)
	}
	if month > 0 {
		query = query.Where("EXTRACT(MONTH FROM wallet_transactions.created_at) = ?", month)
	}
	if day > 0 {
		query = query.Where("EXTRACT(DAY FROM wallet_transactions.created_at) = ?", day)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	if err := query.Order("wallet_transactions.created_at DESC").
		Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *transactionRepository) GetAvailableYears() ([]int, error) {
	var years []int
	err := r.db.Model(&models.WalletTransaction{}).
		Distinct("EXTRACT(YEAR FROM created_at)").
		Order("EXTRACT(YEAR FROM created_at) DESC").
		Find(&years).Error
	if err != nil {
		return nil, err
	}
	return years, nil
}
