package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	"github.com/fathimasithara01/tradeverse/pkg/auth"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

type UpdateUserRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone"`
	IsBlocked  bool   `json:"is_blocked"`
	IsVerified bool   `json:"is_verified"`
	RoleID     uint   `json:"role_id"`
}

type AssignRoleRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

type AdminUpdateProfileRequest struct {
	Name       string                `form:"name"`
	Email      string                `form:"email"`
	Phone      string                `form:"phone"`
	ProfilePic *multipart.FileHeader `form:"profile_pic"`
}

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

	GetAdminProfile(userID uint) (models.User, error)
	UpdateAdminProfile(userID uint, req AdminUpdateProfileRequest) error
	ChangeAdminPassword(userID uint, oldPassword, newPassword string) error
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
func (s *UserService) GetAdminProfile(userID uint) (models.User, error) {
	user, err := s.UserRepo.GetAdminProfile(userID)
	if err != nil {
		return models.User{}, err
	}
	// Important: Do not send password hash to the client
	user.Password = ""
	return user, nil
}
func (s *UserService) UpdateAdminProfile(userID uint, req AdminUpdateProfileRequest) error {
	const uploadDir = "./static/images/profile_pics"

	log.Printf("[INFO] UpdateAdminProfile: Initiating update for admin user ID %d", userID)

	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("admin user not found: %w", err)
	}

	if user.Role != models.RoleAdmin {
		return errors.New("access denied: user is not an admin")
	}

	// Email duplication check
	if req.Email != "" && user.Email != req.Email { // Only check if email is provided and changed
		existingUserByEmail, err := s.UserRepo.FindByEmail(req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check email availability: %w", err)
		}
		if existingUserByEmail != nil && existingUserByEmail.ID != userID {
			return errors.New("email already in use by another account")
		}
		user.Email = req.Email // Update email only if unique and valid
	}

	// Update other user fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	// Handle profile picture upload ONLY IF a file was actually provided AND has a filename
	if req.ProfilePic != nil && req.ProfilePic.Filename != "" { // Explicit check for filename
		absUploadDir, err := filepath.Abs(uploadDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for upload directory: %w", err)
		}

		if err := os.MkdirAll(absUploadDir, 0755); err != nil {
			return fmt.Errorf("failed to create upload directory '%s': %w", absUploadDir, err)
		}

		safeFilename := filepath.Base(req.ProfilePic.Filename)
		// Generate a unique filename to prevent clashes
		filename := fmt.Sprintf("%d-%s%s", userID, time.Now().Format("20060102150405"), filepath.Ext(safeFilename))
		filePath := filepath.Join(absUploadDir, filename)

		// Call the now more robust SaveUploadedFile
		if err := SaveUploadedFile(req.ProfilePic, filePath); err != nil {
			// The specific error from SaveUploadedFile will now propagate
			return fmt.Errorf("failed to save profile picture: %w", err)
		}

		user.ProfilePic = "/static/images/profile_pics/" + filename
		log.Printf("[INFO] UpdateAdminProfile: Updated profile picture for admin user ID %d to %s", userID, user.ProfilePic)
	} else {
		log.Printf("[INFO] UpdateAdminProfile: No valid profile picture provided for admin user ID %d. Skipping file upload.", userID)
	}

	return s.UserRepo.UpdateUser(user) 
}

func (s *UserService) ChangeAdminPassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("admin user not found")
	}

	if user.Role != models.RoleAdmin {
		return errors.New("access denied: user is not an admin")
	}

	// Verify old password
	if !checkPasswordHash(oldPassword, user.Password) { // Use checkPasswordHash
		return errors.New("current password is incorrect")
	}

	// Validate new password strength
	if !IsValidPassword(newPassword) {
		return errors.New("new password does not meet strength requirements (min 8 chars, 1 uppercase, 1 lowercase, 1 number, 1 special char)")
	}

	hashedPassword, err := hashPassword(newPassword) // Use hashPassword
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}
	user.Password = hashedPassword

	return s.UserRepo.UpdateUser(user) // Pass pointer to UpdateUser
}

func IsValidPassword(password string) bool {

	var (
		minLen     = 8
		hasUpper   = regexp.MustCompile(`[A-Z]`)
		hasLower   = regexp.MustCompile(`[a-z]`)
		hasNumber  = regexp.MustCompile(`[0-9]`)
		hasSpecial = regexp.MustCompile(`[!@#$%^&*()_+=\-{}\[\]:;<>,.?~\\|]`) // Added - and |
	)

	if len(password) < minLen {
		return false
	}
	if !hasUpper.MatchString(password) {
		return false
	}
	if !hasLower.MatchString(password) {
		return false
	}
	if !hasNumber.MatchString(password) {
		return false
	}
	if !hasSpecial.MatchString(password) {
		return false
	}

	return true
}

// hashPassword hashes a password using bcrypt.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash compares a plaintext password with a bcrypt hash.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SaveUploadedFile saves a multipart.FileHeader to the specified destination.
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	if file == nil {
		return errors.New("cannot save: file header is nil")
	}
	if file.Filename == "" {
		return errors.New("cannot save: file has an empty filename")
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file '%s': %w", file.Filename, err)
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file '%s': %w", dst, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return fmt.Errorf("failed to copy file to '%s': %w", dst, err)
	}

	return nil
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
		if err := s.UserRepo.UpdateUser(user); err != nil {
			return "", models.User{}, fmt.Errorf("failed to update user role info")
		}
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, string(user.Role), *user.RoleID, s.JWTSecret)
	if err != nil {
		return "", models.User{}, fmt.Errorf("failed to generate JWT: %w", err)
	}
	return token, *user, nil

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

	if profile.Phone != "" {
		existingUser.CustomerProfile.Phone = profile.Phone
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
