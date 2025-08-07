package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type TraderRepository struct{}

func (r *TraderRepository) GetAllTraders() ([]models.Trader, error) {
	var traders []models.Trader
	err := db.DB.Find(&traders).Error
	return traders, err
}

func (r *TraderRepository) ToggleBanStatus(id uint) error {
	var trader models.Trader
	if err := db.DB.First(&trader, id).Error; err != nil {
		return err
	}

	trader.IsBanned = !trader.IsBanned
	return db.DB.Save(&trader).Error
}
