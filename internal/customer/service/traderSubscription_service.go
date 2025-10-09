package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	adminrepo "github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	walletrepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/walletrepo"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrTraderNotFound           = errors.New("trader not found")
	ErrSubscriptionPlanNotFound = errors.New("subscription plan not found")
	ErrCustomerWalletNotFound   = errors.New("customer wallet not found")
	ErrAlreadySubscribed        = errors.New("customer already has an active subscription to this trader")
	ErrNotAuthorized            = errors.New("not authorized to access this subscription")
	ErrSubscriptionNotFound     = errors.New("trader subscription not found")
	ErrTraderWalletNotFound     = errors.New("trader wallet not found")
)

type ITraderSubscriptionService interface {
	// SubscribeCustomerToTrader(ctx context.Context, req models.TraderSubscriptionRequest) (*models.TraderSubscriptionResponse, error)
	SubscribeCustomer(req models.TraderSubscriptionRequest) (*models.TraderSubscriptionResponse, error)
	DeactivateExpiredTraderSubscriptions(ctx context.Context) error
	GetCustomerTraderSubscriptions(ctx context.Context, customerID uint) ([]models.TraderSubscription, error)
	GetCustomerTraderSubscriptionByID(ctx context.Context, customerID, subscriptionID uint) (*models.TraderSubscription, error)
}

type traderSubscriptionService struct {
	db                   *gorm.DB
	customerWalletRepo   walletrepo.WalletRepository
	adminWalletRepo      adminrepo.IAdminWalletRepository
	traderSubRepo        customerrepo.ITraderSubscriptionRepository
	userRepo             customerrepo.IUserRepository
	subscriptionPlanRepo customerrepo.ISubscriptionPlanRepository
}

func NewTraderSubscriptionService(
	db *gorm.DB,
	customerWalletRepo walletrepo.WalletRepository,
	adminWalletRepo adminrepo.IAdminWalletRepository,
	traderSubRepo customerrepo.ITraderSubscriptionRepository,
	userRepo customerrepo.IUserRepository,
	subscriptionPlanRepo customerrepo.ISubscriptionPlanRepository,
) ITraderSubscriptionService {
	return &traderSubscriptionService{
		db:                   db,
		customerWalletRepo:   customerWalletRepo,
		adminWalletRepo:      adminWalletRepo,
		traderSubRepo:        traderSubRepo,
		userRepo:             userRepo,
		subscriptionPlanRepo: subscriptionPlanRepo,
	}
}

func (s *traderSubscriptionService) SubscribeCustomer(req models.TraderSubscriptionRequest) (*models.TraderSubscriptionResponse, error) {
	plan, err := s.traderSubRepo.GetPlanByID(req.TraderSubscriptionPlanID)
	if err != nil {
		return nil, ErrSubscriptionPlanNotFound
	}

	// 2️⃣ Check existing active subscription
	existing, _ := s.traderSubRepo.CheckExistingSubscription(req.CustomerID, req.TraderID)
	if existing != nil {
		return nil, ErrAlreadySubscribed
	}

	startDate := time.Now()
	endDate := startDate.Add(plan.Duration)
	adminCommission := plan.Price * plan.CommissionRate
	traderShare := plan.Price - adminCommission
	transactionID := fmt.Sprintf("TXN-%d-%d-%d", req.CustomerID, req.TraderID, time.Now().Unix())

	// 3️⃣ Perform transaction: debit customer, credit admin & trader, create subscription
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Customer wallet
		customerWallet, err := s.customerWalletRepo.GetUserWallet(req.CustomerID)
		if err != nil {
			return fmt.Errorf("customer wallet not found: %w", err)
		}
		if customerWallet.Balance < plan.Price {
			return errors.New("insufficient customer wallet balance")
		}
		customerWallet.Balance -= plan.Price
		if err := s.customerWalletRepo.UpdateWallet(customerWallet); err != nil {
			return err
		}

		// Admin wallet
		adminWallet, err := s.adminWalletRepo.GetAdminWallet()
		if err != nil {
			return fmt.Errorf("admin wallet not found: %w", err)
		}
		adminWallet.Balance += adminCommission
		if err := s.adminWalletRepo.UpdateWalletBalance(tx, adminWallet); err != nil {
			return err
		}

		// Trader wallet
		traderWallet, err := s.customerWalletRepo.GetUserWallet(req.TraderID)
		if err != nil {
			return fmt.Errorf("trader wallet not found: %w", err)
		}
		traderWallet.Balance += traderShare
		if err := s.customerWalletRepo.UpdateWallet(traderWallet); err != nil {
			return err
		}

		// Create subscription
		sub := &models.TraderSubscription{
			UserID:                   req.CustomerID,
			TraderID:                 req.TraderID,
			TraderSubscriptionPlanID: plan.ID,
			StartDate:                startDate,
			EndDate:                  endDate,
			IsActive:                 true,
			PaymentStatus:            "Paid",
			AmountPaid:               plan.Price,
			TraderShare:              traderShare,
			AdminCommission:          adminCommission,
			TransactionID:            transactionID,
		}
		if err := s.traderSubRepo.CreateSubscription(sub); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 4️⃣ Return response
	resp := &models.TraderSubscriptionResponse{
		TraderName:      fmt.Sprintf("Trader %d", req.TraderID),
		PlanName:        plan.Name,
		AmountPaid:      plan.Price,
		TraderShare:     traderShare,
		AdminCommission: adminCommission,
		PaymentStatus:   "Paid",
		TransactionID:   transactionID,
		StartDate:       startDate.Format("2006-01-02"),
		EndDate:         endDate.Format("2006-01-02"),
		IsActive:        true,
	}
	return resp, nil
}

// func (s *traderSubscriptionService) SubscribeCustomerToTrader(ctx context.Context, req models.TraderSubscriptionRequest) (*models.TraderSubscriptionResponse, error) {
// 	// 1. Validate customer and trader
// 	customer, err := s.userRepo.GetUserByID(ctx, req.CustomerID)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid customer: %w", err)
// 	}

// 	trader, err := s.userRepo.GetUserByID(ctx, req.TraderID)
// 	if err != nil {
// 		return nil, ErrTraderNotFound
// 	}

// 	// 2. Validate Subscription Plan
// 	subPlan, err := s.subscriptionPlanRepo.GetSubscriptionPlanByID(ctx, req.TraderSubscriptionPlanID)
// 	if err != nil {
// 		return nil, ErrSubscriptionPlanNotFound
// 	}

// 	if subPlan.UserID == nil || *subPlan.UserID != req.TraderID {
// 		return nil, errors.New("subscription plan does not belong to this trader")
// 	}

// 	if !subPlan.IsActive {
// 		return nil, errors.New("subscription plan is not active")
// 	}

// 	// 3. Check for existing subscription
// 	existingSub, err := s.traderSubRepo.GetActiveTraderSubscriptionForCustomer(ctx, req.CustomerID, req.TraderID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if existingSub != nil {
// 		return nil, ErrAlreadySubscribed
// 	}

// 	// 4. Wallet validation
// 	customerWallet, err := s.customerWalletRepo.GetUserWallet(req.CustomerID)
// 	if err != nil {
// 		return nil, ErrCustomerWalletNotFound
// 	}
// 	if customerWallet.Balance < subPlan.Price {
// 		return nil, ErrInsufficientFunds
// 	}

// 	adminWallet, err := s.adminWalletRepo.GetAdminWallet()
// 	if err != nil {
// 		return nil, ErrAdminWalletNotFound
// 	}

// 	traderWallet, err := s.customerWalletRepo.GetUserWallet(req.TraderID)
// 	if err != nil {
// 		return nil, ErrTraderWalletNotFound
// 	}

// 	const adminCommissionRate = 0.10
// 	adminCommission := subPlan.Price * adminCommissionRate
// 	traderAmount := subPlan.Price - adminCommission

// 	var newSubscription models.TraderSubscription

// 	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		// Debit Customer
// 		if err := s.customerWalletRepo.DebitWallet(tx, customerWallet.ID, subPlan.Price,
// 			models.TxTypeSubscription,
// 			fmt.Sprintf("SUB_%d_%d", req.TraderID, req.CustomerID),
// 			fmt.Sprintf("Subscription to trader %s", trader.Name),
// 		); err != nil {
// 			return err
// 		}

// 		// Credit Admin
// 		adminWallet.Balance += adminCommission
// 		if err := s.adminWalletRepo.UpdateWalletBalance(tx, adminWallet); err != nil {
// 			return err
// 		}

// 		// Credit Trader
// 		if err := s.customerWalletRepo.CreditWallet(tx, traderWallet.ID, traderAmount,
// 			models.TxTypeSubscription,
// 			fmt.Sprintf("REV_%d_%d", req.TraderID, req.CustomerID),
// 			fmt.Sprintf("Trader revenue from customer %s", customer.Name),
// 		); err != nil {
// 			return err
// 		}

// 		// Create Subscription
// 		newSubscription = models.TraderSubscription{
// 			UserID:                   req.CustomerID,
// 			TraderID:                 req.TraderID,
// 			TraderSubscriptionPlanID: req.TraderSubscriptionPlanID,
// 			StartDate:                time.Now(),
// 			EndDate:                  time.Now().Add(subPlan.Duration),
// 			IsActive:                 true,
// 			PaymentStatus:            string(models.TxStatusSuccess),
// 			AmountPaid:               subPlan.Price,
// 		}
// 		return s.traderSubRepo.CreateTraderSubscription(ctx, &newSubscription)
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("transaction failed: %w", err)
// 	}

// 	return &models.TraderSubscriptionResponse{
// 		TraderSubscriptionID: newSubscription.ID,
// 		Message:              "Successfully subscribed to trader",
// 		Status:               string(models.TxStatusSuccess),
// 	}, nil
// }

func (s *traderSubscriptionService) DeactivateExpiredTraderSubscriptions(ctx context.Context) error {
	return s.traderSubRepo.DeactivateExpiredTraderSubscriptions(ctx)
}

func (s *traderSubscriptionService) GetCustomerTraderSubscriptions(ctx context.Context, customerID uint) ([]models.TraderSubscription, error) {
	return s.traderSubRepo.GetCustomerTraderSubscriptions(ctx, customerID)
}

func (s *traderSubscriptionService) GetCustomerTraderSubscriptionByID(ctx context.Context, customerID, subscriptionID uint) (*models.TraderSubscription, error) {
	sub, err := s.traderSubRepo.GetTraderSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	if sub.UserID != customerID {
		return nil, ErrNotAuthorized
	}
	return sub, nil
}
