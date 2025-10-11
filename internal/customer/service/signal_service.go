package service

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"time"

// 	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
// 	"github.com/fathimasithara01/tradeverse/pkg/models"
// 	"gorm.io/gorm"
// )

// var (
// 	ErrNotSubscribed        = errors.New("customer is not subscribed to this trader")
// 	ErrSignalNotFound       = errors.New("signal not found")
// 	ErrSignalTraderMismatch = errors.New("signal does not belong to the specified trader")
// )

// type ICustomerSignalService interface {
// 	GetTraderSignalsForCustomer(ctx context.Context, customerID, traderID uint) ([]models.Signal, error)
// 	GetSignalCardForCustomer(ctx context.Context, customerID, traderID, signalID uint) (*models.Signal, error)
// 	// Add other customer-specific signal methods here
// }

// type customerSignalService struct {
// 	signalRepo    customerrepo.ICustomerSignalRepository
// 	traderSubRepo customerrepo.ITraderSubscriptionRepository
// }

// func NewCustomerSignalService(
// 	signalRepo customerrepo.ICustomerSignalRepository,
// 	traderSubRepo customerrepo.ITraderSubscriptionRepository,
// ) ICustomerSignalService {
// 	return &customerSignalService{
// 		signalRepo:    signalRepo,
// 		traderSubRepo: traderSubRepo,
// 	}
// }

// // GetTraderSignalsForCustomer retrieves all signals for a specific trader,
// // but only if the customer has an active subscription to that trader.
// func (s *customerSignalService) GetTraderSignalsForCustomer(ctx context.Context, customerID, traderID uint) ([]models.Signal, error) {
// 	// 1. Check for active subscription
// 	activeSub, err := s.traderSubRepo.GetActiveTraderSubscriptionForCustomer(ctx, customerID, traderID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to check subscription status: %w", err)
// 	}
// 	// Check if activeSub is nil OR if it's not active OR if its end date is in the past
// 	if activeSub == nil || !activeSub.IsActive || activeSub.EndDate.Before(time.Now()) {
// 		return nil, ErrNotSubscribed
// 	}

// 	// 2. If subscribed, retrieve signals for the trader
// 	signals, err := s.signalRepo.GetSignalsByTraderID(ctx, traderID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to retrieve signals for trader %d: %w", traderID, err)
// 	}

// 	return signals, nil
// }

// // GetSignalCardForCustomer retrieves a single signal card,
// // but only if the customer has an active subscription to the trader and the signal belongs to that trader.
// func (s *customerSignalService) GetSignalCardForCustomer(ctx context.Context, customerID, traderID, signalID uint) (*models.Signal, error) {
// 	// 1. Check for active subscription
// 	activeSub, err := s.traderSubRepo.GetActiveTraderSubscriptionForCustomer(ctx, customerID, traderID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to check subscription status: %w", err)
// 	}
// 	if activeSub == nil || !activeSub.IsActive || activeSub.EndDate.Before(time.Now()) {
// 		return nil, ErrNotSubscribed
// 	}

// 	// 2. If subscribed, retrieve the specific signal
// 	signal, err := s.signalRepo.GetSignalByID(ctx, signalID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrSignalNotFound
// 		}
// 		return nil, fmt.Errorf("failed to retrieve signal %d: %w", signalID, err)
// 	}

// 	// 3. Ensure the signal belongs to the specified trader
// 	if signal.TraderID != traderID {
// 		return nil, ErrSignalTraderMismatch
// 	}

// 	return signal, nil
// }
