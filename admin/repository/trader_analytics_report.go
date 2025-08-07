package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
)

type TraderAnalyticsRepository struct{}

type TraderStats struct {
	TraderID      uint
	Name          string
	Followers     int64
	SignalsPosted int64
	TotalPnL      float64
	Total         int
	Won           int
	Lost          int
	WinRate       float64
	Badge         string
}

func (r *TraderAnalyticsRepository) GetTraderStats() ([]TraderStats, error) {
	var stats []TraderStats

	rows, err := db.DB.Raw(`
		SELECT t.id as trader_id, t.name,
			(SELECT COUNT(*) FROM followers f WHERE f.trader_id = t.id) as followers,
			(SELECT COUNT(*) FROM signals s WHERE s.trader_id = t.id) as signals_posted,
			t.total_pnl
		FROM traders t
	`).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stat TraderStats
		rows.Scan(&stat.TraderID, &stat.Name, &stat.Followers, &stat.SignalsPosted, &stat.TotalPnL)
		stats = append(stats, stat)
	}

	return stats, nil
}

func (r *TraderAnalyticsRepository) GetTopRankedTraders(limit int) ([]TraderStats, error) {
	var stats []TraderStats
	db.DB.Raw(`
		SELECT t.id as trader_id, t.name,
			(SELECT COUNT(*) FROM followers f WHERE f.trader_id = t.id) as followers,
			(SELECT COUNT(*) FROM signals s WHERE s.trader_id = t.id) as signals_posted,
			t.total_pnl
		FROM traders t
		ORDER BY t.total_pnl DESC
		LIMIT ?
	`, limit).Scan(&stats)

	return stats, nil
}
