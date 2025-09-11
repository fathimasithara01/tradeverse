package service

import (
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type KYCServicer interface {
	SubmitKYCDocuments(userID uint, docType, docURL string) error
	GetKYCStatus(userID uint) (*models.KYCStatusResponse, error)
}

type kycService struct {
	kycRepo repository.KYCRepository
}

func NewKYCService(kycRepo repository.KYCRepository) KYCServicer {
	return &kycService{kycRepo: kycRepo}
}

func (s *kycService) SubmitKYCDocuments(userID uint, docType, docURL string) error {
	// 1. Create the new KYC document record
	newDoc := &models.KYCDocument{
		UserID:             userID,
		DocumentType:       docType,
		DocumentURL:        docURL,
		VerificationStatus: models.KYCStatusPending, // New documents are always pending review
	}
	if err := s.kycRepo.CreateKYCDocument(newDoc); err != nil {
		return errors.New("failed to save KYC document: " + err.Error())
	}

	// 2. Update/Create the user's overall KYC status
	userKYCStatus, err := s.kycRepo.FindUserKYCStatus(userID)
	if err != nil {
		return errors.New("failed to retrieve user KYC status: " + err.Error())
	}

	if userKYCStatus == nil {
		// No existing status, create a new one
		newUserStatus := &models.UserKYCStatus{
			UserID:          userID,
			Status:          models.KYCStatusPending,
			LastUpdatedBy:   userID, // The user themselves submitted it
			LastUpdatedDate: time.Now(),
		}
		if err := s.kycRepo.CreateUserKYCStatus(newUserStatus); err != nil {
			return errors.New("failed to create user KYC status: " + err.Error())
		}
	} else if userKYCStatus.Status != models.KYCStatusApproved {
		// If status exists and is not already approved, set it to PENDING
		userKYCStatus.Status = models.KYCStatusPending
		userKYCStatus.Reason = "" // Clear any previous rejection reasons
		userKYCStatus.LastUpdatedBy = userID
		userKYCStatus.LastUpdatedDate = time.Now()
		if err := s.kycRepo.UpdateUserKYCStatus(userKYCStatus); err != nil {
			return errors.New("failed to update user KYC status: " + err.Error())
		}
	}
	// If it's already approved, submitting more documents doesn't change the overall APPROVED status,
	// though the individual document will be pending. An admin would re-review.

	return nil
}

// GetKYCStatus retrieves the user's current overall KYC verification status.
func (s *kycService) GetKYCStatus(userID uint) (*models.KYCStatusResponse, error) {
	userKYCStatus, err := s.kycRepo.FindUserKYCStatus(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve KYC status from repository: " + err.Error())
	}

	if userKYCStatus == nil {
		// If no status entry exists, it means nothing was ever submitted.
		return &models.KYCStatusResponse{
			Status:      models.KYCStatusNotSubmitted,
			LastUpdated: time.Now(), // Default to now as a "not found" time
		}, nil
	}

	return &models.KYCStatusResponse{
		Status:        userKYCStatus.Status,
		Reason:        userKYCStatus.Reason,
		LastUpdatedBy: userKYCStatus.LastUpdatedBy,
		LastUpdated:   userKYCStatus.LastUpdatedDate,
	}, nil
}
