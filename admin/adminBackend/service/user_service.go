package service

import (
	"errors"
	"strings"

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

func (s *UserService) Register(user models.User) (models.User, error) {
	user.Email = strings.TrimSpace(user.Email)
	if user.Email == "" || user.Password == "" {
		return models.User{}, errors.New("email and password are required")
	}

	_, err := s.Repo.FindByEmail(user.Email)
	if err == nil {
		return models.User{}, errors.New("email already registered")
	}

	// Password hashing is now handled by the BeforeCreate hook in the user model.
	// We just need to call the repository's Create method.
	err = s.Repo.Create(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserService) Login(email, password string) (string, models.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", models.User{}, errors.New("invalid email or password")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", models.User{}, errors.New("invalid email or password")
	}

	// Generate JWT
	token, err := auth.GenerateJWT(user.ID, user.Email, string(user.Role))
	if err != nil {
		return "", models.User{}, errors.New("failed to generate token")
	}
	return token, user, nil
}

func (s *UserService) CreateCustomer(user models.User, profile models.CustomerProfile) error {
	user.Role = models.RoleCustomer
	user.CustomerProfile = profile

	// Check if email already exists before trying to create.
	_, err := s.Repo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email already registered")
	}

	// The BeforeCreate hook in the User model will hash the password.
	return s.Repo.Create(&user)
}

func (s *UserService) GetUsersByRole(role models.UserRole) ([]models.User, error) {
	return s.Repo.FindByRole(role)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.Repo.Delete(id)
}
