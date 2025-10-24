package service

import (
	"context"
	"fmt"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ITraderSubscriptionService interface {
	CreateTraderSubscriptionPlan(ctx context.Context, traderID uint, input models.CreateTraderSubscriptionPlanInput) (*models.TraderSignalSubscriptionPlan, error)
	GetTraderSubscriptionPlanByID(ctx context.Context, planID uint) (*models.TraderSignalSubscriptionPlan, error)
	GetMyTraderSubscriptionPlans(ctx context.Context, traderID uint) ([]models.TraderSignalSubscriptionPlan, error)
	UpdateTraderSubscriptionPlan(ctx context.Context, traderID uint, planID uint, input models.CreateTraderSubscriptionPlanInput) (*models.TraderSignalSubscriptionPlan, error)
	DeleteTraderSubscriptionPlan(ctx context.Context, traderID uint, planID uint) error

	SubscribeToTraderPlan(ctx context.Context, customerID uint, traderID uint, planID uint) error

	GetAllTraderUpgradePlans(ctx context.Context) ([]models.AdminTraderSubscriptionPlan, error)
	SubscribeToTraderUpgradePlan(ctx context.Context, userID uint, planID uint) error
}

type TraderSubscriptionService struct {
	repo repository.ITraderSubscriptionRepository
	db   *gorm.DB // For transactions
}

func NewTraderSubscriptionService(repo repository.ITraderSubscriptionRepository, db *gorm.DB) ITraderSubscriptionService {
	return &TraderSubscriptionService{repo: repo, db: db}
}

func (s *TraderSubscriptionService) CreateTraderSubscriptionPlan(ctx context.Context, traderID uint, input models.CreateTraderSubscriptionPlanInput) (*models.TraderSignalSubscriptionPlan, error) {

	user, err := s.repo.GetUserByID(ctx, traderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user %d: %w", traderID, err)
	}
	if user.Role != models.RoleTrader {
		return nil, fmt.Errorf("user is not a trader and cannot create subscription plans")
	}

	traderShareAmount := input.Price * (1 - (input.AdminCommissionPercentage / 100.0))
	if traderShareAmount < 0 {
		traderShareAmount = 0 // Ensure it doesn't go negative
	}

	plan := &models.TraderSignalSubscriptionPlan{
		TraderID:        traderID,
		Name:            input.Name,
		Description:     input.Description,
		Price:           input.Price,
		Currency:        input.Currency,
		DurationDays:    input.DurationDays,
		IsActive:        true, // Default to active upon creation
		AdminCommission: input.AdminCommissionPercentage,
		// --- SET THE NEW FIELD ---
		TraderShare: traderShareAmount,
		// --- END SET ---
	}

	return s.repo.CreateTraderSubscriptionPlan(ctx, plan)
}
func (s *TraderSubscriptionService) GetTraderSubscriptionPlanByID(ctx context.Context, planID uint) (*models.TraderSignalSubscriptionPlan, error) {
	return s.repo.GetTraderSubscriptionPlanByID(ctx, planID)
}

func (s *TraderSubscriptionService) GetMyTraderSubscriptionPlans(ctx context.Context, traderID uint) ([]models.TraderSignalSubscriptionPlan, error) {
	// This method needs to query the repository for plans belonging to this traderID
	// For example:
	plans, err := s.repo.GetTraderSubscriptionPlansByTraderID(ctx, traderID) // <--- This repository method is key
	if err != nil {
		return nil, fmt.Errorf("failed to get trader subscription plans for trader %d: %w", traderID, err)
	}
	return plans, nil
}
func (s *TraderSubscriptionService) UpdateTraderSubscriptionPlan(ctx context.Context, traderID uint, planID uint, input models.CreateTraderSubscriptionPlanInput) (*models.TraderSignalSubscriptionPlan, error) {
	existingPlan, err := s.repo.GetTraderSubscriptionPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	if existingPlan.TraderID != traderID {
		return nil, fmt.Errorf("unauthorized: plan does not belong to this trader")
	}

	existingPlan.Name = input.Name
	existingPlan.Description = input.Description
	existingPlan.Price = input.Price
	existingPlan.Currency = input.Currency
	existingPlan.DurationDays = input.DurationDays
	existingPlan.AdminCommission = input.AdminCommissionPercentage
	// Optionally update IsActive, but generally controlled by a separate endpoint

	if err := s.repo.UpdateTraderSubscriptionPlan(ctx, existingPlan); err != nil {
		return nil, err
	}
	return existingPlan, nil
}

func (s *TraderSubscriptionService) DeleteTraderSubscriptionPlan(ctx context.Context, traderID uint, planID uint) error {
	return s.repo.DeleteTraderSubscriptionPlan(ctx, planID, traderID)
}

// --- Customer subscription to a specific Trader's plan ---
func (s *TraderSubscriptionService) SubscribeToTraderPlan(ctx context.Context, customerID uint, traderID uint, planID uint) error {
	// 1. Fetch the Trader's plan
	plan, err := s.repo.GetTraderSubscriptionPlanByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("trader subscription plan not found: %w", err)
	}
	if !plan.IsActive {
		return fmt.Errorf("trader subscription plan is not active")
	}
	if plan.TraderID != traderID {
		return fmt.Errorf("trader subscription plan does not belong to the specified trader")
	}
	if customerID == traderID {
		return fmt.Errorf("customer cannot subscribe to their own plan")
	}

	// 2. Check if customer is already subscribed to this specific plan and it's active
	alreadySubscribed, err := s.repo.CheckIfCustomerIsSubscribedToTraderPlan(ctx, customerID, planID)
	if err != nil {
		return fmt.Errorf("failed to check existing subscription: %w", err)
	}
	if alreadySubscribed {
		return fmt.Errorf("user is already subscribed to this plan")
	}

	// 3. Start a database transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 4. Deduct money from customer's wallet
	customerWallet, err := s.repo.GetUserWallet(ctx, customerID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get customer wallet: %w", err)
	}

	if customerWallet.Balance < plan.Price {
		tx.Rollback()
		return fmt.Errorf("insufficient funds in wallet. current balance: %.2f, required: %.2f", customerWallet.Balance, plan.Price)
	}

	if err := s.repo.UpdateWalletBalance(ctx, customerID, -plan.Price, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deduct funds from customer wallet: %w", err)
	}
	customerBalanceAfter := customerWallet.Balance - plan.Price

	// Calculate commission and actual amount for trader
	adminCommissionAmount := plan.Price * (plan.AdminCommission / 100.0)
	traderReceiveAmount := plan.Price - adminCommissionAmount

	// 5. Credit commission to Admin wallet
	adminWallet, err := s.repo.GetAdminWallet(ctx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get admin wallet: %w", err)
	}
	if err := s.repo.UpdateWalletBalance(ctx, adminWallet.UserID, adminCommissionAmount, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to credit admin commission: %w", err)
	}
	adminBalanceAfter := adminWallet.Balance + adminCommissionAmount

	// 6. Credit remaining amount to Trader's wallet
	traderWallet, err := s.repo.GetUserWallet(ctx, traderID) // Re-using GetUserWallet for trader
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get trader wallet: %w", err)
	}
	if err := s.repo.UpdateWalletBalance(ctx, traderID, traderReceiveAmount, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to credit trader wallet: %w", err)
	}
	traderBalanceAfter := traderWallet.Balance + traderReceiveAmount

	// 7. Create wallet transactions
	customerTx := models.WalletTransaction{
		WalletID:        customerWallet.ID,
		UserID:          customerID,
		Type:            models.TxTypeSignalPayment,
		TransactionType: models.TxTypeDebit,
		Name:            "Debit for Trader Subscription",
		Amount:          plan.Price,
		Currency:        plan.Currency,
		Status:          models.TxStatusSuccess,
		Description:     fmt.Sprintf("Subscription to trader %d's plan '%s'", traderID, plan.Name),
		BalanceBefore:   customerWallet.Balance,
		BalanceAfter:    customerBalanceAfter,
		ReferenceID:     fmt.Sprintf("TRADER_SUB_%d_PLAN_%d", customerID, plan.ID),
		TransactionID:   fmt.Sprintf("TRADER_SUB_CUST_%d_%d_%d", customerID, plan.ID, time.Now().UnixNano()),
	}

	if err := s.repo.CreateWalletTransaction(ctx, &customerTx, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record customer transaction for trader subscription: %w", err)
	}

	adminCommissionTx := models.WalletTransaction{
		WalletID:        adminWallet.ID,
		UserID:          adminWallet.UserID,
		Type:            models.TxTypeCommission,
		TransactionType: models.TxTypeCredit,
		Name:            "Credit from Trader Plan Commission",
		Amount:          adminCommissionAmount,
		Currency:        plan.Currency,
		Status:          models.TxStatusSuccess,
		Description:     fmt.Sprintf("Commission from customer %d subscribing to trader %d's plan '%s'", customerID, traderID, plan.Name),
		BalanceBefore:   adminWallet.Balance,
		BalanceAfter:    adminBalanceAfter,
		ReferenceID:     fmt.Sprintf("TRADER_SUB_ADMIN_COMM_%d_PLAN_%d", customerID, plan.ID),
		TransactionID:   fmt.Sprintf("TRADER_SUB_ADM_%d_%d_%d", customerID, plan.ID, time.Now().UnixNano()),
	}
	if err := s.repo.CreateWalletTransaction(ctx, &adminCommissionTx, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record admin commission transaction: %w", err)
	}

	traderCreditTx := models.WalletTransaction{
		WalletID:        traderWallet.ID,
		UserID:          traderID,
		Type:            models.TxTypeSignalPayment,
		TransactionType: models.TxTypeCredit,
		Name:            "Credit from Customer Subscription",
		Amount:          traderReceiveAmount,
		Currency:        plan.Currency,
		Status:          models.TxStatusSuccess,
		Description:     fmt.Sprintf("Credit from customer %d subscribing to plan '%s'", customerID, plan.Name),
		BalanceBefore:   traderWallet.Balance,
		BalanceAfter:    traderBalanceAfter,
		ReferenceID:     fmt.Sprintf("TRADER_SUB_TRADER_CREDIT_%d_PLAN_%d", customerID, plan.ID),
		TransactionID:   fmt.Sprintf("TRADER_SUB_TRD_%d_%d_%d", customerID, plan.ID, time.Now().UnixNano()),
	}
	if err := s.repo.CreateWalletTransaction(ctx, &traderCreditTx, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record trader credit transaction: %w", err)
	}

	// 8. Create CustomerTraderSubscription record
	startDate := time.Now()
	endDate := startDate.Add(time.Duration(plan.DurationDays) * 24 * time.Hour)

	customerTraderSubscription := &models.CustomerTraderSignalSubscription{
		CustomerID:               customerID,
		TraderID:                 traderID,
		TraderSubscriptionPlanID: plan.ID,
		StartDate:                startDate,
		EndDate:                  endDate,
		IsActive:                 true,
		TransactionID:            customerTx.ID,
	}

	if err := s.repo.CreateCustomerTraderSubscription(ctx, customerTraderSubscription); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create customer trader subscription record: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// --- Trader Upgrade Subscription related (for users to take a plan to become trader) ---

func (s *TraderSubscriptionService) GetAllTraderUpgradePlans(ctx context.Context) ([]models.AdminTraderSubscriptionPlan, error) {
	return s.repo.GetAllTraderUpgradePlans(ctx)
}

func (s *TraderSubscriptionService) SubscribeToTraderUpgradePlan(ctx context.Context, userID uint, planID uint) error {
	plan, err := s.repo.GetTraderUpgradePlanByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("invalid subscription plan: %w", err)
	}
	if !plan.IsActive {
		return fmt.Errorf("subscription plan is not active")
	}
	if !plan.IsUpgradeToTrader {
		return fmt.Errorf("this plan is not for upgrading to a trader role")
	}

	// Check if user is already an active trader via this (or any other relevant) upgrade plan
	isAlreadyTrader, err := s.repo.CheckIfUserIsActiveTrader(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user's trader status: %w", err)
	}
	if isAlreadyTrader {
		// More specific check: if already subscribed to *this specific* upgrade plan
		existingUpgradeSub, err := s.repo.GetUserActiveUpgradeSubscription(ctx, userID, planID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to check for existing upgrade subscription: %w", err)
		}
		if existingUpgradeSub != nil && existingUpgradeSub.IsActive && existingUpgradeSub.EndDate.After(time.Now()) {
			return fmt.Errorf("user is already an active trader with this upgrade plan")
		}
	}

	// Start a database transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Deduct money from user's wallet
	userWallet, err := s.repo.GetUserWallet(ctx, userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get user wallet: %w", err)
	}

	if userWallet.Balance < plan.Price {
		tx.Rollback()
		return fmt.Errorf("insufficient funds in wallet. current balance: %.2f, required: %.2f", userWallet.Balance, plan.Price)
	}

	if err := s.repo.UpdateWalletBalance(ctx, userID, -plan.Price, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deduct funds from user wallet: %w", err)
	}
	userBalanceAfter := userWallet.Balance - plan.Price

	// 2. Credit money to Admin wallet (full price, as this is an admin plan)
	adminWallet, err := s.repo.GetAdminWallet(ctx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get admin wallet: %w", err)
	}
	if err := s.repo.UpdateWalletBalance(ctx, adminWallet.UserID, plan.Price, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to credit admin wallet: %w", err)
	}
	adminBalanceAfter := adminWallet.Balance + plan.Price

	// 3. Create wallet transactions
	userTx := models.WalletTransaction{
		WalletID:        userWallet.ID,
		UserID:          userID,
		Type:            models.TxTypeSubscription,
		TransactionType: models.TxTypeDebit,
		Name:            "Debit for Trader Upgrade Subscription",
		Amount:          plan.Price,
		Currency:        plan.Currency,
		Status:          models.TxStatusSuccess,
		Description:     fmt.Sprintf("Subscription to admin plan '%s' to become a trader", plan.Name),
		BalanceBefore:   userWallet.Balance,
		BalanceAfter:    userBalanceAfter,
		ReferenceID:     fmt.Sprintf("ADMIN_SUB_%d_PLAN_%d", userID, plan.ID),
		TransactionID:   fmt.Sprintf("ADMIN_SUB_USER_%d_%d_%d", userID, plan.ID, time.Now().UnixNano()),
	}
	if err := s.repo.CreateWalletTransaction(ctx, &userTx, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record user transaction for admin subscription: %w", err)
	}

	adminTx := models.WalletTransaction{
		WalletID:        adminWallet.ID,
		UserID:          adminWallet.UserID,
		Type:            models.TxTypeSubscription,
		TransactionType: models.TxTypeCredit,
		Name:            "Credit from Trader Upgrade Subscription",
		Amount:          plan.Price,
		Currency:        plan.Currency,
		Status:          models.TxStatusSuccess,
		Description:     fmt.Sprintf("Credit from user %d subscribing to admin plan '%s'", userID, plan.Name),
		BalanceBefore:   adminWallet.Balance,
		BalanceAfter:    adminBalanceAfter,
		ReferenceID:     fmt.Sprintf("ADMIN_SUB_CREDIT_%d_PLAN_%d", userID, plan.ID),
		TransactionID:   fmt.Sprintf("ADMIN_SUB_ADMIN_%d_%d_%d", userID, plan.ID, time.Now().UnixNano()),
	}
	if err := s.repo.CreateWalletTransaction(ctx, &adminTx, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record admin transaction for admin subscription: %w", err)
	}

	// 4. Create UserSubscription record (for the admin plan)
	startDate := time.Now()
	endDate := startDate.Add(time.Duration(plan.Duration) * 24 * time.Hour)

	userSubscription := &models.UserSubscription{
		UserID:             userID,
		SubscriptionPlanID: plan.ID,
		StartDate:          startDate,
		EndDate:            endDate,
		IsActive:           true,
		TransactionID:      userTx.ID, // Link to the user's payment transaction
	}

	if err := s.repo.CreateUserSubscription(ctx, userSubscription); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create user subscription record: %w", err)
	}

	// 5. Update user's Role to Trader
	if err := s.repo.SetUserRole(ctx, userID, models.RoleTrader, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set user as trader: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
