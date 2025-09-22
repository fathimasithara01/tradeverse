package service

import (
	"errors"
	"time"

	"github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"
	"github.com/fathimasithara01/tradeverse/pkg/models"
)

type KYCServicer interface {
	SubmitKYCDocuments(userID uint, docType, docURL string) error
	GetKYCStatus(userID uint) (*models.KYCStatusResponse, error)
}

type kycService struct {
	kycRepo customerrepo.KYCRepository
}

func NewKYCService(kycRepo customerrepo.KYCRepository) KYCServicer {
	return &kycService{kycRepo: kycRepo}
}

func (s *kycService) SubmitKYCDocuments(userID uint, docType, docURL string) error {
	newDoc := &models.KYCDocument{
		UserID:             userID,
		DocumentType:       docType,
		DocumentURL:        docURL,
		VerificationStatus: models.KYCStatusPending,
	}
	if err := s.kycRepo.CreateKYCDocument(newDoc); err != nil {
		return errors.New("failed to save KYC document: " + err.Error())
	}

	userKYCStatus, err := s.kycRepo.FindUserKYCStatus(userID)
	if err != nil {
		return errors.New("failed to retrieve user KYC status: " + err.Error())
	}

	if userKYCStatus == nil {
		newUserStatus := &models.UserKYCStatus{
			UserID:          userID,
			Status:          models.KYCStatusPending,
			LastUpdatedBy:   userID,
			LastUpdatedDate: time.Now(),
		}
		if err := s.kycRepo.CreateUserKYCStatus(newUserStatus); err != nil {
			return errors.New("failed to create user KYC status: " + err.Error())
		}
	} else if userKYCStatus.Status != models.KYCStatusApproved {
		userKYCStatus.Status = models.KYCStatusPending
		userKYCStatus.Reason = ""
		userKYCStatus.LastUpdatedBy = userID
		userKYCStatus.LastUpdatedDate = time.Now()
		if err := s.kycRepo.UpdateUserKYCStatus(userKYCStatus); err != nil {
			return errors.New("failed to update user KYC status: " + err.Error())
		}
	}

	return nil
}

func (s *kycService) GetKYCStatus(userID uint) (*models.KYCStatusResponse, error) {
	userKYCStatus, err := s.kycRepo.FindUserKYCStatus(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve KYC status from repository: " + err.Error())
	}

	if userKYCStatus == nil {
		return &models.KYCStatusResponse{
			Status:      models.KYCStatusNotSubmitted,
			LastUpdated: time.Now(), 
		}, nil
	}

	return &models.KYCStatusResponse{
		Status:        userKYCStatus.Status,
		Reason:        userKYCStatus.Reason,
		LastUpdatedBy: userKYCStatus.LastUpdatedBy,
		LastUpdated:   userKYCStatus.LastUpdatedDate,
	}, nil
}
