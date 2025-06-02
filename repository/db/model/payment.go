package model

import "time"

// Payment represents a payment record for an order.
type Payment struct {
	ID        uint      `gorm:"primaryKey"`
	OrderID   string    // Associated Order ID (assuming OrderID in Order model is string, adjust if necessary)
	Amount    float64   // Payment amount
	Status    string    // Payment status, e.g., "UNPAID", "PAID", "FAILED"
	CreatedAt time.Time // Timestamp of payment creation
	// UpdatedAt time.Time // Optionally, if payment status can change
	// TransactionID string // Optional: payment gateway transaction ID
}
