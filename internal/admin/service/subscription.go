package service

import (
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
}

type SubscriptionService struct {
	subscriptionRepo   repository.ISubscriptionRepository
	planRepo           repository.ISubscriptionPlanRepository
	userRepo           repository.IUserRepository
	adminWalletService IAdminWalletService // Add AdminWalletService
	DB                 *gorm.DB            // Add DB for transaction management
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

func (s *SubscriptionService) UpgradeUserToTrader(userID uint) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	traderRole, err := s.userRepo.GetRoleByName(models.RoleTrader) // Assuming GetRoleByName exists
	if err != nil {
		log.Printf("Error: Trader role not found: %v", err)
		return err
	}

	user.RoleID = &traderRole.ID
	user.Role = models.RoleTrader

	// If a trader profile doesn't exist, create one
	if user.TraderProfile.UserID == 0 { // Check if default/zero value, indicating no profile
		user.TraderProfile = models.TraderProfile{
			UserID: user.ID,
			Status: models.StatusApproved, // Automatically approved if upgraded via subscription
		}
	} else {
		user.TraderProfile.Status = models.StatusApproved // Ensure status is approved
	}

	// Ensure to save the TraderProfile as well
	// Depending on your GORM setup, this might require a separate call or could be handled by `UpdateUser` if relations are set for `Save`.
	// For explicit handling:
	if user.TraderProfile.ID == 0 { // New profile
		if err := s.DB.Create(&user.TraderProfile).Error; err != nil {
			return fmt.Errorf("failed to create trader profile for user %d: %w", userID, err)
		}
	} else { // Existing profile
		if err := s.DB.Save(&user.TraderProfile).Error; err != nil {
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
			endDate = startDate.AddDate(0, 1, 0) // Default to 1 month
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
		subscription = newSubscription // Assign to outer variable

		// Credit the admin's wallet with the subscription amount
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
	return s.subscriptionRepo.GetAllSubscriptions()
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
