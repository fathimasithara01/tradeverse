package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRole defines the possible roles a user can have.
type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
	RoleTrader   UserRole = "trader"
)

// User represents a user in the system.
type User struct {
	gorm.Model
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"` // Stored hashed, excluded from JSON output
	Phone    string `json:"phone"`

	Role UserRole `gorm:"type:varchar(20);not null;default:'customer'" json:"role"`

	RoleID    *uint `json:"role_id"`
	RoleModel Role  `gorm:"foreignKey:RoleID" json:"role_model,omitempty"` // Assuming Role struct exists

	IsBlocked  bool   `gorm:"default:false" json:"is_blocked"`
	IsVerified bool   `gorm:"default:false" json:"is_verified"`
	ProfilePic string `json:"profile_pic"`
	// ReferralCode string `gorm:"unique" json:"referral_code"`

	CustomerProfile CustomerProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"customer_profile,omitempty"`
	TraderProfile   *TraderProfile  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"trader_profile,omitempty"`

	Subscriptions []UserSubscription `gorm:"foreignKey:UserID" json:"subscriptions,omitempty"`

	TraderSubscriptionPlans []TraderSubscriptionPlan `gorm:"foreignKey:TraderID;constraint:OnDelete:CASCADE;" json:"trader_subscription_plans,omitempty"` // Renamed json tag for clarity

	CustomerTraderSubscriptions []CustomerTraderSubscription `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;" json:"customer_trader_subscriptions,omitempty"`

	Wallet Wallet `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"wallet,omitempty"`

	Trades []Trade `gorm:"foreignKey:TraderID;constraint:OnDelete:SET NULL;" json:"trades,omitempty"`

	TraderPerformance *TraderPerformance `gorm:"foreignKey:TraderID;constraint:OnDelete:CASCADE;" json:"trader_performance,omitempty"`

	Notifications []Notification `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"notifications,omitempty"`

	Referrals  []Referral `gorm:"foreignKey:ReferrerID;constraint:OnDelete:SET NULL;" json:"referrals,omitempty"`
	ReferredBy *Referral  `gorm:"foreignKey:RefereeID;constraint:OnDelete:SET NULL;" json:"referred_by,omitempty"`

	AdminActionLogs []AdminActionLog `gorm:"foreignKey:AdminID;constraint:OnDelete:SET NULL;" json:"admin_action_logs,omitempty"`
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	wallet := Wallet{
		UserID:   u.ID,
		Balance:  0,
		Currency: "USD",
	}
	if err := tx.Create(&wallet).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsTrader() bool {
	return u.Role == RoleTrader
}

func (u *User) IsCustomer() bool {
	return u.Role == RoleCustomer
}
