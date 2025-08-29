package service

import (
	"sync"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/fathimasithara01/tradeverse/pkg/repository"
)

type IDashboardService interface {
	GetDashboardStats() (DashboardStats, error)
	GetChartData() (ChartData, error)
	GetTopTraders() ([]models.User, error)
	GetLatestSignups() ([]models.User, error)
}

type GrowthData struct {
	Labels    []string `json:"labels"`
	Followers []int    `json:"followers"`
	Traders   []int    `json:"traders"`
}
type DistributionData struct {
	Labels []string `json:"labels"`
	Data   []int64  `json:"data"`
}
type ChartData struct {
	Growth       GrowthData       `json:"growth"`
	Distribution DistributionData `json:"distribution"`
}
type DashboardStats struct {
	MRR       int64 `json:"mrr"`
	Followers int64 `json:"followers"`
	Traders   int64 `json:"traders"`
	Sessions  int64 `json:"sessions"`
}

type DashboardService struct {
	Repo repository.IDashboardRepository
}

func NewDashboardService(repo repository.IDashboardRepository) IDashboardService {
	return &DashboardService{Repo: repo}
}

func (s *DashboardService) GetDashboardStats() (DashboardStats, error) {
	var stats DashboardStats
	var err error
	var wg sync.WaitGroup
	var errChan = make(chan error, 4)

	wg.Add(4)

	go func() {
		defer wg.Done()
		stats.MRR, _ = s.Repo.GetMonthlyRecurringRevenue() // Ignoring error for mock data
	}()

	go func() {
		defer wg.Done()
		stats.Followers, err = s.Repo.GetCustomerCount() // Corrected function call
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		stats.Traders, err = s.Repo.GetApprovedTraderCount()
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		stats.Sessions, err = s.Repo.GetActiveSessionCount()
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for e := range errChan {
		if e != nil {
			return DashboardStats{}, e
		}
	}

	return stats, nil
}

func (s *DashboardService) GetChartData() (ChartData, error) {
	var chartData ChartData

	followerCount, err := s.Repo.GetCustomerCount()
	if err != nil {
		return ChartData{}, err
	}
	traderCount, err := s.Repo.GetApprovedTraderCount()
	if err != nil {
		return ChartData{}, err
	}
	chartData.Distribution.Labels = []string{"Followers", "Traders"}
	chartData.Distribution.Data = []int64{followerCount, traderCount}

	followerSignups, err := s.Repo.GetMonthlySignups(models.RoleCustomer)
	if err != nil {
		return ChartData{}, err
	}
	traderSignups, err := s.Repo.GetMonthlySignups(models.RoleTrader)
	if err != nil {
		return ChartData{}, err
	}

	labels, followerData := processSignupStats(followerSignups)
	_, traderData := processSignupStats(traderSignups)

	chartData.Growth.Labels = labels
	chartData.Growth.Followers = followerData
	chartData.Growth.Traders = traderData

	return chartData, nil
}

func processSignupStats(stats []repository.SignupStat) ([]string, []int) {
	labels := make([]string, 6)
	data := make([]int, 6)
	now := time.Now()

	statsMap := make(map[time.Month]int)
	for _, s := range stats {
		statsMap[s.Month.Month()] = s.Count
	}

	for i := 5; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		labels[5-i] = month.Format("Jan")
		if count, ok := statsMap[month.Month()]; ok {
			data[5-i] = count
		} else {
			data[5-i] = 0
		}
	}
	return labels, data
}

func (s *DashboardService) GetTopTraders() ([]models.User, error) {
	return s.Repo.GetTopTraders()
}

func (s *DashboardService) GetLatestSignups() ([]models.User, error) {
	return s.Repo.GetLatestSignups()
}
