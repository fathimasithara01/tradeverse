package service

// import (
// 	"errors"

// 	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
// 	"github.com/fathimasithara01/tradeverse/pkg/models"
// 	"gorm.io/gorm"
// )

// type SubscriptionService interface {
// 	ListSubscribers(traderID uint) ([]models.TraderSubscription, error)
// 	GetSubscriberDetails(traderID, subscriptionID uint) (*models.TraderSubscription, error)
// }

// type subscriptionService struct {
// 	traderRepo repository.SubscriptionRepository
// }

// func NewSubscriptionService(traderRepo repository.SubscriptionRepository) SubscriptionService {
// 	return &subscriptionService{traderRepo: traderRepo}
// }

// func (s *subscriptionService) ListSubscribers(traderID uint) ([]models.TraderSubscription, error) {
// 	subscriptions, err := s.traderRepo.ListTraderSubscribers(traderID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return subscriptions, nil
// }

// func (s *subscriptionService) GetSubscriberDetails(traderID, subscriptionID uint) (*models.TraderSubscription, error) {
// 	subscription, err := s.traderRepo.GetTraderSubscriberDetails(traderID, subscriptionID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("subscriber not found or not associated with this trader")
// 		}
// 		return nil, err
// 	}
// 	return subscription, nil
// }
