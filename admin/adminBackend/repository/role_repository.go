package repository

import (
	"errors"

	"github.com/fathimasithara01/tradeverse/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{DB: db}
}

func (r *RoleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	if err := r.DB.Preload("CreatedBy").Order("id asc").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.DB.Create(role).Error
}

func (r *RoleRepository) FindByID(id uint) (models.Role, error) {
	var role models.Role
	if err := r.DB.Preload("CreatedBy").First(&role, id).Error; err != nil {
		return models.Role{}, errors.New("role not found")
	}
	return role, nil
}

func (r *RoleRepository) Update(role *models.Role) error {
	return r.DB.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Role{}, id).Error
}
