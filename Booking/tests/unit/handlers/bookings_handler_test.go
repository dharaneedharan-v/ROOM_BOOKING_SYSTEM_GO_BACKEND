
package handlers_test

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"BookingSystem/Booking/internal/config"
	"BookingSystem/Booking/internal/handlers"
	"BookingSystem/Booking/internal/loggers"
	"BookingSystem/Booking/pkg/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Helper function to setup mock GORM DB
func setupTestDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := sqlserver.New(sqlserver.Config{
		Conn:       db,
		DriverName: "sqlmock",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

// Helper function to setup mock SQL DB
func SetupMockSqlDB() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

func TestSetupRoutes_Complete(t *testing.T) {
	// 1. Create temporary directory structure for OpenAPI Spec
	tmpDir, err := os.MkdirTemp("", "routes_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	docsDir := filepath.Join(tmpDir, "docs")
	require.NoError(t, os.MkdirAll(filepath.Join(docsDir, "swaggerui"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(docsDir, "openapi.yaml"), []byte("url: {BASE_URL}"), 0644))

	// 2. Initialize real router and mock DBs
	router := mux.NewRouter()
	logger := loggers.NewTestLogger()
	gormDB, mock, _ := setupTestDB()
	sqlDB, _, _ := SetupMockSqlDB()
	db := &database.Db{Gorm: gormDB, SqlDb: sqlDB}

	mockConfig := &config.Config{
		DatabaseURL: "mock_db_url",
		Port:        "8080",
		BaseUrl:     "http://localhost:8080",
		ProjectRoot: tmpDir,
	}

	// 3. Register your real routes mapping
	handlers.SetupRoutes(router, db, logger, mockConfig)

	// 4. Define real table-driven test cases covering ALL methods
	testCases := []struct {
		name           string
		method         string
		path           string
		payload        string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:   "Health Check Success",
			method: "GET",
			path:   "/health",
			setupMocks: func() {}, // No DB action needed
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Create Booking - Invalid JSON Payload",
			method: "POST",
			path:   "/api/v1/booking",
			payload: `{"customer_uuid": "some-uuid", "room_uuid": }`, // Bad JSON formatting
			setupMocks: func() {}, // Fails in handler validation, no DB hit
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Get Booking - Record Not Found",
			method: "GET",
			path:   "/api/v1/bookings/bad-uuid",
			setupMocks: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM")).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:    "Update Booking - Invalid JSON Payload",
			method:  "PUT",
			path:    "/api/v1/bookings/valid-uuid",
			payload: `{"booking_date": "invalid-json"`, // Incomplete json block
			setupMocks: func() {},
			expectedStatus: http.StatusUnprocessableEntity, 
		},
		{
			name:   "Delete Booking - Record Not Found Error",
			method: "DELETE",
			path:   "/api/v1/bookings/missing-uuid",
			setupMocks: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE")).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	// 5. Run the test engine loop
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Trigger specific DB expectation rules for this test case
			tc.setupMocks()

			var req *http.Request
			if tc.method == "POST" || tc.method == "PUT" {
				req, err = http.NewRequest(tc.method, tc.path, bytes.NewBuffer([]byte(tc.payload)))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tc.method, tc.path, nil)
			}
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// Assertions confirm the request went completely through the routing layer
			assert.Equal(t, tc.expectedStatus, rr.Code, "Failed on endpoint: %s %s", tc.method, tc.path)
		})
	}
}


func TestSetupRoutes_YAML_File_NOT_FOUND(t *testing.T) {
	router := mux.NewRouter()
	db := &database.Db{}
	logger := loggers.NewTestLogger()

	mockConfig := &config.Config{
		DatabaseURL: "mock_db_url",
		Port:        "8080",
		BaseUrl:     "http://localhost:8080",
		ProjectRoot: "/completely/invalid/path/to/trigger/not/found",
	}

	handlers.SetupRoutes(router, db, logger, mockConfig)

	req, err := http.NewRequest("GET", "/docs/openapi.yaml", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "Should return 404 when the openapi.yaml filesystem lookup breaks")
}

type simpleErrorWriter struct{}

func (s *simpleErrorWriter) Header() http.Header        { return make(http.Header) }
func (s *simpleErrorWriter) WriteHeader(statusCode int) {}
func (s *simpleErrorWriter) Write(b []byte) (int, error) {
	// This instantly forces writeErr != nil inside your routes.go
	return 0, os.ErrPermission 
}

// 2. Run the test case
func TestSetupRoutes_OpenAPI_WriteError(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")
	_ = os.MkdirAll(docsDir, 0755)
	
	// Create a dummy file
	yamlPath := filepath.Join(docsDir, "openapi.yaml")
	_ = os.WriteFile(yamlPath, []byte("url: {BASE_URL}"), 0644)

	router := mux.NewRouter()
	logger := loggers.NewTestLogger()
	cfg := &config.Config{BaseUrl: "http://localhost", ProjectRoot: tmpDir}

	handlers.SetupRoutes(router, &database.Db{}, logger, cfg)

	// Send a request into our fake error writer
	req, _ := http.NewRequest("GET", "/docs/openapi.yaml", nil)
	rw := &simpleErrorWriter{}
	
	router.ServeHTTP(rw, req)
}




















// package handlers_test

// import (
// 	"bytes"
// 	"database/sql"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"path/filepath"
// 	"regexp"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/gorilla/mux"
// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/driver/sqlserver"
// 	"gorm.io/gorm"

// 	"BookingSystem/Booking/internal/config"
// 	handlers "BookingSystem/Booking/internal/handlers"
// 	loggers "BookingSystem/Booking/internal/loggers"
// 	"BookingSystem/Booking/pkg/database"
// )

// func setupTestDB() (*gorm.DB, sqlmock.Sqlmock, error) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	dialector := sqlserver.New(sqlserver.Config{
// 		Conn:       db,
// 		DriverName: "sqlmock",
// 	})

// 	gormDB, err := gorm.Open(dialector, &gorm.Config{})
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return gormDB, mock, nil
// }

// // SetupMockSqlDB creates and returns a mocked *sql.DB and sqlmock.Sqlmock for testing raw SQL operations.
// func SetupMockSqlDB() (*sql.DB, sqlmock.Sqlmock, error) {
// 	// Create a new sqlmock database connection and mock controller
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return db, mock, nil
// }

// func TestSetupRoutes(t *testing.T) {
// 	// Setup mock router and config
// 	router := mux.NewRouter()
// 	logger := loggers.NewTestLogger()
// 	gormDB, mock, _ := setupTestDB()
// 	sqlDB, _, _ := SetupMockSqlDB()
// 	db := &database.Db{Gorm: gormDB, SqlDb: sqlDB}
// 	mockConfig := &config.Config{
// 		DatabaseURL: "mock_db_url",
// 		Port:        "8080",
// 		BaseUrl:     "http://localhost:8080",
// 		ProjectRoot: filepath.Join("..", "..", ".."),
// 	}

// 	// Setup routes with our mocked dependencies
// 	handlers.SetupRoutes(router, db, logger, mockConfig)

// 	// Test cases for different endpoints
// 	t.Run("Health Check Endpoint", func(t *testing.T) {
// 		req, err := http.NewRequest(http.MethodGet, "/health", nil)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "healthy")
// 		assert.Contains(t, rr.Body.String(), "training-service")
// 	})

// 	t.Run("Create Booking - Invalid JSON", func(t *testing.T) {
// 		payload := `{"customer_uuid": "6ba7b810-9dad-11d1-80b4-00c04fd430c8", "room_uuid": }` // Invalid JSON

// 		req, err := http.NewRequest(http.MethodPost, "/api/v1/booking", bytes.NewBuffer([]byte(payload)))
// 		assert.NoError(t, err)
// 		req.Header.Set("Content-Type", "application/json")

// 		rr := httptest.NewRecorder()
// 		router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusBadRequest, rr.Code)
// 	})
// 	t.Run("Get Booking - Booking Not Found", func(t *testing.T) {
// 		// Expect database query to run, then simulate record not found error
// 		mock.ExpectQuery(regexp.QuoteMeta("SELECT")).WillReturnError(gorm.ErrRecordNotFound)

// 		req, err := http.NewRequest(http.MethodGet, "/api/v1/bookings/non-existent-uuid", nil)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		router.ServeHTTP(rr, req)

// 		// Evaluates that the status code bubbles up correctly through the routing layers
// 		assert.Equal(t, http.StatusNotFound, rr.Code)
// 	})

// 	t.Run("Verify Booking PUT Registration", func(t *testing.T) {
// 		req, _ := http.NewRequest(http.MethodPut, "/api/v1/bookings/abc-123", bytes.NewBuffer([]byte(`{}`)))
// 		rr := httptest.NewRecorder()
// 		router.ServeHTTP(rr, req)

// 		assert.NotEqual(t, http.StatusNotFound, rr.Code)
// 	})

// 	t.Run("Delete Booking - Record Not Found", func(t *testing.T) {
// 		// Expect the query to fail because the target row does not exist
// 		mock.ExpectBegin()
// 		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WillReturnError(gorm.ErrRecordNotFound)
// 		// mock.Rollback()

// 		req, err := http.NewRequest(http.MethodDelete, "/api/v1/bookings/non-existent-uuid", nil)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusNotFound, rr.Code)
// 	})

// 		t.Run("Serve OpenAPI YAML - Success", func(t *testing.T) {
// 		// Setup temporary openapi.yaml
// 		tmpDir := t.TempDir()
// 		openapiPath := filepath.Join(tmpDir, "docs", "openapi.yaml")

// 		err := os.MkdirAll(filepath.Dir(openapiPath), os.ModePerm)
// 		assert.NoError(t, err)

// 		originalContent := "openapi: 3.0.0\nservers:\n  - url: {BASE_URL}"
// 		err = os.WriteFile(openapiPath, []byte(originalContent), 0644)
// 		assert.NoError(t, err)

// 		// Setup config
// 		cfg := &config.Config{
// 			BaseUrl:     "http://localhost:8080",
// 			ProjectRoot: tmpDir,
// 		}

// 		// Setup router
// 		router := mux.NewRouter()
// 		handlers.SetupRoutes(router, nil, logger, cfg)

// 		req, err := http.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		assert.Equal(t, "application/x-yaml", rr.Header().Get("Content-Type"))
// 		assert.Contains(t, rr.Body.String(), cfg.BaseUrl)
// 	})

// 	t.Run("Serve OpenAPI YAML - File Not Found", func(t *testing.T) {
// 		// Setup config with non-existent file path
// 		cfg := &config.Config{
// 			BaseUrl:     "http://localhost:8080",
// 			ProjectRoot: t.TempDir(), // no openapi.yaml created
// 		}

// 		// Setup router
// 		router := mux.NewRouter()
// 		logger := loggers.NewTestLogger()

// 		handlers.SetupRoutes(router, nil, logger, cfg)

// 		req, err := http.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusNotFound, rr.Code)
// 	})
// }

// func TestSetupRoutes_YAML_File_NOT_FOUND(t *testing.T) {
// 	// Create a new router
// 	router := mux.NewRouter()

// 	// Create a mock database and logger
// 	db := &database.Db{} // Mock the database object if needed
// 	logger := loggers.NewTestLogger()
// 	// Mock configuration
// 	mockConfig := &config.Config{
// 		DatabaseURL: "mock_db_url",
// 		Port:        "8080",
// 		BaseUrl:     "http://localhost:8080",
// 		ProjectRoot: "invalid_path",
// 	}

// 	// Call the SetupRoutes function with the proper mock services
// 	handlers.SetupRoutes(router, db, logger, mockConfig)

// 	// Define test cases
// 	testCases := []struct {
// 		name           string
// 		method         string
// 		path           string
// 		expectedStatus int
// 		payload        string // Adding payload for POST requests
// 	}{
// 		{
// 			name:           "OpenAPI Spec Read Failure",
// 			method:         "GET",
// 			path:           "/docs/openapi.yaml",
// 			expectedStatus: http.StatusNotFound,
// 		},
// 		{
// 			name:           "Health Check Success",
// 			method:         "GET",
// 			path:           "/health",
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:           "Create User Invalid JSON",
// 			method:         "POST",
// 			path:           "/api/v1/users",
// 			payload:        `{"name": "Test", "age": }`, // Invalid JSON
// 			expectedStatus: http.StatusBadRequest,
// 		},
// 	}

// 	// Run test cases
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var req *http.Request
// 			var err error

// 			// Handle POST requests with a payload
// 			if tc.method == "POST" {
// 				req, err = http.NewRequest(tc.method, tc.path, bytes.NewBuffer([]byte(tc.payload)))
// 				req.Header.Set("Content-Type", "application/json") // Setting the content type to JSON
// 			} else {
// 				req, err = http.NewRequest(tc.method, tc.path, nil)
// 			}

// 			assert.NoError(t, err)

// 			rr := httptest.NewRecorder()
// 			router.ServeHTTP(rr, req)

// 			assert.Equal(t, tc.expectedStatus, rr.Code, "Unexpected status code for %s", tc.path)
// 		})
// 	}
// }
 

