package service

import (
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type AnnouncementService struct {
	Repo repository.AnnouncementRepository
}

func (s *AnnouncementService) Create(a models.Announcement) (models.Announcement, error) {
	return s.Repo.Create(a)
}

func (s *AnnouncementService) GetAll() ([]models.Announcement, error) {
	return s.Repo.GetAll()
}
