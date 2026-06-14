package services_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"BookingSystem/Booking/internal/dtos"
	"BookingSystem/Booking/internal/loggers"
	"BookingSystem/Booking/internal/models"
	"BookingSystem/Booking/internal/services"
	repoMock "BookingSystem/Booking/tests/unit/repository/mock" 

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	
)

func TestCreateBooking(t *testing.T) {
	mockRepo := new(repoMock.UserRepositoryMock)
	logger := loggers.NewTestLogger()
	userService := services.NewUserService(mockRepo, logger)

	requestID := "req-test-123"

	loc, _ := time.LoadLocation("Asia/Kolkata")
	futureDateStr := time.Now().In(loc).AddDate(0, 0, 2).Format("2006-01-02")
	pastDateStr := time.Now().In(loc).AddDate(0, 0, -2).Format("2006-01-02")

	todayDateStr := time.Now().In(loc).Format("2006-01-02")
	twoHoursAgoStr := time.Now().In(loc).Add(-2 * time.Hour).Format("15:04")


	testCases := []struct {
		name          string
		request       dtos.BookingRequest
		mockErrorMsg  string
		conflictFound bool
		expectedError bool
		expectedMsg   string
	}{
		{
			name: "Success-RoomBookedSuccessfully",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "",
			conflictFound: false,
			expectedError: false,
			expectedMsg:   "Room booked successfully",
		},
		{
			name: "Success-MidnightCrossoverRescued",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "23:30",
				EndTime:      "01:30",
			},
			mockErrorMsg:  "",
			conflictFound: false,
			expectedError: true,
			expectedMsg:   "End time must be grater than the start time",
		},
		{
			name: "Failure-InvalidBookingDateFormat",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  "25/05/2026",
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "",
			expectedError: true,
			expectedMsg:   "Invalid booking_date format. Use YYYY-MM-DD",
		},
		{
			name: "Failure-ZeroDurationSlotRejected",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "14:00",
			},
			mockErrorMsg:  "",
			expectedError: true,
			expectedMsg:   "Start time and End time Are Same......!!!!!",
		},
		{
			name: "Failure-PastDateNotAllowed",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  pastDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "",
			expectedError: true,
			expectedMsg:   "Booking rejected due to the past date....!!!",
		},
		{
			name: "Failure-CustomerNotFound",
			request: dtos.BookingRequest{
				CustomerUUID: "invalid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "Customer not found",
			expectedError: true,
			expectedMsg:   "Customer not found",
		},
		{
			name: "Failure-RoomNotFound",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "invalid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "Room not found",
			expectedError: true,
			expectedMsg:   "Room not found",
		},
		{
			name: "Failure-RoomAlreadyBookedConflict",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "",
			conflictFound: true,
			expectedError: true,
			expectedMsg:   "Room already booked for this time slot",
		},
		{
			name: "Failure-CreateBookingDBError",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "Failed to create booking",
			conflictFound: false,
			expectedError: true,
			expectedMsg:   "Failed to create booking",
		},

		// ----
		{
			name: "Success-RoomBookedSuccessfully",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "",
			conflictFound: false,
			expectedError: false,
			expectedMsg:   "Room booked successfully",
		},
		{
			name: "Negative-InvalidStartTimeFormat",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "02:00PM", // Invalid: Has PM indicator
				EndTime:      "16:00",
			},
			mockErrorMsg:  "",
			conflictFound: false,
			expectedError: true,
			expectedMsg:   "Invalid start_time format. Use HH:MM (24-hour style)",
		},
		{
			name: "Negative-InvalidEndTimeFormat",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "99:99", // Invalid: Impossible hours and minutes
			},
			mockErrorMsg:  "",
			conflictFound: false,
			expectedError: true,
			expectedMsg:   "Invalid end_time format. Use HH:MM (24-hour style)",
		},
		{
			name: "Negative-PastTimeValidationForToday",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,       
				StartTime:    pastDateStr,  
				EndTime:      "23:59",           
			},
			mockErrorMsg:  "",
			conflictFound: false,
			expectedError: true,
			expectedMsg:   "Booking rejected due to Time Slot has passed.....!!!",
		},
		{
			name: "Negative-PastTimeValidationForToday",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  todayDateStr,    
				StartTime:    twoHoursAgoStr, 
				EndTime:      "23:59",
			},
			mockErrorMsg:  "",
			conflictFound: false,
			expectedError: true,
			expectedMsg:   "Booking rejected due to Time Slot has passed.....!!!",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			bDate, errDate := time.Parse("2006-01-02", tc.request.BookingDate)
			sTime, errStart := time.Parse("15:04", tc.request.StartTime)
			eTime, errEnd := time.Parse("15:04", tc.request.EndTime)

			if errDate == nil && errStart == nil && errEnd == nil {
				start := time.Date(bDate.Year(), bDate.Month(), bDate.Day(), sTime.Hour(), sTime.Minute(), 0, 0, loc)
				end := time.Date(bDate.Year(), bDate.Month(), bDate.Day(), eTime.Hour(), eTime.Minute(), 0, 0, loc)

				serverNow := time.Now().In(loc)
				serverDate := time.Date(serverNow.Year(), serverNow.Month(), serverNow.Day(), 0, 0, 0, 0, loc)
				requestedDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)

				isZeroDuration := sTime.Equal(eTime)
				isEndBeforeStart := !end.After(start)
				isPastDate := requestedDate.Before(serverDate)
				isPastTimeToday := requestedDate.Equal(serverDate) && start.Before(serverNow)

				shouldMock := !isZeroDuration && !isEndBeforeStart && !isPastDate && !isPastTimeToday

				if shouldMock {
					// 1. Customer Error Mock Path
					if tc.mockErrorMsg == "Customer not found" {
						mockRepo.On("GetCustomerByUUID", mock.Anything, tc.request.CustomerUUID).
							Return(nil, &dtos.Error{Message: tc.mockErrorMsg}).Once()

					// 2. Room Error Mock Path
					} else if tc.mockErrorMsg == "Room not found" {
						mockRepo.On("GetCustomerByUUID", mock.Anything, tc.request.CustomerUUID).
							Return(&models.Customer{ID: 1, CustomerUUID: tc.request.CustomerUUID}, nil).Once()
						mockRepo.On("GetRoomByUUID", mock.Anything, tc.request.RoomUUID).
							Return(nil, &dtos.Error{Message: tc.mockErrorMsg}).Once()

					// 3. Time Overlap Conflict Mock Path
					} else if tc.conflictFound {
						mockRepo.On("GetCustomerByUUID", mock.Anything, tc.request.CustomerUUID).
							Return(&models.Customer{ID: 1, CustomerUUID: tc.request.CustomerUUID}, nil).Once()
						mockRepo.On("GetRoomByUUID", mock.Anything, tc.request.RoomUUID).
							Return(&models.Room{ID: 5, RoomUUID: tc.request.RoomUUID}, nil).Once()
						mockRepo.On("CheckRoomAvailability", mock.Anything, uint(5), mock.Anything, mock.Anything).
							Return(true, "conflicting-booking-uuid-abc").Once()

					// 4. Persistence DB Error Mock Path
					} else if tc.mockErrorMsg == "Failed to create booking" {
						mockRepo.On("GetCustomerByUUID", mock.Anything, tc.request.CustomerUUID).
							Return(&models.Customer{ID: 1, CustomerUUID: tc.request.CustomerUUID}, nil).Once()
						mockRepo.On("GetRoomByUUID", mock.Anything, tc.request.RoomUUID).
							Return(&models.Room{ID: 5, RoomUUID: tc.request.RoomUUID}, nil).Once()
						mockRepo.On("CheckRoomAvailability", mock.Anything, uint(5), mock.Anything, mock.Anything).
							Return(false, "").Once()
						mockRepo.On("CreateBooking", mock.Anything, mock.AnythingOfType("*models.Booking")).
							Return(&dtos.Error{Message: tc.mockErrorMsg}).Once()

					// 5. Complete Success Mock Path
					} else {
						mockRepo.On("GetCustomerByUUID", mock.Anything, tc.request.CustomerUUID).
							Return(&models.Customer{ID: 1, CustomerUUID: tc.request.CustomerUUID}, nil).Once()
						mockRepo.On("GetRoomByUUID", mock.Anything, tc.request.RoomUUID).
							Return(&models.Room{ID: 5, RoomUUID: tc.request.RoomUUID}, nil).Once()
						mockRepo.On("CheckRoomAvailability", mock.Anything, uint(5), mock.Anything, mock.Anything).
							Return(false, "").Once()
						mockRepo.On("CreateBooking", mock.Anything, mock.AnythingOfType("*models.Booking")).
							Return(nil).Once()
					}
				}
			}

			// Call the service
			response := userService.CreateBooking(context.Background(), tc.request, requestID)

			// Assertions
			assert.NotNil(t, response)
			if tc.expectedError {
				assert.Equal(t, "Error", response.Status)
				// Conflict message uses fmt.Sprintf with UUID so use Contains
				if tc.conflictFound {
					assert.Contains(t, response.Message, tc.expectedMsg)
				} else {
					assert.Equal(t, tc.expectedMsg, response.Message)
				}
			} else {
				assert.Equal(t, "Success", response.Status)
				assert.Equal(t, http.StatusCreated, response.Code)
				assert.Equal(t, tc.expectedMsg, response.Message)
				assert.NotNil(t, response.Data)

				resResp, ok := response.Data.(*dtos.BookingResponse)
				if assert.True(t, ok) {
					assert.Equal(t, tc.request.BookingDate, resResp.BookingDate)
					assert.Equal(t, tc.request.StartTime, resResp.StartTime)
					assert.Equal(t, tc.request.EndTime, resResp.EndTime)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}


func TestGetBookingByID(t *testing.T) {
	// Create mock and logger dependencies
	mockRepo := new(repoMock.UserRepositoryMock)

	// logger := loggers.NewLogger(loggers.LogConfig{Level: "info", ServiceName: "test-service"})
			logger := loggers.NewTestLogger()

	userService := services.NewUserService(mockRepo, logger)

	requestID := "req-test-456"
	bookingUUID := "94841ade-468c-480f-8b69-ee911e6fcbdb"
	
	// Create a mock model layout matching a successful database preload result
	mockBooking := &models.Booking{
		ID:           1,
		BookingUUID:  bookingUUID,
		CustomerUUID: "customer-uuid-123", 
		RoomUUID:     "room-uuid-456",     
		BookingDate: time.Date(2027, 5, 25, 0, 0, 0, 0, time.UTC),
		StartTime:   time.Date(2027, 5, 25, 18, 45, 0, 0, time.UTC),
		EndTime:     time.Date(2027, 5, 25, 20, 45, 0, 0, time.UTC),
	}

	testCases := []struct {
		name          string
		bookingUUID   string
		mockRepoError *dtos.Error
		expectedError bool
		expectedCode  int
		expectedMsg   string
	}{
		{
			name:          "Success-BookingRetrievedSuccessfully",
			bookingUUID:   bookingUUID,
			mockRepoError: nil,
			expectedError: false,
			expectedCode:  http.StatusOK,
			expectedMsg:   "Booking retrieved successfully",
		},
		{
			name:         "Failure-BookingNotFound",
			bookingUUID:  "non-existent-uuid",
			mockRepoError: &dtos.Error{Message: "Booking not found"},
			expectedError: true,
			expectedCode:  http.StatusNotFound,
			expectedMsg:   "Booking not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear expectation traces before each individual loop run
			mockRepo.ExpectedCalls = nil

			// Setup Mock expectations
			if tc.mockRepoError != nil {
				mockRepo.On("GetBookingByUUID", mock.Anything, tc.bookingUUID).
					Return(nil, tc.mockRepoError).Once()
			} else {
				mockRepo.On("GetBookingByUUID", mock.Anything, tc.bookingUUID).
					Return(mockBooking, nil).Once()
			}

			// Call the service method
			response := userService.GetBookingByID(context.Background(), tc.bookingUUID, requestID)

			// Assertions matching your structural reference pattern
			assert.NotNil(t, response)
			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedError {
				assert.Equal(t, "Error", response.Status)
				assert.Equal(t, tc.expectedMsg, response.Message)
			} else {
				assert.Equal(t, "Success", response.Status)
				assert.Equal(t, tc.expectedMsg, response.Message)
				assert.NotNil(t, response.Data)

				resResp, ok := response.Data.(*dtos.BookingResponse)
				if assert.True(t, ok) {
					assert.Equal(t, bookingUUID, resResp.BookingUUID)
					assert.Equal(t, "customer-uuid-123", resResp.CustomerUUID)
					assert.Equal(t, "room-uuid-456", resResp.RoomUUID)
					assert.Equal(t, "2027-05-25", resResp.BookingDate)
					assert.Equal(t, "18:45", resResp.StartTime)
					assert.Equal(t, "20:45", resResp.EndTime)
				}
			}

			// Verify mock assertions
			mockRepo.AssertExpectations(t)
		})
	}
}



func TestDeleteBooking(t *testing.T) {
	// Create mock and logger dependencies
	mockRepo := new(repoMock.UserRepositoryMock)
	// logger := loggers.NewLogger(loggers.LogConfig{Level: "info", ServiceName: "test-service"})
			logger := loggers.NewTestLogger()

	userService := services.NewUserService(mockRepo, logger)

	requestID := "req-test-789"
	bookingUUID := "94841ade-468c-480f-8b69-ee911e6fcbdb"

	testCases := []struct {
		name          string
		bookingUUID   string
		mockRepoError *dtos.Error
		expectedError bool
		expectedCode  int
		expectedMsg   string
	}{
		{
			name:          "Success-BookingDeletedSuccessfully",
			bookingUUID:   bookingUUID,
			mockRepoError: nil,
			expectedError: false,
			expectedCode:  http.StatusOK,
			expectedMsg:   "Booking deleted successfully",
		},
		{
			name:         "Failure-BookingNotFound",
			bookingUUID:  "missing-uuid",
			mockRepoError: &dtos.Error{Message: "Booking not found"},
			expectedError: true,
			expectedCode:  http.StatusNotFound, // Evaluates to 404 based on status code conditions
			expectedMsg:   "Booking not found",
		},
		{
			name:         "Failure-BookingAlreadyCancelled",
			bookingUUID:  bookingUUID,
			mockRepoError: &dtos.Error{Message: "Booking Already Cancelled...!!"},
			expectedError: true,
			expectedCode:  http.StatusBadRequest, // Evaluates to 400 based on message match condition
			expectedMsg:   "Booking Already Cancelled...!!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear expectation traces before each individual loop run
			mockRepo.ExpectedCalls = nil

			// Setup Mock expectations
			if tc.mockRepoError != nil {
				mockRepo.On("SoftDeleteBooking", mock.Anything, tc.bookingUUID).
					Return(tc.mockRepoError).Once()
			} else {
				mockRepo.On("SoftDeleteBooking", mock.Anything, tc.bookingUUID).
					Return(nil).Once()
			}

			// Call the service method
			response := userService.DeleteBooking(context.Background(), tc.bookingUUID, requestID)

			// Assertions matching your structural reference pattern
			assert.NotNil(t, response)
			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedError {
				assert.Equal(t, "Error", response.Status)
				assert.Equal(t, tc.expectedMsg, response.Message)
			} else {
				assert.Equal(t, "Success", response.Status)
				assert.Equal(t, tc.expectedMsg, response.Message)
			}

			// Verify mock assertions
			mockRepo.AssertExpectations(t)
		})
	}
}


func TestUpdateBooking(t *testing.T) {
	mockRepo := new(repoMock.UserRepositoryMock)
	logger := loggers.NewTestLogger()
	userService := services.NewUserService(mockRepo, logger)

	requestID := "req-test-update-123"
	bookingUUID := "existing-booking-uuid"

	loc, _ := time.LoadLocation("Asia/Kolkata")
	futureDateStr := time.Now().In(loc).AddDate(0, 0, 2).Format("2006-01-02")
	pastDateStr := time.Now().In(loc).AddDate(0, 0, -2).Format("2006-01-02")

	// A future StartTime for the existing booking (not yet started)
	futureStart := time.Now().In(loc).AddDate(0, 0, 2)
	// A past StartTime for the existing booking (already started)
	pastStart := time.Now().In(loc).AddDate(0, 0, -1)

	testCases := []struct {
		name              string
		bookingUUID       string
		request           dtos.BookingRequest
		mockErrorMsg      string
		bookingNotFound   bool
		existingIsPast    bool
		conflictFound     bool
		expectedError     bool
		expectedMsg       string
	}{
		{
			name:        "Success-BookingUpdatedSuccessfully",
			bookingUUID: bookingUUID,
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			expectedError: false,
			expectedMsg:   "Booking updated successfully",
		},
		{
			name:            "Failure-BookingNotFound",
			bookingUUID:     "non-existent-uuid",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			bookingNotFound: true,
			expectedError:   true,
			expectedMsg:     "Original booking not found",
		},
		{
			name:           "Failure-ExistingBookingAlreadyStarted",
			bookingUUID:    bookingUUID,
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			existingIsPast: true,
			expectedError:  true,
			expectedMsg:    "Cannot modify a booking that has already started or passed",
		},
		{
			name:        "Failure-InvalidBookingDateFormat",
			bookingUUID: bookingUUID,
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  "25/05/2026",
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			expectedError: true,
			expectedMsg:   "Invalid booking_date format. Use YYYY-MM-DD",
		},
		{
			name:        "Failure-EndTimeNotAfterStartTime",
			bookingUUID: bookingUUID,
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "16:00",
				EndTime:      "14:00",
			},
			expectedError: true,
			expectedMsg:   "End time must be strictly after the start time",
		},
		{
			name:        "Failure-ZeroDurationSlot",
			bookingUUID: bookingUUID,
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "14:00",
			},
			expectedError: true,
			expectedMsg:   "End time must be strictly after the start time",
		},
		{
			name:        "Failure-RescheduleToPastDateTime",
			bookingUUID: bookingUUID,
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  pastDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			expectedError: true,
			expectedMsg:   "Cannot reschedule a booking to a past date or time slot",
		},
		{
			name:        "Failure-CustomerNotFound",
			bookingUUID: bookingUUID,
			request: dtos.BookingRequest{
				CustomerUUID: "invalid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "Customer not found",
			expectedError: true,
			expectedMsg:   "Customer not found",
		},
		{
			name:        "Failure-RoomNotFound",
			bookingUUID: bookingUUID,
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "invalid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "Room not found",
			expectedError: true,
			expectedMsg:   "Room not found",
		},
		{
			name:        "Failure-RoomConflictWithAnotherBooking",
			bookingUUID: bookingUUID,
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			conflictFound: true,
			expectedError: true,
			expectedMsg:   "Room is already occupied by another booking during this time frame",
		},
		{
			name:        "Failure-UpdateBookingDBError",
			bookingUUID: bookingUUID,
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "Failed to update booking",
			expectedError: true,
			expectedMsg:   "Failed to update booking",
		},

		{
			name: "Negative-InvalidStartTimeFormat",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "02:00PM", // Invalid: Has PM indicator
				EndTime:      "16:00",
			},
			mockErrorMsg:  "",
			conflictFound: false,
			expectedError: true,
			expectedMsg:   "Invalid start_time format. Use HH:MM (24-hour style)",
		},
		{
			name: "Negative-InvalidEndTimeFormat",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  futureDateStr,
				StartTime:    "14:00",
				EndTime:      "99:99", // Invalid: Impossible hours and minutes
			},
			mockErrorMsg:  "",
			conflictFound: false,
			expectedError: true,
			expectedMsg:   "Invalid end_time format. Use HH:MM (24-hour style)",
		},

		{
			name: "Negative-TimezoneLoadFailure",
			request: dtos.BookingRequest{
				CustomerUUID: "valid-cust-uuid",
				RoomUUID:     "valid-room-uuid",
				BookingDate:  "2026-06-15",
				StartTime:    "14:00",
				EndTime:      "16:00",
			},
			mockErrorMsg:  "",
			conflictFound: false,
			expectedError: true,
			expectedMsg:   "Internal server configuration error: invalid timezone data",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			// Determine existing booking start time
			existingStart := futureStart
			if tc.existingIsPast {
				existingStart = pastStart
			}

			existingBooking := &models.Booking{
				ID:          1,
				BookingUUID: tc.bookingUUID,
				CustomerID:  1,
				RoomID:      5,
				StartTime:   existingStart,
				EndTime:     existingStart.Add(2 * time.Hour),
			}

			// Step 1: GetBookingByUUID
			if tc.bookingNotFound {
				mockRepo.On("GetBookingByUUID", mock.Anything, tc.bookingUUID).
					Return(nil, &dtos.Error{Message: "Original booking not found"}).Once()
				// No further mocks needed
				goto runTest
			}

			mockRepo.On("GetBookingByUUID", mock.Anything, tc.bookingUUID).
				Return(existingBooking, nil).Once()

			// Step 2: If existing booking already started, no further mocks needed
			if tc.existingIsPast {
				goto runTest
			}

			// Step 3: Only mock further if date/time parsing and validation will pass
			{
				bDate, errDate := time.Parse("2006-01-02", tc.request.BookingDate)
				sTime, errStart := time.Parse("15:04", tc.request.StartTime)
				eTime, errEnd := time.Parse("15:04", tc.request.EndTime)

				if errDate != nil || errStart != nil || errEnd != nil {
					goto runTest
				}

				start := time.Date(bDate.Year(), bDate.Month(), bDate.Day(), sTime.Hour(), sTime.Minute(), 0, 0, loc)
				end := time.Date(bDate.Year(), bDate.Month(), bDate.Day(), eTime.Hour(), eTime.Minute(), 0, 0, loc)

				// Skip mocking if chronological or past checks will fail
				if !end.After(start) {
					goto runTest
				}
				if start.Before(time.Now().In(loc)) {
					goto runTest
				}

				// Step 4: Customer lookup
				if tc.mockErrorMsg == "Customer not found" {
					mockRepo.On("GetCustomerByUUID", mock.Anything, tc.request.CustomerUUID).
						Return(nil, &dtos.Error{Message: tc.mockErrorMsg}).Once()
					goto runTest
				}

				mockRepo.On("GetCustomerByUUID", mock.Anything, tc.request.CustomerUUID).
					Return(&models.Customer{ID: 1, CustomerUUID: tc.request.CustomerUUID}, nil).Once()

				// Step 5: Room lookup
				if tc.mockErrorMsg == "Room not found" {
					mockRepo.On("GetRoomByUUID", mock.Anything, tc.request.RoomUUID).
						Return(nil, &dtos.Error{Message: tc.mockErrorMsg}).Once()
					goto runTest
				}

				mockRepo.On("GetRoomByUUID", mock.Anything, tc.request.RoomUUID).
					Return(&models.Room{ID: 5, RoomUUID: tc.request.RoomUUID}, nil).Once()

				// Step 6: Conflict check
				if tc.conflictFound {
					mockRepo.On("CheckRoomAvailabilityForUpdate", mock.Anything, uint(5), mock.Anything, mock.Anything, tc.bookingUUID).
						Return(true).Once()
					goto runTest
				}

				mockRepo.On("CheckRoomAvailabilityForUpdate", mock.Anything, uint(5), mock.Anything, mock.Anything, tc.bookingUUID).
					Return(false).Once()

				// Step 7: UpdateBooking
				if tc.mockErrorMsg == "Failed to update booking" {
					mockRepo.On("UpdateBooking", mock.Anything, mock.AnythingOfType("*models.Booking")).
						Return(&dtos.Error{Message: tc.mockErrorMsg}).Once()
					goto runTest
				}

				mockRepo.On("UpdateBooking", mock.Anything, mock.AnythingOfType("*models.Booking")).
					Return(nil).Once()
			}

		runTest:
			response := userService.UpdateBooking(context.Background(), tc.bookingUUID, tc.request, requestID)

			assert.NotNil(t, response)
			if tc.expectedError {
				assert.Equal(t, "Error", response.Status)
				assert.Equal(t, tc.expectedMsg, response.Message)
			} else {
				assert.Equal(t, "Success", response.Status)
				assert.Equal(t, http.StatusOK, response.Code)
				assert.Equal(t, tc.expectedMsg, response.Message)
				assert.NotNil(t, response.Data)

				resResp, ok := response.Data.(*dtos.BookingResponse)
				if assert.True(t, ok) {
					assert.Equal(t, tc.request.BookingDate, resResp.BookingDate)
					assert.Equal(t, tc.request.StartTime, resResp.StartTime)
					assert.Equal(t, tc.request.EndTime, resResp.EndTime)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}




//  go test ./... -coverpkg=BookingSystem/Booking/internal/services -coverprofile="coverage.out"

// go tool cover -html="coverage.out"


// go test ./... -coverpkg=BookingSystem/Booking/internal/handler -coverprofile="coverage.out"



// Mockery : 


// Go to the Root Folder.. 
// 1) go install github.com/vektra/mockery/v2@latest
// 2)  Test-Path "$env:USERPROFILE\go\bin\mockery.exe"
// 3) $env:Path += ";$env:USERPROFILE\go\bin"
// 4) mockery --version
// 5) mockery --name=UserRepositoryInterface --dir=internal/repository --output=tests/unit/repository/mock --filename=bookings_repository_mock.go --structname=UserRepositoryMock --outpkg=mock

// [ change the File and folder name according to that ] 

