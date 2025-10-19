package repository

import (
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type IDashboardRepository interface {
	GetCustomerCount() (int64, error)
	GetApprovedTraderCount() (int64, error)
	GetMonthlyRecurringRevenue() (int64, error)
	GetTotalSignalCount() (int64, error)
	GetMonthlySignups(role models.UserRole) ([]SignupStat, error)
	GetLatestSignups() ([]models.User, error)
	GetTopTraders() ([]models.User, error)
}

type DashboardRepository struct {
	DB *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) IDashboardRepository {
	return &DashboardRepository{DB: db}
}

// ✅ Count total customers
func (r *DashboardRepository) GetCustomerCount() (int64, error) {
	var count int64
	err := r.DB.Model(&models.User{}).
		Where("role = ?", models.RoleCustomer).
		Count(&count).Error
	return count, err
}

// ✅ Count total traders (approved only if profile exists)
func (r *DashboardRepository) GetApprovedTraderCount() (int64, error) {
	var count int64

	// safer for cases where table might not exist yet
	if !r.DB.Migrator().HasTable("trader_profiles") {
		return 0, nil
	}

	err := r.DB.Table("users").
		Joins("LEFT JOIN trader_profiles ON users.id = trader_profiles.user_id").
		Where("users.role = ?", models.RoleTrader).
		Where("(trader_profiles.status = ? OR trader_profiles.status IS NULL)", models.StatusApproved).
		Count(&count).Error

	return count, err
}

// ✅ Total signals created
func (r *DashboardRepository) GetTotalSignalCount() (int64, error) {
	var count int64
	err := r.DB.Model(&models.Signal{}).Count(&count).Error
	return count, err
}

// ✅ Total revenue from subscriptions (all customers)
func (r *DashboardRepository) GetMonthlyRecurringRevenue() (int64, error) {
	var total float64
	err := r.DB.Model(&models.WalletTransaction{}).
		Where("type = ? AND status = ?", models.TxTypeSubscription, models.TxStatusSuccess).
		Select("COALESCE(SUM(amount),0)").Scan(&total).Error
	if err != nil {
		return 0, fmt.Errorf("failed to calculate MRR: %w", err)
	}
	return int64(total), nil
}

// ✅ Monthly signups (for charts)
type SignupStat struct {
	Month time.Time
	Count int
}

func (r *DashboardRepository) GetMonthlySignups(role models.UserRole) ([]SignupStat, error) {
	var stats []SignupStat
	err := r.DB.Model(&models.User{}).
		Select("DATE_TRUNC('month', created_at) AS month, COUNT(*) AS count").
		Where("role = ?", role).
		Group("month").
		Order("month ASC").
		Scan(&stats).Error
	return stats, err
}

// Optional for table widgets
func (r *DashboardRepository) GetLatestSignups() ([]models.User, error) {
	var users []models.User
	err := r.DB.Where("role = ?", models.RoleCustomer).
		Order("created_at DESC").Limit(5).Find(&users).Error
	return users, err
}

func (r *DashboardRepository) GetTopTraders() ([]models.User, error) {
	var traders []models.User
	err := r.DB.Where("role = ?", models.RoleTrader).
		Order("created_at DESC").Limit(5).Find(&traders).Error
	return traders, err
}
