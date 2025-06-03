package dao_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"douyin/repository/db/dao" // Your DAO package
	"douyin/repository/db/model" // Your model package
	// "gorm.io/gorm" // Already available via testDB in common_dao_test.go
)

// TestCreateUser tests the CreateUser and GetUserByID DAO functions.
func TestUserDAO_CreateAndGet(t *testing.T) {
	require.NotNil(t, testDB, "testDB is not initialized. TestMain in common_dao_test.go should have initialized it.")

	// Clear users table before this specific test execution for isolation.
	// Note: TestMain already runs, this is for per-test cleanup if needed.
	// clearTables(testDB, "users") // users is the GORM default table name for model.User

	userDAO := dao.NewUserDAO(context.Background()) // NewUserDAO doesn't take DB

	uniqueSuffix := time.Now().UnixNano()
	testUsername := fmt.Sprintf("testuser_%d", uniqueSuffix)
	testEmail := fmt.Sprintf("testuser_%d@example.com", uniqueSuffix)

	newUser := &model.User{
		Username:      testUsername,
		PasswordHash:  "hashedpassword_test", // Example
		Email:         testEmail,
		FollowCount:   0, // Default value
		FollowerCount: 0, // Default value
		// Initialize other required fields for your User model if any
		// Avatar: "default_avatar.png",
		// BackgroundImage: "default_bg.jpg",
		// Signature: "Test signature",
	}

	// 1. Test CreateUser
	createdUser, err := userDAO.CreateUser(testDB, newUser)
	require.NoError(t, err, "CreateUser should not return an error for a new user")
	require.NotNil(t, createdUser, "CreateUser result should not be nil")
	assert.Greater(t, createdUser.ID, uint(0), "CreatedUser ID should be positive after creation")
	assert.Equal(t, newUser.Username, createdUser.Username, "Username of created user should match input")
	assert.Equal(t, newUser.Email, createdUser.Email, "Email of created user should match input")

	// 2. Test GetUserByID
	fetchedUserByID, err := userDAO.GetUserByID(testDB, createdUser.ID)
	require.NoError(t, err, "GetUserByID should not return an error for an existing user ID")
	require.NotNil(t, fetchedUserByID, "GetUserByID result should not be nil for an existing user ID")
	assert.Equal(t, createdUser.ID, fetchedUserByID.ID, "Fetched user ID by ID should match created user's ID")
	assert.Equal(t, testUsername, fetchedUserByID.Username, "Fetched user username by ID should match")

	// 3. Test GetUserByUsername
	fetchedUserByUsername, err := userDAO.GetUserByUsername(testDB, testUsername)
	require.NoError(t, err, "GetUserByUsername should not return an error for an existing username")
	require.NotNil(t, fetchedUserByUsername, "GetUserByUsername result should not be nil for an existing username")
	assert.Equal(t, createdUser.ID, fetchedUserByUsername.ID, "Fetched user ID by username should match created user's ID")
	assert.Equal(t, testUsername, fetchedUserByUsername.Username, "Fetched user username by username should match")

	// 4. Test creating a user with a duplicate username (should fail due to unique constraint)
	duplicateUser := &model.User{
		Username:      testUsername, // Same username
		PasswordHash:  "anotherpassword",
		Email:         fmt.Sprintf("testuser_dup_%d@example.com", uniqueSuffix),
	}
	_, err = userDAO.CreateUser(testDB, duplicateUser)
	assert.Error(t, err, "CreateUser with duplicate username should return an error")
	// Note: The exact error type can be checked if GORM/driver provides a way to identify unique constraint violations.
	// For example, errors.Is(err, gorm.ErrDuplicatedKey) or checking underlying driver error.

	// 5. Test GetUserByID for a non-existent user
	nonExistentID := uint(9999999)
	_, err = userDAO.GetUserByID(testDB, nonExistentID)
	assert.Error(t, err, "GetUserByID for a non-existent ID should return an error")
	// assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "Error should be gorm.ErrRecordNotFound for non-existent ID") // More specific check

	// 6. Test GetUserByUsername for a non-existent username
	nonExistentUsername := "iamnotauserthatshouldexist"
	_, err = userDAO.GetUserByUsername(testDB, nonExistentUsername)
	assert.Error(t, err, "GetUserByUsername for a non-existent username should return an error")
	// assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "Error should be gorm.ErrRecordNotFound for non-existent username")

	// Cleanup: Delete the created user (optional, as TestMain might drop/recreate DB or tables)
	// If not, manual cleanup is good practice for test isolation.
	// err = testDB.Unscoped().Delete(&model.User{}, createdUser.ID).Error
	// require.NoError(t, err, "Failed to clean up test user")
}
