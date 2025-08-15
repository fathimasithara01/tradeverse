package repository

import (
	"github.com/fathimasithara01/tradeverse/models"
	"gorm.io/gorm"
)

type ActivityRepository struct{ DB *gorm.DB }

func NewActivityRepository(db *gorm.DB) *ActivityRepository { return &ActivityRepository{DB: db} }

func (r *ActivityRepository) GetActiveSessions() ([]models.CopySession, error) {
	var sessions []models.CopySession
	err := r.DB.Where("is_active = ?", true).
		Preload("Follower").
		Preload("Master").
		Order("id desc").
		Find(&sessions).Error
	return sessions, err
}

func (r *ActivityRepository) GetTradeLogs(limit int) ([]models.TradeLog, error) {
	var logs []models.TradeLog
	err := r.DB.Order("status asc, timestamp desc").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}
