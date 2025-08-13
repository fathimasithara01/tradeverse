package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID uint            `json:"customerId"`
	Customer   CustomerProfile `gorm:"foreignKey:CustomerID"`
	OrderItems []OrderItem     `gorm:"foreignKey:OrderID" json:"orderItems"` // Has Many relationship

	ShippingAddressID uint    `json:"shippingAddressId"`
	ShippingAddress   Address `gorm:"foreignKey:ShippingAddressID"` // Has One relationship

	OrderDate   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"orderDate"`
	TotalAmount float64   `gorm:"type:decimal(10,2)" json:"totalAmount"`

	Status string `gorm:"default:'Pending'" json:"status"`

	PaymentMethod string `json:"paymentMethod"`                          // e.g., "COD", "Stripe"
	PaymentStatus string `gorm:"default:'Pending'" json:"paymentStatus"` // e.g., "Pending", "Paid", "Failed"
	PaymentID     string `json:"paymentId,omitempty"`                    // Optional: To store the transaction ID from a payment gateway
}

type OrderItem struct {
	gorm.Model

	OrderID uint `json:"orderId"`

	ProductID uint    `json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID"`

	Quantity int `json:"quantity"`

	PriceAtPurchase float64 `gorm:"type:decimal(10,2)" json:"priceAtPurchase"`
}

type Address struct {
	gorm.Model
	CustomerID uint
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}
