package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
)

type TraderStatsRepository struct{}

func (r *TraderStatsRepository) GetAllRankings() ([]TraderStats, error) {
	query := `
		SELECT trader_id,
			COUNT(*) AS total,
			SUM(CASE WHEN status = 'won' THEN 1 ELSE 0 END) AS won,
			SUM(CASE WHEN status = 'lost' THEN 1 ELSE 0 END) AS lost
		FROM signals
		GROUP BY trader_id;
	`

	rows, err := db.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []TraderStats
	for rows.Next() {
		var t TraderStats
		rows.Scan(&t.TraderID, &t.Total, &t.Won, &t.Lost)

		t.WinRate = float64(t.Won) / float64(t.Total) * 100
		t.Badge = assignBadge(t.WinRate)

		stats = append(stats, t)
	}

	return stats, nil
}

func (r *TraderStatsRepository) GetTraderBadge(traderID uint) (TraderStats, error) {
	query := `
		SELECT trader_id,
			COUNT(*) AS total,
			SUM(CASE WHEN status = 'won' THEN 1 ELSE 0 END) AS won,
			SUM(CASE WHEN status = 'lost' THEN 1 ELSE 0 END) AS lost
		FROM signals
		WHERE trader_id = ?
		GROUP BY trader_id;
	`

	var t TraderStats
	err := db.DB.Raw(query, traderID).Row().Scan(&t.TraderID, &t.Total, &t.Won, &t.Lost)
	if err != nil {
		return t, err
	}

	t.WinRate = float64(t.Won) / float64(t.Total) * 100
	t.Badge = assignBadge(t.WinRate)
	return t, nil
}

// badge logic
func assignBadge(rate float64) string {
	switch {
	case rate >= 90:
		return "🏅 Elite Trader"
	case rate >= 75:
		return "🥈 Pro Trader"
	case rate >= 60:
		return "🥉 Intermediate"
	default:
		return "❌ Needs Review"
	}
}
