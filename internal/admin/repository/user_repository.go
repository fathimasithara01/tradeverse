package repository

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *models.User) error
	CreateCustomerWithProfile(user *models.User, profile *models.CustomerProfile) error
	CreateTraderWithProfile(user *models.User, profile *models.TraderProfile) error

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

	GetUserByIDWithProfile(id uint) (*models.User, error)
	UpdateUserAndProfile(user *models.User) error
	DeleteUser(id uint) error

	GetUserByID(id uint) (*models.User, error)
	GetRoleByName(name models.UserRole) (*models.Role, error)
	UpdateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)

	GetUsersByRole(role models.UserRole) ([]models.User, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUsersByRole(role models.UserRole) ([]models.User, error) {
	var users []models.User
	// Make sure the query correctly targets the 'role' field
	err := r.DB.Where("role = ?", role).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("repo: failed to find users by role: %w", err)
	}
	return users, nil
}

func (r *UserRepository) CreateCustomerWithProfile(user *models.User, profile *models.CustomerProfile) error {
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	profile.UserID = user.ID
	if err := tx.Create(profile).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create customer profile: %w", err)
	}

	return tx.Commit().Error
}

func (r *UserRepository) CreateTraderWithProfile(user *models.User, profile *models.TraderProfile) error {
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	profile.UserID = user.ID
	if err := tx.Create(profile).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create trader profile: %w", err)
	}

	return tx.Commit().Error
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("RoleModel").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[DEBUG Repo] GetUserByID: User ID %d not found.", id)
			return nil, errors.New("user not found")
		}
		log.Printf("[ERROR Repo] GetUserByID: Failed to find user ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	log.Printf("[DEBUG Repo] GetUserByID: Found user ID %d, Email '%s'.", id, user.Email)
	return &user, nil
}

func (r *UserRepository) GetRoleByName(name models.UserRole) (*models.Role, error) {
	var role models.Role
	err := r.DB.Where("name = ?", string(name)).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[DEBUG Repo] GetRoleByName: Role '%s' not found.", name)
			return nil, errors.New("role not found")
		}
		log.Printf("[ERROR Repo] GetRoleByName: Failed to find role '%s': %v", name, err)
		return nil, fmt.Errorf("failed to find role by name: %w", err)
	}
	log.Printf("[DEBUG Repo] GetRoleByName: Found role '%s', ID %d.", name, role.ID)
	return &role, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("RoleModel").Where("LOWER(email) = LOWER(?)", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[DEBUG Repo] GetUserByEmail: User email '%s' not found.", email)
			return nil, errors.New("user not found")
		}
		log.Printf("[ERROR Repo] GetUserByEmail: Failed to find user by email '%s': %v", email, err)
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	log.Printf("[DEBUG Repo] GetUserByEmail: Found user with email '%s', ID %d.", email, user.ID)
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (models.User, error) {
	var user models.User
	err := r.DB.Preload("CustomerProfile").Preload("TraderProfile").Preload("RoleModel").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[DEBUG Repo] FindByID: User ID %d not found.", id)
			return models.User{}, errors.New("user not found")
		}
		log.Printf("[ERROR Repo] FindByID: Failed to find user ID %d: %v", id, err)
		return models.User{}, fmt.Errorf("database error finding user by ID: %w", err)
	}
	log.Printf("[DEBUG Repo] FindByID: Found user ID %d, Email '%s'.", id, user.Email)
	return user, nil
}

func (r *UserRepository) GetUserByIDWithProfile(id uint) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("CustomerProfile").Preload("TraderProfile").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserAndProfile(user *models.User) error {
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update user: %w", err)
	}

	if user.CustomerProfile.UserID != 0 {
		if user.CustomerProfile.ID != 0 {
			if err := tx.Save(&user.CustomerProfile).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update customer profile: %w", err)
			}
		} else {
			user.CustomerProfile.UserID = user.ID
			if err := tx.Create(&user.CustomerProfile).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create customer profile: %w", err)
			}
		}
	}

	if user.TraderProfile.UserID != 0 {
		if user.TraderProfile.ID != 0 {
			if err := tx.Save(&user.TraderProfile).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update trader profile: %w", err)
			}
		} else {
			user.TraderProfile.UserID = user.ID
			if err := tx.Create(&user.TraderProfile).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create trader profile: %w", err)
			}
		}
	}

	return tx.Commit().Error
}

func (r *UserRepository) DeleteUser(id uint) error {
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Unscoped().Delete(&models.User{}, id).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return tx.Commit().Error
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (models.User, error) {
	var user models.User

	err := r.DB.Preload("CustomerProfile").Preload("TraderProfile").Preload("RoleModel").Where("LOWER(email) = LOWER(?)", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[DEBUG Repo] FindByEmail: User email '%s' not found.", email)
			return models.User{}, errors.New("user not found")
		}
		log.Printf("[ERROR Repo] FindByEmail: Failed to find user by email '%s': %v", email, err)
		return models.User{}, fmt.Errorf("database error finding user by email: %w", err)
	}
	log.Printf("[DEBUG Repo] FindByEmail: Found user with email '%s', ID %d.", email, user.ID)
	return user, nil
}

func (r *UserRepository) FindByRole(role models.UserRole) ([]models.User, error) {
	var users []models.User
	err := r.DB.
		Preload("CustomerProfile").
		Preload("TraderProfile").
		Preload("RoleModel").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", string(role)).
		Order("users.id asc").
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find users by role '%s': %w", role, err)
	}
	return users, nil
}

func (r *UserRepository) FindAllNonAdmins() ([]models.User, error) {
	var users []models.User
	err := r.DB.
		Preload("CustomerProfile").
		Preload("TraderProfile").
		Preload("RoleModel").
		Where("role <> ?", models.RoleAdmin).
		Order("id asc").
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find non-admin users: %w", err)
	}
	return users, nil
}

func (r *UserRepository) FindTradersByStatus(status models.TraderStatus) ([]models.User, error) {
	var users []models.User
	err := r.DB.Joins("JOIN trader_profiles ON users.id = trader_profiles.user_id").
		Where("users.role = ? AND trader_profiles.status = ?", models.RoleTrader, status).
		Preload("TraderProfile").
		Preload("RoleModel").
		Order("users.id asc").
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find traders by status '%s': %w", status, err)
	}
	return users, nil
}

func (r *UserRepository) FindByIDs(ids []uint) ([]models.User, error) {
	var users []models.User
	if len(ids) == 0 {
		return users, nil
	}
	if err := r.DB.Preload("RoleModel").Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find users by IDs: %w", err)
	}
	return users, nil
}

func (r *UserRepository) UpdateTraderStatus(userID uint, newStatus models.TraderStatus) error {
	res := r.DB.Model(&models.TraderProfile{}).Where("user_id = ?", userID).Update("status", newStatus)
	if res.Error != nil {
		return fmt.Errorf("failed to update trader status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return errors.New("trader profile not found or status already set")
	}
	return nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	res := r.DB.Unscoped().Delete(&models.User{}, id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete user (hard delete): %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return errors.New("user not found for deletion")
	}
	return nil
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

	query := r.DB.Model(&models.User{}).Preload("RoleModel")

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
		return PaginatedUsers{}, fmt.Errorf("failed to count users for advanced query: %w", err)
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
		return PaginatedUsers{}, fmt.Errorf("failed to fetch paginated users for advanced query: %w", err)
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
		Preload("RoleModel").
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

	res := r.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	if res.Error != nil {
		return fmt.Errorf("failed to assign role to user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return errors.New("user not found for role assignment")
	}
	return nil
}

func (r *UserRepository) GetLatestOpenTradeForUser(userID uint) (models.Trade, error) {
	var trade models.Trade
	err := r.DB.Where("master_user_id = ? AND status = ?", userID, "open").
		Order("opened_at desc").
		First(&trade).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Trade{}, errors.New("no open trade found for user")
		}
		return models.Trade{}, fmt.Errorf("failed to get latest open trade: %w", err)
	}
	return trade, nil
}
