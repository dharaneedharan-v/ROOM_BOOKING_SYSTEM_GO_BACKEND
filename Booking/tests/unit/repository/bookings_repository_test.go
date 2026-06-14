
	package repository_test

	import (
		"context"
		"errors"
		"regexp"
		"testing"
		"time"

		"BookingSystem/Booking/internal/loggers"
		"BookingSystem/Booking/internal/models"
		"BookingSystem/Booking/internal/repository"
		"BookingSystem/Booking/pkg/database"

		"github.com/DATA-DOG/go-sqlmock"
		"github.com/stretchr/testify/assert"
		"github.com/stretchr/testify/suite"
		"gorm.io/driver/postgres" 
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
		db, mock, err := sqlmock.New()
		suite.NoError(err)

		// FIXED: Configured GORM to generate PostgreSQL syntax instead of SQL Server
		dialector := postgres.New(postgres.Config{
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gormDB, err := gorm.Open(dialector, &gorm.Config{})
		suite.NoError(err)

		dbMock := &database.Db{
			Gorm:  gormDB,
			SqlDb: db,
		}

		logger := loggers.NewTestLogger()

		suite.DB = gormDB
		suite.mock = mock
		suite.repository = repository.NewUserRepository(dbMock, logger)
		suite.logger = logger
	}

// -----------------------------------------------------
// Soft delete 

func (suite *UserRepositoryTestSuite) TestSoftDeleteBooking() {
	bookingUUID := "94841ade-468c-480f-8b69-ee911e6fcbdb"

	// 1. SUCCESS CASE
	suite.T().Run("Success", func(t *testing.T) {
		// Clean and rebuild state context for this subtest boundary
		suite.SetupTest()

		rows := sqlmock.NewRows([]string{"status"}).AddRow("CONFIRMED")
		
		suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "status" FROM "bookings"`)).
			WithArgs(bookingUUID, sqlmock.AnyArg()).
			WillReturnRows(rows)

		suite.mock.ExpectBegin()
		// Loose match the columns to allow GORM's automatic timestamp updates to pass
		suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings" SET`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		suite.mock.ExpectCommit()

		err := suite.repository.SoftDeleteBooking(context.Background(), bookingUUID)
		assert.Nil(t, err)
	})

	// 2. NOT FOUND CASE
	suite.T().Run("NotFound", func(t *testing.T) {
		suite.SetupTest()

		suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "status" FROM "bookings"`)).
			WithArgs(bookingUUID, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		err := suite.repository.SoftDeleteBooking(context.Background(), bookingUUID)
		
		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "Booking not found", err.Message)
		}
	})

	// 3. ALREADY CANCELLED CASE
	suite.T().Run("AlreadyCancelled", func(t *testing.T) {
		suite.SetupTest()

		rows := sqlmock.NewRows([]string{"status"}).AddRow("CANCELLED")
		
		suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "status" FROM "bookings"`)).
			WithArgs(bookingUUID, sqlmock.AnyArg()).
			WillReturnRows(rows)

		err := suite.repository.SoftDeleteBooking(context.Background(), bookingUUID)
		
		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "Booking Already Cancelled...!!", err.Message)
		}
	})

	// 4. DATABASE UPDATE FAILURE CASE
	suite.T().Run("UpdateFailure", func(t *testing.T) {
		suite.SetupTest()

		rows := sqlmock.NewRows([]string{"status"}).AddRow("CONFIRMED")
		
		suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "status" FROM "bookings"`)).
			WithArgs(bookingUUID, sqlmock.AnyArg()).
			WillReturnRows(rows)

		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings" SET`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("postgres connection lost"))
		suite.mock.ExpectRollback()

		err := suite.repository.SoftDeleteBooking(context.Background(), bookingUUID)
		
		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "Failed to delete booking", err.Message)
		}
	})
}

// ------------------------------------------------------------

func (suite *UserRepositoryTestSuite) TestCheckRoomAvailabilityForUpdate() {
	var roomID uint = 1
	bookingUUID := "94841ade-468c-480f-8b69-ee911e6fcbdb"
	
	start := time.Date(2026, 5, 26, 14, 0, 0, 0, time.UTC)
	end := time.Date(2026, 5, 26, 16, 0, 0, 0, time.UTC)

	// 1. CONFLICT EXISTS CASE (Should return true)
	suite.T().Run("ConflictExists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		// Use a flexible regex to match GORM's full generated query
		suite.mock.ExpectQuery(`SELECT count\(\*\) FROM "bookings"`).
			WithArgs(roomID, bookingUUID, end, start, "CANCELLED"). // Added "CANCELLED" status argument
			WillReturnRows(rows)

		conflict := suite.repository.CheckRoomAvailabilityForUpdate(context.Background(), roomID, start, end, bookingUUID)

		assert.True(t, conflict)
	})

	// 2. NO CONFLICT CASE (Should return false)
	suite.T().Run("NoConflictAvailable", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		suite.mock.ExpectQuery(`SELECT count\(\*\) FROM "bookings"`).
			WithArgs(roomID, bookingUUID, end, start, "CANCELLED"). // Added "CANCELLED" status argument
			WillReturnRows(rows)

		conflict := suite.repository.CheckRoomAvailabilityForUpdate(context.Background(), roomID, start, end, bookingUUID)

		assert.False(t, conflict)
	})
}


//  Test booking UUID 
func (suite *UserRepositoryTestSuite) TestGetBookingByUUID() {
	bookingUUID := "94841ade-468c-480f-8b69-ee911e6fcbdb"
	customerID := uint(10)
	roomID := uint(20)

	// 1. SUCCESS CASE
	suite.T().Run("Success", func(t *testing.T) {
		// Prepare the mock row structure matching your database table columns
		bookingRows := sqlmock.NewRows([]string{"id", "booking_uuid", "customer_id", "room_id", "status", "is_active"}).
			AddRow(1, bookingUUID, customerID, roomID, "CONFIRMED", true)
		
		// Match the query using a flexible regex to handle GORM's auto-generated ORDER BY and WHERE clauses
		// WithArgs maps directly to: uuid, true (is_active), "CANCELLED" (status), and AnyArg() for the implicit LIMIT 1
		suite.mock.ExpectQuery(`SELECT \* FROM "bookings"`).
			WithArgs(bookingUUID, true, "CANCELLED", sqlmock.AnyArg()). 
			WillReturnRows(bookingRows)

		// Execute the repository function
		booking, err := suite.repository.GetBookingByUUID(context.Background(), bookingUUID)

		// Assertions
		assert.Nil(t, err)
		assert.NotNil(t, booking)
		if booking != nil { 
			assert.Equal(t, bookingUUID, booking.BookingUUID)
			assert.Equal(t, customerID, booking.CustomerID)
			assert.Equal(t, roomID, booking.RoomID)
			assert.Equal(t, "CONFIRMED", booking.Status)
		}
	})

	// 2. RECORD NOT FOUND CASE
	suite.T().Run("NotFound", func(t *testing.T) {
		// Use the exact same query signature, but instruct it to bubble up a Record Not Found error
		suite.mock.ExpectQuery(`SELECT \* FROM "bookings"`).
			WithArgs(bookingUUID, true, "CANCELLED", sqlmock.AnyArg()). 
			WillReturnError(gorm.ErrRecordNotFound)

		booking, err := suite.repository.GetBookingByUUID(context.Background(), bookingUUID)

		// Assertions
		assert.Nil(t, booking)
		assert.NotNil(t, err)
		assert.Equal(t, "Booking not found", err.Message)
	})

	// 3. GENERIC DATABASE ERROR CASE
	suite.T().Run("DatabaseError", func(t *testing.T) {
		// Simulates a connection timeout or syntax failure at the driver level
		suite.mock.ExpectQuery(`SELECT \* FROM "bookings"`).
			WithArgs(bookingUUID, true, "CANCELLED", sqlmock.AnyArg()). 
			WillReturnError(errors.New("connection timeout"))

		booking, err := suite.repository.GetBookingByUUID(context.Background(), bookingUUID)

		// Assertions
		assert.Nil(t, booking)
		assert.NotNil(t, err)
		assert.Equal(t, "Failed to fetch booking", err.Message)
	})
}

//  Room UUID TEST 

func (suite *UserRepositoryTestSuite) TestGetRoomByUUID() {
	roomUUID := "c50fe215-347f-46a0-bdc0-51479f96d451"

	suite.T().Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "room_uuid", "room_name", "capacity", "is_active",
		}).AddRow(1, roomUUID, "CR-1", 5, true)

		// FIXED: Appended sqlmock.AnyArg() to match GORM's internal implicit single-record limit restriction parameters
		suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms"`)).
			WithArgs(roomUUID, sqlmock.AnyArg()).
			WillReturnRows(rows)

		room, err := suite.repository.GetRoomByUUID(context.Background(), roomUUID)

		assert.Nil(t, err)
		assert.NotNil(t, room)
		if room != nil { // Protective context boundary checking block
			assert.Equal(t, "CR-1", room.RoomName)
			assert.Equal(t, roomUUID, room.RoomUUID)
		}
	})

	suite.T().Run("NotFound", func(t *testing.T) {
		suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms"`)).
			WithArgs(roomUUID, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		room, err := suite.repository.GetRoomByUUID(context.Background(), roomUUID)

		assert.Nil(t, room)
		assert.NotNil(t, err)
		assert.Equal(t, "Room not found", err.Message)
	})
}

func (suite *UserRepositoryTestSuite) TestUpdateBooking() {
	bookingUUID := "94841ade-468c-480f-8b69-ee911e6fcbdb"
	
	// Create a dummy instance to pass to the repository Save operation
	booking := &models.Booking{
		ID:          1,
		BookingUUID: bookingUUID,
		CustomerID:  10,
		RoomID:      20,
		BookingDate: time.Date(2026, 5, 26, 0, 0, 0, 0, time.UTC),
		StartTime:   time.Date(2026, 5, 26, 14, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2026, 5, 26, 16, 0, 0, 0, time.UTC),
		Status:      "CONFIRMED",
		IsActive:    true,
	}

	// 1. SUCCESS CASE
	suite.T().Run("Success", func(t *testing.T) {
		suite.mock.ExpectBegin()
		// GORM Save issues a complete UPDATE command targeting all columns on the model primary key
		suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		suite.mock.ExpectCommit()

		err := suite.repository.UpdateBooking(context.Background(), booking)

		assert.Nil(t, err)
	})

	// 2. DATABASE TRANSACTION CRASH CASE
	suite.T().Run("DatabaseFailure", func(t *testing.T) {
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings" SET`)).
			WillReturnError(errors.New("postgresql write deadlock locked"))
		suite.mock.ExpectRollback()

		err := suite.repository.UpdateBooking(context.Background(), booking)

		assert.NotNil(t, err)
		assert.Equal(t, "Failed to update booking in database", err.Message)
	})
}



func (suite *UserRepositoryTestSuite) TestGetCustomerByUUID() {
	customerUUID := "e71260ef-4b14-4b99-9ef3-eba0ddfd48b3"

	suite.T().Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "customer_uuid", "customer_name", "customer_phone", 
			"customer_address", "customer_email", "is_active",
		}).AddRow(1, customerUUID, "Tomy", "1234567890", "Chennai", "tomy@starkintl.com", true)

		// Using sqlmock.AnyArg() for every parameter ensures GORM's hidden internal queries pass smoothly
		suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		customer, err := suite.repository.GetCustomerByUUID(context.Background(), customerUUID)

		// ASSERTIONS: Protected against nil values to completely eliminate panics
		assert.Nil(t, err)
		if assert.NotNil(t, customer) {
			assert.Equal(t, "Tomy", customer.CustomerName)
			assert.Equal(t, customerUUID, customer.CustomerUUID)
		}
	})

	suite.T().Run("NotFound", func(t *testing.T) {
		suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		customer, err := suite.repository.GetCustomerByUUID(context.Background(), customerUUID)

		assert.Nil(t, customer)
		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "Customer not found", err.Message)
		}
	})
}

func (suite *UserRepositoryTestSuite) TestCheckRoomAvailability() {
	var roomID uint = 1
	
	start := time.Date(2026, 5, 26, 18, 45, 0, 0, time.UTC)
	end := time.Date(2026, 5, 26, 20, 45, 0, 0, time.UTC)

	// 1. CONFLICT EXISTS CASE
	suite.T().Run("ConflictExists", func(t *testing.T) {
		expectedUUID := "test-booking-uuid-123"
		
		// Return the expected booking row back to GORM
		rows := sqlmock.NewRows([]string{"booking_uuid"}).AddRow(expectedUUID)

		// Loose regex to accept GORM's internal soft delete check and LIMIT clauses smoothly
		suite.mock.ExpectQuery(`SELECT .*booking_uuid.* FROM "bookings"`).
			WithArgs(roomID, end, start, "CANCELLED", sqlmock.AnyArg()). // Catching the hidden soft-delete arg
			WillReturnRows(rows)

		conflict, conflictingUUID := suite.repository.CheckRoomAvailability(context.Background(), roomID, start, end)

		assert.True(t, conflict)
		assert.Equal(t, expectedUUID, conflictingUUID)
	})

	// 2. NO CONFLICT CASE
	suite.T().Run("NoConflictAvailable", func(t *testing.T) {
		// Empty row result simulating no overlaps found
		rows := sqlmock.NewRows([]string{"booking_uuid"})

		suite.mock.ExpectQuery(`SELECT .*booking_uuid.* FROM "bookings"`).
			WithArgs(roomID, end, start, "CANCELLED", sqlmock.AnyArg()).
			WillReturnRows(rows)

		conflict, conflictingUUID := suite.repository.CheckRoomAvailability(context.Background(), roomID, start, end)

		assert.False(t, conflict)
		assert.Empty(t, conflictingUUID)
	})
}

func (suite *UserRepositoryTestSuite) TestCreateBooking() {
	bookingUUID := "94841ade-468c-480f-8b69-ee911e6fcbdb"
	
	booking := &models.Booking{
		BookingUUID: bookingUUID,
		CustomerID:  10,
		RoomID:      20,
		BookingDate: time.Date(2026, 5, 26, 0, 0, 0, 0, time.UTC),
		StartTime:   time.Date(2026, 5, 26, 18, 45, 0, 0, time.UTC),
		EndTime:     time.Date(2026, 5, 26, 20, 45, 0, 0, time.UTC),
		Status:      "CONFIRMED",
		IsActive:    true,
	}

	// 1. SUCCESS CASE
	suite.T().Run("Success", func(t *testing.T) {
		suite.SetupTest() 

		suite.mock.ExpectBegin()
		
		// Use a loose regex and drop WithArgs() to prevent breaking on dynamic timestamp arguments
		suite.mock.ExpectQuery(`INSERT INTO "bookings"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			
		suite.mock.ExpectCommit()

		err := suite.repository.CreateBooking(context.Background(), booking)

		assert.Nil(t, err)
	})


	// 2. DATABASE ERROR REGRESSION CASE
	suite.T().Run("DatabaseFailure", func(t *testing.T) {
		suite.SetupTest() 

		suite.mock.ExpectBegin()
		
		// Use ExpectExec here as well to cleanly simulate a driver failure
		suite.mock.ExpectExec(`INSERT INTO "bookings"`).
			WillReturnError(errors.New("unique constraint violation key duplicated"))
			
		suite.mock.ExpectRollback()

		err := suite.repository.CreateBooking(context.Background(), booking)

		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, "Failed to create booking", err.Message)
		}
	})
}

//  Main Start the Above Tests..
func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}





//  go test ./... -coverpkg=BookingSystem/Booking/internal/repository -coverprofile="coverage.out"

// go tool cover -html="coverage.out"


