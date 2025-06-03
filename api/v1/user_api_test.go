package v1_test

import (
	"bytes"
	"context" // Required for dao.NewUserDAO if it keeps this param
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// Adjust these imports based on your actual project structure
	"douyin/api/v1" // To potentially access handler functions directly if needed, or types
	"douyin/routes" // For routes.NewRouter
	"douyin/conf"   // For conf.InitConfig (if needed and test-friendly)
	"douyin/global" // For global.DB
	"douyin/middleware" // To add other middlewares if NewRouter doesn't include all globals
	"douyin/pkg/utils/log" // Project's logger
	daoTest "douyin/repository/db/dao_test" // To ensure TestMain runs via blank import (convention)

	// These are needed if you directly interact with DAO/model in your test setup for this API test
	// "douyin/repository/db/dao"
	// "douyin/repository/db/model"

	// Required for blank import to trigger TestMain in common_dao_test.go
	_ "douyin/repository/db/dao_test"
)

var testRouter *gin.Engine

// setupTestRouter initializes the Gin engine with all routes for testing.
func setupTestRouter(t *testing.T) *gin.Engine {
	// Ensure logger is test-friendly.
	// Assuming log.InitLogger() is safe or a test-specific version is available.
	log.InitLogger() // Or log.InitLoggerForTest() if you create one

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Configuration: APP_ENV is often used to load different configs.
	// Ensure your conf.InitConfig() respects this or uses a test-specific file.
	// os.Setenv("APP_ENV", "test") // Handled by CI.yml, ensure local test setup matches.
	// if err := conf.InitConfig(); err != nil {
	//     t.Fatalf("Failed to init config for test router: %v", err)
	// }
    // mylog.InitLogger() // Re-init logger if config changed logging settings.

	// Database:
	// The blank import `_ "douyin/repository/db/dao_test"` ensures that TestMain
	// in common_dao_test.go runs. That TestMain calls setupTestDB, which initializes
	// testDB and assigns it to global.DB. So, global.DB should be the test DB here.
	if global.DB == nil {
		// This might happen if the import side-effect isn't working as expected or
		// if dao_test.TestMain didn't complete successfully.
		// For robustness, explicitly call the setup if global.DB is nil,
		// but this means dao_test.TestMain might run its setup twice if not careful.
		// A better way is a getter for testDB from dao_test package or rely on TestMain.
		t.Log("global.DB is nil, attempting to trigger dao_test.TestMain or setup manually")
		// This is tricky. Let's rely on the blank import for now.
		// If this fails, the test will fail on require.NotNil(global.DB).
		// One could expose `GetTestDB()` from `dao_test` package, but that's not standard for _test packages.
	}
	require.NotNil(t, global.DB, "global.DB (testDB) not initialized. Check dao_test.TestMain.")

	// Redis:
	// Similar to DB, if routes depend on Redis (e.g., cache.Rdb for rate limiting),
	// it should be initialized to a test Redis instance.
	// For now, RateLimitMiddleware has a nil check for cache.Rdb.
	// cache.InitCacheForTest() // You'd need a test initializer for cache.Rdb

	engine := gin.New()

	// Apply global middlewares that are typically added in main.go before NewRouter
	// This order should mimic your main.go setup.
	engine.Use(gin.Recovery()) // Standard recovery
	engine.Use(middleware.RequestID())
	engine.Use(middleware.LoggerMiddleware()) // Your custom logger middleware
	engine.Use(middleware.Cors())
	engine.Use(middleware.SecurityHeadersMiddleware())
	engine.Use(middleware.ContextTimeout(15 * time.Second)) // Example timeout
	engine.Use(middleware.PrometheusMiddleware())
	engine.Use(middleware.DBInjectorMiddleware()) // Injects global.DB into context as "db"
	engine.Use(middleware.I18nMiddleware())
	// AuthMiddleware is usually route-specific, not global, unless all routes are auth'd.
	// ErrorHandler should be late.
	engine.Use(middleware.ErrorHandler())


	// Initialize routes using the test DB (global.DB)
	routes.NewRouter(engine, global.DB) // This function sets up all app routes

	return engine
}

// TestMain for api_v1_test package
func TestMain(m *testing.M) {
	// Perform setup for the API test package
	// This ensures that the router is set up once for all tests in this package.
	// Note: DB setup is handled by TestMain in dao_test package due to blank import.
	// If you need specific setup for API tests beyond DB, do it here.

	// Set environment to test if not already set (e.g. by CI)
	if os.Getenv("APP_ENV") == "" {
		os.Setenv("APP_ENV", "test")
	}
	// Initialize config if it's light and test-friendly
	// Or rely on env vars set by CI/test scripts for DB/Redis DSNs.
	// conf.InitConfig()
	// log.InitLogger()


	// Run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}


func TestUserRegisterAPI(t *testing.T) {
	// Get the router instance. setupTestRouter is called once per package via TestMain logic (implicitly)
	// or explicitly if testRouter is nil.
	// For simplicity and ensuring it's always initialized for each test run if needed:
	if testRouter == nil {
		testRouter = setupTestRouter(t) // Pass `t` for failing fast if setup fails
	}
	require.NotNil(t, testRouter, "Test router could not be initialized")

	// Table cleanup for this specific test.
	// It's better to clean up data created by this test to ensure isolation.
	// You can use a helper from dao_test or write specific cleanup SQL.
	// For User model, table name is likely "users".
	defer func() {
		// This is a simple cleanup. In real tests, you might want more targeted deletion.
		// Or, if each test runs in a transaction that's rolled back (more complex setup).
		// global.DB.Exec("DELETE FROM users WHERE username LIKE 'testapiuser_%'")
	}()


	uniqueSuffix := time.Now().UnixNano()
	username := fmt.Sprintf("testapiuser_%d", uniqueSuffix)
	email := fmt.Sprintf("testapiuser_%d@example.com", uniqueSuffix)
	password := "strongpassword123"

	payload := gin.H{ // Using gin.H for convenience, similar to map[string]interface{}
		"username": username,
		"password": password,
		"email":    email,
		// "confirm_password": password, // Add if your API expects it
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// Print response for debugging if test fails
	if w.Code != http.StatusOK {
		t.Logf("Register API Response Code: %d", w.Code)
		t.Logf("Register API Response Body: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK for successful registration")

	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err, "Failed to unmarshal response body for successful registration")

	// Assertions on response structure (adjust to your actual API response)
	// Example: {"code":0,"message":"OK","data":{"user_id":1,"username":"testapiuser_123","token":"ey..."}}
	assert.Contains(t, responseBody, "code", "Response should have a code")
	// Assuming your response.Success sets code to 0 or a specific success code
	// Let's assume response.Success uses 0 for success code as per some conventions.
	// This depends on your `douyin/pkg/utils/response` package.
	// If response.Success doesn't embed code directly in data, adjust assertions.
	// For now, let's assume the structure from previous examples where data is the main object.
	// assert.Equal(t, 0, int(responseBody["code"].(float64)), "Response code should be 0 for success")


	// Test duplicate registration
	wDup := httptest.NewRecorder()
	// Use the same body for request
	reqDup, _ := http.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewBuffer(body))
	reqDup.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(wDup, reqDup)

	if wDup.Code == http.StatusOK {
		t.Logf("Duplicate Register API Response Code: %d", wDup.Code)
		t.Logf("Duplicate Register API Response Body: %s", wDup.Body.String())
	}
	assert.NotEqual(t, http.StatusOK, wDup.Code, "Duplicate registration should not return HTTP 200 OK")
	// A common response for duplicate entry might be 400 (Bad Request) or 409 (Conflict)
	// or even 500 if not handled gracefully. Assert a specific non-OK code you expect.
	// For example, if your API returns 400 for "username already exists":
	// assert.Equal(t, http.StatusBadRequest, wDup.Code, "Expected HTTP 400 for duplicate registration")

	// Unmarshal error response for duplicate registration
	var errResponseBody map[string]interface{}
	err = json.Unmarshal(wDup.Body.Bytes(), &errResponseBody)
	require.NoError(t, err, "Failed to unmarshal error response body for duplicate registration")
	// Assert something about the error message if your API provides one
	// assert.Contains(t, errResponseBody, "message", "Error response should have a message")
	// message, _ := errResponseBody["message"].(string)
	// assert.Contains(t, strings.ToLower(message), "username already exists", "Error message should indicate duplicate username")
}
