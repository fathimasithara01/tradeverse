package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.WalletTransaction) error
	UpdateTransaction(transaction *models.WalletTransaction) error
	// ... other transaction methods
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(transaction *models.WalletTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) UpdateTransaction(transaction *models.WalletTransaction) error {
	return r.db.Save(transaction).Error
}
