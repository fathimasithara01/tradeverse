package service

import (
	"errors"
	"log"

	"github.com/fathimasithara01/tradeverse/auth"
	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Login(email, password string) (string, models.User, error) {
	log.Printf("[SERVICE-LOGIN] Attempting to find user by email: %s", email)
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		log.Printf("[SERVICE-LOGIN-ERROR] User not found or DB error for email '%s': %v\n", email, err)
		return "", models.User{}, errors.New("invalid credentials")
	}

	log.Printf("[SERVICE-LOGIN-SUCCESS] Found user in database. UserID: %d, Email: %s\n", user.ID, user.Email)

	if user.ID == 0 {
		return "", models.User{}, errors.New("internal error: user record found but ID is zero")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("[SERVICE-LOGIN-ERROR] Password mismatch for user '%s'\n", email)
		return "", models.User{}, errors.New("invalid credentials")
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, string(user.Role))
	if err != nil {
		log.Printf("[SERVICE-LOGIN-ERROR] Failed to generate JWT for user '%s': %v\n", email, err)
		return "", models.User{}, errors.New("failed to generate token")
	}

	return token, user, nil
}

func (s *UserService) RegisterCustomer(user models.User, profile models.CustomerProfile) error {
	_, err := s.Repo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email already registered")
	}

	user.Role = models.RoleCustomer
	user.CustomerProfile = profile
	return s.Repo.Create(&user)
}
func (s *UserService) RegisterTrader(user models.User, profile models.TraderProfile) error {
	_, err := s.Repo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email already registered")
	}

	user.Role = models.RoleTrader
	user.TraderProfile = profile
	return s.Repo.Create(&user)
}

func (s *UserService) RegisterAdmin(user models.User) (models.User, error) {
	_, err := s.Repo.FindByEmail(user.Email)
	if err == nil {
		return models.User{}, errors.New("email already registered")
	}

	user.Role = models.RoleAdmin
	err = s.Repo.Create(&user)
	return user, err
}

func (s *UserService) CreateTrader(user models.User, profile models.TraderProfile) error {
	_, err := s.Repo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email is already registered")
	}

	user.Role = models.RoleTrader
	user.TraderProfile = profile

	return s.Repo.Create(&user)
}

func (s *UserService) CreateTraderByAdmin(user models.User, profile models.TraderProfile) error {
	_, err := s.Repo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email is already registered")
	}

	user.Role = models.RoleTrader // Set the role to 'trader'

	profile.Status = models.StatusApproved

	user.TraderProfile = profile // Attach the trader-specific profile

	return s.Repo.Create(&user)
}

func (s *UserService) CreateCustomer(user models.User, profile models.CustomerProfile) error {
	user.Role = models.RoleCustomer
	user.CustomerProfile = profile

	_, err := s.Repo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email already registered")
	}

	return s.Repo.Create(&user)
}

func (s *UserService) GetUsersByRole(role models.UserRole) ([]models.User, error) {
	return s.Repo.FindByRole(role)
}

func (s *UserService) GetTradersByStatus(status models.TraderStatus) ([]models.User, error) {
	return s.Repo.FindTradersByStatus(status)
}

func (s *UserService) ApproveTrader(traderID uint) error {
	return s.Repo.UpdateTraderStatus(traderID, models.StatusApproved)
}

func (s *UserService) RejectTrader(traderID uint) error {
	return s.Repo.UpdateTraderStatus(traderID, models.StatusRejected)
}
