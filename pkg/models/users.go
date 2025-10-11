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

	IsBlocked    bool   `gorm:"default:false" json:"is_blocked"`
	IsVerified   bool   `gorm:"default:false" json:"is_verified"`
	ProfilePic   string `json:"profile_pic"`
	ReferralCode string `gorm:"unique" json:"referral_code"`

	CustomerProfile CustomerProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"customer_profile,omitempty"`
	// Profile for a trader (pointer because not all users are traders)
	TraderProfile *TraderProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"trader_profile,omitempty"`

	// Generic subscriptions a user might have (e.g., to platform features)
	Subscriptions []UserSubscription `gorm:"foreignKey:UserID" json:"subscriptions,omitempty"`

	TraderSubscriptionPlans []TraderSubscriptionPlan `gorm:"foreignKey:TraderID;constraint:OnDelete:CASCADE;" json:"trader_subscription_plans,omitempty"` // Renamed json tag for clarity

	CustomerTraderSubscriptions []CustomerTraderSubscription `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;" json:"customer_trader_subscriptions,omitempty"`

	Wallet Wallet `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"wallet,omitempty"`

	// Trades initiated by this user if they are a trader
	Trades []Trade `gorm:"foreignKey:TraderID;constraint:OnDelete:SET NULL;" json:"trades,omitempty"`

	// Performance metrics if this user is a trader
	TraderPerformance *TraderPerformance `gorm:"foreignKey:TraderID;constraint:OnDelete:CASCADE;" json:"trader_performance,omitempty"`

	Notifications []Notification `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"notifications,omitempty"`

	// Referrals made by this user (this user is the referrer)
	Referrals []Referral `gorm:"foreignKey:ReferrerID;constraint:OnDelete:SET NULL;" json:"referrals,omitempty"`
	// If this user was referred by someone (this user is the referee)
	ReferredBy *Referral `gorm:"foreignKey:RefereeID;constraint:OnDelete:SET NULL;" json:"referred_by,omitempty"`

	AdminActionLogs []AdminActionLog `gorm:"foreignKey:AdminID;constraint:OnDelete:SET NULL;" json:"admin_action_logs,omitempty"`
}

// AfterCreate is a GORM hook that runs after a new User record is created.
// It automatically creates a new Wallet for the user.
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	wallet := Wallet{
		UserID:   u.ID,
		Balance:  0,
		Currency: "USD", // Default currency
	}
	if err := tx.Create(&wallet).Error; err != nil {
		return err
	}
	return nil
}

// SetPassword hashes the given password and stores it in the User struct.
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compares a plaintext password with the stored hashed password.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil // Returns true if passwords match, false otherwise
}

// IsAdmin returns true if the user's role is RoleAdmin.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsTrader() bool {
	return u.Role == RoleTrader
}

// IsCustomer returns true if the user's role is RoleCustomer.
func (u *User) IsCustomer() bool {
	return u.Role == RoleCustomer
}
