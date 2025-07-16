package repository

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/stretchr/testify/mock"
)

// MockCPSActionRepo is a mock implementation of the CPSUserRepo interface.
type MockCPSActionRepo struct {
	mock.Mock
}

func (m *MockCPSActionRepo) CreateUserRequest(ctx context.Context, cpsAction model.CPSAction) (*model.CPSAction, error) {
	args := m.Called(ctx, cpsAction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CPSAction), args.Error(1)
}

func (m *MockCPSActionRepo) UpdateUserRequest(ctx context.Context, cpsAction model.CPSAction, userCode string) (*model.CPSAction, error) {
	args := m.Called(ctx, cpsAction, userCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CPSAction), args.Error(1)
}

func (m *MockCPSActionRepo) GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.CPSAction), args.Error(1)
}

func (m *MockCPSActionRepo) ApproveUserAction(ctx context.Context, actionCode string, checker model.CPSAction) (*model.CPSAction, error) {
	args := m.Called(ctx, actionCode, checker)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CPSAction), args.Error(1)
}

func (m *MockCPSActionRepo) FetchUserByUserCode(ctx context.Context, userCode string) (*model.CPSUser, error) {
	args := m.Called(ctx, userCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CPSUser), args.Error(1)
}
