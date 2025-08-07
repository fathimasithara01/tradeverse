package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type LogRepository struct{}

func (r *LogRepository) CreateLog(log models.Log) error {
	return db.DB.Create(&log).Error
}

func (r *LogRepository) GetAllLogs() ([]models.Log, error) {
	var logs []models.Log
	err := db.DB.Order("timestamp DESC").Find(&logs).Error
	return logs, err
}
