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
	Password string `gorm:"default:false;size:255;not null" json:"-"`

	Role   UserRole `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
	RoleID *uint    `gorm:"index" json:"role_id"`

	IsBlocked bool `gorm:"default:false" json:"is_blocked"`

	CustomerProfile CustomerProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"customer_profile,omitempty"`
	TraderProfile   TraderProfile   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"trader_profile,omitempty"`

	Subscriptions []Subscription `gorm:"foreignKey:UserID" json:"subscriptions,omitempty"`

	TraderSubscriptions []TraderSubscription `gorm:"foreignKey:UserID" json:"trader_subscriptions,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password != "" { // Only hash if password is provided
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	if u.Role == "" {
		u.Role = RoleCustomer
	}
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
