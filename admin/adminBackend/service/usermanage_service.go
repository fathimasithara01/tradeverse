package service

import (
	"github.com/fathimasithara01/tradeverse/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) GetUserByID(id uint) (models.User, error) {
	return s.Repo.FindByID(id)
}

func (s *UserService) UpdateUser(userToUpdate *models.User) error {
	originalUser, err := s.Repo.FindByID(userToUpdate.ID)
	if err != nil {
		return err
	}

	if userToUpdate.Password == "" {
		userToUpdate.Password = originalUser.Password
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userToUpdate.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		userToUpdate.Password = string(hashedPassword)
	}

	return s.Repo.Update(userToUpdate)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.Repo.Delete(id)
}
