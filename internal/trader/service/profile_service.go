package service

import (
	"errors"

	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

var (
	ErrTraderProfileNotFound = errors.New("trader profile not found")
	ErrTraderProfileExists   = errors.New("trader profile already exists for this user")
	ErrUnauthorized          = errors.New("unauthorized action")
	ErrInvalidInput          = errors.New("invalid input for trader profile")
	ErrPermissionDenied      = errors.New("permission denied to perform this action")
)

// ITraderProfileService defines the interface for TraderProfile business logic.
type ITraderProfileService interface {
	GetProfile(userID uint) (*models.TraderProfile, error)
	CreateProfile(userID uint, name, companyName, bio string) (*models.TraderProfile, error)
	UpdateProfile(userID uint, profileID uint, name, companyName, bio *string) (*models.TraderProfile, error)
	DeleteProfile(userID uint, profileID uint) error // Assuming only the user or admin can delete their profile
}

// TraderProfileService implements ITraderProfileService.
type TraderProfileService struct {
	traderRepo repository.ITraderProfileRepository
	userRepo   repository.IUserRepository // Assuming you have a User repository for user-related checks
}

// NewTraderProfileService creates a new TraderProfileService.
func NewTraderProfileService(traderRepo repository.ITraderProfileRepository, userRepo repository.IUserRepository) *TraderProfileService {
	return &TraderProfileService{traderRepo: traderRepo, userRepo: userRepo}
}

// GetProfile retrieves a trader's profile.
func (s *TraderProfileService) GetProfile(userID uint) (*models.TraderProfile, error) {
	profile, err := s.traderRepo.GetTraderProfileByUserID(userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, ErrTraderProfileNotFound
	}
	return profile, nil
}

// CreateProfile creates a new trader profile for a given user.
func (s *TraderProfileService) CreateProfile(userID uint, name, companyName, bio string) (*models.TraderProfile, error) {
	// First, check if a profile already exists for this user
	existingProfile, err := s.traderRepo.GetTraderProfileByUserID(userID)
	if err != nil {
		return nil, err
	}
	if existingProfile != nil {
		return nil, ErrTraderProfileExists
	}

	// You might want to check if the user is actually designated as a 'trader' role in the User model
	user, err := s.userRepo.GetUserByID(userID) // Assuming this method exists
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsTrader() {
		return nil, ErrPermissionDenied // Or a more specific error
	}

	profile := &models.TraderProfile{
		UserID:      userID,
		Name:        name,
		CompanyName: companyName,
		Bio:         bio,
		Status:      models.StatusPending, // New profiles are pending by default
		TotalPnL:    0.0,
		IsVerified:  false,
	}

	if err := s.traderRepo.CreateTraderProfile(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *TraderProfileService) UpdateProfile(userID uint, profileID uint, name, companyName, bio *string) (*models.TraderProfile, error) {
	profile, err := s.traderRepo.GetTraderProfileByUserID(userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, ErrTraderProfileNotFound
	}

	if profile.UserID != userID {
		return nil, ErrUnauthorized
	}

	if name != nil {
		profile.Name = *name
	}
	if companyName != nil {
		profile.CompanyName = *companyName
	}
	if bio != nil {
		profile.Bio = *bio
	}

	if err := s.traderRepo.UpdateTraderProfile(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *TraderProfileService) DeleteProfile(userID uint, profileID uint) error {
	profile, err := s.traderRepo.GetTraderProfileByUserID(userID)
	if err != nil {
		return err
	}
	if profile == nil {
		return ErrTraderProfileNotFound
	}

	// Only the owner of the profile can delete it, or an admin
	if profile.UserID != userID {
		// Here you would also check if the requesting user is an admin
		// For simplicity, we'll just check ownership for now.
		return ErrUnauthorized
	}

	// Before deleting the profile, consider cascading deletes for related data
	// (e.g., subscriptions specific to this trader if not handled by DB cascades)
	// GORM's `constraint:OnDelete:CASCADE` on the User struct for TraderProfile helps here.

	return s.traderRepo.DeleteTraderProfile(profile.ID)
}
