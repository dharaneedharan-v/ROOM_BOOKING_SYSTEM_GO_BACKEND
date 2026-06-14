
package config_test

import (
	"os"
	"testing"

	config "lynxis-gate/training-service/internal/config"
	loggers "lynxis-gate/training-service/internal/loggers"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Setup code, if needed
	exitCode := m.Run()
	os.Exit(exitCode)
}

// A dummy logger to inject for testing purposes
func getDummyLogger() *loggers.Logger {
	return loggers.NewLogger("test-service")
}

func TestLoadConfig(t *testing.T) {
	// Test case: Success case
	t.Run("Successful Load", func(t *testing.T) {
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
}