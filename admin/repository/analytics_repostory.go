package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type AnalyticsRepository struct{}

func (r *AnalyticsRepository) CountSignals() (total, won, lost int64) {
	db.DB.Model(&models.Signal{}).Count(&total)
	db.DB.Model(&models.Signal{}).Where("status = ?", "won").Count(&won)
	db.DB.Model(&models.Signal{}).Where("status = ?", "lost").Count(&lost)
	return
}

func (r *AnalyticsRepository) GetTraderStats() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `
		SELECT trader_id,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'won' THEN 1 ELSE 0 END) as won,
			SUM(CASE WHEN status = 'lost' THEN 1 ELSE 0 END) as lost
		FROM signals
		GROUP BY trader_id
	`
	rows, err := db.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var traderID, total, won, lost int
		_ = rows.Scan(&traderID, &total, &won, &lost)
		entry := map[string]interface{}{
			"trader_id": traderID,
			"total":     total,
			"won":       won,
			"lost":      lost,
			"win_rate":  float64(won) / float64(total) * 100,
		}
		results = append(results, entry)
	}
	return results, nil
}
