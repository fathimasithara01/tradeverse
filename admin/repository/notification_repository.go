package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type NotificationRepository struct{}

func (r *NotificationRepository) Create(n models.Notification) error {
	return db.DB.Create(&n).Error
}

func (r *NotificationRepository) GetByUser(userID uint) ([]models.Notification, error) {
	var notifs []models.Notification
	err := db.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifs).Error
	return notifs, err
}

func (r *NotificationRepository) MarkAllRead(userID uint) error {
	return db.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read = false", userID).
		Update("read", true).Error
}
