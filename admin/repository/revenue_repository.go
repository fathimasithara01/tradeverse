package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
)

type MonthlyRevenue struct {
	Month string  `json:"month"`
	Total float64 `json:"total"`
}

type RevenueRepository struct{}

func (r *RevenueRepository) GetMonthlyRevenue() ([]MonthlyRevenue, error) {
	var result []MonthlyRevenue

	err := db.DB.
		Raw(`
			SELECT
				TO_CHAR(paid_at, 'YYYY-MM') AS month,
				SUM(amount) as total
			FROM payments
			GROUP BY month
			ORDER BY month DESC
		`).Scan(&result).Error

	return result, err
}
