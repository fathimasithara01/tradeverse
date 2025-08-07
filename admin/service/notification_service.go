package service

import (
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type NotificationService struct {
	Repo repository.NotificationRepository
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>
// func NotifyUser(userID uint, msg string) {
// 	ws.WSManager.SendTo(userID, msg)
// }

func (s *NotificationService) Send(userID uint, title, message, notifType string) error {
	n := models.Notification{
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    notifType,
		Read:    false,
	}
	return s.Repo.Create(n)
}

func (s *NotificationService) GetUserNotifications(userID uint) ([]models.Notification, error) {
	return s.Repo.GetByUser(userID)
}

func (s *NotificationService) MarkAsRead(userID uint) error {
	return s.Repo.MarkAllRead(userID)
}
