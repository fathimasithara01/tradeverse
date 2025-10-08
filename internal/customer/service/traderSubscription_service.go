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
	SubscribeCustomer(req *models.TraderSubscriptionRequest) (*models.TraderSubscriptionResponse, error)
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

func (s *traderSubscriptionService) SubscribeCustomer(req *models.TraderSubscriptionRequest) (*models.TraderSubscriptionResponse, error) {
	ctx := context.Background()

	// 1️⃣ Validate and fetch plan
	plan, err := s.subscriptionPlanRepo.GetSubscriptionPlanByID(ctx, req.TraderSubscriptionPlanID)
	if err != nil {
		return nil, fmt.Errorf("invalid subscription plan ID: %w", err)
	}

	// 2️⃣ Validate customer wallet
	customerWallet, err := s.customerWalletRepo.GetUserWallet(req.CustomerID)
	if err != nil {
		return nil, errors.New("customer wallet not found")
	}

	if customerWallet.Balance < plan.Price {
		return nil, errors.New("insufficient wallet balance")
	}

	// 3️⃣ Deduct from customer wallet
	prevBalance := customerWallet.Balance
	customerWallet.Balance -= plan.Price
	customerWallet.LastUpdated = time.Now()
	if err := s.customerWalletRepo.UpdateWallet(customerWallet); err != nil {
		return nil, err
	}

	// 4️⃣ Credit admin wallet (assuming admin user_id = 1)
	adminWallet, err := s.customerWalletRepo.GetUserWallet(1)
	if err != nil {
		adminWallet, _ = s.customerWalletRepo.GetOrCreateWallet(1)
	}
	adminPrevBalance := adminWallet.Balance
	adminWallet.Balance += plan.Price
	adminWallet.LastUpdated = time.Now()
	if err := s.customerWalletRepo.UpdateWallet(adminWallet); err != nil {
		return nil, err
	}

	// 5️⃣ Create subscription record
	transactionID := fmt.Sprintf("ADMIN_SUB_%d_%d_%d", req.CustomerID, plan.ID, time.Now().UnixNano())

	sub := &models.TraderSubscription{
		UserID:                   req.CustomerID,
		TraderID:                 req.TraderID,
		TraderSubscriptionPlanID: req.TraderSubscriptionPlanID,
		StartDate:                time.Now(),
		EndDate:                  time.Now().AddDate(0, 0, int(plan.Duration)), // Duration in days
		IsActive:                 true,
		PaymentStatus:            "SUCCESS",
		AmountPaid:               plan.Price,
		TraderShare:              0,
		AdminCommission:          plan.Price,
		TransactionID:            transactionID,
	}

	if err := s.traderSubRepo.CreateTraderSubscription(ctx, sub); err != nil {
		return nil, err
	}

	// 6️⃣ Record wallet transactions
	customerTx := &models.WalletTransaction{
		WalletID:        customerWallet.ID,
		UserID:          req.CustomerID,
		Name:            "Admin Subscription Plan",
		Type:            models.TxTypeSubscription,
		TransactionType: models.TxTypeSubscription,
		Amount:          plan.Price,
		Currency:        customerWallet.Currency,
		Status:          models.TxStatusSuccess,
		BalanceBefore:   prevBalance,
		BalanceAfter:    customerWallet.Balance,
		TransactionID:   transactionID,
	}
	_ = s.customerWalletRepo.CreateTransaction(customerTx)

	adminTx := &models.WalletTransaction{
		WalletID:        adminWallet.ID,
		UserID:          1,
		Name:            "Customer Subscription Payment",
		Type:            models.TxTypeDeposit,
		TransactionType: models.TxTypeDeposit,
		Amount:          plan.Price,
		Currency:        adminWallet.Currency,
		Status:          models.TxStatusSuccess,
		BalanceBefore:   adminPrevBalance,
		BalanceAfter:    adminWallet.Balance,
		TransactionID:   transactionID,
	}
	_ = s.customerWalletRepo.CreateTransaction(adminTx)

	// 7️⃣ Build and return response
	resp := &models.TraderSubscriptionResponse{
		TraderName:      "Admin",
		PlanName:        plan.Name,
		AmountPaid:      plan.Price,
		TraderShare:     0,
		AdminCommission: plan.Price,
		PaymentStatus:   sub.PaymentStatus,
		TransactionID:   transactionID,
		StartDate:       sub.StartDate.Format(time.RFC3339),
		EndDate:         sub.EndDate.Format(time.RFC3339),
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
