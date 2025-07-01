package mocks

import (
	"context"

	"cbe-super-app-cps-action/internal/domain/action"
	"cbe-super-app-cps-action/internal/domain/service"

	"github.com/stretchr/testify/mock"
)

type ServiceDetailsRepository struct {
	mock.Mock
}

func (m *ServiceDetailsRepository) GetAllServiceDetails(ctx context.Context) ([]*service.Service, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*service.Service), args.Error(1)
}

func (m *ServiceDetailsRepository) GetOneServiceDetail(ctx context.Context, id string) (service.Service, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(service.Service), args.Error(1)
}

func (m *ServiceDetailsRepository) UpdateOneServiceDetail(ctx context.Context, id string, update service.Service) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *ServiceDetailsRepository) FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error) {
	args := m.Called(ctx, uniqueID)
	return args.Get(0).([]action.CPSAction), args.Error(1)
}
