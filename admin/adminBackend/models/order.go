package models

import (
	"time"

	"gorm.io/gorm"
)

// Order represents a customer's order in the system.
type Order struct {
	gorm.Model // Includes ID, CreatedAt, UpdatedAt, DeletedAt

	// Foreign key to the Customer who placed the order.
	CustomerID uint            `json:"customerId"`
	Customer   CustomerProfile `gorm:"foreignKey:CustomerID"` // Belongs to relationship

	// A list of all products in this order.
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"orderItems"` // Has Many relationship

	// Foreign key for the shipping address.
	ShippingAddressID uint    `json:"shippingAddressId"`
	ShippingAddress   Address `gorm:"foreignKey:ShippingAddressID"` // Has One relationship

	OrderDate   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"orderDate"`
	TotalAmount float64   `gorm:"type:decimal(10,2)" json:"totalAmount"`

	// Status of the order (e.g., "Pending", "Processing", "Shipped", "Delivered", "Cancelled")
	Status string `gorm:"default:'Pending'" json:"status"`

	// Payment details
	PaymentMethod string `json:"paymentMethod"`                          // e.g., "COD", "Stripe"
	PaymentStatus string `gorm:"default:'Pending'" json:"paymentStatus"` // e.g., "Pending", "Paid", "Failed"
	PaymentID     string `json:"paymentId,omitempty"`                    // Optional: To store the transaction ID from a payment gateway
}

// OrderItem represents a single product within an order.
type OrderItem struct {
	gorm.Model // Includes its own ID, CreatedAt, etc.

	// Foreign key back to the Order.
	OrderID uint `json:"orderId"`

	// Foreign key to the Product.
	ProductID uint    `json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID"` // Belongs to relationship

	Quantity int `json:"quantity"`

	// Price of the product at the time of purchase.
	// It's crucial to store this here so historical orders are not affected by future price changes.
	PriceAtPurchase float64 `gorm:"type:decimal(10,2)" json:"priceAtPurchase"`
}

// You would also need an Address model like this for the relationship to work.
type Address struct {
	gorm.Model
	CustomerID uint
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// And ensure your Customer and Product models are defined.
/*
type Customer struct {
    gorm.Model
    // ... other customer fields
}

type Product struct {
    gorm.Model
    // ... other product fields
}
*/
