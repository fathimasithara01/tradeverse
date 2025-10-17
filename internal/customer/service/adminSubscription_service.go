// internal/customer/service/customer_subscription_service.go
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

	// Get all active trader subscriptions
	activeSubscriptions, err := s.customerSubscriptionRepo.GetActiveTraderSubscriptions()
	if err != nil {
		return err
	}

	now := time.Now()
	var deactivatedCount int

	for _, sub := range activeSubscriptions {
		if sub.EndDate.Before(now) && sub.IsActive { // Assuming an EndDate field and an IsActive boolean
			log.Printf("Deactivating trader subscription ID: %d for user: %d, expired on: %s", sub.ID, sub.UserID, sub.EndDate.Format(time.RFC3339))
			sub.IsActive = false
			sub.DeactivatedAt = &now // Assuming a DeactivatedAt field
			err := s.customerSubscriptionRepo.UpdateTraderSubscription(&sub)
			if err != nil {
				log.Printf("Error deactivating trader subscription ID %d: %v", sub.ID, err)
				// Depending on your error handling strategy, you might return here or continue
				// and log all errors. For a cron job, logging and continuing might be better.
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
		plan, err := s.adminSubscriptionPlanRepo.GetSubscriptionPlanByID(planID)
		if err != nil {
			return fmt.Errorf("subscription plan not found: %w", err)
		}

		startDate := time.Now()
		var endDate time.Time

		switch plan.Interval {
		case "days":
			endDate = startDate.AddDate(0, 0, int(plan.Duration))
		case "monthly":
			endDate = startDate.AddDate(0, int(plan.Duration), 0)
		case "yearly":
			endDate = startDate.AddDate(int(plan.Duration), 0, 0)
		default:
			endDate = startDate.AddDate(0, 1, 0) // Default to 1 month if interval is unknown
		}

		newSubscription := &models.CustomerToTraderSub{
			UserID:             userID,
			SubscriptionPlanID: planID,
			StartDate:          startDate,
			EndDate:            endDate,
			IsActive:           true,
			PaymentStatus:      "paid",
			AmountPaid:         amount,
			TransactionID:      transactionID,
		}

		if err := s.customerSubscriptionRepo.CreateSubscription(newSubscription); err != nil {
			return fmt.Errorf("failed to create customer subscription: %w", err)
		}
		subscription = newSubscription

		creditDescription := fmt.Sprintf("Subscription payment from User %d for Plan %d (Amount: %.2f)", userID, planID, amount)
		if err := s.adminWalletService.CreditAdminWallet(tx, amount, plan.Currency, creditDescription); err != nil {
			log.Printf("Error crediting admin wallet for subscription: %v", err)
			return fmt.Errorf("failed to credit admin wallet for subscription: %w", err)
		}

		// If it's a trader upgrade plan, update user role
		if plan.IsUpgradeToTrader {
			user, err := s.userRepo.GetUserByID(userID)
			if err != nil {
				return fmt.Errorf("failed to get user for upgrade: %w", err)
			}
			traderRole, err := s.userRepo.GetRoleByName(models.RoleTrader)
			if err != nil {
				return fmt.Errorf("trader role not found: %w", err)
			}
			user.RoleID = &traderRole.ID
			user.Role = models.RoleTrader

			if user.TraderProfile == nil || user.TraderProfile.UserID == 0 {
				user.TraderProfile = &models.TraderProfile{
					UserID: user.ID,
					Status: models.StatusApproved,
				}
				if err := tx.Create(user.TraderProfile).Error; err != nil {
					return fmt.Errorf("failed to create trader profile for user %d: %w", userID, err)
				}
			} else {
				user.TraderProfile.Status = models.StatusApproved
				if err := tx.Save(user.TraderProfile).Error; err != nil {
					return fmt.Errorf("failed to update trader profile for user %d: %w", userID, err)
				}
			}

			if err := s.userRepo.UpdateUser(user); err != nil {
				return fmt.Errorf("failed to upgrade user to trader role: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return subscription, nil
}

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
	subscription.EndDate = time.Now() // Set end date to now upon cancellation

	return s.customerSubscriptionRepo.UpdateSubscription(subscription)
}
