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

	RoleID    *uint
	IsBlocked bool `gorm:"default:false" json:"is_blocked"`

	CustomerProfile     CustomerProfile      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"customer_profile,omitempty"`
	TraderProfile       TraderProfile        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"trader_profile,omitempty"`
	Subscriptions       []Subscription       `gorm:"foreignKey:UserID" json:"subscriptions,omitempty"`
	TraderSubscriptions []TraderSubscription `gorm:"foreignKey:UserID" json:"trader_subscriptions,omitempty"`
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	// Ensure a wallet is created for every new user
	wallet := Wallet{
		UserID:   u.ID,
		Balance:  0,
		Currency: "INR", // Default currency
	}
	if err := tx.Create(&wallet).Error; err != nil {
		return err
	}

	// Create a CustomerProfile by default for new users
	if u.Role == RoleCustomer {
		customerProfile := CustomerProfile{
			UserID: u.ID,
			// ShippingAddress: "", // Can be filled later
			// PhoneNumber:     "", // Can be filled later
		}
		if err := tx.Create(&customerProfile).Error; err != nil {
			return err
		}
	}
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
func (u *User) IsAdmin() bool    { return u.Role == RoleAdmin }
func (u *User) IsTrader() bool   { return u.Role == RoleTrader }
func (u *User) IsCustomer() bool { return u.Role == RoleCustomer }
