package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	adminRepo "github.com/fathimasithara01/tradeverse/internal/admin/repository"
	adminService "github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ICustomerSubscriptionService interface {
	CreateSubscription(userID, planID uint, amount float64, transactionID string) (*models.CustomerToTraderSub, error)
	GetSubscriptionsByUserID(userID uint) ([]models.CustomerToTraderSub, error)
	CancelSubscription(userID, subscriptionID uint) error
	DeactivateExpiredTraderSubscriptions() error
}

type CustomerSubscriptionService struct {
	customerSubscriptionRepo  customerrepo.ICustomerSubscriptionRepository
	adminSubscriptionPlanRepo adminRepo.ISubscriptionPlanRepository
	adminWalletService        adminService.IAdminWalletService
	userRepo                  adminRepo.IUserRepository
	DB                        *gorm.DB
}

func NewCustomerSubscriptionService(
	customerSubscriptionRepo customerrepo.ICustomerSubscriptionRepository,
	adminSubscriptionPlanRepo adminRepo.ISubscriptionPlanRepository,
	adminWalletService adminService.IAdminWalletService,
	userRepo adminRepo.IUserRepository,
	db *gorm.DB,
) *CustomerSubscriptionService {
	return &CustomerSubscriptionService{
		customerSubscriptionRepo:  customerSubscriptionRepo,
		adminSubscriptionPlanRepo: adminSubscriptionPlanRepo,
		adminWalletService:        adminWalletService,
		userRepo:                  userRepo,
		DB:                        db,
	}
}

func (s *CustomerSubscriptionService) DeactivateExpiredTraderSubscriptions() error {
	log.Println("Running DeactivateExpiredTraderSubscriptions cron job...")

	activeSubscriptions, err := s.customerSubscriptionRepo.GetActiveTraderSubscriptions()
	if err != nil {
		return err
	}

	now := time.Now()
	var deactivatedCount int

	for _, sub := range activeSubscriptions {
		if sub.EndDate.Before(now) && sub.IsActive {
			log.Printf("Deactivating trader subscription ID: %d for user: %d, expired on: %s", sub.ID, sub.UserID, sub.EndDate.Format(time.RFC3339))
			sub.IsActive = false
			sub.DeactivatedAt = &now // Assuming a DeactivatedAt field
			err := s.customerSubscriptionRepo.UpdateTraderSubscription(&sub)
			if err != nil {
				log.Printf("Error deactivating trader subscription ID %d: %v", sub.ID, err)
			} else {
				deactivatedCount++
			}
		}
	}

	log.Printf("Deactivated %d expired trader subscriptions.", deactivatedCount)
	return nil
}

func (s *CustomerSubscriptionService) CreateSubscription(userID, planID uint, amount float64, transactionID string) (*models.CustomerToTraderSub, error) {
	var subscription *models.CustomerToTraderSub

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1️⃣ Get plan from admin repo
		plan, err := s.adminSubscriptionPlanRepo.GetSubscriptionPlanByID(planID)
		if err != nil {
			return fmt.Errorf("plan not found: %w", err)
		}

		// 2️⃣ Fetch wallets
		var customerWallet, adminWallet models.Wallet
		if err := tx.Where("user_id = ?", userID).First(&customerWallet).Error; err != nil {
			return fmt.Errorf("customer wallet not found: %w", err)
		}
		if err := tx.Where("user_id = ?", 1).First(&adminWallet).Error; err != nil { // admin ID = 1
			return fmt.Errorf("admin wallet not found: %w", err)
		}

		if customerWallet.Balance < amount {
			return fmt.Errorf("insufficient wallet balance")
		}

		// 3️⃣ Wallet Balances Before Transfer
		beforeCustomer := customerWallet.Balance
		beforeAdmin := adminWallet.Balance

		// 4️⃣ Perform Transfer
		customerWallet.Balance -= amount
		adminWallet.Balance += amount
		customerWallet.LastUpdated = time.Now()
		adminWallet.LastUpdated = time.Now()

		if err := tx.Save(&customerWallet).Error; err != nil {
			return fmt.Errorf("failed to update customer wallet: %w", err)
		}
		if err := tx.Save(&adminWallet).Error; err != nil {
			return fmt.Errorf("failed to update admin wallet: %w", err)
		}

		// 5️⃣ Create Transactions
		customerTx := models.WalletTransaction{
			WalletID:        customerWallet.ID,
			UserID:          userID,
			Type:            models.TxTypeSubscription,
			TransactionType: models.TxTypeDebit,
			Amount:          amount,
			Currency:        plan.Currency,
			Status:          models.TxStatusSuccess,
			Description:     fmt.Sprintf("Subscription to %s", plan.Name),
			BalanceBefore:   beforeCustomer,
			BalanceAfter:    customerWallet.Balance,
			TransactionID:   transactionID,
		}
		if err := tx.Create(&customerTx).Error; err != nil {
			return fmt.Errorf("failed to record customer transaction: %w", err)
		}

		adminTx := models.WalletTransaction{
			WalletID:        adminWallet.ID,
			UserID:          1, // Admin user ID
			Type:            models.TxTypeSubscription,
			TransactionType: models.TxTypeCredit,
			Amount:          amount,
			Currency:        plan.Currency,
			Status:          models.TxStatusSuccess,
			Description:     fmt.Sprintf("Received subscription payment from user %d", userID),
			BalanceBefore:   beforeAdmin,
			BalanceAfter:    adminWallet.Balance,
			TransactionID:   transactionID,
		}
		if err := tx.Create(&adminTx).Error; err != nil {
			return fmt.Errorf("failed to record admin transaction: %w", err)
		}

		// 6️⃣ Create Subscription Record
		startDate := time.Now()
		endDate := startDate.AddDate(0, int(plan.Duration), 0)
		subscription = &models.CustomerToTraderSub{
			UserID:             userID,
			SubscriptionPlanID: plan.ID,
			StartDate:          startDate,
			EndDate:            endDate,
			IsActive:           true,
			PaymentStatus:      "paid",
			AmountPaid:         amount,
			TransactionID:      transactionID,
		}
		if err := tx.Create(subscription).Error; err != nil {
			return fmt.Errorf("failed to create subscription: %w", err)
		}

		// 7️⃣ Upgrade user to trader
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return fmt.Errorf("user not found: %w", err)
		}
		user.Role = "trader"
		user.RoleID = uintPtr(3) // trader role_id = 3

		if err := tx.Save(&user).Error; err != nil {
			return fmt.Errorf("failed to update user role: %w", err)
		}

		// 8️⃣ Ensure trader profile exists
		var existingProfile models.TraderProfile
		err = tx.Where("user_id = ?", user.ID).First(&existingProfile).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newProfile := models.TraderProfile{
				UserID: user.ID,
				Status: models.StatusApproved,
			}
			if err := tx.Create(&newProfile).Error; err != nil {
				return fmt.Errorf("failed to create trader profile: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func uintPtr(v uint) *uint { return &v }

func (s *CustomerSubscriptionService) GetSubscriptionsByUserID(userID uint) ([]models.CustomerToTraderSub, error) {
	return s.customerSubscriptionRepo.GetSubscriptionsByUserID(userID)
}

func (s *CustomerSubscriptionService) CancelSubscription(userID, subscriptionID uint) error {
	subscription, err := s.customerSubscriptionRepo.GetSubscriptionByID(subscriptionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("subscription not found")
		}
		return fmt.Errorf("failed to retrieve subscription: %w", err)
	}

	if subscription.UserID != userID {
		return fmt.Errorf("unauthorized to cancel this subscription")
	}

	if !subscription.IsActive {
		return fmt.Errorf("subscription is already inactive or cancelled")
	}

	subscription.IsActive = false
	subscription.PaymentStatus = "cancelled"
	subscription.EndDate = time.Now()

	return s.customerSubscriptionRepo.UpdateSubscription(subscription)
}
