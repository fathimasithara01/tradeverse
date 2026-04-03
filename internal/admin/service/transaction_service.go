package service

import (
	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type ITransactionService interface {
	GetTransactions(page, limit int, search string, year, month, day int) ([]models.WalletTransaction, int64, error)
	GetAvailableYears() ([]int, error)
}

type transactionService struct {
	repo repository.ITransactionRepository
}

func NewTransactionService(repo repository.ITransactionRepository) ITransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) GetTransactions(page, limit int, search string, year, month, day int) ([]models.WalletTransaction, int64, error) {
	return s.repo.GetAllTransactions(page, limit, search, year, month, day)
}

func (s *transactionService) GetAvailableYears() ([]int, error) {
	return s.repo.GetAvailableYears()
}
