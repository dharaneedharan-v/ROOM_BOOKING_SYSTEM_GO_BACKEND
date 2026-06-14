


package server_test

import (
	"os"
	"testing"
	"time"

	"BookingSystem/Booking/internal/config"
	log "BookingSystem/Booking/internal/loggers"
	"BookingSystem/Booking/pkg/database"
	"BookingSystem/Booking/pkg/server"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDBService struct {
	mock.Mock
}

func (m *MockDBService) EstablishConnection(dbURL string) (*database.Db, error) {
	args := m.Called(dbURL)
	if db, ok := args.Get(0).(*database.Db); ok {
		return db, args.Error(1)
	}
	return nil, args.Error(1)
}

// Helper function to create a temporary encrypted configuration file with proper encryption
func createTempConfigFile(t *testing.T, content string) string {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "config-*.bin")
	assert.NoError(t, err)

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}

func TestRunServer(t *testing.T) {
	logger := log.NewTestLogger()
	router := mux.NewRouter()
	cfg := &config.Config{
		Port: "0", // Use port 0 to let the OS assign a free port
	}

	app := &server.Application{
		Logger: logger,
		Router: router,
		Config: cfg,
	}

	// Start server in a goroutine
	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- server.RunServer(app)
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)
}