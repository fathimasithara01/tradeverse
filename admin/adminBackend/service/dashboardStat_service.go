package service

import (
	"github.com/fathimasithara01/tradeverse/repository" // Make sure this import path is correct
)

// DashboardStats defines the structure for our API response.
// The `json` tags are critical and must match what the JavaScript expects.
type DashboardStats struct {
	CustomerCount int64 `json:"CustomerCount"`
	TraderCount   int64 `json:"TraderCount"`
	ProductCount  int64 `json:"ProductCount"`
	OrderCount    int64 `json:"OrderCount"`
}

type DashboardService struct {
	Repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{Repo: repo}
}

// GetDashboardStats calls the repository to fetch all required counts.
func (s *DashboardService) GetDashboardStats() (DashboardStats, error) {
	var stats DashboardStats
	var err error

	stats.CustomerCount, err = s.Repo.GetCustomerCount()
	if err != nil {
		return DashboardStats{}, err
	}

	stats.TraderCount, err = s.Repo.GetTraderCount()
	if err != nil {
		return DashboardStats{}, err
	}

	stats.ProductCount, err = s.Repo.GetProductCount()
	if err != nil {
		return DashboardStats{}, err
	}

	stats.OrderCount, err = s.Repo.GetOrderCount()
	if err != nil {
		return DashboardStats{}, err
	}

	return stats, nil
}

type MonthlyOrdersStats struct {
	Labels []string `json:"labels"`
	Data   []int    `json:"data"`
}

// GetMonthlyOrderStats processes the raw data from the repository into a chart-friendly format.
func (s *DashboardService) GetMonthlyOrderStats(year int) (MonthlyOrdersStats, error) {
	// Get the raw counts from the repository
	monthlyCounts, err := s.Repo.GetMonthlyOrderCounts(year)
	if err != nil {
		return MonthlyOrdersStats{}, err
	}

	// Initialize a full year of data with zeros. This ensures all 12 months are present.
	labels := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	data := make([]int, 12)

	// Populate the data array with the counts from the database.
	// The month from the DB is 1-based (January=1), so we subtract 1 for the 0-based array index.
	for _, result := range monthlyCounts {
		if result.Month >= 1 && result.Month <= 12 {
			data[result.Month-1] = result.Count
		}
	}

	return MonthlyOrdersStats{
		Labels: labels,
		Data:   data,
	}, nil
}
