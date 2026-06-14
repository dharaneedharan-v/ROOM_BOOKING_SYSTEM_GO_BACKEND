package config_test

import (
	"os"
	"path/filepath"
	"testing"

	config "BookingSystem/Booking/internal/config"
	loggers "BookingSystem/Booking/internal/loggers"


	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func getDummyLogger() *loggers.Logger {
	return loggers.NewTestLogger()
}

func TestLoadConfig(t *testing.T) {
	// Test case: Success case
	t.Run("Successful Load", func(t *testing.T) {
		// Clean up env variables after this subtest runs
		defer os.Unsetenv("DATABASE_URL")
		defer os.Unsetenv("PORT")
		defer os.Unsetenv("PROJECT_ROOT")
		defer os.Unsetenv("BASE_URL")

		os.Setenv("DATABASE_URL", "db_test_url")
		os.Setenv("PORT", "9090")
		os.Setenv("PROJECT_ROOT", "/test/root")
		os.Setenv("BASE_URL", "http://localhost:9090")

		cfg, err := config.LoadConfig("testService", "development")

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "9090", cfg.Port)
		assert.Equal(t, "db_test_url", cfg.DatabaseURL)
		assert.Equal(t, "/test/root", cfg.ProjectRoot)
		assert.Equal(t, "http://localhost:9090", cfg.BaseUrl)
	})

	// Test case: Empty DATABASE_URL block error
	t.Run("Missing DATABASE_URL Error", func(t *testing.T) {
		defer os.Unsetenv("DATABASE_URL")
		
		// Force DATABASE_URL to be completely empty
		os.Setenv("DATABASE_URL", "")

		cfg, err := config.LoadConfig("testService", "nonexistent")

		assert.Nil(t, cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DATABASE_URL environment variable is required")
	})

	// Test case: Godotenv Load else block error (unreadable file)
	t.Run("Godotenv Load Else Block Error", func(t *testing.T) {
		defer os.Unsetenv("DATABASE_URL")
		os.Setenv("DATABASE_URL", "db_test_url")

		// Create a temporary directory structure
		tmpDir, err := os.MkdirTemp("", "config_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		
		envsDir := filepath.Join(tmpDir, "envs")
		err = os.MkdirAll(envsDir, 0755)
		require.NoError(t, err)
		
		brokenFileDir := filepath.Join(envsDir, ".env.broken")
		err = os.Mkdir(brokenFileDir, 0755)
		require.NoError(t, err)

		deepDir := filepath.Join(tmpDir, "internal", "config")
		err = os.MkdirAll(deepDir, 0755)
		require.NoError(t, err)

		oldWd, err := os.Getwd()
		require.NoError(t, err)
		
		err = os.Chdir(deepDir)
		require.NoError(t, err)
		defer func() {
			_ = os.Chdir(oldWd)
		}()

		cfg, err := config.LoadConfig("testService", "broken")

		assert.Nil(t, cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load")
	})
}


//  go test ./... -coverpkg=BookingSystem/Booking/internal/config -coverprofile="coverage.out"


// go tool cover -html="coverage.out"                
 

