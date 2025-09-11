package service

import (
	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type IActivityService interface {
	GetActiveSessions() ([]models.CopySession, error)
	GetRecentTradeLogs() ([]models.TradeLog, error)
}

type ActivityService struct {
	Repo repository.IActivityRepository
}

func NewActivityService(repo repository.IActivityRepository) IActivityService {
	return &ActivityService{Repo: repo}
}

func (s *ActivityService) GetActiveSessions() ([]models.CopySession, error) {
	return s.Repo.GetActiveSessions()
}

func (s *ActivityService) GetRecentTradeLogs() ([]models.TradeLog, error) {
	return s.Repo.GetTradeLogs(100)
}
