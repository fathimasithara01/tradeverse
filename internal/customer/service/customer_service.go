package service

import (
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"gorm.io/gorm"
)

type SubscriptionService interface {
	SubscribeCustomerToTrader(userID, planID uint) (*models.TraderSubscription, error)
}

type subscriptionService struct {
	db   *gorm.DB
	repo customerrepo.SubscriptionRepository
}

func NewSubscriptionService(db *gorm.DB, repo customerrepo.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{db: db, repo: repo}
}

func (s *subscriptionService) SubscribeCustomerToTrader(userID, planID uint) (*models.TraderSubscription, error) {
	var plan models.SubscriptionPlan
	if err := s.db.First(&plan, planID).Error; err != nil {
		return nil, errors.New("subscription plan not found")
	}

	existing, _ := s.repo.GetActiveSubscription(userID, planID)
	if existing != nil {
		return nil, errors.New("already subscribed to this trader plan")
	}

	var adminUser models.User
	if err := s.db.Where("role = ?", models.RoleAdmin).First(&adminUser).Error; err != nil {
		return nil, errors.New("admin not found")
	}

	var adminWallet models.Wallet
	if err := s.db.Where("user_id = ?", adminUser.ID).First(&adminWallet).Error; err != nil {
		return nil, errors.New("admin wallet not found")
	}

	sub := &models.TraderSubscription{
		UserID:                   userID,
		TraderSubscriptionPlanID: plan.ID,
		StartDate:                time.Now(),
		EndDate:                  time.Now().AddDate(0, 1, 0), // Example: 1 month
		IsActive:                 true,
		PaymentStatus:            "paid",
		AmountPaid:               plan.Price,
		TransactionID:            "", // Can fill with PG TX ID later
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(sub).Error; err != nil {
			return err
		}

		before := adminWallet.Balance
		adminWallet.Balance += plan.Price
		if err := tx.Save(&adminWallet).Error; err != nil {
			return err
		}

		wtx := models.WalletTransaction{
			WalletID:        adminWallet.ID,
			UserID:          adminUser.ID,
			TransactionType: models.TxTypeSubscription,
			Amount:          plan.Price,
			Currency:        plan.Currency,
			Status:          models.TxStatusSuccess,
			Description:     "Customer subscription payment",
			BalanceBefore:   before,
			BalanceAfter:    adminWallet.Balance,
			SubscriptionID:  &sub.ID,
		}
		if err := tx.Create(&wtx).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return sub, nil
}
