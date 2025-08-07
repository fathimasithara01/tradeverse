package service

import (
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type TraderService struct {
	Repo repository.TraderRepository
}

func (s *TraderService) GetAllTraders() ([]models.Trader, error) {
	return s.Repo.GetAllTraders()
}

func (s *TraderService) ToggleBan(id uint) error {
	return s.Repo.ToggleBanStatus(id)
}
