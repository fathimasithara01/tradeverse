package repository

import (
	"github.com/fathimasithara01/tradeverse/models" // Make sure this import path is correct
	"gorm.io/gorm"
)

type DashboardRepository struct {
	DB *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{DB: db}
}

// GetCustomerCount counts all records in the 'customers' table.
func (r *DashboardRepository) GetCustomerCount() (int64, error) {
	var count int64
	// Assumes you have a models.Customer struct
	if err := r.DB.Model(&models.CustomerProfile{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetTraderCount counts all records in the 'traders' table.
func (r *DashboardRepository) GetTraderCount() (int64, error) {
	var count int64
	// Assumes you have a models.Trader struct
	if err := r.DB.Model(&models.TraderProfile{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetProductCount counts all records in the 'products' table.
func (r *DashboardRepository) GetProductCount() (int64, error) {
	var count int64
	// Assumes you have a models.Product struct
	if err := r.DB.Model(&models.Product{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetOrderCount counts all records in the 'orders' table.
func (r *DashboardRepository) GetOrderCount() (int64, error) {
	var count int64
	// Uses the models.Order struct you asked for
	if err := r.DB.Model(&models.Order{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

type MonthlyCountResult struct {
	Month int `json:"month"`
	Count int `json:"count"`
}

// GetMonthlyOrderCounts retrieves the count of orders for each month of a given year.
func (r *DashboardRepository) GetMonthlyOrderCounts(year int) ([]MonthlyCountResult, error) {
	var results []MonthlyCountResult

	// This query groups orders by the month they were created and counts them.
	// It's designed to work with PostgreSQL. See note below for MySQL.
	err := r.DB.Model(&models.Order{}).
		Select("EXTRACT(MONTH FROM created_at) as month, COUNT(id) as count").
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Group("month").
		Order("month asc").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
