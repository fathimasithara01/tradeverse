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

	// FIXED: no foreignKey here (it broke GORM)
	Role   UserRole `gorm:"type:varchar(20);not null;default:'customer'" json:"role"`
	RoleID *uint    `gorm:"index" json:"role_id"`

	IsBlocked bool `gorm:"default:false" json:"is_blocked"`

	CustomerProfile     CustomerProfile      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"customer_profile,omitempty"`
	TraderProfile       TraderProfile        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"trader_profile,omitempty"`
	Subscriptions       []Subscription       `gorm:"foreignKey:UserID" json:"subscriptions,omitempty"`
	TraderSubscriptions []TraderSubscription `gorm:"foreignKey:UserID" json:"trader_subscriptions,omitempty"`
}

// Before creating user → hash password
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password != "" {
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

// After creating user → auto-create wallet
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	wallet := Wallet{
		UserID:   u.ID,
		Balance:  0,
		Currency: "INR",
	}
	if err := tx.Create(&wallet).Error; err != nil {
		return err
	}
	return nil
}

// Helpers
func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
func (u *User) IsAdmin() bool    { return u.Role == RoleAdmin }
func (u *User) IsTrader() bool   { return u.Role == RoleTrader }
func (u *User) IsCustomer() bool { return u.Role == RoleCustomer }
