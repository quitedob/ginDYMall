package dao

import (
	"context"
	"douyin/repository/db/model" // Model package
	"gorm.io/gorm"
)

// UserDAO handles database operations for User model.
type UserDAO struct {
	// No internal state like *gorm.DB; pass DB instance to methods for testability
	// and to allow use with transactions if needed.
	// If your project convention is to store DB in DAO, adjust accordingly.
}

// NewUserDAO creates a new UserDAO.
// ctx is included for consistency if other DAO methods might need it for cancellation/deadlines.
func NewUserDAO(_ context.Context) *UserDAO {
	return &UserDAO{}
}

// CreateUser creates a new user record in the database.
// It's better to pass *gorm.DB to DAO methods for testability and flexibility.
func (d *UserDAO) CreateUser(db *gorm.DB, user *model.User) (*model.User, error) {
	err := db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (d *UserDAO) GetUserByID(db *gorm.DB, id uint) (*model.User, error) {
	var user model.User
	err := db.First(&user, id).Error // First finds record by primary key
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username.
func (d *UserDAO) GetUserByUsername(db *gorm.DB, username string) (*model.User, error) {
	var user model.User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
