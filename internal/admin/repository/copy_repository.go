package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ICopyRepository interface {
	GetSessionStatus(followerID, masterID uint) (models.CopySession, error)
	StartOrUpdateSession(followerID, masterID uint) error
	StopSession(followerID, masterID uint) error
}

type CopyRepository struct{ DB *gorm.DB }

func NewCopyRepository(db *gorm.DB) ICopyRepository { return &CopyRepository{DB: db} }

func (r *CopyRepository) GetSessionStatus(followerID, masterID uint) (models.CopySession, error) {
	var session models.CopySession
	err := r.DB.Where("follower_id = ? AND master_id = ?", followerID, masterID).First(&session).Error
	return session, err
}

func (r *CopyRepository) StartOrUpdateSession(followerID, masterID uint) error {
	session := models.CopySession{
		FollowerID: followerID,
		MasterID:   masterID,
		IsActive:   true,
	}

	return r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "follower_id"}, {Name: "master_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_active"}),
	}).Create(&session).Error
}

func (r *CopyRepository) StopSession(followerID, masterID uint) error {
	return r.DB.Model(&models.CopySession{}).
		Where("follower_id = ? AND master_id = ?", followerID, masterID).
		Update("is_active", false).Error
}
