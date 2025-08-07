package service

import (
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type WithdrawalService struct {
	Repo repository.WithdrawalRepository
}

func (s *WithdrawalService) GetPending() ([]models.Withdrawal, error) {
	return s.Repo.GetPending()
}

func (s *WithdrawalService) Approve(id uint) error {
	return s.Repo.UpdateStatus(id, "approved", "Approved by admin")
}

func (s *WithdrawalService) Reject(id uint, note string) error {
	return s.Repo.UpdateStatus(id, "rejected", note)
}
