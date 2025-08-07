package service

import (
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type WalletService struct {
	Repo repository.WalletRepository
}

func (s *WalletService) GetWallet(userID uint) (models.Wallet, []models.WalletTransaction, error) {
	return s.Repo.GetWalletByUserID(userID)
}

func (s *WalletService) Credit(userID uint, amount float64, desc string) error {
	return s.Repo.Credit(userID, amount, desc)
}

func (s *WalletService) Debit(userID uint, amount float64, desc string) error {
	return s.Repo.Debit(userID, amount, desc)
}
