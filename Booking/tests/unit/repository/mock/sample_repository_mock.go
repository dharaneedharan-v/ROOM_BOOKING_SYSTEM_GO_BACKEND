
package mock

import (
	"context"

	"lynxis-gate/training-service/internal/dtos"
	"lynxis-gate/training-service/internal/models"

	"github.com/stretchr/testify/mock"
)

// UserRepositoryMock is a mock of UserRepositoryInterface
type UserRepositoryMock struct {
	mock.Mock
}

// CreateUser mocks the CreateUser method
func (m *UserRepositoryMock) CreateUser(ctx context.Context, user *models.User) *dtos.Error {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dtos.Error)
}

// GetUserByUUID mocks the GetUserByUUID method
func (m *UserRepositoryMock) GetUserByUUID(ctx context.Context, userUUID string) (*models.User, *dtos.Error) {
	args := m.Called(ctx, userUUID)
	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}
	var err *dtos.Error
	if args.Get(1) != nil {
		err = args.Get(1).(*dtos.Error)
	}
	return user, err
}

// UpdateUser mocks the UpdateUser method
func (m *UserRepositoryMock) UpdateUser(ctx context.Context, userUUID string, updateData map[string]interface{}) *dtos.Error {
	args := m.Called(ctx, userUUID, updateData)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dtos.Error)
}

// DeleteUser mocks the DeleteUser method
func (m *UserRepositoryMock) DeleteUser(ctx context.Context, userUUID string) *dtos.Error {
	args := m.Called(ctx, userUUID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dtos.Error)
}

// GetAllUsers mocks the GetAllUsers method
func (m *UserRepositoryMock) GetAllUsers(ctx context.Context) ([]models.User, int64, *dtos.Error) {
	args := m.Called(ctx)
	var users []models.User
	if args.Get(0) != nil {
		users = args.Get(0).([]models.User)
	}
	count := args.Get(1).(int64)
	var err *dtos.Error
	if args.Get(2) != nil {
		err = args.Get(2).(*dtos.Error)
	}
	return users, count, err
}