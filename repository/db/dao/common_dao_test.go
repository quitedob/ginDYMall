package dao_test // Using _test package to avoid circular dependencies if DAO uses models from its own package directly

import (
	"fmt"
	"os"
	"testing"
	// "time" // Not directly used in this common setup, but often useful

	"gorm.io/driver/mysql"
	// "gorm.io/driver/sqlite" // Alternative for faster, simpler tests if compatible
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // GORM's logger
	"douyin/global"       // For global.DB
	"douyin/conf"         // Assuming access to config for DSN or project structure
	"douyin/repository/db/model" // For migrating all models
	"douyin/pkg/utils/log" // Project's logger utility, assuming InitLoggerForTest exists or can be added
)

var testDB *gorm.DB

// InitLoggerForTest is a placeholder for a function that would initialize
// the project's logger specifically for test environments (e.g., different format, level).
// If your mylog package doesn't have this, you might need to add it or adjust.
func InitLoggerForTest() {
    // Example: If your logger has a SetLevel or similar configuration
    // mylog.SetLevel("debug")
    // mylog.SetOutput(os.Stdout) // Or a test-specific log file
    // For now, assume mylog.InitLogger() is safe to call or adapts.
    // If log.InitLogger is already called by conf.InitConfig, this might not be needed.
    // Let's assume `log.InitLogger()` is sufficient or a specific test init is added to the log package.
    // For the purpose of this task, if `log.InitLoggerForTest()` is not defined in `douyin/pkg/utils/log`,
    // we'll call the standard `log.InitLogger()`.
    // It's better if the logger package itself handles the test context.
    // pkgLog.InitLogger() // Assuming pkgLog is an alias to "douyin/pkg/utils/log"
    // This will be called from setupTestDB
}


// setupTestDB initializes a test database connection.
func setupTestDB() (func(), error) { // Return error for better handling in TestMain
	// Attempt to initialize logger for test context.
    // This might need to be adapted based on your actual logger package.
    // If InitLogger is part of your conf.InitConfig, ensure that's test-friendly.
    // For now, assuming a generic Init or one that can be called multiple times.
    // Let's assume `log.InitLogger()` can be called or there's a test-specific one.
    // log.InitLogger() // This should ideally be idempotent or test-specific.
    // The task used mylog.InitLoggerForTest(), so we'll assume that exists.
    // If not, this would be a point of failure or require adjustment.
    // For now, let's assume it's available in "douyin/pkg/utils/log"
    log.InitLogger() // Using the existing logger init. Best to have a test-specific one.


	dsn := os.Getenv("MYSQL_TEST_DSN")
	if dsn == "" {
		// Fallback to project config if available and DSN is for a test DB
		// This logic is complex and error-prone for generic projects.
		// Prioritizing environment variable for CI/testing is best.
		if conf.GlobalConfig != nil && conf.GlobalConfig.MySql.Default.DbName != "" && !strings.Contains(conf.GlobalConfig.MySql.Default.DbName, "test") {
			mylog.Warnf("Project's default DB '%s' does not seem to be a test database. Using hardcoded fallback DSN.", conf.GlobalConfig.MySql.Default.DbName)
			// Hardcoded fallback DSN for local testing if MYSQL_TEST_DSN is not set.
			// Replace 'yourpassword' and 'douyin_test' as appropriate.
			dsn = "root:testpassword@tcp(127.0.0.1:3306)/douyin_test?charset=utf8mb4&parseTime=True&loc=Local"
		} else if conf.GlobalConfig != nil && conf.GlobalConfig.MySql.Default.DSN() != "" { // Assuming DSN() method exists
            dsn = conf.GlobalConfig.MySql.Default.DSN() // Use configured DSN, assuming it's for test
            if !strings.Contains(dsn, "test") { // Simple check
                 mylog.Warnf("Using DSN from config: %s. Ensure this is a TEST database.", dsn)
            }
        } else {
			// Final fallback if no other DSN is found
			dsn = "root:testpassword@tcp(127.0.0.1:3306)/douyin_test?charset=utf8mb4&parseTime=True&loc=Local"
			mylog.Warnf("MYSQL_TEST_DSN and project config DSN not suitable/found. Using hardcoded fallback DSN for testing: %s", dsn)
		}
	}

	var err error
	gormLogger := logger.New(
        stdLog.New(os.Stdout, "\r\n", stdLog.LstdFlags), // io writer
        logger.Config{
            SlowThreshold:             time.Second,   // Slow SQL threshold
            LogLevel:                  logger.Silent, // Log level (Silent, Error, Warn, Info)
            IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
            Colorful:                  false,         // Disable color
        },
    )

	testDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger, // Use customized logger
	})
	if err != nil {
		mylog.Fatalf("Failed to connect to test database using DSN '%s': %v", dsn, err)
		return nil, err
	}

	// Assign to global.DB as per task's suggestion for API tests later
	global.DB = testDB
	mylog.Info("Test database connected and global.DB assigned.")

	// Auto-migrate all known models
	allModels := []interface{}{
		&model.User{},
		&model.Address{},
		&model.Cart{},
		&model.Category{},
		// &model.Checkout{}, // Assuming Checkout might not be a direct GORM table model based on typical naming
		&model.Order{},
		&model.OrderItem{},
		&model.Payment{},
		&model.Product{},
		&model.ProductCategory{},
		// RBAC models are added next
	}
    // Add RBAC models (Role, Permission, RolePermission, UserRole)
    allModels = append(allModels, model.GetRBACModels()...)
    // Note: If GetRBACModels() already includes User, it might be duplicated, but AutoMigrate handles this.

	err = testDB.AutoMigrate(allModels...)
	if err != nil {
		mylog.Fatalf("Failed to auto-migrate tables for test: %v. Models: %v", err, allModels)
		return nil, err
	}
	mylog.Info("Test database schema migrated.")

	// Teardown function to close DB connection
	teardown := func() {
		sqlDB, _ := testDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		mylog.Info("Test database connection closed.")
	}
	return teardown, nil
}

// TestMain can be used for package-level setup/teardown
func TestMain(m *testing.M) {
    // Call conf.InitConfig() if it's necessary for obtaining DB DSN or other settings
    // Ensure it's test-friendly (e.g., loads a test config or uses test env vars)
    // if err := conf.InitConfig(); err != nil {
    //     stdLog.Fatalf("Failed to initialize config for tests: %v", err)
    // }
    // mylog.InitLogger() // Initialize logger based on loaded config (if any test specific settings)


	teardown, err := setupTestDB()
	if err != nil {
		// Error already logged by setupTestDB
		os.Exit(1)
	}

	exitCode := m.Run()

	if teardown != nil {
		teardown()
	}
	os.Exit(exitCode)
}

// Helper to clear tables (basic version, might need adjustment for foreign keys)
func clearTables(db *gorm.DB, tableNames ...string) {
	if db == nil {
		mylog.Error("clearTables called with nil db instance.")
		return
	}
	for _, tableName := range tableNames {
		// Using Exec to run raw SQL for TRUNCATE or DELETE
		// TRUNCATE is faster but might be blocked by foreign key constraints
		// DELETE FROM is safer with foreign keys but can be slower
		// For simplicity, trying TRUNCATE first, then DELETE as a fallback concept
		// Note: Specific model structs are not used here, relying on table names.
		err := db.Exec(fmt.Sprintf("DELETE FROM %s", tableName)).Error
		// err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error // Alternative
		if err != nil {
			mylog.Warnf("Failed to clear table %s: %v. Data might persist between tests.", tableName, err)
		}
	}
}

// Ensure imports for standard log and strings are present
import (
	stdLog "log"
	"strings"
	"time"
)
