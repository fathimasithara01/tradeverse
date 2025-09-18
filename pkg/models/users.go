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
	Password string `gorm:"size:255;not null" json:"-"` // This field stores the HASHED password

	Role UserRole `gorm:"type:varchar(20);not null;default:'customer'" json:"role"`

	RoleID *uint `json:"role_id"`

	RoleModel Role `gorm:"foreignKey:RoleID" json:"role_model,omitempty"` // Renamed from RoleAssociation to RoleModel for clarity

	IsBlocked bool `gorm:"default:false" json:"is_blocked"`

	CustomerProfile     CustomerProfile      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"customer_profile,omitempty"`
	TraderProfile       TraderProfile        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"trader_profile,omitempty"`
	Subscriptions       []Subscription       `gorm:"foreignKey:UserID" json:"subscriptions,omitempty"`
	TraderSubscriptions []TraderSubscription `gorm:"foreignKey:UserID" json:"trader_subscriptions,omitempty"`
	Wallet              Wallet               `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"wallet,omitempty"`
}

// AfterCreate hook to create a wallet for every new user.
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	wallet := Wallet{
		UserID:   u.ID,
		Balance:  0,
		Currency: "INR", // Default currency
	}
	if err := tx.Create(&wallet).Error; err != nil {
		return err
	}
	return nil
}

// SetPassword hashes the user's password using bcrypt.
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword) // Store the HASH
	return nil
}

// CheckPassword verifies the provided password against the stored hashed password.
func (u *User) CheckPassword(password string) bool {
	// Debug logging for password comparison
	// fmt.Printf("[DEBUG-PASSWORD] User %s: Received plain password: '%s'\n", u.Email, password)
	// fmt.Printf("[DEBUG-PASSWORD] User %s: Stored hashed password: '%s'\n", u.Email, u.Password)

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// fmt.Printf("[DEBUG-PASSWORD-ERROR] bcrypt comparison failed for user '%s': %v \n", u.Email, err)
		return false
	}
	return true
}

func (u *User) IsAdmin() bool    { return u.Role == RoleAdmin }
func (u *User) IsTrader() bool   { return u.Role == RoleTrader }
func (u *User) IsCustomer() bool { return u.Role == RoleCustomer }
