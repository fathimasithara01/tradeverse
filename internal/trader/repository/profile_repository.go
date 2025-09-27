package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ITraderProfileRepository interface {
	GetTraderProfileByUserID(userID uint) (*models.TraderProfile, error)
	CreateTraderProfile(profile *models.TraderProfile) error
	UpdateTraderProfile(profile *models.TraderProfile) error
	DeleteTraderProfile(profileID uint) error
}

type TraderProfileRepository struct {
	db *gorm.DB
}

func NewTraderProfileRepository(db *gorm.DB) *TraderProfileRepository {
	return &TraderProfileRepository{db: db}
}

func (r *TraderProfileRepository) GetTraderProfileByUserID(userID uint) (*models.TraderProfile, error) {
	var profile models.TraderProfile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *TraderProfileRepository) CreateTraderProfile(profile *models.TraderProfile) error {
	return r.db.Create(profile).Error
}

func (r *TraderProfileRepository) UpdateTraderProfile(profile *models.TraderProfile) error {
	return r.db.Save(profile).Error
}

func (r *TraderProfileRepository) DeleteTraderProfile(profileID uint) error {
	return r.db.Delete(&models.TraderProfile{}, profileID).Error
}
