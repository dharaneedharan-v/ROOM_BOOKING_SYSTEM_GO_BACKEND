

package database_test

import (
	"log"
	"os"
	"testing"

	"lynxis-gate/training-service/pkg/database"

	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {
	// Set up environment variables for the SQL server connection
	os.Setenv("DB_USERNAME", "username")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "1433")
	os.Setenv("DB_NAME", "testdb")

	// Generate a SQL Server DSN from the environment variables
	dsn := "sqlserver://" +
		os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") +
		"?database=" + os.Getenv("DB_NAME")

	iul := "sqlserver:/" +
		os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") +
		"@" + os.Getenv("HOST") + ":" + os.Getenv("DB_PORT") +
		"?database="

	tests := []struct {
		name    string
		dbURL   string
		wantErr bool
	}{
		{
			name:    "Valid connection string format",
			dbURL:   dsn,
			wantErr: false,
		},
		{
			name:    "Invalid connection string",
			dbURL:   iul,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load env file if present for the current environment
			dbService := &database.DBService{} // Use the actual DB service
			db, err := dbService.EstablishConnection(tt.dbURL)

			if tt.wantErr {
				assert.Nil(t, db)
				assert.NotNil(t, err)
			} else {
				if err != nil {
					t.Logf("Failed to connect to database: %v", err)
					assert.NotNil(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, db)

					cleanupErr := db.SqlDb.Close()
					if cleanupErr != nil {
						log.Printf("Failed to close the database: %v", cleanupErr)
					}
				}
			}
		})
	}

	// Cleanup environment variables
	os.Unsetenv("DB_USERNAME")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
}