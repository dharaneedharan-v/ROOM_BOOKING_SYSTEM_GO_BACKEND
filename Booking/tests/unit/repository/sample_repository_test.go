


package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"lynxis-gate/training-service/internal/loggers"
	"lynxis-gate/training-service/internal/models"
	"lynxis-gate/training-service/internal/repository"
	"lynxis-gate/training-service/pkg/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository repository.UserRepositoryInterface
	logger     *loggers.Logger
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	suite.NoError(err)

	// Configure the dialector for GORM
	dialector := sqlserver.New(sqlserver.Config{
		Conn:       db,
		DriverName: "sqlmock",
	})

	// Open a GORM connection
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	suite.NoError(err)

	// Create a new Database struct
	dbMock := &database.Db{
		Gorm:  gormDB,
		SqlDb: db,
	}

	// Create a logger
	logger := loggers.NewLogger("test-service")

	// Initialize the repository with the mock database
	suite.DB = gormDB
	suite.mock = mock
	suite.repository = repository.NewUserRepository(dbMock, logger)
	suite.logger = logger
}

func (suite *UserRepositoryTestSuite) TestCreateUser() {
	user := &models.User{
		UserUUID:  "test-uuid",
		Name:      "John Doe",
		Age:       30,
		Address:   "123 Main St",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Success case
	suite.T().Run("Success", func(t *testing.T) {
		// Setup expectations - match any query that contains INSERT INTO
		suite.mock.ExpectBegin()
		suite.mock.ExpectQuery("INSERT INTO").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		suite.mock.ExpectCommit()

		// Call repository method
		err := suite.repository.CreateUser(context.Background(), user)

		// Assertions
		assert.Nil(t, err)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})

	// Error case
	suite.T().Run("Error", func(t *testing.T) {
		// Setup expectations
		suite.mock.ExpectBegin()
		suite.mock.ExpectQuery("INSERT INTO").
			WillReturnError(sqlmock.ErrCancelled)
		suite.mock.ExpectRollback()

		// Call repository method
		err := suite.repository.CreateUser(context.Background(), user)

		// Assertions
		assert.NotNil(t, err)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})
}

func (suite *UserRepositoryTestSuite) TestGetUserByUUID() {
	uuid := "test-uuid"
	columns := []string{"id", "user_uuid", "name", "age", "address", "created_at", "updated_at", "deleted_at"}

	// Success case
	suite.T().Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).
			AddRow(1, uuid, "John Doe", 30, "123 Main St", time.Now(), time.Now(), nil)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WillReturnRows(rows)

		// Call repository method
		user, err := suite.repository.GetUserByUUID(context.Background(), uuid)

		// Assertions
		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uuid, user.UserUUID)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})

	// Not found case
	suite.T().Run("Not Found", func(t *testing.T) {
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WillReturnError(gorm.ErrRecordNotFound)

		// Call repository method
		user, err := suite.repository.GetUserByUUID(context.Background(), "non-existent-uuid")

		// Assertions
		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})

	// Database error case
	suite.T().Run("Database Error", func(t *testing.T) {
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WillReturnError(sqlmock.ErrCancelled)

		// Call repository method
		user, err := suite.repository.GetUserByUUID(context.Background(), uuid)

		// Assertions
		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})
}

func (suite *UserRepositoryTestSuite) TestUpdateUser() {
	uuid := "test-uuid"
	updateData := map[string]interface{}{
		"name":    "Updated Name",
		"age":     35,
		"address": "456 New Address",
	}

	// Success case
	suite.T().Run("Success", func(t *testing.T) {
		// Setup expectations
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE")).
			WillReturnResult(sqlmock.NewResult(0, 1))
		suite.mock.ExpectCommit()

		// Call repository method
		err := suite.repository.UpdateUser(context.Background(), uuid, updateData)

		// Assertions
		assert.Nil(t, err)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})

	// User not found case
	suite.T().Run("Not Found", func(t *testing.T) {
		// Setup expectations
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE")).
			WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
		suite.mock.ExpectCommit()

		// Call repository method
		err := suite.repository.UpdateUser(context.Background(), "non-existent-uuid", updateData)

		// Assertions
		assert.NotNil(t, err)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})

	// Database error case
	suite.T().Run("Database Error", func(t *testing.T) {
		// Setup expectations
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE")).
			WillReturnError(sqlmock.ErrCancelled)
		suite.mock.ExpectRollback()

		// Call repository method
		err := suite.repository.UpdateUser(context.Background(), uuid, updateData)

		// Assertions
		assert.NotNil(t, err)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})
}

func (suite *UserRepositoryTestSuite) TestDeleteUser() {
	uuid := "test-uuid"

	// Success case
	suite.T().Run("Success", func(t *testing.T) {
		// Setup expectations - GORM uses soft deletes with UPDATE statements, not actual DELETE
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec("UPDATE").
			WillReturnResult(sqlmock.NewResult(0, 1))
		suite.mock.ExpectCommit()

		// Call repository method
		err := suite.repository.DeleteUser(context.Background(), uuid)

		// Assertions
		assert.Nil(t, err)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})

	// User not found case
	suite.T().Run("Not Found", func(t *testing.T) {
		// Setup expectations
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec("UPDATE").
			WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
		suite.mock.ExpectCommit()

		// Call repository method
		err := suite.repository.DeleteUser(context.Background(), "non-existent-uuid")

		// Assertions
		assert.NotNil(t, err)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})

	// Database error case
	suite.T().Run("Database Error", func(t *testing.T) {
		// Setup expectations
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec("UPDATE").
			WillReturnError(sqlmock.ErrCancelled)
		suite.mock.ExpectRollback()

		// Call repository method
		err := suite.repository.DeleteUser(context.Background(), uuid)

		// Assertions
		assert.NotNil(t, err)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})
}

func (suite *UserRepositoryTestSuite) TestGetAllUsers() {
	columns := []string{"id", "user_uuid", "name", "age", "address", "created_at", "updated_at", "deleted_at"}
	now := time.Now()

	// Success case with users
	suite.T().Run("Success - With Users", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*)")).
			WillReturnRows(countRows)

		rows := sqlmock.NewRows(columns).
			AddRow(1, "uuid-1", "User 1", 30, "Address 1", now, now, nil).
			AddRow(2, "uuid-2", "User 2", 35, "Address 2", now, now, nil)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WillReturnRows(rows)

		// Call repository method
		users, count, err := suite.repository.GetAllUsers(context.Background())

		// Assertions
		assert.Nil(t, err)
		assert.Equal(t, int64(2), count)
		assert.Equal(t, 2, len(users))
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})

	// Success case with no users
	suite.T().Run("Success - Empty List", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*)")).
			WillReturnRows(countRows)

		rows := sqlmock.NewRows(columns)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WillReturnRows(rows)

		// Call repository method
		users, count, err := suite.repository.GetAllUsers(context.Background())

		// Assertions
		assert.Nil(t, err)
		assert.Equal(t, int64(0), count)
		assert.Equal(t, 0, len(users))
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})

	// Count query error
	suite.T().Run("Count Error", func(t *testing.T) {
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*)")).
			WillReturnError(sqlmock.ErrCancelled)

		// Call repository method
		users, count, err := suite.repository.GetAllUsers(context.Background())

		// Assertions
		assert.NotNil(t, err)
		assert.Equal(t, int64(0), count)
		assert.Nil(t, users)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})

	// Select query error
	suite.T().Run("Select Error", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*)")).
			WillReturnRows(countRows)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WillReturnError(sqlmock.ErrCancelled)

		// Call repository method
		users, count, err := suite.repository.GetAllUsers(context.Background())

		// Assertions
		assert.NotNil(t, err)
		assert.Equal(t, int64(0), count)
		assert.Nil(t, users)
		assert.NoError(t, suite.mock.ExpectationsWereMet())
	})
}

// Run the test suite
func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
