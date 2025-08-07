package service

import (
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type StatsService struct {
	Repo repository.StatsRepository
}

func (s *StatsService) GetPlanStats() ([]repository.PlanUsageStats, error) {
	return s.Repo.GetPlanUsageStats()
}
