package repository

import (
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type TraderRepository interface {
	CreateTrader(user *models.User, profile *models.TraderProfile) error
	GetByEmail(email string) (*models.User, error)
}

type traderRepository struct {
	db *gorm.DB
}

func NewTraderRepository(db *gorm.DB) TraderRepository {
	return &traderRepository{db: db}
}

func (r *traderRepository) CreateTrader(user *models.User, profile *models.TraderProfile) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		profile.UserID = user.ID
		if err := tx.Create(profile).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *traderRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("TraderProfile").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
