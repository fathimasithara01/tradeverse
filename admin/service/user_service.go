package service

import (
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
)

type UserService struct {
	Repo repository.UserRepository
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAllUsers()
}

func (s *UserService) BanUser(userID uint) error {
	return s.Repo.UpdateStatus(userID, "banned")
}

func (s *UserService) UnbanUser(userID uint) error {
	return s.Repo.UpdateStatus(userID, "active")
}

func (s *UserService) GetUserByID(id uint) (models.User, error) {
	return s.Repo.GetByID(id)
}
