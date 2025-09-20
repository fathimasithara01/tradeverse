package repository

import (
	"errors"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type SubscriberRepository interface {
	GetAll(traderID uint) ([]models.Subscriber, error)
	GetByID(id uint) (*models.Subscriber, error)
}

type subscriberRepo struct {
	db *gorm.DB
}

func NewSubscriberRepository(db *gorm.DB) SubscriberRepository {
	return &subscriberRepo{db: db}
}

func (r *subscriberRepo) GetAll(traderID uint) ([]models.Subscriber, error) {
	var subs []models.Subscriber
	if err := r.db.Where("trader_id = ?", traderID).Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *subscriberRepo) GetByID(id uint) (*models.Subscriber, error) {
	var sub models.Subscriber
	if err := r.db.First(&sub, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}
