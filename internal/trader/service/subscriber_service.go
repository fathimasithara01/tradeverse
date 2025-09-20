package service

import (
	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type SubscriberService interface {
	ListSubscribers(traderID uint) ([]models.Subscriber, error)
	GetSubscriber(id uint) (*models.Subscriber, error)
}

type subscriberService struct {
	repo repository.SubscriberRepository
}

func NewSubscriberService(repo repository.SubscriberRepository) SubscriberService {
	return &subscriberService{repo: repo}
}

func (s *subscriberService) ListSubscribers(traderID uint) ([]models.Subscriber, error) {
	return s.repo.GetAll(traderID)
}

func (s *subscriberService) GetSubscriber(id uint) (*models.Subscriber, error) {
	return s.repo.GetByID(id)
}
