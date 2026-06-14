
package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"lynxis-gate/training-service/internal/dtos"
	"lynxis-gate/training-service/internal/errorcodes"
	"lynxis-gate/training-service/internal/handlers"
	"lynxis-gate/training-service/internal/loggers"
	serviceMock "lynxis-gate/training-service/tests/unit/services/mock"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	// Create mocks
	mockUserService := new(serviceMock.UserServiceMock)
	logger := loggers.NewLogger("test-service")

	// Create handler with mock
	userHandler := handlers.NewUserHandler(mockUserService, logger)

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		mockResponse   *dtos.APIResponse
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Success - Create User",
			requestBody: map[string]interface{}{
				"name":    "John Doe",
				"age":     30,
				"address": "123 Main St",
			},
			mockResponse: &dtos.APIResponse{
				Status:  "Success",
				Code:    http.StatusOK,
				Message: "User created successfully",
				Data: &dtos.UserResponse{
					UserUUID:  "test-uuid",
					Name:      "John Doe",
					Age:       30,
					Address:   "123 Main St",
					CreatedAt: "2023-01-01T12:00:00Z",
					UpdatedAt: "2023-01-01T12:00:00Z",
				},
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "Failure - Invalid Request",
			requestBody: map[string]interface{}{
				"name": "", // Empty name is invalid
				"age":  30,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert request body to JSON
			reqBody, _ := json.Marshal(tc.requestBody)

			// Create a request
			req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Record response
			rr := httptest.NewRecorder()

			// Setup mock expectations
			if !tc.expectedError {
				mockUserService.On("CreateUser", mock.Anything, mock.Anything).Return(tc.mockResponse)
			}

			// Call the handler
			userHandler.CreateUser(rr, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, rr.Code)

			if !tc.expectedError {
				// Verify mock was called
				mockUserService.AssertExpectations(t)

				// Check response
				var response dtos.APIResponse
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tc.mockResponse.Status, response.Status)
				assert.Equal(t, tc.mockResponse.Message, response.Message)
			}
		})
	}
}

func TestGetUserByUUID(t *testing.T) {
	// Create mocks
	mockUserService := new(serviceMock.UserServiceMock)
	logger := loggers.NewLogger("test-service")

	// Create handler with mock
	userHandler := handlers.NewUserHandler(mockUserService, logger)

	// Test cases
	testCases := []struct {
		name           string
		userUUID       string
		mockResponse   *dtos.APIResponse
		expectedStatus int
		expectedError  bool
	}{
		{
			name:     "Success - Get User",
			userUUID: "valid-uuid",
			mockResponse: &dtos.APIResponse{
				Status:  "Success",
				Code:    http.StatusOK,
				Message: "User fetched successfully",
				Data: &dtos.UserResponse{
					UserUUID:  "valid-uuid",
					Name:      "John Doe",
					Age:       30,
					Address:   "123 Main St",
					CreatedAt: "2023-01-01T12:00:00Z",
					UpdatedAt: "2023-01-01T12:00:00Z",
				},
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:     "Failure - User Not Found",
			userUUID: "invalid-uuid",
			mockResponse: &dtos.APIResponse{
				Status:  "Error",
				Code:    http.StatusNotFound,
				Message: "User not found",
				Errors: []dtos.Error{
					{
						Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
						Message: "User not found",
					},
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("GET", "/api/v1/users/"+tc.userUUID, nil)
			assert.NoError(t, err)

			// Add gorilla/mux vars
			vars := map[string]string{
				"uuid": tc.userUUID,
			}
			req = mux.SetURLVars(req, vars)

			// Record response
			rr := httptest.NewRecorder()

			// Setup mock expectations
			mockUserService.On("GetUserByUUID", mock.Anything, tc.userUUID).Return(tc.mockResponse)

			// Call the handler
			userHandler.GetUserByUUID(rr, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, rr.Code)

			// Verify mock was called
			mockUserService.AssertExpectations(t)

			// Check response
			var response dtos.APIResponse
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tc.mockResponse.Status, response.Status)
			assert.Equal(t, tc.mockResponse.Message, response.Message)
		})
	}
}

func TestGetUserByUUID_MissingUUID(t *testing.T) {
	// Create mocks
	mockUserService := new(serviceMock.UserServiceMock)
	logger := loggers.NewLogger("test-service")

	// Create handler with mock
	userHandler := handlers.NewUserHandler(mockUserService, logger)

	// Create request WITHOUT setting mux vars
	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	assert.NoError(t, err)

	// Note: do NOT call mux.SetURLVars here, so uuid == ""
	rr := httptest.NewRecorder()

	// Call handler
	userHandler.GetUserByUUID(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Ensure service was NOT called
	mockUserService.AssertNotCalled(t, "GetUserByUUID", mock.Anything, mock.Anything)

	// Assert response body
	var resp dtos.APIResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "Error", resp.Status)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, "Validation failed for the request", resp.Message)
	if assert.Len(t, resp.Errors, 1) {
		assert.Equal(t, "body", resp.Errors[0].Field)
		assert.Equal(t, "Invalid JSON format", resp.Errors[0].Message)
		assert.Equal(t,
			errorcodes.ErrorCodeStatus[errorcodes.InvalidJSONFormatErrorCode],
			resp.Errors[0].Code,
		)
	}
}

func TestUpdateUser(t *testing.T) {
	// Create mocks
	mockUserService := new(serviceMock.UserServiceMock)
	logger := loggers.NewLogger("test-service")

	// Create handler with mock
	userHandler := handlers.NewUserHandler(mockUserService, logger)

	// Test cases
	testCases := []struct {
		name           string
		userUUID       string
		requestBody    map[string]interface{}
		mockResponse   *dtos.APIResponse
		expectedStatus int
		expectedError  bool
	}{
		{
			name:     "Success - Update User",
			userUUID: "valid-uuid",
			requestBody: map[string]interface{}{
				"name":    "Updated Name",
				"age":     35,
				"address": "456 New Address",
			},
			mockResponse: &dtos.APIResponse{
				Status:  "Success",
				Code:    http.StatusOK,
				Message: "User updated successfully",
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:     "Failure - Invalid Request",
			userUUID: "valid-uuid",
			requestBody: map[string]interface{}{
				"name": "", // Empty name is invalid
				"age":  35,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:     "Failure - User Not Found",
			userUUID: "invalid-uuid",
			requestBody: map[string]interface{}{
				"name":    "Updated Name",
				"age":     35,
				"address": "456 New Address",
			},
			mockResponse: &dtos.APIResponse{
				Status:  "Error",
				Code:    http.StatusNotFound,
				Message: "User not found",
				Errors: []dtos.Error{
					{
						Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
						Message: "User not found",
					},
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert request body to JSON
			reqBody, _ := json.Marshal(tc.requestBody)

			// Create a request
			req, err := http.NewRequest("PUT", "/api/v1/users/"+tc.userUUID, bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Add gorilla/mux vars
			vars := map[string]string{
				"uuid": tc.userUUID,
			}
			req = mux.SetURLVars(req, vars)

			// Record response
			rr := httptest.NewRecorder()

			// Setup mock expectations
			if !tc.expectedError {
				mockUserService.On("UpdateUser", mock.Anything, tc.userUUID, mock.Anything).Return(tc.mockResponse)
			}

			// Call the handler
			userHandler.UpdateUser(rr, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, rr.Code)

			if !tc.expectedError {
				// Verify mock was called
				mockUserService.AssertExpectations(t)

				// Check response
				var response dtos.APIResponse
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tc.mockResponse.Status, response.Status)
				assert.Equal(t, tc.mockResponse.Message, response.Message)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	// Create mocks
	mockUserService := new(serviceMock.UserServiceMock)
	logger := loggers.NewLogger("test-service")

	// Create handler with mock
	userHandler := handlers.NewUserHandler(mockUserService, logger)

	// Test cases
	testCases := []struct {
		name           string
		userUUID       string
		mockResponse   *dtos.APIResponse
		expectedStatus int
	}{
		{
			name:     "Success - Delete User",
			userUUID: "valid-uuid",
			mockResponse: &dtos.APIResponse{
				Status:  "Success",
				Code:    http.StatusOK,
				Message: "User deleted successfully",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Failure - User Not Found",
			userUUID: "invalid-uuid",
			mockResponse: &dtos.APIResponse{
				Status:  "Error",
				Code:    http.StatusNotFound,
				Message: "User not found",
				Errors: []dtos.Error{
					{
						Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
						Message: "User not found",
					},
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("DELETE", "/api/v1/users/"+tc.userUUID, nil)
			assert.NoError(t, err)

			// Add gorilla/mux vars
			vars := map[string]string{
				"uuid": tc.userUUID,
			}
			req = mux.SetURLVars(req, vars)

			// Record response
			rr := httptest.NewRecorder()

			// Setup mock expectations
			mockUserService.On("DeleteUser", mock.Anything, tc.userUUID).Return(tc.mockResponse)

			// Call the handler
			userHandler.DeleteUser(rr, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, rr.Code)

			// Verify mock was called
			mockUserService.AssertExpectations(t)

			// Check response
			var response dtos.APIResponse
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tc.mockResponse.Status, response.Status)
			assert.Equal(t, tc.mockResponse.Message, response.Message)
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	// Create mocks
	mockUserService := new(serviceMock.UserServiceMock)
	logger := loggers.NewLogger("test-service")

	// Create handler with mock
	userHandler := handlers.NewUserHandler(mockUserService, logger)

	// Test case 1: Get all users - success case with users
	t.Run("Get All Users - With Data", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("GET", "/api/v1/users", nil)
		assert.NoError(t, err)

		// Record response
		rr := httptest.NewRecorder()

		// Mock data
		usersList := &dtos.UserListResponse{
			Users: []dtos.UserResponse{
				{
					UserUUID:  "uuid-1",
					Name:      "User 1",
					Age:       30,
					Address:   "Address 1",
					CreatedAt: "2023-01-01T12:00:00Z",
					UpdatedAt: "2023-01-01T12:00:00Z",
				},
				{
					UserUUID:  "uuid-2",
					Name:      "User 2",
					Age:       35,
					Address:   "Address 2",
					CreatedAt: "2023-01-02T12:00:00Z",
					UpdatedAt: "2023-01-02T12:00:00Z",
				},
			},
			Total: 2,
		}

		// Create mock response
		mockResponse := &dtos.APIResponse{
			Status:  "Success",
			Code:    http.StatusOK,
			Message: "All users fetched successfully",
			Data:    usersList,
		}

		// Setup mock expectations
		mockUserService.On("GetAllUsers", mock.Anything).Return(mockResponse)

		// Call the handler
		userHandler.GetAllUsers(rr, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify mock was called
		mockUserService.AssertExpectations(t)

		// Check response
		var response dtos.APIResponse
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, mockResponse.Status, response.Status)
		assert.Equal(t, mockResponse.Message, response.Message)
	})

	// Test case 2: Get all users - empty list
	t.Run("Get All Users - Empty List", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("GET", "/api/v1/users", nil)
		assert.NoError(t, err)

		// Record response
		rr := httptest.NewRecorder()

		// Mock data
		emptyList := &dtos.UserListResponse{
			Users: []dtos.UserResponse{},
			Total: 0,
		}

		// Create mock response
		mockResponse := &dtos.APIResponse{
			Status:  "Success",
			Code:    http.StatusOK,
			Message: "All users fetched successfully",
			Data:    emptyList,
		}

		// Setup mock expectations
		mockUserService.On("GetAllUsers", mock.Anything).Return(mockResponse)

		// Call the handler
		userHandler.GetAllUsers(rr, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify mock was called
		mockUserService.AssertExpectations(t)

		// Check response
		var response dtos.APIResponse
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, mockResponse.Status, response.Status)
		assert.Equal(t, mockResponse.Message, response.Message)
	})
}