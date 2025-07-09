package mock

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
)

type MockRepository struct {
	GetPortalCardFunc func(ctx context.Context) ([]model.Card, error)
}

func (m *MockRepository) GetAllPortalCard(ctx context.Context, id string) ([]model.Card, error) {
	if m.GetPortalCardFunc != nil {
		return m.GetPortalCardFunc(ctx)
	}
	return []model.Card{}, nil
}
