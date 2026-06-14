
package mock

import (
	"context"

	"lynxis-gate/training-service/internal/dtos"

	"github.com/stretchr/testify/mock"
)

// UserServiceMock is a mock for UserServiceInterface
type UserServiceMock struct {
	mock.Mock
}

// CreateUser mocks the CreateUser method
func (m *UserServiceMock) CreateUser(ctx context.Context, req dtos.CreateUserRequest) *dtos.APIResponse {
	args := m.Called(ctx, req)
	return args.Get(0).(*dtos.APIResponse)
}

// GetUserByUUID mocks the GetUserByUUID method
func (m *UserServiceMock) GetUserByUUID(ctx context.Context, userUUID string) *dtos.APIResponse {
	args := m.Called(ctx, userUUID)
	return args.Get(0).(*dtos.APIResponse)
}

// UpdateUser mocks the UpdateUser method
func (m *UserServiceMock) UpdateUser(ctx context.Context, userUUID string, req dtos.UpdateUserRequest) *dtos.APIResponse {
	args := m.Called(ctx, userUUID, req)
	return args.Get(0).(*dtos.APIResponse)
}

// DeleteUser mocks the DeleteUser method
func (m *UserServiceMock) DeleteUser(ctx context.Context, userUUID string) *dtos.APIResponse {
	args := m.Called(ctx, userUUID)
	return args.Get(0).(*dtos.APIResponse)
}

// GetAllUsers mocks the GetAllUsers method
func (m *UserServiceMock) GetAllUsers(ctx context.Context) *dtos.APIResponse {
	args := m.Called(ctx)
	return args.Get(0).(*dtos.APIResponse)
}