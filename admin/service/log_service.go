package service

import (
	"time"

	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type LogService struct {
	Repo repository.LogRepository
}

func (s *LogService) Create(actorRole string, userID *uint, action, details string) error {
	log := models.Log{
		UserID:    userID,
		ActorRole: actorRole,
		Action:    action,
		Details:   details,
		Timestamp: time.Now(),
	}
	return s.Repo.CreateLog(log)
}

func (s *LogService) GetAll() ([]models.Log, error) {
	return s.Repo.GetAllLogs()
}
