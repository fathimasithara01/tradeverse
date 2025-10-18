package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type IDashboardRepository interface {
	GetCustomerCount() (int64, error)
	GetApprovedTraderCount() (int64, error)
	// GetActiveSessionCount() (int64, error)
	GetMonthlyRecurringRevenue() (int64, error)
	GetTotalSignalCount() (int64, error)
	GetMonthlySignups(role models.UserRole) ([]SignupStat, error)
	GetLatestSignups() ([]models.User, error)
	GetTopTraders() ([]models.User, error)
}

type DashboardRepository struct{ DB *gorm.DB }

func NewDashboardRepository(db *gorm.DB) IDashboardRepository {
	return &DashboardRepository{DB: db}
}

func (r *DashboardRepository) GetCustomerCount() (int64, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("role = ?", models.RoleCustomer).Count(&count).Error
	return count, err
}

func (r *DashboardRepository) GetApprovedTraderCount() (int64, error) {
	var count int64
	err := r.DB.Model(&models.User{}).
		Joins("JOIN trader_profiles ON users.id = trader_profiles.user_id").
		Where("users.role = ? AND trader_profiles.status = ?", models.RoleTrader, models.StatusApproved).
		Count(&count).Error
	return count, err
}

// func (r *DashboardRepository) GetActiveSessionCount() (int64, error) {
// 	var count int64
// 	err := r.DB.Model(&models.CopySession{}).Where("is_active = ?", true).Count(&count).Error
// 	return count, err
// }

func (r *DashboardRepository) GetTotalSignalCount() (int64, error) {
	var count int64
	err := r.DB.Model(&models.Signal{}).Count(&count).Error
	return count, err
}
func (r *DashboardRepository) FindAdminUser() (*models.User, error) {
	var adminUser models.User
	err := r.DB.Where("role = ?", models.RoleAdmin).First(&adminUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin user not found")
		}
		return nil, fmt.Errorf("failed to find admin user: %w", err)
	}
	return &adminUser, nil
}
func (r *DashboardRepository) GetMonthlyRecurringRevenue() (int64, error) {
	adminUser, err := r.FindAdminUser()
	if err != nil {
		return 0, err
	}

	var total float64
	err = r.DB.Model(&models.WalletTransaction{}).
		Where("user_id = ? AND type = ? AND status = ?", adminUser.ID, models.TxTypeSubscription, models.TxStatusSuccess).
		Select("COALESCE(SUM(amount),0)").
		Scan(&total).Error
	if err != nil {
		return 0, fmt.Errorf("failed to calculate MRR: %w", err)
	}

	return int64(total), nil
}

type SignupStat struct {
	Month time.Time
	Count int
}

func (r *DashboardRepository) GetMonthlySignups(role models.UserRole) ([]SignupStat, error) {
	var results []SignupStat
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)

	err := r.DB.Model(&models.User{}).
		Select("DATE_TRUNC('month', created_at) as month, COUNT(id) as count").
		Where("role = ? AND created_at >= ?", role, sixMonthsAgo).
		Group("month").
		Order("month ASC").
		Scan(&results).Error

	return results, err
}

func (r *DashboardRepository) GetLatestSignups() ([]models.User, error) {
	var users []models.User
	err := r.DB.Order("created_at DESC").Limit(5).Find(&users).Error
	return users, err
}

// Top 5 traders by PNL
func (r *DashboardRepository) GetTopTraders() ([]models.User, error) {
	var users []models.User
	err := r.DB.Joins("JOIN trader_profiles ON users.id = trader_profiles.user_id").
		Where("users.role = ? AND trader_profiles.status = ?", models.RoleTrader, models.StatusApproved).
		Order("trader_profiles.total_pnl DESC").
		Limit(5).
		Preload("TraderProfile").
		Find(&users).Error
	return users, err
}
