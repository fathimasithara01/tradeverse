package repository

import (
	"time"

	"github.com/fathimasithara01/tradeverse/models"
	"gorm.io/gorm"
)

type DashboardRepository struct{ DB *gorm.DB }

func NewDashboardRepository(db *gorm.DB) *DashboardRepository { return &DashboardRepository{DB: db} }

// --- KPI Card Queries ---
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

func (r *DashboardRepository) GetActiveSessionCount() (int64, error) {
	var count int64
	err := r.DB.Model(&models.CopySession{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

func (r *DashboardRepository) GetMonthlyRecurringRevenue() (int64, error) {
	// MOCK DATA: Replace with a real query when you have a subscriptions table.
	// Example real query:
	// var total int64
	// r.DB.Table("subscriptions").Where("status = ?", "active").Select("SUM(price)").Row().Scan(&total)
	return 12450, nil
}

// --- Chart & Table Queries ---

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

func (r *DashboardRepository) GetLatestSignups() ([]models.User, error) {
	var users []models.User
	err := r.DB.Order("created_at DESC").Limit(5).Find(&users).Error
	return users, err
}
