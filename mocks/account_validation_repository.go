package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/account_validation"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type AccountValidationRepository struct {
	mock.Mock
}

func (m *AccountValidationRepository) GetAccountValidationByID(ctx context.Context, id string) (account_validation.ValidationRule, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(account_validation.ValidationRule), args.Error(1)
}

func (m *AccountValidationRepository) UpdateAccountValidation(ctx context.Context, id string, update account_validation.ValidationRule) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *AccountValidationRepository) FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error) {
	args := m.Called(ctx, uniqueID)
	return args.Get(0).([]action.CPSAction), args.Error(1)
}
