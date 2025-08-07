package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type UserRepository struct{}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := db.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) UpdateStatus(userID uint, status string) error {
	return db.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", status).Error
}

func (r *UserRepository) GetByID(id uint) (models.User, error) {
	var user models.User
	err := db.DB.First(&user, id).Error
	return user, err
}
