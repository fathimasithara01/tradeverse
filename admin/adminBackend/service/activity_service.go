package service

import (
	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/repository"
)

type ActivityService struct {
	Repo *repository.ActivityRepository
}

func NewActivityService(repo *repository.ActivityRepository) *ActivityService {
	return &ActivityService{Repo: repo}
}

func (s *ActivityService) GetActiveSessions() ([]models.CopySession, error) {
	return s.Repo.GetActiveSessions()
}

func (s *ActivityService) GetRecentTradeLogs() ([]models.TradeLog, error) {
	// Let's fetch the last 100 logs for the admin panel.
	return s.Repo.GetTradeLogs(100)
}
