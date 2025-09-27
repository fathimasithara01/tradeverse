package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/auth"
	"github.com/fathimasithara01/tradeverse/pkg/models"
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

	UpdateCustomerProfile(userID uint, user models.User, profile models.CustomerProfile) error
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserService) Login(email, password string) (string, models.User, error) {
	user, err := s.UserRepo.FindByEmail(email) 
	if err != nil {
		log.Printf("[LOGIN SERVICE] User '%s' not found or other error: %v", email, err)
		return "", models.User{}, errors.New("invalid credentials")
	}

	if user.IsBlocked {
		return "", models.User{}, errors.New("account is blocked")
	}

	if !checkPasswordHash(password, user.Password) {
		log.Printf("[LOGIN SERVICE] Password mismatch for user '%s'", email)
		return "", models.User{}, errors.New("invalid credentials")
	}

	if user.RoleID == nil {
		log.Printf("[LOGIN SERVICE] User '%s' has nil RoleID, fixing...", email)
		role, err := s.UserRepo.GetRoleByName(user.Role)
		if err != nil {
			return "", models.User{}, fmt.Errorf("failed to get role for user %s: %w", email, err)
		}
		user.RoleID = &role.ID
		if err := s.UserRepo.UpdateUser(&user); err != nil {
			return "", models.User{}, fmt.Errorf("failed to update user role info")
		}
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, string(user.Role), *user.RoleID, s.JWTSecret)
	if err != nil {
		return "", models.User{}, fmt.Errorf("failed to generate JWT: %w", err)
	}
	return token, user, nil

}
func (s *UserService) UpdateCustomerProfile(userID uint, user models.User, profile models.CustomerProfile) error {
	existingUser, err := s.UserRepo.GetUserByIDWithProfile(userID)
	if err != nil {
		return fmt.Errorf("user not found for update: %w", err)
	}

	if existingUser.Role != models.RoleCustomer {
		return errors.New("cannot update customer profile for a non-customer user")
	}

	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Email != "" && existingUser.Email != user.Email {
		_, err := s.UserRepo.FindByEmail(user.Email)
		if err == nil {
			return errors.New("email already registered by another user")
		}
		existingUser.Email = user.Email
	}

	if profile.PhoneNumber != "" {
		existingUser.CustomerProfile.PhoneNumber = profile.PhoneNumber
	}

	return s.UserRepo.UpdateUserAndProfile(existingUser)
}

func (s *UserService) RegisterCustomer(user models.User, profile models.CustomerProfile) error {
	_, err := s.UserRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email already registered")
	}

	customerRole, err := s.UserRepo.GetRoleByName(models.RoleCustomer)
	if err != nil {
		return fmt.Errorf("failed to retrieve customer role: %w", err)
	}

	user.Role = models.RoleCustomer
	user.RoleID = &customerRole.ID
	profile.UserID = user.ID

	return s.UserRepo.CreateCustomerWithProfile(&user, &profile)
}

func (s *UserService) RegisterTrader(user models.User, profile models.TraderProfile) error {
	_, err := s.UserRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email already registered")
	}

	traderRole, err := s.UserRepo.GetRoleByName(models.RoleTrader)
	if err != nil {
		return fmt.Errorf("failed to retrieve trader role: %w", err)
	}

	user.Role = models.RoleTrader
	user.RoleID = &traderRole.ID
	profile.Status = models.StatusPending
	profile.UserID = user.ID

	return s.UserRepo.CreateTraderWithProfile(&user, &profile)
}

func (s *UserService) CreateTraderByAdmin(user models.User, profile models.TraderProfile) error {
	_, err := s.UserRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email is already registered")
	}

	traderRole, err := s.UserRepo.GetRoleByName(models.RoleTrader)
	if err != nil {
		return fmt.Errorf("failed to retrieve trader role: %w", err)
	}

	user.Role = models.RoleTrader
	user.RoleID = &traderRole.ID
	profile.Status = models.StatusApproved
	profile.UserID = user.ID

	return s.UserRepo.CreateTraderWithProfile(&user, &profile)
}

func (s *UserService) CreateInternalUser(user models.User) (models.User, error) {
	_, err := s.UserRepo.FindByEmail(user.Email)
	if err == nil {
		return models.User{}, errors.New("a user with this email already exists")
	}

	if user.Role != "" {
		role, err := s.UserRepo.GetRoleByName(user.Role)
		if err != nil {
			return models.User{}, fmt.Errorf("failed to retrieve role '%s': %w", user.Role, err)
		}
		user.RoleID = &role.ID
	}

	err = s.UserRepo.Create(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create internal user: %w", err)
	}
	return user, nil
}

func (s *UserService) GetUserByID(id uint) (models.User, error) {
	user, err := s.UserRepo.FindByID(id)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user by ID %d: %w", id, err)
	}
	return user, nil
}

func (s *UserService) GetUsersByRole(role models.UserRole) ([]models.User, error) {
	log.Printf("Fetching users with role: %s", role)
	users, err := s.UserRepo.GetUsersByRole(role)
	if err != nil {
		log.Printf("Error fetching users by role %s: %v", role, err)
		return nil, fmt.Errorf("service: failed to get users by role %s: %w", role, err)
	}
	log.Printf("Found %d users with role %s", len(users), role)
	return users, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.UserRepo.FindAllNonAdmins()
	if err != nil {
		return nil, fmt.Errorf("failed to get all non-admin users: %w", err)
	}
	return users, nil
}

func (s *UserService) DeleteUser(id uint) error {
	err := s.UserRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete user %d: %w", id, err)
	}
	return nil
}

func (s *UserService) UpdateUser(userToUpdate *models.User) error {
	err := s.UserRepo.UpdateUserAndProfile(userToUpdate)
	if err != nil {
		return fmt.Errorf("failed to update user %d: %w", userToUpdate.ID, err)
	}
	return nil
}

func (s *UserService) GetAllUsersAdvanced(options repository.UserQueryOptions) (repository.PaginatedUsers, error) {
	paginatedUsers, err := s.UserRepo.FindAllAdvanced(options)
	if err != nil {
		return repository.PaginatedUsers{}, fmt.Errorf("failed to get advanced paginated users: %w", err)
	}
	return paginatedUsers, nil
}

func (s *UserService) GetTradersByStatus(status models.TraderStatus) ([]models.User, error) {
	traders, err := s.UserRepo.FindTradersByStatus(status)
	if err != nil {
		return nil, fmt.Errorf("failed to get traders by status '%s': %w", status, err)
	}
	return traders, nil
}
func (s *UserService) ApproveTrader(traderID uint) error {
	err := s.UserRepo.UpdateTraderStatus(traderID, models.StatusApproved)
	if err != nil {
		return fmt.Errorf("failed to approve trader %d: %w", traderID, err)
	}
	return nil
}
func (s *UserService) RejectTrader(traderID uint) error {
	err := s.UserRepo.UpdateTraderStatus(traderID, models.StatusRejected)
	if err != nil {
		return fmt.Errorf("failed to reject trader %d: %w", traderID, err)
	}
	return nil
}

func (s *UserService) GetAllUsersWithRole() ([]models.User, error) {
	users, err := s.UserRepo.FindAllWithRole()
	if err != nil {
		return nil, fmt.Errorf("failed to get all users with role: %w", err)
	}
	return users, nil
}

func (s *UserService) AssignRoleToUser(userID, roleID uint) error {
	role, err := s.RoleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("invalid role selected: role not found in database")
	}

	log.Printf("==> ATTEMPTING TO ASSIGN ROLE: UserID=%d, RoleID=%d, RoleName='%s'\n", userID, roleID, role.Name)

	err = s.UserRepo.AssignRoleToUser(userID, roleID, models.UserRole(role.Name))
	if err != nil {
		return fmt.Errorf("failed to assign role %s to user %d: %w", role.Name, userID, err)
	}
	log.Printf("==> SUCCESSFULLY ASSIGNED ROLE: UserID=%d, RoleID=%d, RoleName='%s'\n", userID, roleID, role.Name)
	return nil
}
