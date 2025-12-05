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

type ITraderProfileService interface {
	GetProfile(userID uint) (*models.TraderProfile, error)
	CreateProfile(userID uint, name, companyName, bio string) (*models.TraderProfile, error)
	UpdateProfile(userID uint, profileID uint, name, companyName, bio *string) (*models.TraderProfile, error)
	DeleteProfile(userID uint, profileID uint) error
}

type TraderProfileService struct {
	traderRepo repository.ITraderProfileRepository
}

func NewTraderProfileService(traderRepo repository.ITraderProfileRepository) *TraderProfileService {
	return &TraderProfileService{traderRepo: traderRepo}
}

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

func (s *TraderProfileService) CreateProfile(userID uint, name, companyName, bio string) (*models.TraderProfile, error) {
	existingProfile, err := s.traderRepo.GetTraderProfileByUserID(userID)
	if err != nil {
		return nil, err
	}
	if existingProfile != nil {
		return nil, ErrTraderProfileExists
	}

	user, err := s.traderRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsTrader() {
		return nil, ErrPermissionDenied
	}

	profile := &models.TraderProfile{
		UserID:      userID,
		Name:        name,
		CompanyName: companyName,
		Bio:         bio,
		Status:      models.StatusPending,
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

	if profile.UserID != userID {
		return ErrUnauthorized
	}

	return s.traderRepo.DeleteTraderProfile(profile.ID)
}
