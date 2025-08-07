package repository

import (
	"errors"

	"github.com/fathimasithara01/tradeverse/admin/models"
	"gorm.io/gorm"
)

type AdminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{DB: db}
}

func (r *AdminRepository) Create(admin models.Admin) (models.Admin, error) {
	if err := r.DB.Create(&admin).Error; err != nil {
		return models.Admin{}, err
	}
	return admin, nil
}

func (r *AdminRepository) GetByEmail(email string) (models.Admin, error) {
	var admin models.Admin
	if err := r.DB.Where("email = ?", email).First(&admin).Error; err != nil {
		return models.Admin{}, errors.New("admin not found")
	}
	return admin, nil
}
