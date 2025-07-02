package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type ActionRepository struct {
	mock.Mock
}

func (m *ActionRepository) CreateCpsAction(ctx context.Context, cpsAction action.CPSAction) (action.CPSAction, error) {
	args := m.Called(ctx, cpsAction)
	return args.Get(0).(action.CPSAction), args.Error(1)
}

func (m *ActionRepository) FetchCpsActionById(ctx context.Context, id string) (action.CPSAction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(action.CPSAction), args.Error(1)
}

func (m *ActionRepository) UpdateCpsAction(ctx context.Context, cpsAction action.CPSAction) error {
	args := m.Called(ctx, cpsAction)
	return args.Error(0)
}

func (m *ActionRepository) FetchAccountsByAccountNumber(ctx context.Context, accountNumber string) ([]action.LinkedAccount, error) {
	args := m.Called(ctx, accountNumber)
	return args.Get(0).([]action.LinkedAccount), args.Error(1)
}

func (m *ActionRepository) FetchLastCpsActionByMakerID(ctx context.Context, makerId string) (action.CPSAction, error) {
	args := m.Called(ctx, makerId)
	return args.Get(0).(action.CPSAction), args.Error(1)
}

func (m *ActionRepository) GetAllHqServices(ctx context.Context) ([]action.ServiceDetails, error) {
	args := m.Called(ctx)
	return args.Get(0).([]action.ServiceDetails), args.Error(1)
}

func (m *ActionRepository) GetAllHqServicesPaginated(ctx context.Context, offset, limit int) ([]action.ServiceDetails, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]action.ServiceDetails), args.Error(1)
}

func (m *ActionRepository) GetHqServiceById(ctx context.Context, id string) (action.ServiceDetails, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(action.ServiceDetails), args.Error(1)
}

func (m *ActionRepository) UpdateHqService(ctx context.Context, service action.ServiceDetails) error {
	args := m.Called(ctx, service)
	return args.Error(0)
}

func (m *ActionRepository) FetchLinkedAccountById(ctx context.Context, id []string) ([]action.LinkedAccount, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]action.LinkedAccount), args.Error(1)
}

func (m *ActionRepository) UpdateAccounts(ctx context.Context, linkedAccounts []action.LinkedAccount) ([]action.LinkedAccount, error) {
	args := m.Called(ctx, linkedAccounts)
	return args.Get(0).([]action.LinkedAccount), args.Error(1)
}

func (m *ActionRepository) UpdateAccount(ctx context.Context, linkedAccount action.LinkedAccount) (action.LinkedAccount, error) {
	args := m.Called(ctx, linkedAccount)
	return args.Get(0).(action.LinkedAccount), args.Error(1)
}
