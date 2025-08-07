package service

import (
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type SubscriptionService struct {
	Repo repository.SubscriptionRepository
}

func (s *SubscriptionService) GetAll() ([]models.Subscription, error) {
	return s.Repo.GetAll()
}
