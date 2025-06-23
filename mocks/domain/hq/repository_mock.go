package mock

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/hq"
)

type MockRepository struct {
	GetHQByIDFunc                     func(ctx context.Context, id string) (hq.HQ, error)
	UpdateHQFunc                      func(ctx context.Context, id string, hq hq.HQ) error
	FetchPendingActionsByUniqueIDFunc func(ctx context.Context, uniqueID string) ([]action.CPSAction, error)
}

func (m *MockRepository) GetHQByID(ctx context.Context, id string) (hq.HQ, error) {
	if m.GetHQByIDFunc != nil {
		return m.GetHQByIDFunc(ctx, id)
	}
	return hq.HQ{}, nil
}

func (m *MockRepository) UpdateHQ(ctx context.Context, id string, hq hq.HQ) error {
	if m.UpdateHQFunc != nil {
		return m.UpdateHQFunc(ctx, id, hq)
	}
	return nil
}

func (m *MockRepository) FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error) {
	if m.FetchPendingActionsByUniqueIDFunc != nil {
		return m.FetchPendingActionsByUniqueIDFunc(ctx, uniqueID)
	}
	return nil, nil
}
