package mocks

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/stretchr/testify/mock"
)

type Account struct {
	mock.Mock
}

func (m *Account) LookupAccountByPhone(ctx context.Context, phone string) (bool, error) {
	args := m.Called(ctx, phone)
	return args.Bool(0), args.Error(1)
}

func (m *Account) LookupAccountByAccountNumber(ctx context.Context, account model.AccountLookUpRequest) (*model.AccountInfo, error) {
	args := m.Called(ctx, account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AccountInfo), args.Error(1)
}
