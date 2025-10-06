package service

import (
	"context"
	"errors"
	"fmt"

	tradersignalrepo "github.com/fathimasithara01/tradeverse/internal/trader/repository" // Reuse trader's signal repository
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

var (
	ErrNotSubscribed = errors.New("customer is not subscribed to this trader")
)

type ICustomerSignalService interface {
	GetTraderSignalsForCustomer(ctx context.Context, customerID, traderID uint) ([]models.Signal, error)
	GetSignalCardForCustomer(ctx context.Context, customerID, traderID, signalID uint) (*models.Signal, error)
}

type customerSignalService struct {
	signalRepo            tradersignalrepo.ISignalRepository
	traderSubscriptionSvc ITraderSubscriptionService
}

func NewCustomerSignalService(signalRepo tradersignalrepo.ISignalRepository, traderSubscriptionSvc ITraderSubscriptionService) ICustomerSignalService {
	return &customerSignalService{
		signalRepo:            signalRepo,
		traderSubscriptionSvc: traderSubscriptionSvc,
	}
}

// GetSignalCardForCustomer retrieves a single signal from a specific trader for a subscribed customer.
func (s *customerSignalService) GetSignalCardForCustomer(ctx context.Context, customerID, traderID, signalID uint) (*models.Signal, error) {
	// 1. Check if the customer has an active subscription to this trader
	isSubscribed, err := s.traderSubscriptionSvc.IsCustomerSubscribedToTrader(ctx, customerID, traderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check subscription status: %w", err)
	}
	if !isSubscribed {
		return nil, ErrNotSubscribed
	}

	// 2. If subscribed, retrieve the specific signal
	signal, err := s.signalRepo.GetSignalByID(ctx, signalID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve signal: %w", err)
	}

	// 3. Ensure the signal actually belongs to the specified trader
	if signal.TraderID == nil || *signal.TraderID != traderID {
		return nil, errors.New("signal does not belong to the specified trader")
	}

	return signal, nil
}

// internal/customer/service/signal_service.go (modified)
func (s *customerSignalService) GetTraderSignalsForCustomer(ctx context.Context, customerID, traderID uint) ([]models.Signal, error) {
	// ... (subscription check remains the same) ...

	isSubscribed, err := s.traderSubscriptionSvc.IsCustomerSubscribedToTrader(ctx, customerID, traderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check subscription status: %w", err)
	}
	if !isSubscribed {
		return nil, ErrNotSubscribed
	}

	// Use the new repository method
	signals, err := s.signalRepo.GetSignalsByTraderID(ctx, traderID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve signals for trader: %w", err)
	}

	return signals, nil
}
