package service

import "github.com/fathimasithara01/tradeverse/admin/repository"

type TraderStatsService struct {
	Repo repository.TraderStatsRepository
}

func (s *TraderStatsService) GetAllTraderRankings() ([]repository.TraderStats, error) {
	return s.Repo.GetAllRankings()
}

func (s *TraderStatsService) GetBadgeForTrader(id uint) (repository.TraderStats, error) {
	return s.Repo.GetTraderBadge(id)
}
