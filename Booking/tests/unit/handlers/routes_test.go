
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

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"lynxis-gate/training-service/internal/config"
	handlers "lynxis-gate/training-service/internal/handlers"
	loggers "lynxis-gate/training-service/internal/loggers"
	"lynxis-gate/training-service/pkg/database"
)

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

// SetupMockSqlDB creates and returns a mocked *sql.DB and sqlmock.Sqlmock for testing raw SQL operations.
func SetupMockSqlDB() (*sql.DB, sqlmock.Sqlmock, error) {
	// Create a new sqlmock database connection and mock controller
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	return db, mock, nil
}

func TestSetupRoutes(t *testing.T) {
	// Setup mock router and config
	router := mux.NewRouter()
	logger := loggers.NewLogger("test-service")
	gormDB, mock, _ := setupTestDB()
	sqlDB, _, _ := SetupMockSqlDB()
	db := &database.Db{Gorm: gormDB, SqlDb: sqlDB}
	mockConfig := &config.Config{
		DatabaseURL: "mock_db_url",
		Port:        "8080",
		BaseUrl:     "http://localhost:8080",
		ProjectRoot: filepath.Join("..", "..", ".."),
	}

	// Setup routes with our mocked dependencies
	handlers.SetupRoutes(router, db, logger, mockConfig)

	// Test cases for different endpoints
	t.Run("Health Check Endpoint", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/health", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "healthy")
		assert.Contains(t, rr.Body.String(), "training-service")
	})

	t.Run("Create User - Invalid JSON", func(t *testing.T) {
		payload := `{"name": "Test User", "age": }` // Invalid JSON

		req, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer([]byte(payload)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Get User - User Not Found", func(t *testing.T) {
		// Expect query for user - just expect any SELECT query and return not found
		mock.ExpectQuery(regexp.QuoteMeta("SELECT")).WillReturnError(gorm.ErrRecordNotFound)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/users/non-existent-uuid", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Should return NotFound
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Serve OpenAPI YAML - Success", func(t *testing.T) {
		// Setup temporary openapi.yaml
		tmpDir := t.TempDir()
		openapiPath := filepath.Join(tmpDir, "docs", "openapi.yaml")

		err := os.MkdirAll(filepath.Dir(openapiPath), os.ModePerm)
		assert.NoError(t, err)

		originalContent := "openapi: 3.0.0\nservers:\n  - url: {BASE_URL}"
		err = os.WriteFile(openapiPath, []byte(originalContent), 0644)
		assert.NoError(t, err)

		// Setup config
		cfg := &config.Config{
			BaseUrl:     "http://localhost:8080",
			ProjectRoot: tmpDir,
		}

		// Setup router
		router := mux.NewRouter()
		handlers.SetupRoutes(router, nil, logger, cfg)

		req, err := http.NewRequest(http.MethodGet, "/training_service/docs/openapi.yaml", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/x-yaml", rr.Header().Get("Content-Type"))
		assert.Contains(t, rr.Body.String(), cfg.BaseUrl)
	})

	t.Run("Serve OpenAPI YAML - File Not Found", func(t *testing.T) {
		// Setup config with non-existent file path
		cfg := &config.Config{
			BaseUrl:     "http://localhost:8080",
			ProjectRoot: t.TempDir(), // no openapi.yaml created
		}

		// Setup router
		router := mux.NewRouter()
		logger := loggers.NewLogger("test-service")

		handlers.SetupRoutes(router, nil, logger, cfg)

		req, err := http.NewRequest(http.MethodGet, "/training_service/docs/openapi.yaml", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestSetupRoutes_YAML_File_NOT_FOUND(t *testing.T) {
	// Create a new router
	router := mux.NewRouter()

	// Create a mock database and logger
	db := &database.Db{} // Mock the database object if needed
	logger := loggers.NewLogger("test-service")

	// Mock configuration
	mockConfig := &config.Config{
		DatabaseURL: "mock_db_url",
		Port:        "8080",
		BaseUrl:     "http://localhost:8080",
		ProjectRoot: "invalid_path",
	}

	// Call the SetupRoutes function with the proper mock services
	handlers.SetupRoutes(router, db, logger, mockConfig)

	// Define test cases
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		payload        string // Adding payload for POST requests
	}{
		{
			name:           "OpenAPI Spec Read Failure",
			method:         "GET",
			path:           "/training_service/docs/openapi.yaml",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Health Check Success",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Create User Invalid JSON",
			method:         "POST",
			path:           "/api/v1/users",
			payload:        `{"name": "Test", "age": }`, // Invalid JSON
			expectedStatus: http.StatusBadRequest,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			var err error

			// Handle POST requests with a payload
			if tc.method == "POST" {
				req, err = http.NewRequest(tc.method, tc.path, bytes.NewBuffer([]byte(tc.payload)))
				req.Header.Set("Content-Type", "application/json") // Setting the content type to JSON
			} else {
				req, err = http.NewRequest(tc.method, tc.path, nil)
			}

			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "Unexpected status code for %s", tc.path)
		})
	}
}
