package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type SignalRepository struct{}

func (r *SignalRepository) GetAllSignals() ([]models.Signal, error) {
	var signals []models.Signal
	err := db.DB.Find(&signals).Error
	return signals, err
}

func (r *SignalRepository) DeactivateSignal(signalID uint) error {
	return db.DB.Model(&models.Signal{}).
		Where("id = ?", signalID).
		Update("status", "inactive").Error
}

func (r *SignalRepository) GetPendingSignals() ([]models.Signal, error) {
	var signals []models.Signal
	err := db.DB.Where("status = ?", "pending").Find(&signals).Error
	return signals, err
}

func (r *SignalRepository) UpdateStatus(signalID uint, status string) error {
	return db.DB.Model(&models.Signal{}).Where("id = ?", signalID).Update("status", status).Error
}
