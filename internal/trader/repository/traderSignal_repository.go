package repository

import (
	"context"
	"fmt"

	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ISignalRepository interface {
	CreateSignal(ctx context.Context, signal *models.Signal) (*models.Signal, error)
	GetAllSignals(ctx context.Context) ([]models.Signal, error)
	GetMarketDataBySymbol(ctx context.Context, symbol string) (*models.MarketData, error)
	UpdateSignalCurrentPrice(ctx context.Context, signalID uint, price float64) error
	UpdateSignalStatus(ctx context.Context, signalID uint, status string) error
	GetPendingSignals(ctx context.Context) ([]models.Signal, error)
	GetActiveSignals(ctx context.Context) ([]models.Signal, error)
	GetSignalByID(ctx context.Context, id uint) (*models.Signal, error)
	UpdateSignal(ctx context.Context, signal *models.Signal) error
	DeleteSignal(ctx context.Context, id uint) error
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
		return nil, err
	}
	return signals, nil
}

func (r *SignalRepository) GetMarketDataBySymbol(ctx context.Context, symbol string) (*models.MarketData, error) {
	var md models.MarketData
	if err := r.db.WithContext(ctx).Where("symbol = ?", symbol).First(&md).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &md, nil
}

func (r *SignalRepository) UpdateSignalCurrentPrice(ctx context.Context, signalID uint, price float64) error {
	return r.db.WithContext(ctx).Model(&models.Signal{}).Where("id = ?", signalID).Update("current_price", price).Error
}

func (r *SignalRepository) UpdateSignalStatus(ctx context.Context, signalID uint, status string) error {
	return r.db.WithContext(ctx).Model(&models.Signal{}).Where("id = ?", signalID).Update("status", status).Error
}

func (r *SignalRepository) GetPendingSignals(ctx context.Context) ([]models.Signal, error) {
	var signals []models.Signal
	err := r.db.WithContext(ctx).Where("status = ?", "Pending").Find(&signals).Error
	return signals, err
}

func (r *SignalRepository) GetActiveSignals(ctx context.Context) ([]models.Signal, error) {
	var signals []models.Signal
	err := r.db.WithContext(ctx).Where("status = ?", "Active").Find(&signals).Error
	return signals, err
}

func (r *SignalRepository) GetSignalByID(ctx context.Context, id uint) (*models.Signal, error) {
	var signal models.Signal
	if err := r.db.WithContext(ctx).First(&signal, id).Error; err != nil {
		return nil, err
	}
	return &signal, nil
}

func (r *SignalRepository) UpdateSignal(ctx context.Context, signal *models.Signal) error {
	return r.db.WithContext(ctx).Save(signal).Error
}

func (r *SignalRepository) DeleteSignal(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Signal{}, id).Error
}
