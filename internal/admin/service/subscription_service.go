package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type ISubscriptionService interface {
	CreateSubscription(userID, planID uint, amount float64, transactionID string) (*models.Subscription, error)
	GetAllSubscriptions() ([]models.Subscription, error)
	GetSubscriptionByID(id uint) (*models.Subscription, error)
	GetSubscriptionsByUserID(userID uint) ([]models.Subscription, error)
	UpdateSubscription(subscription *models.Subscription) error
	DeleteSubscription(id uint) error
	GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)
	UpgradeUserToTrader(userID uint) error

	DeactivateExpiredSubscriptions() error
	UpdateUserTraderStatus(userID uint, status string) error
}

type SubscriptionService struct {
	subscriptionRepo   repository.ISubscriptionRepository
	planRepo           repository.ISubscriptionPlanRepository
	userRepo           repository.IUserRepository
	adminWalletService IAdminWalletService
	DB                 *gorm.DB
}

func NewSubscriptionService(subRepo repository.ISubscriptionRepository, planRepo repository.ISubscriptionPlanRepository, userRepo repository.IUserRepository, adminWalletService IAdminWalletService, db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo:   subRepo,
		planRepo:           planRepo,
		userRepo:           userRepo,
		adminWalletService: adminWalletService,
		DB:                 db,
	}
}

func (s *SubscriptionService) UpdateUserTraderStatus(userID uint, status string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %d not found", userID)
		}
		return fmt.Errorf("failed to get user by ID %d: %w", userID, err)
	}

	// Initialize TraderProfile if it's nil or has a zero UserID
	if user.TraderProfile == nil || user.TraderProfile.UserID == 0 {
		return fmt.Errorf("user %d does not have an initialized trader profile", userID)
	}

	// Validate the status string against models.TraderProfileStatus
	switch status {
	case string(models.StatusApproved):
		user.TraderProfile.Status = models.StatusApproved
	case string(models.StatusRejected):
		user.TraderProfile.Status = models.StatusRejected
	default:
		return fmt.Errorf("invalid trader status provided: %s", status)
	}

	// Using DB.Save to update the nested TraderProfile
	if err := s.DB.Save(&user.TraderProfile).Error; err != nil {
		return fmt.Errorf("failed to update trader profile for user %d: %w", userID, err)
	}

	// If status is approved, ensure the user has the trader role
	if user.TraderProfile.Status == models.StatusApproved && user.Role != models.RoleTrader {
		traderRole, err := s.userRepo.GetRoleByName(models.RoleTrader)
		if err != nil {
			log.Printf("Warning: Trader role not found when approving user %d: %v", userID, err)
			// Decide if this should be a critical error or just a warning.
			// For now, let's just log and continue, but ideally, the role should exist.
		} else {
			user.RoleID = &traderRole.ID
			user.Role = models.RoleTrader
			if err := s.userRepo.UpdateUser(user); err != nil {
				log.Printf("Warning: Failed to update user %d to trader role after profile approval: %v", userID, err)
			}
		}
	}

	return nil
}

func (s *SubscriptionService) DeactivateExpiredSubscriptions() error {
	log.Println("Running cron job: Deactivating expired subscriptions...")
	expiredSubs, err := s.subscriptionRepo.GetExpiredActiveSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to get expired active subscriptions: %w", err)
	}

	if len(expiredSubs) == 0 {
		log.Println("No expired subscriptions found to deactivate.")
		return nil
	}

	for _, sub := range expiredSubs {
		sub.IsActive = false
		sub.PaymentStatus = "expired"
		if err := s.subscriptionRepo.UpdateSubscription(&sub); err != nil {
			log.Printf("Error deactivating subscription ID %d for user %d: %v", sub.ID, sub.UserID, err)
		} else {
			log.Printf("Deactivated subscription ID %d for user %d (Plan: %s)", sub.ID, sub.UserID, sub.SubscriptionPlan.Name)
		}
	}
	log.Printf("Cron job finished: Deactivated %d subscriptions.", len(expiredSubs))
	return nil
}

func (s *SubscriptionService) UpgradeUserToTrader(userID uint) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	traderRole, err := s.userRepo.GetRoleByName(models.RoleTrader)
	if err != nil {
		log.Printf("Error: Trader role not found: %v", err)
		return err
	}

	user.RoleID = &traderRole.ID
	user.Role = models.RoleTrader

	// Initialize TraderProfile if it's nil or has a zero UserID
	if user.TraderProfile == nil || user.TraderProfile.UserID == 0 {
		user.TraderProfile = &models.TraderProfile{ // Make sure to initialize as pointer if it's a pointer in the model
			UserID: user.ID,
			Status: models.StatusApproved,
		}
		if err := s.DB.Create(user.TraderProfile).Error; err != nil {
			return fmt.Errorf("failed to create trader profile for user %d: %w", userID, err)
		}
	} else {
		user.TraderProfile.Status = models.StatusApproved
		if err := s.DB.Save(user.TraderProfile).Error; err != nil {
			return fmt.Errorf("failed to update trader profile for user %d: %w", userID, err)
		}
	}

	return s.userRepo.UpdateUser(user)
}

func (s *SubscriptionService) CreateSubscription(userID, planID uint, amount float64, transactionID string) (*models.Subscription, error) {
	var subscription *models.Subscription
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		plan, err := s.planRepo.GetSubscriptionPlanByID(planID)
		if err != nil {
			return err
		}

		startDate := time.Now()
		var endDate time.Time

		switch plan.Interval {
		case "days":
			endDate = startDate.AddDate(0, 0, plan.Duration)
		case "monthly":
			endDate = startDate.AddDate(0, plan.Duration, 0)
		case "yearly":
			endDate = startDate.AddDate(plan.Duration, 0, 0)
		default:
			endDate = startDate.AddDate(0, 1, 0) // Default to 1 month if interval is not specified/recognized
		}

		newSubscription := &models.Subscription{
			UserID:             userID,
			SubscriptionPlanID: planID,
			StartDate:          startDate,
			EndDate:            endDate,
			IsActive:           true,
			PaymentStatus:      "paid",
			AmountPaid:         amount,
			TransactionID:      transactionID,
		}

		if err := s.subscriptionRepo.CreateSubscription(newSubscription); err != nil {
			return err
		}
		subscription = newSubscription

		creditDescription := fmt.Sprintf("Subscription payment from User %d for Plan %d (Amount: %.2f)", userID, planID, amount)
		if err := s.adminWalletService.CreditAdminWallet(tx, amount, "INR", creditDescription); err != nil { // Assuming INR as default currency
			log.Printf("Error crediting admin wallet for subscription: %v", err)
			return fmt.Errorf("failed to credit admin wallet for subscription: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *SubscriptionService) GetAllSubscriptions() ([]models.Subscription, error) {
	log.Println("DEBUG: SubscriptionService.GetAllSubscriptions was called.") // Add this
	subs, err := s.subscriptionRepo.GetAllSubscriptions()
	if err != nil {
		log.Printf("ERROR: SubscriptionService.GetAllSubscriptions failed: %v", err)
		return nil, fmt.Errorf("failed to retrieve all subscriptions: %w", err)
	}
	log.Printf("DEBUG: SubscriptionService fetched %d subscriptions from repo.", len(subs)) // Add this
	return subs, nil
}

func (s *SubscriptionService) GetSubscriptionByID(id uint) (*models.Subscription, error) {
	return s.subscriptionRepo.GetSubscriptionByID(id)
}

func (s *SubscriptionService) GetSubscriptionsByUserID(userID uint) ([]models.Subscription, error) {
	return s.subscriptionRepo.GetSubscriptionsByUserID(userID)
}

func (s *SubscriptionService) UpdateSubscription(subscription *models.Subscription) error {
	return s.subscriptionRepo.UpdateSubscription(subscription)
}

func (s *SubscriptionService) DeleteSubscription(id uint) error {
	return s.subscriptionRepo.DeleteSubscription(id)
}

func (s *SubscriptionService) GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error) {
	return s.planRepo.GetSubscriptionPlanByID(id)
}
