package action_manual_mock

import (
	"context"

	"cbe-super-app-cps-action/internal/domain/action"
)

type MockActionRepository struct {
	CreateCpsActionFunc              func(ctx context.Context, a action.CPSAction) (action.CPSAction, error)
	FetchCpsActionByIdFunc           func(ctx context.Context, id string) (action.CPSAction, error)
	UpdateCpsActionFunc              func(ctx context.Context, a action.CPSAction) error
	FetchAccountsByAccountNumberFunc func(ctx context.Context, accountNumber string) ([]action.LinkedAccount, error)
	FetchLinkedAccountByIdFunc       func(ctx context.Context, id []string) ([]action.LinkedAccount, error)
	UpdateAccountsFunc               func(ctx context.Context, linkedAccounts []action.LinkedAccount) ([]action.LinkedAccount, error)
	UpdateAccountFunc                func(ctx context.Context, linkedAccount action.LinkedAccount) (action.LinkedAccount, error)
	GetAllHqServicesFunc             func(ctx context.Context) ([]action.ServiceDetails, error)
	GetAllHqServicesPaginatedFunc    func(ctx context.Context, offset, limit int) ([]action.ServiceDetails, error)
	GetHqServiceByIdFunc             func(ctx context.Context, id string) (action.ServiceDetails, error)
	UpdateHqServiceFunc              func(ctx context.Context, service action.ServiceDetails) error
	FetchLastCpsActionByMakerIDFunc  func(ctx context.Context, makerId string) (action.CPSAction, error)
}

func (m *MockActionRepository) CreateCpsAction(ctx context.Context, a action.CPSAction) (action.CPSAction, error) {
	if m.CreateCpsActionFunc != nil {
		return m.CreateCpsActionFunc(ctx, a)
	}
	return action.CPSAction{}, nil
}

func (m *MockActionRepository) FetchCpsActionById(ctx context.Context, id string) (action.CPSAction, error) {
	if m.FetchCpsActionByIdFunc != nil {
		return m.FetchCpsActionByIdFunc(ctx, id)
	}
	return action.CPSAction{}, nil
}

func (m *MockActionRepository) UpdateCpsAction(ctx context.Context, a action.CPSAction) error {
	if m.UpdateCpsActionFunc != nil {
		return m.UpdateCpsActionFunc(ctx, a)
	}
	return nil
}

func (m *MockActionRepository) FetchAccountsByAccountNumber(ctx context.Context, accountNumber string) ([]action.LinkedAccount, error) {
	if m.FetchAccountsByAccountNumberFunc != nil {
		return m.FetchAccountsByAccountNumberFunc(ctx, accountNumber)
	}
	return nil, nil
}

func (m *MockActionRepository) FetchLinkedAccountById(ctx context.Context, id []string) ([]action.LinkedAccount, error) {
	if m.FetchLinkedAccountByIdFunc != nil {
		return m.FetchLinkedAccountByIdFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockActionRepository) UpdateAccounts(ctx context.Context, linkedAccounts []action.LinkedAccount) ([]action.LinkedAccount, error) {
	if m.UpdateAccountsFunc != nil {
		return m.UpdateAccountsFunc(ctx, linkedAccounts)
	}
	return nil, nil
}

func (m *MockActionRepository) UpdateAccount(ctx context.Context, linkedAccount action.LinkedAccount) (action.LinkedAccount, error) {
	if m.UpdateAccountFunc != nil {
		return m.UpdateAccountFunc(ctx, linkedAccount)
	}
	return action.LinkedAccount{}, nil
}

func (m *MockActionRepository) GetAllHqServices(ctx context.Context) ([]action.ServiceDetails, error) {
	if m.GetAllHqServicesFunc != nil {
		return m.GetAllHqServicesFunc(ctx)
	}
	return nil, nil
}

func (m *MockActionRepository) GetAllHqServicesPaginated(ctx context.Context, offset, limit int) ([]action.ServiceDetails, error) {
	if m.GetAllHqServicesPaginatedFunc != nil {
		return m.GetAllHqServicesPaginatedFunc(ctx, offset, limit)
	}
	return nil, nil
}

func (m *MockActionRepository) GetHqServiceById(ctx context.Context, id string) (action.ServiceDetails, error) {
	if m.GetHqServiceByIdFunc != nil {
		return m.GetHqServiceByIdFunc(ctx, id)
	}
	return action.ServiceDetails{}, nil
}

func (m *MockActionRepository) UpdateHqService(ctx context.Context, service action.ServiceDetails) error {
	if m.UpdateHqServiceFunc != nil {
		return m.UpdateHqServiceFunc(ctx, service)
	}
	return nil
}

func (m *MockActionRepository) FetchLastCpsActionByMakerID(ctx context.Context, makerId string) (action.CPSAction, error) {
	if m.FetchLastCpsActionByMakerIDFunc != nil {
		return m.FetchLastCpsActionByMakerIDFunc(ctx, makerId)
	}
	return action.CPSAction{}, nil
}
