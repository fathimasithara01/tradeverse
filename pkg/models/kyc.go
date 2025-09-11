package models

import (
	"time"

	"gorm.io/gorm"
)

type KYCDocument struct {
	gorm.Model
	UserID             uint   `gorm:"index;not null"`
	DocumentType       string `gorm:"size:50;not null"`
	DocumentURL        string `gorm:"size:255;not null"`
	VerificationStatus string `gorm:"size:20;default:'PENDING'"`
	AdminNotes         string `gorm:"type:text"`
}

type UserKYCStatus struct {
	gorm.Model
	UserID          uint   `gorm:"uniqueIndex;not null"`
	Status          string `gorm:"size:20;default:'NOT_SUBMITTED'"`
	Reason          string `gorm:"type:text"`
	LastUpdatedBy   uint
	LastUpdatedDate time.Time
}

const (
	KYCStatusNotSubmitted = "NOT_SUBMITTED"
	KYCStatusPending      = "PENDING"
	KYCStatusApproved     = "APPROVED"
)

type SubmitKYCRequest struct {
	DocumentType string `json:"document_type" binding:"required"`
	DocumentURL  string `json:"document_url" binding:"required,url"`
}

type KYCStatusResponse struct {
	Status        string    `json:"status"`
	Reason        string    `json:"reason,omitempty"`
	LastUpdatedBy uint      `json:"last_updated_by,omitempty"`
	LastUpdated   time.Time `json:"last_updated"`
}
