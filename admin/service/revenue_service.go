package service

import (
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type RevenueService struct {
	Repo repository.RevenueRepository
}

func (s *RevenueService) GetMonthly() ([]repository.MonthlyRevenue, error) {
	return s.Repo.GetMonthlyRevenue()
}
