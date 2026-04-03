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

type SignalRepository struct {
	db *gorm.DB
}

func NewSignalRepository(db *gorm.DB) ISignalRepository {
	return &SignalRepository{db: db}
}

func (r *SignalRepository) CreateSignal(ctx context.Context, signal *models.Signal) (*models.Signal, error) {
	if err := r.db.WithContext(ctx).Create(signal).Error; err != nil {
		return nil, fmt.Errorf("failed to create signal: %w", err)
	}
	return signal, nil
}

func (r *SignalRepository) GetAllSignals(ctx context.Context) ([]models.Signal, error) {
	var signals []models.Signal
	if err := r.db.WithContext(ctx).Find(&signals).Error; err != nil {
		return nil, fmt.Errorf("failed to get all signals: %w", err)
	}
	return signals, nil
}

func (r *SignalRepository) GetActiveAndPendingSignals(ctx context.Context) ([]models.Signal, error) {
	var signals []models.Signal
	if err := r.db.WithContext(ctx).Where("status = ? OR status = ?", "Active", "Pending").Find(&signals).Error; err != nil {
		return nil, fmt.Errorf("failed to get active/pending signals: %w", err)
	}
	return signals, nil
}

func (r *SignalRepository) GetMarketDataBySymbol(ctx context.Context, symbol string) (*models.MarketData, error) {
	var marketData models.MarketData

	symbol = strings.ToUpper(symbol)

	err := r.db.WithContext(ctx).
		Where("UPPER(symbol) = ?", symbol).
		First(&marketData).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get market data for symbol %s: %w", symbol, err)
	}

	return &marketData, nil
}

func (r *SignalRepository) UpdateSignalCurrentPrice(ctx context.Context, signalID uint, newPrice float64) error {
	result := r.db.WithContext(ctx).Model(&models.Signal{}).Where("id = ?", signalID).Update("current_price", newPrice)
	if result.Error != nil {
		return fmt.Errorf("failed to update signal current price: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SignalRepository) UpdateSignalStatus(ctx context.Context, signalID uint, newStatus string) error {
	result := r.db.WithContext(ctx).Model(&models.Signal{}).Where("id = ?", signalID).Update("status", newStatus)
	if result.Error != nil {
		return fmt.Errorf("failed to update signal status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
