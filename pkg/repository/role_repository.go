package repository

import (
	"errors"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	Create(user *models.Role) error
	FindAll() ([]models.Role, error)
	FindByID(id uint) (models.Role, error)
	Update(role *models.Role) error
	Delete(id uint) error
	FindByIDWithPermissions(id uint) (models.Role, error)
	UpdatePermissions(role *models.Role, permissions []models.Permission) error
	RoleHasPermission(roleID uint, permissionName string) (bool, error)
}

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &RoleRepository{DB: db}
}

func (r *RoleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	if err := r.DB.Order("id asc").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.DB.Create(role).Error
}

func (r *RoleRepository) FindByID(id uint) (models.Role, error) {
	var role models.Role
	if err := r.DB.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Role{}, errors.New("role not found")
		}
		return models.Role{}, err
	}
	return role, nil
}

func (r *RoleRepository) Update(role *models.Role) error {
	return r.DB.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Role{}, id).Error
}

func (r *RoleRepository) FindByIDWithPermissions(id uint) (models.Role, error) {
	var role models.Role
	if err := r.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		return models.Role{}, errors.New("role not found")
	}
	return role, nil
}

func (r *RoleRepository) UpdatePermissions(role *models.Role, permissions []models.Permission) error {
	return r.DB.Model(role).Association("Permissions").Replace(permissions)
}

func (r *RoleRepository) RoleHasPermission(roleID uint, permissionName string) (bool, error) {
	var count int64
	err := r.DB.Table("role_permissions").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND permissions.name = ?", roleID, permissionName).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
