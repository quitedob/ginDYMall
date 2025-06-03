package dao

import (
	"context"
	"douyin/repository/db/model"

	"gorm.io/gorm"
)

// AddressDao defines the data access object for address operations.
type AddressDao struct {
	db *gorm.DB
}

// NewAddressDao creates a new AddressDao instance.
func NewAddressDao(db *gorm.DB) *AddressDao {
	return &AddressDao{
		db: db,
	}
}

// GetAddressByID retrieves a specific address by its ID and UserID.
// This ensures a user can only fetch their own addresses.
func (dao *AddressDao) GetAddressByID(ctx context.Context, userID uint, addressID uint) (*model.Address, error) {
	var address model.Address
	err := dao.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", addressID, userID).
		First(&address).Error

	if err != nil {
		// gorm.ErrRecordNotFound is returned if no record is found.
		// Other errors could indicate database issues.
		return nil, err
	}
	return &address, nil
}

// TODO: Add other CRUD operations for addresses as needed, e.g.:
// CreateAddress(ctx context.Context, address *model.Address) error
// ListAddressesByUserID(ctx context.Context, userID uint) ([]model.Address, error)
// UpdateAddress(ctx context.Context, address *model.Address) error
// DeleteAddress(ctx context.Context, userID uint, addressID uint) error
