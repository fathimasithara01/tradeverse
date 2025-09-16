package service

import (
	"errors"

	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/pkg/auth"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type TraderService interface {
	Signup(name, email, password, companyName, bio, jwtSecret string) (string, error)
	Login(email, password, jwtSecret string) (string, error)
}

type traderService struct {
	repo repository.TraderRepository
}

func NewTraderService(repo repository.TraderRepository) TraderService {
	return &traderService{repo: repo}
}

func (s *traderService) Signup(name, email, password, companyName, bio, jwtSecret string) (string, error) {
	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     models.RoleTrader,
	}

	profile := &models.TraderProfile{
		CompanyName: companyName,
		Bio:         bio,
		Status:      models.StatusPending,
		IsVerified:  false,
	}

	if err := s.repo.CreateTrader(user, profile); err != nil {
		return "", err
	}

	// create token
	return auth.GenerateJWT(user.ID, user.Email, string(user.Role), *user.RoleID, jwtSecret)
}

func (s *traderService) Login(email, password, jwtSecret string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !user.CheckPassword(password) {
		return "", errors.New("invalid email or password")
	}

	return auth.GenerateJWT(user.ID, user.Email, string(user.Role), *user.RoleID, jwtSecret)
}
