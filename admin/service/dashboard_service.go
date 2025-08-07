package service

import "github.com/fathimasithara01/tradeverse/admin/repository"

type DashboardService struct {
	Repo repository.DashboardRepository
}

func (s *DashboardService) GetDashboardStats() map[string]int64 {
	return map[string]int64{
		"users_count":          s.Repo.GetUserCount(),
		"traders_count":        s.Repo.GetTraderCount(),
		"active_subscriptions": s.Repo.GetActiveSubscriptionCount(),
		"total_signals":        s.Repo.GetSignalCount(),
	}
}
