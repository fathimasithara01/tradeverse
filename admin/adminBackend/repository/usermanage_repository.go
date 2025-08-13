package repository

import (
	"github.com/fathimasithara01/tradeverse/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}
