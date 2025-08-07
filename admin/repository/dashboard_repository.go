package repository

import (
	"github.com/fathimasithara01/tradeverse/admin/db"
	"github.com/fathimasithara01/tradeverse/admin/models"
)

type DashboardRepository struct{}

func (r *DashboardRepository) GetUserCount() int64 {
	var count int64
	// db.DB.Models(&user{}).Count(&count)
	db.DB.Model(&models.User{}).Count(&count)
	return count
}

func (r *DashboardRepository) GetTraderCount() int64 {
	var count int64
	db.DB.Model(&models.Trader{}).Count(&count)
	return count
}

func (r *DashboardRepository) GetActiveSubscriptionCount() int64 {
	var count int64
	db.DB.Model(&models.Subscription{}).Where("active = ?", true).Count(&count)
	return count
}

func (r *DashboardRepository) GetSignalCount() int64 {
	var count int64
	db.DB.Model(&models.Signal{}).Count(&count)
	return count
}
