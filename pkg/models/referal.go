package models

import "gorm.io/gorm"

type Referral struct {
	gorm.Model
	ReferrerID       uint    `gorm:"not null;index" json:"referrer_id"`
	Referrer         User    `gorm:"foreignKey:ReferrerID" json:"referrer"`
	RefereeID        *uint   `gorm:"uniqueIndex" json:"referee_id"`
	Referee          *User   `gorm:"foreignKey:RefereeID" json:"referee,omitempty"`
	ReferralCode     string  `gorm:"size:50;uniqueIndex;not null" json:"referral_code"`
	Status           string  `gorm:"size:50;not null;default:'pending'" json:"status"`
	CommissionEarned float64 `gorm:"type:numeric(18,4);default:0.00" json:"commission_earned"`
	CommissionRate   float64 `gorm:"type:numeric(5,2);default:0.00" json:"commission_rate"`
}
