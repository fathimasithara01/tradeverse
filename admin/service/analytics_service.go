package service

import "github.com/fathimasithara01/tradeverse/admin/repository"

type AnalyticsService struct {
	Repo repository.AnalyticsRepository
}

func (s *AnalyticsService) GetSignalStats() map[string]interface{} {
	total, won, lost := s.Repo.CountSignals()
	winRate := float64(won) / float64(total) * 100
	return map[string]interface{}{
		"total":   total,
		"won":     won,
		"lost":    lost,
		"winRate": winRate,
	}
}

func (s *AnalyticsService) GetTraderStats() ([]map[string]interface{}, error) {
	return s.Repo.GetTraderStats()
}
