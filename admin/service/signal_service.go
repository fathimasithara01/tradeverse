package service

import (
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type SignalService struct {
	Repo repository.SignalRepository
}

func (s *SignalService) GetAllSignals() ([]models.Signal, error) {
	return s.Repo.GetAllSignals()
}

func (s *SignalService) Deactivate(id uint) error {
	return s.Repo.DeactivateSignal(id)
}

func (s *SignalService) GetPending() ([]models.Signal, error) {
	return s.Repo.GetPendingSignals()
}

func (s *SignalService) Approve(id uint) error {
	return s.Repo.UpdateStatus(id, "approved")
}

func (s *SignalService) Reject(id uint) error {
	return s.Repo.UpdateStatus(id, "rejected")
}
