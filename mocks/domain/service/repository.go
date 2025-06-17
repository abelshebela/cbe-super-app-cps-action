package mock

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/service"
)

type MockServiceRepository struct {
	GetServiceDetailsByIDFunc         func(ctx context.Context, id string) (service.ServiceDetails, error)
	UpdateServiceDetailsFunc          func(ctx context.Context, id string, update service.ServiceDetails) error
	FetchPendingActionsByUniqueIDFunc func(ctx context.Context, uniqueID string) ([]action.CPSAction, error)
	GetAllServiceDetailsFunc          func(ctx context.Context) ([]*service.Service, error)
	GetOneServiceDetailFunc           func(ctx context.Context, id string) (service.Service, error)
	UpdateOneServiceDetailFunc        func(ctx context.Context, id string, update service.Service) error
}

func (m *MockServiceRepository) GetServiceDetailsByID(ctx context.Context, id string) (service.ServiceDetails, error) {
	return m.GetServiceDetailsByIDFunc(ctx, id)
}

func (m *MockServiceRepository) UpdateServiceDetails(ctx context.Context, id string, update service.ServiceDetails) error {
	return m.UpdateServiceDetailsFunc(ctx, id, update)
}

func (m *MockServiceRepository) FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error) {
	return m.FetchPendingActionsByUniqueIDFunc(ctx, uniqueID)
}

func (m *MockServiceRepository) GetAllServiceDetails(ctx context.Context) ([]*service.Service, error) {
	return m.GetAllServiceDetailsFunc(ctx)
}

func (m *MockServiceRepository) GetOneServiceDetail(ctx context.Context, id string) (service.Service, error) {
	return m.GetOneServiceDetailFunc(ctx, id)
}

func (m *MockServiceRepository) UpdateOneServiceDetail(ctx context.Context, id string, update service.Service) error {
	return m.UpdateOneServiceDetailFunc(ctx, id, update)
}
