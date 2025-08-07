package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
)

type PlanUsageStats struct {
	PlanID            uint `json:"plan_id"`
	SubscriptionCount int  `json:"subscription_count"`
}

type StatsRepository struct{}

func (r *StatsRepository) GetPlanUsageStats() ([]PlanUsageStats, error) {
	var stats []PlanUsageStats

	err := db.DB.Raw(`
		SELECT plan_id, COUNT(*) AS subscription_count
		FROM subscriptions
		GROUP BY plan_id
		ORDER BY subscription_count DESC
	`).Scan(&stats).Error

	return stats, err
}
