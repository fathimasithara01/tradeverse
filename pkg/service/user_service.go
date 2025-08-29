package service

import (
	"errors"
	"fmt"

	"github.com/fathimasithara01/tradeverse/pkg/auth"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/fathimasithara01/tradeverse/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Login(email, password string) (string, models.User, error)
	RegisterCustomer(user models.User, profile models.CustomerProfile) error
	RegisterTrader(user models.User, profile models.TraderProfile) error
	CreateTraderByAdmin(user models.User, profile models.TraderProfile) error
	GetUserByID(id uint) (models.User, error)
	CreateInternalUser(user models.User) (models.User, error)
	GetUsersByRole(role models.UserRole) ([]models.User, error)
	GetAllUsers() ([]models.User, error)
	DeleteUser(id uint) error
	UpdateUser(userToUpdate *models.User) error
	GetAllUsersAdvanced(options repository.UserQueryOptions) (repository.PaginatedUsers, error)
	GetTradersByStatus(status models.TraderStatus) ([]models.User, error)
	ApproveTrader(traderID uint) error
	RejectTrader(traderID uint) error
	GetAllUsersWithRole() ([]models.User, error)
	AssignRoleToUser(userID, roleID uint) error
}

type UserService struct {
	UserRepo  repository.IUserRepository
	RoleRepo  repository.IRoleRepository
	JWTSecret string
}

func NewUserService(userRepo repository.IUserRepository, roleRepo repository.IRoleRepository, jwtSecret string) IUserService {
	return &UserService{
		UserRepo:  userRepo,
		RoleRepo:  roleRepo,
		JWTSecret: jwtSecret,
	}
}

type UserWithRoleName struct {
	models.User
	RoleName string `json:"role_name"`
}

func (s *UserService) Login(email, password string) (string, models.User, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", models.User{}, errors.New("invalid credentials")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", models.User{}, errors.New("invalid credentials")
	}

	var roleID uint
	if user.RoleID != nil {
		roleID = *user.RoleID
	} else {
		fmt.Printf("[WARN] User '%s' (ID: %d) has a nil RoleID during login.\n", user.Email, user.ID)
	}

	fmt.Printf("[DEBUG-LOGIN] Generating JWT for UserID: %d, Role: %s, RoleID: %d\n", user.ID, user.Role, roleID)

	token, err := auth.GenerateJWT(user.ID, user.Email, string(user.Role), roleID, s.JWTSecret)
	return token, user, err
}

func (s *UserService) RegisterCustomer(user models.User, profile models.CustomerProfile) error {
	_, err := s.UserRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email already registered")
	}
	user.Role = models.RoleCustomer
	user.CustomerProfile = profile
	return s.UserRepo.Create(&user)
}

func (s *UserService) RegisterTrader(user models.User, profile models.TraderProfile) error {
	_, err := s.UserRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email already registered")
	}
	user.Role = models.RoleTrader
	user.TraderProfile = profile
	return s.UserRepo.Create(&user)
}

func (s *UserService) CreateTraderByAdmin(user models.User, profile models.TraderProfile) error {
	_, err := s.UserRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email is already registered")
	}
	user.Role = models.RoleTrader
	profile.Status = models.StatusApproved
	user.TraderProfile = profile
	return s.UserRepo.Create(&user)
}
func (s *UserService) CreateInternalUser(user models.User) (models.User, error) {
	_, err := s.UserRepo.FindByEmail(user.Email)
	if err == nil {
		return models.User{}, errors.New("a user with this email already exists")
	}

	err = s.UserRepo.Create(&user)
	return user, err
}

func (s *UserService) GetUserByID(id uint) (models.User, error) { return s.UserRepo.FindByID(id) }
func (s *UserService) GetUsersByRole(role models.UserRole) ([]models.User, error) {
	return s.UserRepo.FindByRole(role)
}
func (s *UserService) GetAllUsers() ([]models.User, error) { return s.UserRepo.FindAllNonAdmins() }
func (s *UserService) DeleteUser(id uint) error            { return s.UserRepo.Delete(id) }
func (s *UserService) UpdateUser(userToUpdate *models.User) error {
	return s.UserRepo.Update(userToUpdate)
}
func (s *UserService) GetAllUsersAdvanced(options repository.UserQueryOptions) (repository.PaginatedUsers, error) {
	return s.UserRepo.FindAllAdvanced(options)
}

func (s *UserService) GetTradersByStatus(status models.TraderStatus) ([]models.User, error) {
	return s.UserRepo.FindTradersByStatus(status)
}
func (s *UserService) ApproveTrader(traderID uint) error {
	return s.UserRepo.UpdateTraderStatus(traderID, models.StatusApproved)
}
func (s *UserService) RejectTrader(traderID uint) error {
	return s.UserRepo.UpdateTraderStatus(traderID, models.StatusRejected)

}

func (s *UserService) GetAllUsersWithRole() ([]models.User, error) {
	return s.UserRepo.FindAllNonAdmins()
}

func (s *UserService) AssignRoleToUser(userID, roleID uint) error {
	role, err := s.RoleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("invalid role selected: role not found in database")
	}

	fmt.Printf("==> ATTEMPTING TO ASSIGN ROLE: UserID=%d, RoleID=%d, RoleName='%s'\n", userID, roleID, role.Name)

	return s.UserRepo.AssignRoleToUser(userID, roleID, models.UserRole(role.Name))

}
