package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
	RoleTrader   UserRole = "trader"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`

	Role UserRole `gorm:"type:varchar(20);not null;default:'customer'" json:"role"`

	RoleID *uint `json:"role_id"`

	RoleModel Role `gorm:"foreignKey:RoleID" json:"role_model,omitempty"`

	IsBlocked bool `gorm:"default:false" json:"is_blocked"`

	CustomerProfile     CustomerProfile      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"customer_profile,omitempty"`
	TraderProfile       *TraderProfile        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"trader_profile,omitempty"`
	Subscriptions       []Subscription       `gorm:"foreignKey:UserID" json:"subscriptions,omitempty"`
	TraderSubscriptions []TraderSubscription `gorm:"foreignKey:UserID" json:"trader_subscriptions,omitempty"`
	Wallet              Wallet               `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"wallet,omitempty"`

	Trades []Trade `gorm:"foreignKey:TraderID" json:"trades,omitempty"` // Trades initiated by this user if they are a trader
	// CopiedTrades      []CopyTrade           `gorm:"foreignKey:CustomerID" json:"copied_trades,omitempty"`    // Trades copied by this user if they are a customer
	// CopySettings      []CustomerCopySetting `gorm:"foreignKey:CustomerID" json:"copy_settings,omitempty"`    // Copy settings for various traders
	TraderPerformance *TraderPerformance `gorm:"foreignKey:TraderID" json:"trader_performance,omitempty"` // Performance metrics if this user is a trader
	Notifications     []Notification     `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
	Referrals         []Referral         `gorm:"foreignKey:ReferrerID" json:"referrals,omitempty"`  // Referrals made by this user
	ReferredBy        *Referral          `gorm:"foreignKey:RefereeID" json:"referred_by,omitempty"` // If this user was referred by someone
	AdminActionLogs   []AdminActionLog   `gorm:"foreignKey:AdminID" json:"admin_action_logs,omitempty"`
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
	if err != nil {
		return false
	}
	return true
}

func (u *User) IsAdmin() bool    { return u.Role == RoleAdmin }
func (u *User) IsTrader() bool   { return u.Role == RoleTrader }
func (u *User) IsCustomer() bool { return u.Role == RoleCustomer }
