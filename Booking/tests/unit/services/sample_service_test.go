
package services_test

import (
	"context"
	"testing"
	"time"

	"lynxis-gate/training-service/internal/dtos"
	"lynxis-gate/training-service/internal/errorcodes"
	"lynxis-gate/training-service/internal/loggers"
	"lynxis-gate/training-service/internal/models"
	"lynxis-gate/training-service/internal/services"
	repoMock "lynxis-gate/training-service/tests/unit/repository/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	// Create mocks
	mockRepo := new(repoMock.UserRepositoryMock)
	logger := loggers.NewLogger("test-service")

	// Create service with mock repository
	userService := services.NewUserService(mockRepo, logger)

	// Test cases
	testCases := []struct {
		name          string
		request       dtos.CreateUserRequest
		mockRepoError *dtos.Error
		expectedError bool
	}{
		{
			name: "Success - Create User",
			request: dtos.CreateUserRequest{
				Name:    "John Doe",
				Age:     30,
				Address: "123 Main St",
			},
			mockRepoError: nil,
			expectedError: false,
		},
		{
			name: "Failure - Repository Error",
			request: dtos.CreateUserRequest{
				Name:    "John Doe",
				Age:     30,
				Address: "123 Main St",
			},
			mockRepoError: &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
				Message: "Failed to create user",
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(tc.mockRepoError).Once()

			// Call the service
			response := userService.CreateUser(context.Background(), tc.request)

			// Assertions
			assert.NotNil(t, response)

			if tc.expectedError {
				assert.Equal(t, "Error", response.Status)
				assert.Equal(t, "Failed to create User", response.Message)
			} else {
				assert.Equal(t, "Success", response.Status)
				assert.Equal(t, "User created successfully", response.Message)
				assert.NotNil(t, response.Data)

				userResp, ok := response.Data.(*dtos.UserResponse)
				assert.True(t, ok)
				assert.Equal(t, tc.request.Name, userResp.Name)
				assert.Equal(t, tc.request.Age, userResp.Age)
				assert.Equal(t, tc.request.Address, userResp.Address)
			}

			// Verify mock was called
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetUserByUUID(t *testing.T) {
	// Create mocks
	mockRepo := new(repoMock.UserRepositoryMock)
	logger := loggers.NewLogger("test-service")

	// Create service with mock repository
	userService := services.NewUserService(mockRepo, logger)

	// Mock data
	uuid := "test-uuid"
	now := time.Now()
	mockUser := &models.User{
		ID:        1,
		UserUUID:  uuid,
		Name:      "John Doe",
		Age:       30,
		Address:   "123 Main St",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test cases
	testCases := []struct {
		name          string
		userUUID      string
		mockUser      *models.User
		mockError     *dtos.Error
		expectedError bool
	}{
		{
			name:          "Success - Get User",
			userUUID:      uuid,
			mockUser:      mockUser,
			mockError:     nil,
			expectedError: false,
		},
		{
			name:     "Failure - User Not Found",
			userUUID: "invalid-uuid",
			mockUser: nil,
			mockError: &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
				Message: "User not found",
			},
			expectedError: true,
		},
		{
			name:     "Failure - Database Error",
			userUUID: uuid,
			mockUser: nil,
			mockError: &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
				Message: "Failed to fetch user",
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo.On("GetUserByUUID", mock.Anything, tc.userUUID).Return(tc.mockUser, tc.mockError).Once()

			// Call the service
			response := userService.GetUserByUUID(context.Background(), tc.userUUID)

			// Assertions
			assert.NotNil(t, response)

			if tc.expectedError {
				assert.Equal(t, "Error", response.Status)
				assert.Equal(t, "Failed to fetch user", response.Message)
			} else {
				assert.Equal(t, "Success", response.Status)
				assert.Equal(t, "User fetched successfully", response.Message)
				assert.NotNil(t, response.Data)

				userResp, ok := response.Data.(*dtos.UserResponse)
				assert.True(t, ok)
				assert.Equal(t, tc.mockUser.UserUUID, userResp.UserUUID)
				assert.Equal(t, tc.mockUser.Name, userResp.Name)
				assert.Equal(t, tc.mockUser.Age, userResp.Age)
				assert.Equal(t, tc.mockUser.Address, userResp.Address)
			}

			// Verify mock was called
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	// Create mocks
	mockRepo := new(repoMock.UserRepositoryMock)
	logger := loggers.NewLogger("test-service")

	// Create service with mock repository
	userService := services.NewUserService(mockRepo, logger)

	// Test cases
	testCases := []struct {
		name          string
		userUUID      string
		request       dtos.UpdateUserRequest
		mockError     *dtos.Error
		expectedError bool
	}{
		{
			name:     "Success - Update User",
			userUUID: "valid-uuid",
			request: dtos.UpdateUserRequest{
				Name:    "Updated Name",
				Age:     35,
				Address: "456 New St",
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:     "Failure - User Not Found",
			userUUID: "invalid-uuid",
			request: dtos.UpdateUserRequest{
				Name:    "Updated Name",
				Age:     35,
				Address: "456 New St",
			},
			mockError: &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
				Message: "User not found",
			},
			expectedError: true,
		},
		{
			name:     "Failure - Database Error",
			userUUID: "valid-uuid",
			request: dtos.UpdateUserRequest{
				Name:    "Updated Name",
				Age:     35,
				Address: "456 New St",
			},
			mockError: &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
				Message: "Failed to update user",
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock repository - Verify that updateData map has expected values
			mockRepo.On("UpdateUser", mock.Anything, tc.userUUID, mock.MatchedBy(func(updateData map[string]interface{}) bool {
				return updateData["name"] == tc.request.Name && updateData["age"] == tc.request.Age && updateData["address"] == tc.request.Address
			})).Return(tc.mockError).Once()

			// Call the service
			response := userService.UpdateUser(context.Background(), tc.userUUID, tc.request)

			// Assertions
			assert.NotNil(t, response)

			if tc.expectedError {
				assert.Equal(t, "Error", response.Status)
				assert.Equal(t, "Failed to update user", response.Message)
			} else {
				assert.Equal(t, "Success", response.Status)
				assert.Equal(t, "User updated successfully", response.Message)
			}

			// Verify mock was called
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	// Create mocks
	mockRepo := new(repoMock.UserRepositoryMock)
	logger := loggers.NewLogger("test-service")

	// Create service with mock repository
	userService := services.NewUserService(mockRepo, logger)

	// Mock data
	now := time.Now()
	mockUser := &models.User{
		ID:        1,
		UserUUID:  "valid-uuid",
		Name:      "John Doe",
		Age:       30,
		Address:   "123 Main St",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test cases
	testCases := []struct {
		name             string
		userUUID         string
		mockGetUserError *dtos.Error
		mockUser         *models.User
		mockDeleteError  *dtos.Error
		expectedError    bool
	}{
		{
			name:             "Success - Delete User",
			userUUID:         "valid-uuid",
			mockGetUserError: nil,
			mockUser:         mockUser,
			mockDeleteError:  nil,
			expectedError:    false,
		},
		{
			name:     "Failure - User Not Found",
			userUUID: "invalid-uuid",
			mockGetUserError: &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
				Message: "User not found",
			},
			mockUser:        nil,
			mockDeleteError: nil,
			expectedError:   true,
		},
		{
			name:             "Failure - Delete Error",
			userUUID:         "valid-uuid",
			mockGetUserError: nil,
			mockUser:         mockUser,
			mockDeleteError: &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
				Message: "Failed to delete user",
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo.On("GetUserByUUID", mock.Anything, tc.userUUID).Return(tc.mockUser, tc.mockGetUserError).Once()

			// Only expect DeleteUser to be called if GetUserByUUID didn't return an error
			if tc.mockGetUserError == nil {
				mockRepo.On("DeleteUser", mock.Anything, tc.userUUID).Return(tc.mockDeleteError).Once()
			}

			// Call the service
			response := userService.DeleteUser(context.Background(), tc.userUUID)

			// Assertions
			if tc.expectedError {
				assert.NotNil(t, response)
				assert.Equal(t, "Error", response.Status)
				if tc.mockGetUserError != nil {
					assert.Equal(t, "Fetch to fetch user", response.Message) // This matches the error message in the service
				} else {
					assert.Equal(t, "Failed to delete user", response.Message)
				}
			} else {
				assert.Nil(t, response) // DeleteUser returns nil on success
			}

			// Verify mock was called
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	// Create mocks
	mockRepo := new(repoMock.UserRepositoryMock)
	logger := loggers.NewLogger("test-service")

	// Create service with mock repository
	userService := services.NewUserService(mockRepo, logger)

	// Mock data
	now := time.Now()
	mockUsers := []models.User{
		{
			ID:        1,
			UserUUID:  "uuid-1",
			Name:      "User 1",
			Age:       30,
			Address:   "Address 1",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			UserUUID:  "uuid-2",
			Name:      "User 2",
			Age:       35,
			Address:   "Address 2",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Test cases
	testCases := []struct {
		name          string
		mockUsers     []models.User
		mockCount     int64
		mockError     *dtos.Error
		expectedError bool
		expectedCount int64
	}{
		{
			name:          "Success - Get All Users",
			mockUsers:     mockUsers,
			mockCount:     2,
			mockError:     nil,
			expectedError: false,
			expectedCount: 2,
		},
		{
			name:      "Success - Empty Users List",
			mockUsers: []models.User{},
			mockCount: 0,
			mockError: nil,
			expectedError: false,
			expectedCount: 0,
		},
		{
			name:      "Failure - Database Error",
			mockUsers: nil,
			mockCount: 0,
			mockError: &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
				Message: "Failed to fetch users",
			},
			expectedError: true,
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo.On("GetAllUsers", mock.Anything).Return(tc.mockUsers, tc.mockCount, tc.mockError).Once()

			// Call the service
			response := userService.GetAllUsers(context.Background())

			// Assertions
			assert.NotNil(t, response)

			if tc.expectedError {
				assert.Equal(t, "Error", response.Status)
				assert.Equal(t, "Failed to fetch all user", response.Message)
			} else {
				assert.Equal(t, "Success", response.Status)
				assert.Equal(t, "All users fetched successfully", response.Message)

				userListResp, ok := response.Data.(*dtos.UserListResponse)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCount, userListResp.Total)
				assert.Equal(t, len(tc.mockUsers), len(userListResp.Users))

				// Check user data
				if tc.expectedCount > 0 {
					for i, user := range tc.mockUsers {
						assert.Equal(t, user.UserUUID, userListResp.Users[i].UserUUID)
						assert.Equal(t, user.Name, userListResp.Users[i].Name)
						assert.Equal(t, user.Age, userListResp.Users[i].Age)
						assert.Equal(t, user.Address, userListResp.Users[i].Address)
					}
				}
			}

			// Verify mock was called
			mockRepo.AssertExpectations(t)
		})
	}
} 