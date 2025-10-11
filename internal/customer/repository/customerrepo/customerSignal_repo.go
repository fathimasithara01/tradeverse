package customerrepo

// import (
// 	"context"
// 	"errors"
// 	"fmt"

// 	"github.com/fathimasithara01/tradeverse/pkg/models"
// 	"gorm.io/gorm"
// )

// type ICustomerSignalRepository interface {
// 	GetSignalsByTraderID(ctx context.Context, traderID uint) ([]models.Signal, error)
// 	GetSignalByID(ctx context.Context, signalID uint) (*models.Signal, error)
// }

// type customerSignalRepository struct {
// 	db *gorm.DB
// }

// func NewCustomerSignalRepository(db *gorm.DB) ICustomerSignalRepository {
// 	return &customerSignalRepository{db: db}
// }

// func (r *customerSignalRepository) GetSignalsByTraderID(ctx context.Context, traderID uint) ([]models.Signal, error) {
// 	var signals []models.Signal
// 	// Only fetch active signals that belong to the specified trader
// 	err := r.db.WithContext(ctx).
// 		Where("trader_id = ? AND status IN (?)", traderID, []string{"Active", "Pending"}). // Assuming signals have a status
// 		Order("created_at DESC").
// 		Find(&signals).Error
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get signals for trader %d: %w", traderID, err)
// 	}
// 	return signals, nil
// }

// func (r *customerSignalRepository) GetSignalByID(ctx context.Context, signalID uint) (*models.Signal, error) {
// 	var signal models.Signal
// 	err := r.db.WithContext(ctx).First(&signal, signalID).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, gorm.ErrRecordNotFound // Let the service layer handle custom errors
// 		}
// 		return nil, fmt.Errorf("failed to get signal by ID %d: %w", signalID, err)
// 	}
// 	return &signal, nil
// }
