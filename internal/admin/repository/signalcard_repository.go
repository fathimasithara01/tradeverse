package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ISignalRepository interface {
	CreateSignal(ctx context.Context, signal *models.Signal) (*models.Signal, error)
	GetAllSignals(ctx context.Context) ([]models.Signal, error)
	GetMarketDataBySymbol(ctx context.Context, symbol string) (*models.MarketData, error)
	UpdateSignalCurrentPrice(ctx context.Context, signalID uint, newPrice float64) error
	UpdateSignalStatus(ctx context.Context, signalID uint, newStatus string) error
	GetActiveAndPendingSignals(ctx context.Context) ([]models.Signal, error)
}

type signalRepository struct {
	db *gorm.DB
}

func NewSignalRepository(db *gorm.DB) ISignalRepository {
	return &signalRepository{db: db}
}

func (r *signalRepository) CreateSignal(ctx context.Context, signal *models.Signal) (*models.Signal, error) {
	if err := r.db.WithContext(ctx).Create(signal).Error; err != nil {
		return nil, fmt.Errorf("failed to create signal: %w", err)
	}
	return signal, nil
}

func (r *signalRepository) GetAllSignals(ctx context.Context) ([]models.Signal, error) {
	var signals []models.Signal
	if err := r.db.WithContext(ctx).Find(&signals).Error; err != nil {
		return nil, fmt.Errorf("failed to get all signals: %w", err)
	}
	return signals, nil
}

func (r *signalRepository) GetActiveAndPendingSignals(ctx context.Context) ([]models.Signal, error) {
	var signals []models.Signal
	// Fetch signals that are either "Active" or "Pending" for status checks
	if err := r.db.WithContext(ctx).Where("status = ? OR status = ?", "Active", "Pending").Find(&signals).Error; err != nil {
		return nil, fmt.Errorf("failed to get active/pending signals: %w", err)
	}
	return signals, nil
}

func (r *signalRepository) GetMarketDataBySymbol(ctx context.Context, symbol string) (*models.MarketData, error) {
	var marketData models.MarketData

	// Ensure the symbol is uppercase for case-insensitive matching with the database (if necessary)
	// and to match the stored `Symbol` in `MarketData` which is likely uppercase.
	symbol = strings.ToUpper(symbol)

	err := r.db.WithContext(ctx).
		Where("UPPER(symbol) = ?", symbol). // Use UPPER() for robust case-insensitive search
		First(&marketData).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// If record not found, return nil MarketData and nil error.
			// This indicates that no data exists for the symbol, which is not an "error" in flow.
			return nil, nil
		}
		// For any other DB error, return the error.
		return nil, fmt.Errorf("failed to get market data for symbol %s: %w", symbol, err)
	}

	return &marketData, nil
}

func (r *signalRepository) UpdateSignalCurrentPrice(ctx context.Context, signalID uint, newPrice float64) error {
	result := r.db.WithContext(ctx).Model(&models.Signal{}).Where("id = ?", signalID).Update("current_price", newPrice)
	if result.Error != nil {
		return fmt.Errorf("failed to update signal current price: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Indicate that no record was found/updated
	}
	return nil
}

func (r *signalRepository) UpdateSignalStatus(ctx context.Context, signalID uint, newStatus string) error {
	result := r.db.WithContext(ctx).Model(&models.Signal{}).Where("id = ?", signalID).Update("status", newStatus)
	if result.Error != nil {
		return fmt.Errorf("failed to update signal status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Indicate that no record was found/updated
	}
	return nil
}
