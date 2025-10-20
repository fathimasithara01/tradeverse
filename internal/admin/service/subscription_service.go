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
	CreateSubscription(userID, planID uint, amount float64, transactionID string) (*models.CustomerToTraderSub, error)
	GetAllSubscriptions() ([]models.CustomerToTraderSub, error)
	GetSubscriptionByID(id uint) (*models.CustomerToTraderSub, error)
	GetSubscriptionsByUserID(userID uint) ([]models.CustomerToTraderSub, error)
	UpdateSubscription(subscription *models.CustomerToTraderSub) error
	DeleteSubscription(id uint) error
	GetSubscriptionPlanByID(id uint) (*models.AdminTraderSubscriptionPlan, error)
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

	if user.TraderProfile == nil || user.TraderProfile.UserID == 0 {
		return fmt.Errorf("user %d does not have an initialized trader profile", userID)
	}

	switch status {
	case string(models.StatusApproved):
		user.TraderProfile.Status = models.StatusApproved
	case string(models.StatusRejected):
		user.TraderProfile.Status = models.StatusRejected
	default:
		return fmt.Errorf("invalid trader status provided: %s", status)
	}

	if err := s.DB.Save(&user.TraderProfile).Error; err != nil {
		return fmt.Errorf("failed to update trader profile for user %d: %w", userID, err)
	}

	if user.TraderProfile.Status == models.StatusApproved && user.Role != models.RoleTrader {
		traderRole, err := s.userRepo.GetRoleByName(models.RoleTrader)
		if err != nil {
			log.Printf("Warning: Trader role not found when approving user %d: %v", userID, err)
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

	if user.TraderProfile == nil || user.TraderProfile.UserID == 0 {
		user.TraderProfile = &models.TraderProfile{
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

func (s *SubscriptionService) CreateSubscription(userID, planID uint, amount float64, transactionID string) (*models.CustomerToTraderSub, error) {
	var subscription *models.CustomerToTraderSub
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		plan, err := s.planRepo.GetSubscriptionPlanByID(planID)
		if err != nil {
			return err
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
			endDate = startDate.AddDate(0, 1, 0)
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

func (s *SubscriptionService) GetAllSubscriptions() ([]models.CustomerToTraderSub, error) {
	log.Println("DEBUG: SubscriptionService.GetAllSubscriptions was called.")
	subs, err := s.subscriptionRepo.GetAllSubscriptions()
	if err != nil {
		log.Printf("ERROR: SubscriptionService.GetAllSubscriptions failed: %v", err)
		return nil, fmt.Errorf("failed to retrieve all subscriptions: %w", err)
	}
	log.Printf("DEBUG: SubscriptionService fetched %d subscriptions from repo.", len(subs))
	return subs, nil
}

func (s *SubscriptionService) GetSubscriptionByID(id uint) (*models.CustomerToTraderSub, error) {
	return s.subscriptionRepo.GetSubscriptionByID(id)
}

func (s *SubscriptionService) GetSubscriptionsByUserID(userID uint) ([]models.CustomerToTraderSub, error) {
	return s.subscriptionRepo.GetSubscriptionsByUserID(userID)
}

func (s *SubscriptionService) UpdateSubscription(subscription *models.CustomerToTraderSub) error {
	return s.subscriptionRepo.UpdateSubscription(subscription)
}

func (s *SubscriptionService) DeleteSubscription(id uint) error {
	return s.subscriptionRepo.DeleteSubscription(id)
}

func (s *SubscriptionService) GetSubscriptionPlanByID(id uint) (*models.AdminTraderSubscriptionPlan, error) {
	return s.planRepo.GetSubscriptionPlanByID(id)
}
