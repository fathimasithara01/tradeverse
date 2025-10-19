package customerrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
	GetRoleByName(roleName string) (*models.Role, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetRoleByName(roleName string) (*models.Role, error) {
	var role models.Role
	if err := r.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID %d: %w", userID, err)
	}
	return &user, nil
}
