package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type AnnouncementRepository struct{}

func (r *AnnouncementRepository) Create(a models.Announcement) (models.Announcement, error) {
	err := db.DB.Create(&a).Error
	return a, err
}

func (r *AnnouncementRepository) GetAll() ([]models.Announcement, error) {
	var list []models.Announcement
	err := db.DB.Order("created_at desc").Find(&list).Error
	return list, err
}
