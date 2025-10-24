package models

import (
	"time"

	"gorm.io/gorm"
)

type CommissionSetting struct {
	gorm.Model
	Key         string    `gorm:"uniqueIndex;not null;size:50" json:"key"`
	Value       float64   `gorm:"type:numeric(5,2);not null" json:"value"` // e.g., 10.50 for 10.50%
	Description string    `gorm:"type:text" json:"description,omitempty"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`
	UpdatedBy   uint      `json:"updated_by,omitempty"`
}

type AdminCommissionRequestPayload struct {
	CommissionPercentage float64 `json:"commission_percentage" binding:"required,min=0,max=100"`
}

type AdminCommissionResponsePayload struct {
	ID                   uint      `json:"id"`
	CommissionPercentage float64   `json:"commission_percentage"`
	LastUpdated          time.Time `json:"last_updated"`
	UpdatedBy            uint      `json:"updated_by,omitempty"`
	Description          string    `json:"description,omitempty"`
}
