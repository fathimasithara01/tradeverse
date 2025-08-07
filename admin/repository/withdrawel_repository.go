package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type WithdrawalRepository struct{}

func (r *WithdrawalRepository) GetPending() ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal
	err := db.DB.Where("status = ?", "pending").Find(&withdrawals).Error
	return withdrawals, err
}

func (r *WithdrawalRepository) UpdateStatus(id uint, status string, note string) error {
	return db.DB.Model(&models.Withdrawal{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "note": note}).Error
}
