package service

import (
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type TraderAnalyticsService struct {
	Repo repository.TraderAnalyticsRepository
}

func (s *TraderAnalyticsService) GetAllStats() ([]repository.TraderStats, error) {
	return s.Repo.GetTraderStats()
}

func (s *TraderAnalyticsService) GetRanked(limit int) ([]repository.TraderStats, error) {
	return s.Repo.GetTopRankedTraders(limit)
}
