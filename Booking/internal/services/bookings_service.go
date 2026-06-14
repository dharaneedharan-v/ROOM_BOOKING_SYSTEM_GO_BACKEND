package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"BookingSystem/Booking/internal/dtos"
	"BookingSystem/Booking/internal/loggers"
	"BookingSystem/Booking/internal/models"
	"BookingSystem/Booking/internal/repository"

	"github.com/google/uuid"
	// "go.uber.org/zap"
)

type UserServiceInterface interface {
	CreateBooking(ctx context.Context, req dtos.BookingRequest, requestID string) *dtos.APIResponse
	GetBookingByID(ctx context.Context, bookingUUID, requestID string) *dtos.APIResponse
	//
	// GetAvailableRooms(ctx context.Context, date string, requestID string) *dtos.APIResponse
	UpdateBooking(ctx context.Context, bookingUUID string, req dtos.BookingRequest, requestID string) *dtos.APIResponse
	DeleteBooking(ctx context.Context, bookingUUID string, requestID string) *dtos.APIResponse
}

type UserService struct {
	userRepo repository.UserRepositoryInterface
	logger   *loggers.Logger
}

func NewUserService(userRepo repository.UserRepositoryInterface, logger *loggers.Logger) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Create Booking...
func (s *UserService) CreateBooking(ctx context.Context, req dtos.BookingRequest, requestID string) *dtos.APIResponse {
	s.logger.Info("Received request to get booking by ID---------[CreateBooking---SERVICE]")

	// Prasing the Time and date For the Validation of the date and time.
	bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "Invalid booking_date format. Use YYYY-MM-DD",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "Invalid start_time format. Use HH:MM (24-hour style)",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "Invalid end_time format. Use HH:MM (24-hour style)",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	//Local Time Setup and for the Server Side Time Validations [Example : Github]
	loc, _ := time.LoadLocation("Asia/Kolkata")
	start := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), startTime.Hour(), startTime.Minute(), 0, 0, loc)
	end := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), endTime.Hour(), endTime.Minute(), 0, 0, loc)

	// Negative Test case - Start and End Are same time.
	if end.Equal(start) {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "Start time and End time Are Same......!!!!!",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Negative Test case- Invalid Time Line For the Room Booking [End Time Must be Greater then Start Time] .
	if !end.After(start) {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "End time must be grater than the start time",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	//  Fetch the exact system clock time right 
	serverSystemTime := time.Now().In(loc)

	serverSystemDate := time.Date(serverSystemTime.Year(), serverSystemTime.Month(), serverSystemTime.Day(), 0, 0, 0, 0, loc)
	userRequestedDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)

	// Negative Test case - Checking for the past date.
	if userRequestedDate.Before(serverSystemDate) {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "Booking rejected due to the past date....!!!",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	//  Negative Test case - PAST Time Validation For the Same date. 
	if userRequestedDate.Equal(serverSystemDate) && start.Before(serverSystemTime) {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "Booking rejected due to Time Slot has passed.....!!!",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	//   Negative Test case - Checking Customer UUIDs are in the DB are Not..
	customer, custErr := s.userRepo.GetCustomerByUUID(ctx, req.CustomerUUID)
	if custErr != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusNotFound,
			RequestID: requestID,
			Message:   custErr.Message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	//  MUTUX LOCK.... 
	// =================================================================
	// LOCKING MECHANISM START (ADD THIS BLOCK RIGHT HERE) [ add the sync in the strut.. ]
	// =================================================================
			// actualLock, _ := s.roomLocks.LoadOrStore(req.RoomUUID, &sync.Mutex{})
			// roomMutex := actualLock.(*sync.Mutex)

			// roomMutex.Lock()
			
			// defer roomMutex.Unlock()
	// =================================================================

	//  Negative Test case - Checking the Room UUID are present in the DB are not.

	room, roomErr := s.userRepo.GetRoomByUUID(ctx, req.RoomUUID)
	if roomErr != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusNotFound,
			RequestID: requestID,
			Message:   roomErr.Message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Check for Any Conflicts.
	conflict, conflictingUUID := s.userRepo.CheckRoomAvailability(ctx, room.ID, start, end)
	if conflict {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusConflict,
			RequestID: requestID,
			Message:   fmt.Sprintf("Room already booked for this time slot. Conflicting Booking UUID: %s", conflictingUUID),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Finally Save it in the DB
	booking := &models.Booking{
		BookingUUID:  uuid.NewString(),
		CustomerID:   customer.ID,
		CustomerUUID: customer.CustomerUUID, // FK UUID
		RoomID:       room.ID,
		RoomUUID:     room.RoomUUID, // FK UUID
		BookingDate:  bookingDate,
		StartTime:    start,
		EndTime:      end,
	}

	if err := s.userRepo.CreateBooking(ctx, booking); err != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusInternalServerError,
			RequestID: requestID,
			Message:   err.Message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	return &dtos.APIResponse{
		Status:    "Success",
		Code:      http.StatusCreated,
		RequestID: requestID,
		Message:   "Room booked successfully",
		Data: &dtos.BookingResponse{
			BookingUUID:  booking.BookingUUID,
			CustomerUUID: customer.CustomerUUID,
			RoomUUID:     room.RoomUUID,
			BookingDate:  req.BookingDate,
			StartTime:    req.StartTime,
			EndTime:      req.EndTime,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
// Update Bookings
func (s *UserService) UpdateBooking(ctx context.Context, bookingUUID string, req dtos.BookingRequest, requestID string) *dtos.APIResponse {
	s.logger.Info("Received request to update booking by ID---------[UpdateBooking---SERVICE]")

	// 1. SET THE TARGET LOCAL TIMEZONE
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusInternalServerError,
			RequestID: requestID,
			Message:   "Internal server configuration error: invalid timezone data",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// 2. FETCH EXISTING BOOKING
	existingBooking, errRepo := s.userRepo.GetBookingByUUID(ctx, bookingUUID)
	if errRepo != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusNotFound,
			RequestID: requestID,
			Message:   "Original booking not found",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// 3. BLOCK MODIFICATION OF PAST BOOKINGS
	serverSystemTime := time.Now().In(loc)
	if existingBooking.StartTime.Before(serverSystemTime) {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "Cannot modify a booking that has already started or passed",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// 4. INPUT TIME PARSING
	bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "Invalid booking_date format. Use YYYY-MM-DD",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "Invalid start_time format. Use HH:MM (24-hour style)",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "Invalid end_time format. Use HH:MM (24-hour style)",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// 5. TIMELINE CONSTRUCTION
	start := time.Date(
		bookingDate.Year(), bookingDate.Month(), bookingDate.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, loc,
	)

	end := time.Date(
		bookingDate.Year(), bookingDate.Month(), bookingDate.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, loc,
	)

	// 6. CHRONOLOGICAL SAFEGUARD
	if !end.After(start) {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "End time must be strictly after the start time",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// 7. NEW PAST DATE & TIME CHECK (Combined & simplified)
	if start.Before(serverSystemTime) {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Message:   "Cannot reschedule a booking to a past date or time slot",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// 5. RESOURCE AVAILABILITY LOOKUPS
	customer, custErr := s.userRepo.GetCustomerByUUID(ctx, req.CustomerUUID)
	if custErr != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusNotFound,
			RequestID: requestID,
			Message:   custErr.Message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	room, roomErr := s.userRepo.GetRoomByUUID(ctx, req.RoomUUID)
	if roomErr != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusNotFound,
			RequestID: requestID,
			Message:   roomErr.Message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// 9. CONFLICT CHECK EXCLUDING SELF
	conflict := s.userRepo.CheckRoomAvailabilityForUpdate(ctx, room.ID, start, end, bookingUUID)
	if conflict {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusConflict,
			RequestID: requestID,
			Message:   "Room is already occupied by another booking during this time frame",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// 10. MAP MODIFICATIONS AND PERSIST
	existingBooking.CustomerID = customer.ID
	existingBooking.RoomID = room.ID
	existingBooking.BookingDate = bookingDate
	existingBooking.StartTime = start
	existingBooking.EndTime = end

	if updateErr := s.userRepo.UpdateBooking(ctx, existingBooking); updateErr != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusInternalServerError,
			RequestID: requestID,
			Message:   updateErr.Message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	return &dtos.APIResponse{
		Status:    "Success",
		Code:      http.StatusOK,
		RequestID: requestID,
		Message:   "Booking updated successfully",
		Data: &dtos.BookingResponse{
			BookingUUID:  existingBooking.BookingUUID,
			CustomerUUID: req.CustomerUUID,
			RoomUUID:     room.RoomUUID,
			BookingDate:  req.BookingDate,
			StartTime:    req.StartTime,
			EndTime:      req.EndTime,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
// Get Booking 
func (s *UserService) GetBookingByID(ctx context.Context, bookingUUID string, requestID string) *dtos.APIResponse {
	s.logger.Info("Received request to get booking by ID---------[GetBookingByID---SERVICE]")
	booking, errRepo := s.userRepo.GetBookingByUUID(ctx, bookingUUID)
	if errRepo != nil {
		return &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusNotFound,
			RequestID: requestID,
			Message:   errRepo.Message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}
	response := &dtos.BookingResponse{
		BookingUUID:  booking.BookingUUID,
		CustomerUUID: booking.CustomerUUID, //  FK
		RoomUUID:     booking.RoomUUID,     //  FK
		BookingDate:  booking.BookingDate.Format("2006-01-02"),
		StartTime:    booking.StartTime.Format("15:04"),
		EndTime:      booking.EndTime.Format("15:04"),
	}

	return &dtos.APIResponse{
		Status:    "Success",
		Code:      http.StatusOK,
		RequestID: requestID,
		Message:   "Booking retrieved successfully",
		Data:      response,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func (s *UserService) DeleteBooking(ctx context.Context, bookingUUID string, requestID string) *dtos.APIResponse {

	s.logger.Info("Received request to get booking by ID---------[DeleteBooking ---SERVICE]")
	err := s.userRepo.SoftDeleteBooking(ctx, bookingUUID)
	if err != nil {
		statusCode := http.StatusNotFound
		// If it exists ,  already cancelled.
		if err.Message == "Booking Already Cancelled...!!" {
			statusCode = http.StatusBadRequest
		}

		return &dtos.APIResponse{
			Status:    "Error",
			Code:      statusCode,
			RequestID: requestID,
			Message:   err.Message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	return &dtos.APIResponse{
		Status:    "Success",
		Code:      http.StatusOK,
		RequestID: requestID,
		Message:   "Booking deleted successfully",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}













// func (s *UserService) GetAvailableRooms(ctx context.Context, date string, requestID string) *dtos.APIResponse {
// 	parsedDate, err := time.Parse("2006-01-02", date)
// 	if err != nil {
// 		return &dtos.APIResponse{
// 			Status:    "Error",
// 			Code:      http.StatusBadRequest,
// 			RequestID: requestID,
// 			Message:   "Invalid date format. Use YYYY-MM-DD",
// 			Timestamp: time.Now().UTC().Format(time.RFC3339),
// 		}
// 	}

// 	rooms, repoErr := s.userRepo.GetAvailableRooms(ctx, parsedDate)
// 	if repoErr != nil {
// 		return &dtos.APIResponse{
// 			Status:    "Error",
// 			Code:      http.StatusInternalServerError,
// 			RequestID: requestID,
// 			Message:   repoErr.Message,
// 			Timestamp: time.Now().UTC().Format(time.RFC3339),
// 		}
// 	}

// 	response := make([]dtos.AvailableRoomResponse, 0)
// 	for _, room := range rooms {
// 		response = append(response, dtos.AvailableRoomResponse{
// 			RoomUUID: room.RoomUUID,
// 			RoomName: room.RoomName,
// 			Capacity: room.Capacity,
// 			Status:   "Available",
// 		})
// 	}

// 	return &dtos.APIResponse{
// 		Status:    "Success",
// 		Code:      http.StatusOK,
// 		RequestID: requestID,
// 		Message:   "Rooms retrieved successfully",
// 		Data:      response,
// 		Timestamp: time.Now().UTC().Format(time.RFC3339),
// 	}
// }
