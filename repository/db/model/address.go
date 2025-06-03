package model

import "time"

// Address represents a user's address in the database.
type Address struct {
	ID            uint   `gorm:"primaryKey"`
	UserID        uint   `gorm:"not null;index"` // Belongs to a user, add index for faster lookups
	FirstName     string `gorm:"size:100"`
	LastName      string `gorm:"size:100"`
	StreetAddress string `gorm:"size:255"`
	City          string `gorm:"size:100"`
	State         string `gorm:"size:100"` // Province/State
	Country       string `gorm:"size:100"`
	ZipCode       string `gorm:"size:20"`
	Email         string `gorm:"size:255"` // Optional: Email associated with this specific address record
	PhoneNumber   string `gorm:"size:20"`  // Optional: Phone number for this address
	IsDefault     bool   `gorm:"default:false"` // Optional: Mark as default address for the user
	// UserCurrency is likely a user-level preference, not per address.
	// It was in the previous DAO DTO, so consider if it's truly address-specific.
	// If it is, it can be added here. For now, assuming it's not strictly part of an address model.
	// UserCurrency  string `gorm:"size:10"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
