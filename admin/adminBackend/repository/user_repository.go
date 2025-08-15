package repository

import (
	"errors"

	"github.com/fathimasithara01/tradeverse/models"
	"gorm.io/gorm"
)

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (models.User, error) {
	var user models.User
	if err := r.DB.Preload("CustomerProfile").Preload("TraderProfile").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) FindByRole(role models.UserRole) ([]models.User, error) {
	var users []models.User
	if err := r.DB.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	if err := r.DB.Where("LOWER(email) = LOWER(?)", email).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) FindTradersByStatus(status models.TraderStatus) ([]models.User, error) {
	var users []models.User
	err := r.DB.Joins("JOIN trader_profiles ON users.id = trader_profiles.user_id").
		Where("users.role = ? AND trader_profiles.status = ?", models.RoleTrader, status).
		Preload("TraderProfile"). // IMPORTANT: Preload the profile data.
		Order("users.id asc").
		Find(&users).Error

	return users, err
}

func (r *UserRepository) UpdateTraderStatus(userID uint, newStatus models.TraderStatus) error {
	return r.DB.Model(&models.TraderProfile{}).Where("user_id = ?", userID).Update("status", newStatus).Error
}
