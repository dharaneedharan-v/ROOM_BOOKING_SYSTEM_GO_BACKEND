
package handlers_test

import (
	// "bytes"
	"bytes"
	"context"
	"io"
	"strings"

	// "context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	// "os"
	"testing"
	// "time"

	"BookingSystem/Booking/internal/dtos"
	"BookingSystem/Booking/internal/errorcodes"
	"BookingSystem/Booking/internal/handlers"
	"BookingSystem/Booking/internal/loggers"
	serviceMock "BookingSystem/Booking/tests/unit/services/mock"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// func TestCreateBooking_Handler(t *testing.T) {
// 	// Create mock and logger dependencies
// 	mockUserService := new(serviceMock.UserServiceMock)
// 	logger := loggers.NewTestLogger()

// 	// Create handler with mock service
// 	userHandler := handlers.NewUserHandler(mockUserService, logger)

// 	tomorrowStr := time.Now().Add(24 * time.Hour).Format("2006-01-02")

// 	// Test cases structural layout matching your reference format exactly machi!
// 	testCases := []struct {
// 		name           string
// 		requestBody    interface{} // Allows passing valid structures or bad broken raw strings
// 		mockResponse   *dtos.APIResponse
// 		expectedStatus int
// 		expectedError  bool
// 		isValidation   bool // Distinguishes between parsing breaks vs validation rule triggers
// 	}{
// 		{
// 			name: "Success - Room Booked Successfully",
// 			requestBody: map[string]interface{}{
// 				"customer_uuid": "e71260ef-4b14-4b99-9ef3-eba0ddfd48b3",
// 				"room_uuid":     "c50fe215-347f-46a0-bdc0-51479f96d451",
// 				"booking_date":  tomorrowStr,
// 				"start_time":    "18:45",
// 				"end_time":      "20:45",
// 			},
// 			mockResponse: &dtos.APIResponse{
// 				Status:  "Success",
// 				Code:    http.StatusCreated,
// 				Message: "Room booked successfully",
// 				Data: &dtos.BookingResponse{
// 					BookingUUID:  "generated-booking-uuid-xyz",
// 					CustomerUUID: "e71260ef-4b14-4b99-9ef3-eba0ddfd48b3",
// 					RoomUUID:     "c50fe215-347f-46a0-bdc0-51479f96d451",
// 					BookingDate:  tomorrowStr,
// 					StartTime:    "18:45",
// 					EndTime:      "20:45",
// 				},
// 			},
// 			expectedStatus: http.StatusCreated,
// 			expectedError:  false,
// 			isValidation:   false,
// 		},
// 		{
// 			name: "Failure - Validation Failed",
// 			requestBody: map[string]interface{}{
// 				"customer_uuid": "", // Missing required field triggers validation array output
// 				"room_uuid":     "c50fe215-347f-46a0-bdc0-51479f96d451",
// 				"booking_date":  tomorrowStr,
// 				"start_time":    "18:45",
// 				"end_time":      "20:45",
// 			},
// 			mockResponse:   nil, // Validation blocks flow early before service execution
// 			expectedStatus: http.StatusBadRequest,
// 			expectedError:  true,
// 			isValidation:   true,
// 		},
// 		{
// 			name:           "Failure - Invalid Request Format",
// 			requestBody:    "{invalid-json-body-raw-string", // Triggers json.NewDecoder parsing break block
// 			mockResponse:   nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedError:  true,
// 			isValidation:   false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Clear tracking mock states before testing individual paths
// 			mockUserService.ExpectedCalls = nil

// 			// Convert request body to bytes dynamically based on format pattern type
// 			var reqBody []byte
// 			if strBody, ok := tc.requestBody.(string); ok {
// 				reqBody = []byte(strBody)
// 			} else {
// 				reqBody, _ = json.Marshal(tc.requestBody)
// 			}

// 			// Create an HTTP request
// 			req, err := http.NewRequest("POST", "/bookings", bytes.NewBuffer(reqBody))
// 			assert.NoError(t, err)
// 			req.Header.Set("Content-Type", "application/json")

// 			// Record response
// 			rr := httptest.NewRecorder()

// 			// Setup mock expectations if execution path is expected to hit the business logic service layer
// 			if !tc.expectedError {
// 				// FIXED: Intercept dynamic handler parameters to assign the generated requestID dynamically
// 				mockUserService.On("CreateBooking", mock.Anything, mock.Anything, mock.AnythingOfType("string")).
// 					Return(func(ctx context.Context, bookingReq dtos.BookingRequest, reqID string) *dtos.APIResponse {
// 						if tc.mockResponse != nil {
// 							tc.mockResponse.RequestID = reqID
// 						}
// 						return tc.mockResponse
// 					})
// 			}

// 			// Call the target handler method directly
// 			userHandler.CreateBooking(rr, req)

// 			// Assertions matching your precise reference architecture
// 			assert.Equal(t, tc.expectedStatus, rr.Code)

// 			var response dtos.APIResponse
// 			err = json.Unmarshal(rr.Body.Bytes(), &response)
// 			assert.NoError(t, err)

// 			assert.NotEmpty(t, response.RequestID) // Verifies request ID tracking token is present

// 			if tc.expectedError {
// 				assert.Equal(t, "Error", response.Status)
// 				if tc.isValidation {
// 					assert.Equal(t, "Validation failed", response.Message)
// 				} else {
// 					assert.Equal(t, "Invalid request format", response.Message)
// 				}
// 			} else {
// 				// Verify mock was called
// 				mockUserService.AssertExpectations(t)

// 				assert.Equal(t, tc.mockResponse.Status, response.Status)
// 				assert.Equal(t, tc.mockResponse.Message, response.Message)
// 			}
// 		})
// 	}
// }

// func TestGetBookingByID_Handler(t *testing.T) {
// 	mockUserService := new(serviceMock.UserServiceMock)
// 	logger := loggers.NewTestLogger()
// 	userHandler := handlers.NewUserHandler(mockUserService, logger)

// 	bookingUUID := "94841ade-468c-480f-8b69-ee911e6fcbdb"

// 	testCases := []struct {
// 		name           string
// 		urlParamID     string
// 		mockResponse   *dtos.APIResponse
// 		expectedStatus int
// 		expectedError  bool
// 		expectedMsg    string
// 		isMissingID    bool
// 	}{
// 		{
// 			name:       "Success - Booking Retrieved Successfully",
// 			urlParamID: bookingUUID,
// 			mockResponse: &dtos.APIResponse{
// 				Status:  "Success",
// 				Message: "Booking retrieved successfully",
// 				Data: &dtos.BookingResponse{
// 					BookingUUID:  bookingUUID,
// 					CustomerUUID: "customer-uuid-123",
// 					RoomUUID:     "room-uuid-456",
// 					BookingDate:  "2026-05-26",
// 					StartTime:    "18:45",
// 					EndTime:      "20:45",
// 				},
// 				Errors: []dtos.Error{}, // FIXED: Matches your exact custom dtos.Error type
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedError:  false,
// 			expectedMsg:    "Booking retrieved successfully",
// 		},
// 		{
// 			name:       "Failure - Booking Not Found",
// 			urlParamID: "non-existent-uuid",
// 			mockResponse: &dtos.APIResponse{
// 				Status:  "Error",
// 				Message: "Booking not found",
// 				// FIXED: Populating using your native dtos.Error structure fields
// 				Errors: []dtos.Error{
// 					{
// 						Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
// 						Message: "Booking not found",
// 					},
// 				},
// 			},
// 			expectedStatus: http.StatusNotFound,
// 			expectedError:  true,
// 			expectedMsg:    "Booking not found",
// 		},
// 		{
// 			name:       "Failure - Booking ID Is Empty String",
// 			urlParamID: "",
// 			mockResponse: &dtos.APIResponse{
// 				Status:  "Error",
// 				Message: "Booking ID is required",
// 				// FIXED: Populating using your native dtos.Error structure fields
// 				Errors: []dtos.Error{
// 					{
// 						Code:    errorcodes.ErrorCodeStatus[errorcodes.BadRequestErrorCode],
// 						Message: "Booking ID is required",
// 					},
// 				},
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedError:  true,
// 			expectedMsg:    "Booking ID is required",
// 			isMissingID:    true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			mockUserService.ExpectedCalls = nil

// 			req, err := http.NewRequest("GET", "/bookings/"+tc.urlParamID, nil)
// 			assert.NoError(t, err)

// 			req = mux.SetURLVars(req, map[string]string{"id": tc.urlParamID})
// 			rr := httptest.NewRecorder()

// 			mockUserService.On("GetBookingByID", mock.Anything, tc.urlParamID, mock.AnythingOfType("string")).
// 				Return(func(ctx context.Context, id string, reqID string) *dtos.APIResponse {
// 					if tc.mockResponse != nil {
// 						tc.mockResponse.RequestID = reqID
// 					}
// 					return tc.mockResponse
// 				})

// 			// Call target handler method
// 			userHandler.GetBookingByID(rr, req)

// 			// Assertions
// 			assert.Equal(t, tc.expectedStatus, rr.Code)

// 			var response dtos.APIResponse
// 			err = json.Unmarshal(rr.Body.Bytes(), &response)
// 			assert.NoError(t, err)

// 			assert.NotEmpty(t, response.RequestID)

// 			if tc.expectedError {
// 				assert.Equal(t, "Error", response.Status)
// 				assert.Equal(t, tc.expectedMsg, response.Message)
// 			} else {
// 				mockUserService.AssertExpectations(t)
// 				assert.Equal(t, "Success", response.Status)
// 				assert.Equal(t, tc.expectedMsg, response.Message)
// 			}
// 		})
// 	}
// }

// func TestDeleteBooking_Handler(t *testing.T) {
// 	mockUserService := new(serviceMock.UserServiceMock)
// 	logger := loggers.NewTestLogger()
// 	userHandler := handlers.NewUserHandler(mockUserService, logger)

// 	bookingUUID := "abc-123-delete-uuid"

// 	testCases := []struct {
// 		name           string
// 		urlParamID     string
// 		mockResponse   *dtos.APIResponse
// 		expectedStatus int
// 		expectedError  bool
// 		expectedMsg    string
// 	}{
// 		{
// 			name:       "Success - Booking Cancelled Safely",
// 			urlParamID: bookingUUID,
// 			mockResponse: &dtos.APIResponse{
// 				Status:  "Success",
// 				Code:    http.StatusOK,
// 				Message: "Booking deleted successfully",
// 				Errors:  []dtos.Error{}, // Empty array routes to HTTP 200 via your mapper logic
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedError:  false,
// 			expectedMsg:    "Booking deleted successfully",
// 		},
// 		{
// 			name:       "Failure - Booking Not Found",
// 			urlParamID: "invalid-uuid",
// 			mockResponse: &dtos.APIResponse{
// 				Status:  "Error",
// 				Code:    http.StatusNotFound,
// 				Message: "Booking not found",
// 				Errors: []dtos.Error{
// 					{
// 						Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
// 						Message: "Booking not found",
// 					},
// 				},
// 			},
// 			expectedStatus: http.StatusNotFound,
// 			expectedError:  true,
// 			expectedMsg:    "Booking not found",
// 		},
// 		{
// 			name:       "Failure - Booking Was Already Cancelled Prior",
// 			urlParamID: bookingUUID,
// 			mockResponse: &dtos.APIResponse{
// 				Status:  "Error",
// 				Code:    http.StatusBadRequest,
// 				Message: "Booking Already Cancelled...!!",
// 				Errors: []dtos.Error{
// 					{
// 						Code:    errorcodes.ErrorCodeStatus[errorcodes.BadRequestErrorCode],
// 						Message: "Booking Already Cancelled...!!",
// 					},
// 				},
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedError:  true,
// 			expectedMsg:    "Booking Already Cancelled...!!",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Clear mock tracking registers for clean run isolation
// 			mockUserService.ExpectedCalls = nil

// 			// 1. Build context request target path route string
// 			req, err := http.NewRequest("DELETE", "/bookings/"+tc.urlParamID, nil)
// 			assert.NoError(t, err)

// 			// 2. Inject URL mux param tracking context IDs
// 			req = mux.SetURLVars(req, map[string]string{"id": tc.urlParamID})
// 			rr := httptest.NewRecorder()

// 			// 3. Mock Service execution binding pattern expectations
// 			mockUserService.On("DeleteBooking", mock.Anything, tc.urlParamID, mock.AnythingOfType("string")).
// 				Return(func(ctx context.Context, id string, reqID string) *dtos.APIResponse {
// 					if tc.mockResponse != nil {
// 						tc.mockResponse.RequestID = reqID
// 					}
// 					return tc.mockResponse
// 				})

// 			// 4. Invoke concrete handler pipeline method logic
// 			userHandler.DeleteBooking(rr, req)

// 			// 5. Evaluate returned status parameters against expected profiles
// 			assert.Equal(t, tc.expectedStatus, rr.Code)

// 			var response dtos.APIResponse
// 			err = json.Unmarshal(rr.Body.Bytes(), &response)
// 			assert.NoError(t, err)

// 			assert.NotEmpty(t, response.RequestID)

// 			if tc.expectedError {
// 				assert.Equal(t, "Error", response.Status)
// 				assert.Equal(t, tc.expectedMsg, response.Message)
// 			} else {
// 				mockUserService.AssertExpectations(t)
// 				assert.Equal(t, "Success", response.Status)
// 				assert.Equal(t, tc.expectedMsg, response.Message)
// 			}
// 		})
// 	}
// }

// func TestUpdateBooking_Handler(t *testing.T) {
// 	mockUserService := new(serviceMock.UserServiceMock)
// 	logger := loggers.NewTestLogger()
// 	userHandler := handlers.NewUserHandler(mockUserService, logger)

// 	bookingUUID := "update-123-uuid"
// 	validRequestBody := dtos.BookingRequest{
// 		CustomerUUID: "customer-uuid-123",
// 		RoomUUID:     "room-uuid-456",
// 		BookingDate:  "2026-05-26",
// 		StartTime:    "14:00",
// 		EndTime:      "16:00",
// 	}

// 	testCases := []struct {
// 		name           string
// 		urlParamID     string
// 		requestBody    interface{} // Accept interface to allow invalid strings or real structs
// 		mockResponse   *dtos.APIResponse
// 		expectedStatus int
// 		expectedError  bool
// 		expectedMsg    string
// 		isMalformedJSON bool
// 	}{
// 		{
// 			name:        "Success - Booking Updated Successfully",
// 			urlParamID:  bookingUUID,
// 			requestBody: validRequestBody,
// 			mockResponse: &dtos.APIResponse{
// 				Status:  "Success",
// 				Code:    http.StatusOK,
// 				Message: "Booking updated successfully",
// 				Errors:  []dtos.Error{},
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedError:  false,
// 			expectedMsg:    "Booking updated successfully",
// 		},
// 		{
// 			name:        "Failure - Invalid Request Body (Malformed JSON)",
// 			urlParamID:  bookingUUID,
// 			requestBody: "{malformed-json: clear-text-error-here", // Invalid raw string
// 			mockResponse:   nil, // Terminates early, won't reach the service layer
// 			expectedStatus: http.StatusBadRequest,
// 			expectedError:  true,
// 			expectedMsg:    "Invalid request body",
// 			isMalformedJSON: true,
// 		},
// 		{
// 			name:        "Failure - Business Validation Conflict",
// 			urlParamID:  bookingUUID,
// 			requestBody: validRequestBody,
// 			mockResponse: &dtos.APIResponse{
// 				Status:  "Error",
// 				Code:    http.StatusConflict,
// 				Message: "Room is already booked for this time block",
// 				Errors: []dtos.Error{
// 					{
// 						Code:    errorcodes.ErrorCodeStatus[errorcodes.ValidationErrorCode],
// 						Message: "Room is already booked for this time block",
// 					},
// 				},
// 			},
// 			expectedStatus: http.StatusConflict,
// 			expectedError:  true,
// 			expectedMsg:    "Room is already booked for this time block",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			mockUserService.ExpectedCalls = nil

// 			// 1. Marshall input bodies dynamically based on test case variant type
// 			var bodyBytes []byte
// 			var err error
// 			if strBody, ok := tc.requestBody.(string); ok {
// 				bodyBytes = []byte(strBody)
// 			} else {
// 				bodyBytes, err = json.Marshal(tc.requestBody)
// 				assert.NoError(t, err)
// 			}

// 			// 2. Build HTTP request passing the buffered byte payload stream
// 			req, err := http.NewRequest("PUT", "/bookings/"+tc.urlParamID, bytes.NewBuffer(bodyBytes))
// 			assert.NoError(t, err)

// 			req = mux.SetURLVars(req, map[string]string{"id": tc.urlParamID})
// 			rr := httptest.NewRecorder()

// 			// 3. Set up mock service behavior only if JSON parsing succeeds
// 			if !tc.isMalformedJSON {
// 				mockUserService.On("UpdateBooking", mock.Anything, tc.urlParamID, tc.requestBody.(dtos.BookingRequest), mock.AnythingOfType("string")).
// 					Return(func(ctx context.Context, id string, request dtos.BookingRequest, reqID string) *dtos.APIResponse {
// 						if tc.mockResponse != nil {
// 							tc.mockResponse.RequestID = reqID
// 						}
// 						return tc.mockResponse
// 					})
// 			}

// 			// 4. Invoke target update handler method execution
// 			userHandler.UpdateBooking(rr, req)

// 			// 5. Run structural validations
// 			assert.Equal(t, tc.expectedStatus, rr.Code)

// 			var response dtos.APIResponse
// 			err = json.Unmarshal(rr.Body.Bytes(), &response)
// 			assert.NoError(t, err)

// 			assert.NotEmpty(t, response.RequestID)

// 			if tc.expectedError {
// 				assert.Equal(t, "Error", response.Status)
// 				assert.Equal(t, tc.expectedMsg, response.Message)
// 			} else {
// 				mockUserService.AssertExpectations(t)
// 				assert.Equal(t, "Success", response.Status)
// 				assert.Equal(t, tc.expectedMsg, response.Message)
// 			}
// 		})
// 	}
// }





func TestGetBookingByID_Handler(t *testing.T) {
	mockUserService := new(serviceMock.UserServiceMock)
	logger := loggers.NewTestLogger()
	userHandler := handlers.NewUserHandler(mockUserService, logger)

	bookingUUID := "abc-123-get-uuid"

	testCases := []struct {
		name           string
		urlParamID     string
		mockResponse   *dtos.APIResponse
		expectedStatus int
		expectedError  bool
		expectedMsg    string
	}{
		{
			name:       "Success - Booking Found",
			urlParamID: bookingUUID,
			mockResponse: &dtos.APIResponse{
				Status:  "Success",
				Code:    http.StatusOK,
				Message: "Booking fetched successfully",
				Errors:  []dtos.Error{},
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedMsg:    "Booking fetched successfully",
		},
		{
			name:       "Failure - Booking Not Found",
			urlParamID: "invalid-uuid",
			mockResponse: &dtos.APIResponse{
				Status:  "Error",
				Code:    http.StatusNotFound,
				Message: "Booking not found",
				Errors: []dtos.Error{
					{
						Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
						Message: "Booking not found",
					},
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedMsg:    "Booking not found",
		},
		{
			name:       "Failure - Missing Booking ID",
			urlParamID: "",
			mockResponse: nil, // Service should NOT be called
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedMsg:    "Booking ID is required",
		},
		{
			name:       "Failure - Internal Server Error",
			urlParamID: bookingUUID,
			mockResponse: &dtos.APIResponse{
				Status:  "Error",
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
				Errors: []dtos.Error{
					{
						Code:    errorcodes.ErrorCodeStatus[errorcodes.InternalServerErrorCode],
						Message: "Internal server error",
					},
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedMsg:    "Internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear mock tracking registers for clean run isolation
			mockUserService.ExpectedCalls = nil

			// 1. Build context request target path route string
			req, err := http.NewRequest("GET", "/bookings/"+tc.urlParamID, nil)
			assert.NoError(t, err)

			// 2. Inject URL mux param tracking context IDs
			req = mux.SetURLVars(req, map[string]string{"id": tc.urlParamID})
			rr := httptest.NewRecorder()

			// 3. Mock Service execution binding pattern expectations (skip if ID is empty — handler short-circuits)
			if tc.urlParamID != "" {
				mockUserService.On("GetBookingByID", mock.Anything, tc.urlParamID, mock.AnythingOfType("string")).
					Return(func(ctx context.Context, id string, reqID string) *dtos.APIResponse {
						if tc.mockResponse != nil {
							tc.mockResponse.RequestID = reqID
						}
						return tc.mockResponse
					})
			}

			// 4. Invoke concrete handler pipeline method logic
			userHandler.GetBookingByID(rr, req)

			// 5. Evaluate returned status parameters against expected profiles
			assert.Equal(t, tc.expectedStatus, rr.Code)

			var response dtos.APIResponse
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tc.urlParamID != "" {
				assert.NotEmpty(t, response.RequestID)
			}

			if tc.expectedError {
				assert.Equal(t, "Error", response.Status)
				assert.Equal(t, tc.expectedMsg, response.Message)
			} else {
				mockUserService.AssertExpectations(t)
				assert.Equal(t, "Success", response.Status)
				assert.Equal(t, tc.expectedMsg, response.Message)
			}
		})
	}
}


func TestDeleteBooking_Handler(t *testing.T) {
	mockUserService := new(serviceMock.UserServiceMock)
	logger := loggers.NewTestLogger()
	userHandler := handlers.NewUserHandler(mockUserService, logger)

	bookingUUID := "abc-123-delete-uuid"

	testCases := []struct {
		name           string
		urlParamID     string
		mockResponse   *dtos.APIResponse
		expectedStatus int
		expectedError  bool
		expectedMsg    string
	}{
		{
			name:       "Success - Booking Cancelled Successfully",
			urlParamID: bookingUUID,
			mockResponse: &dtos.APIResponse{
				Status:  "Success",
				Code:    http.StatusOK,
				Message: "Booking deleted successfully",
				Errors:  []dtos.Error{},
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedMsg:    "Booking deleted successfully",
		},
		{
			name:           "Failure - Missing Booking ID",
			urlParamID:     "",
			mockResponse:   nil, // Service should NOT be called
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedMsg:    "Booking ID is required",
		},
		{
			name:       "Failure - Booking Not Found",
			urlParamID: "invalid-uuid",
			mockResponse: &dtos.APIResponse{
				Status:  "Error",
				Code:    http.StatusNotFound,
				Message: "Booking not found",
				Errors: []dtos.Error{
					{
						Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
						Message: "Booking not found",
					},
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedMsg:    "Booking not found",
		},
		{
			name:       "Failure - Booking Already Cancelled",
			urlParamID: bookingUUID,
			mockResponse: &dtos.APIResponse{
				Status:  "Error",
				Code:    http.StatusBadRequest,
				Message: "Booking Already Cancelled...!!",
				Errors: []dtos.Error{
					{
						Code:    errorcodes.ErrorCodeStatus[errorcodes.BadRequestErrorCode],
						Message: "Booking Already Cancelled...!!",
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedMsg:    "Booking Already Cancelled...!!",
		},
		{
			name:       "Failure - Internal Server Error",
			urlParamID: bookingUUID,
			mockResponse: &dtos.APIResponse{
				Status:  "Error",
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
				Errors: []dtos.Error{
					{
						Code:    errorcodes.ErrorCodeStatus[errorcodes.InternalServerErrorCode],
						Message: "Internal server error",
					},
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedMsg:    "Internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear mock tracking registers for clean run isolation
			mockUserService.ExpectedCalls = nil

			// 1. Build context request target path route string
			req, err := http.NewRequest("DELETE", "/bookings/"+tc.urlParamID, nil)
			assert.NoError(t, err)

			// 2. Inject URL mux param tracking context IDs
			req = mux.SetURLVars(req, map[string]string{"id": tc.urlParamID})
			rr := httptest.NewRecorder()

			// 3. Mock Service execution binding (skip if ID is empty — handler short-circuits)
			if tc.urlParamID != "" {
				mockUserService.On("DeleteBooking", mock.Anything, tc.urlParamID, mock.AnythingOfType("string")).
					Return(func(ctx context.Context, id string, reqID string) *dtos.APIResponse {
						if tc.mockResponse != nil {
							tc.mockResponse.RequestID = reqID
						}
						return tc.mockResponse
					})
			}

			// 4. Invoke concrete handler pipeline method logic
			userHandler.DeleteBooking(rr, req)

			// 5. Evaluate returned status parameters against expected profiles
			assert.Equal(t, tc.expectedStatus, rr.Code)

			var response dtos.APIResponse
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tc.urlParamID != "" {
				assert.NotEmpty(t, response.RequestID)
			}

			if tc.expectedError {
				assert.Equal(t, "Error", response.Status)
				assert.Equal(t, tc.expectedMsg, response.Message)
			} else {
				mockUserService.AssertExpectations(t)
				assert.Equal(t, "Success", response.Status)
				assert.Equal(t, tc.expectedMsg, response.Message)
			}
		})
	}
}

func TestUpdateBooking_Handler(t *testing.T) {
	mockUserService := new(serviceMock.UserServiceMock)
	logger := loggers.NewTestLogger()
	userHandler := handlers.NewUserHandler(mockUserService, logger)

	// Valid validation-compliant dummy arguments
	validBookingID := "a3b4c5d6-e7f8-9012-a3b4-c5d6e7f89012"
	validRequestPayload := dtos.BookingRequest{
		CustomerUUID: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", 
		RoomUUID:     "f47ac10b-58cc-4372-a567-0e02b2c3d479", 
		BookingDate:  "2026-06-01",
		StartTime:    "10:00",
		EndTime:      "12:00",
	}

	testCases := []struct {
		name              string
		bookingIDParam    string
		requestBody       interface{}
		mockResponse      *dtos.APIResponse
		expectedStatus    int
		expectedStatusStr string
		expectedMsg       string
	}{
		{
			name:           "1. Success - Booking Updated Successfully",
			bookingIDParam: validBookingID,
			requestBody:    validRequestPayload,
			mockResponse: &dtos.APIResponse{
				Status:  "Success",
				Code:    http.StatusOK,
				Message: "Booking updated successfully",
			},
			expectedStatus:    http.StatusOK,
			expectedStatusStr: "Success",
			expectedMsg:       "Booking updated successfully",
		},
		{
			name:           "2. Failure - Invalid JSON Body",
			bookingIDParam: validBookingID,
			requestBody:    `{ invalid-json }`,
			mockResponse:   nil,
			expectedStatus:    http.StatusUnprocessableEntity, // Your handler explicitly returns 422 here
			expectedStatusStr: "Error",
			expectedMsg:       "Invalid request format",
		},
		{
			name:           "3. Failure - Struct Validation Error",
			bookingIDParam: validBookingID,
			requestBody:    dtos.BookingRequest{}, // Empty fields trigger structural 422
			mockResponse:   nil,
			expectedStatus:    http.StatusUnprocessableEntity,
			expectedStatusStr: "Error",
			expectedMsg:       "Validation failed",
		},
		{
			name:           "4. Failure - Empty Booking ID URL Parameter",
			bookingIDParam: "", // Empties out the URL param router tracking register
			requestBody:    validRequestPayload,
			mockResponse:   nil,
			expectedStatus:    http.StatusBadRequest,
			expectedStatusStr: "Error",
			expectedMsg:       "Booking ID is required",
		},
		{
			name:           "5. Failure - Business Conflict from Service Layer",
			bookingIDParam: validBookingID,
			requestBody:    validRequestPayload,
			mockResponse: &dtos.APIResponse{
				Status:  "Error",
				Code:    http.StatusConflict,
				Message: "Room is already booked for this time block",
			},
			expectedStatus:    http.StatusConflict,
			expectedStatusStr: "Error",
			expectedMsg:       "Room is already booked for this time block",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear all history and mock registries completely between test cycles
			mockUserService.ExpectedCalls = nil
			mockUserService.Calls = nil

			var bodyReader io.Reader
			switch v := tc.requestBody.(type) {
			case string:
				bodyReader = strings.NewReader(v)
			default:
				bodyBytes, err := json.Marshal(tc.requestBody)
				assert.NoError(t, err)
				bodyReader = bytes.NewReader(bodyBytes)
			}

			req, err := http.NewRequest("PUT", "/bookings/"+tc.bookingIDParam, bodyReader)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// FIX 1: Map parameter key exactly to "id" as used by your handler code
			req = mux.SetURLVars(req, map[string]string{"id": tc.bookingIDParam})

			rr := httptest.NewRecorder()

			if tc.mockResponse != nil {
				localMockResponse := *tc.mockResponse

				// FIX 3: Perfect structural match with your UserService.UpdateBooking mock signature
				mockUserService.On("UpdateBooking", 
					mock.Anything, 
					tc.bookingIDParam, 
					mock.AnythingOfType("dtos.BookingRequest"), 
					mock.AnythingOfType("string"),
				).Return(func(ctx context.Context, bID string, bookingReq dtos.BookingRequest, reqID string) *dtos.APIResponse {
					localMockResponse.RequestID = reqID
					return &localMockResponse
				}).Once()
			}

			// Invoke handler pipeline
			userHandler.UpdateBooking(rr, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, rr.Code)

			var response dtos.APIResponse
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatusStr, response.Status)
			assert.Contains(t, response.Message, tc.expectedMsg)

			if tc.mockResponse != nil {
				assert.NotEmpty(t, response.RequestID)
				mockUserService.AssertExpectations(t)
			}
		})
	}
}


func TestCreateBooking_Handler(t *testing.T) {
	// Initialize mocked service and test logger dependencies
	mockUserService := new(serviceMock.UserServiceMock)
	logger := loggers.NewTestLogger()
	userHandler := handlers.NewUserHandler(mockUserService, logger)

	// Valid payload conforming to required field and "uuid" validation rules
	validBookingRequest := dtos.BookingRequest{
		CustomerUUID: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", // Valid UUID v1
		RoomUUID:     "f47ac10b-58cc-4372-a567-0e02b2c3d479", // Valid UUID v4
		BookingDate:  "2026-06-01",
		StartTime:    "10:00",
		EndTime:      "12:00",
	}

	testCases := []struct {
		name           string
		requestBody    interface{}
		mockResponse   *dtos.APIResponse
		expectedStatus int
		expectedError  bool
		expectedMsg    string
	}{
		{
			name:        "Success - Booking Created Successfully",
			requestBody: validBookingRequest,
			mockResponse: &dtos.APIResponse{
				Status:  "Success",
				Code:    http.StatusCreated,
				Message: "Booking created successfully",
				Errors:  []dtos.Error{},
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
			expectedMsg:    "Booking created successfully",
		},
		{
			name:           "Failure - Invalid JSON Body",
			requestBody:    `{ invalid-json }`,
			mockResponse:   nil, // Handler exits before hitting service
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedMsg:    "Invalid request format",
		},
		{
			name:           "Failure - Validation Error Missing Required Fields",
			requestBody:    dtos.BookingRequest{}, // Triggers 400 Bad Request
			mockResponse:   nil,                   // Handler exits before hitting service
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedMsg:    "Validation failed",
		},
		{
			name:        "Failure - Room Not Available",
			requestBody: validBookingRequest,
			mockResponse: &dtos.APIResponse{
				Status:  "Error",
				Code:    http.StatusConflict,
				Message: "Room is not available for the selected time slot",
				Errors: []dtos.Error{
					{
						Message: "Room is not available for the selected time slot",
					},
				},
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
			expectedMsg:    "Room is not available for the selected time slot",
		},
		{
			name:        "Failure - Internal Server Error",
			requestBody: validBookingRequest,
			mockResponse: &dtos.APIResponse{
				Status:  "Error",
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
				Errors: []dtos.Error{
					{
						Message: "Internal server error",
					},
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedMsg:    "Internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear out previous mock execution cycles completely 
			mockUserService.ExpectedCalls = nil
			mockUserService.Calls = nil

			// 1. Convert payload types to bytes
			var bodyReader io.Reader
			switch v := tc.requestBody.(type) {
			case string:
				bodyReader = strings.NewReader(v)
			default:
				bodyBytes, err := json.Marshal(tc.requestBody)
				assert.NoError(t, err)
				bodyReader = bytes.NewReader(bodyBytes)
			}

			// 2. Setup the target HTTP call environment
			req, err := http.NewRequest("POST", "/bookings", bodyReader)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			// 3. Register service layer expectations cleanly
			if tc.mockResponse != nil {
				localMockResponse := *tc.mockResponse 

				mockUserService.On("CreateBooking", mock.Anything, mock.AnythingOfType("dtos.BookingRequest"), mock.AnythingOfType("string")).
					Return(func(ctx context.Context, bookingReq dtos.BookingRequest, reqID string) *dtos.APIResponse {
						localMockResponse.RequestID = reqID
						return &localMockResponse
					}).Once() 
			}

			// 4. Fire the handler function logic execution
			userHandler.CreateBooking(rr, req)

			// 5. Test Assertions
			assert.Equal(t, tc.expectedStatus, rr.Code)

			var response dtos.APIResponse
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tc.mockResponse != nil {
				assert.NotEmpty(t, response.RequestID)
				mockUserService.AssertExpectations(t)
			}

			assert.Equal(t, tc.expectedMsg, response.Message)
			if tc.expectedError {
				assert.Equal(t, "Error", response.Status)
			} else {
				assert.Equal(t, "Success", response.Status)
			}
		})
	}
}

