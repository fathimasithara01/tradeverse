package repository

import (
	"errors"
	"math"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindByRole(role models.UserRole) ([]models.User, error)
	FindAllNonAdmins() ([]models.User, error)
	FindAllAdvanced(options UserQueryOptions) (PaginatedUsers, error)
	Update(user *models.User) error
	Delete(id uint) error
	UpdateTraderStatus(userID uint, newStatus models.TraderStatus) error
	FindAllWithRole() ([]models.User, error)
	AssignRoleToUser(userID uint, roleID uint, roleName models.UserRole) error
	FindByIDs(ids []uint) ([]models.User, error)

	FindTradersByStatus(status models.TraderStatus) ([]models.User, error)
	GetLatestOpenTradeForUser(userID uint) (models.Trade, error)

	GetUserByIDWithProfile(id uint) (*models.User, error) // New
	UpdateUserAndProfile(user *models.User) error         // New
	DeleteUser(id uint) error
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user models.User, profile models.CustomerProfile) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	profile.UserID = user.ID // Link profile to the newly created user
	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByIDWithProfile(id uint) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("CustomerProfile").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserAndProfile(user *models.User) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if user.CustomerProfile.ID != 0 { // Check if a profile exists and has an ID
		if err := tx.Save(&user.CustomerProfile).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else if user.CustomerProfile.UserID != 0 { // If no ID but UserID is set, it might be a new profile for an existing user

		if err := tx.Create(&user.CustomerProfile).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *UserRepository) DeleteUser(id uint) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("user_id = ?", id).Delete(&models.CustomerProfile{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&models.User{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (models.User, error) {
	var user models.User
	err := r.DB.Preload("CustomerProfile").Preload("TraderProfile").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.DB.Where("LOWER(email) = LOWER(?)", email).First(&user).Error

	return user, err
}

func (r *UserRepository) FindByRole(role models.UserRole) ([]models.User, error) {
	var users []models.User
	if err := r.DB.
		Preload("CustomerProfile").
		Preload("TraderProfile").
		Where("role = ?", role).
		Order("id asc").
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) FindAllNonAdmins() ([]models.User, error) {
	var users []models.User
	err := r.DB.
		Where("role <> ?", models.RoleAdmin).
		Order("id asc").
		Find(&users).Error
	return users, err
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

func (r *UserRepository) FindByIDs(ids []uint) ([]models.User, error) {
	var users []models.User
	if len(ids) == 0 {
		return users, nil // Return empty slice if no IDs are provided
	}
	if err := r.DB.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateTraderStatus(userID uint, newStatus models.TraderStatus) error {
	return r.DB.Model(&models.TraderProfile{}).Where("user_id = ?", userID).Update("status", newStatus).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}

type UserQueryOptions struct {
	Search string          `form:"search"`
	Role   models.UserRole `form:"role"`
	Period string          `form:"period"`
	Page   int             `form:"page"`
	Limit  int             `form:"limit"`
}

type PaginatedUsers struct {
	Users      []models.User `json:"users"`
	TotalPages int           `json:"total_pages"`
	Page       int           `json:"page"`
}

func (r *UserRepository) FindAllAdvanced(options UserQueryOptions) (PaginatedUsers, error) {
	var users []models.User
	var totalUsers int64

	query := r.DB.Model(&models.User{})

	if options.Search != "" {
		searchQuery := "%" + options.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery)
	}

	if options.Role != "" {
		query = query.Where("role = ?", options.Role)
	}

	if options.Period != "" {
		now := time.Now()
		var startTime time.Time
		if options.Period == "monthly" {
			startTime = now.AddDate(0, -1, 0)
		} else if options.Period == "yearly" {
			startTime = now.AddDate(-1, 0, 0)
		}
		if !startTime.IsZero() {
			query = query.Where("created_at >= ?", startTime)
		}
	}

	if err := query.Count(&totalUsers).Error; err != nil {
		return PaginatedUsers{}, err
	}

	if options.Page <= 0 {
		options.Page = 1
	}
	if options.Limit <= 0 {
		options.Limit = 10
	}

	offset := (options.Page - 1) * options.Limit

	err := query.Order("id asc").Limit(options.Limit).Offset(offset).Find(&users).Error
	if err != nil {
		return PaginatedUsers{}, err
	}

	totalPages := int(math.Ceil(float64(totalUsers) / float64(options.Limit)))

	return PaginatedUsers{
		Users:      users,
		TotalPages: totalPages,
		Page:       options.Page,
	}, nil
}

func (r *UserRepository) FindAllWithRole() ([]models.User, error) {
	var users []models.User
	err := r.DB.
		Preload("Role").
		Where("role <> ?", models.RoleAdmin).
		Order("id asc").
		Find(&users).Error
	return users, err
}

func (r *UserRepository) AssignRoleToUser(userID uint, roleID uint, roleName models.UserRole) error {
	updates := map[string]interface{}{
		"role_id": roleID,
		"role":    roleName,
	}

	return r.DB.Model(&models.User{}).Where("id = ?", userID).UpdateColumns(updates).Error
}

func (r *UserRepository) GetLatestOpenTradeForUser(userID uint) (models.Trade, error) {
	var trade models.Trade
	err := r.DB.Where("master_user_id = ? AND status = ?", userID, "open").
		Order("opened_at desc").
		First(&trade).Error
	return trade, err
}
