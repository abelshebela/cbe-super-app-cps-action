package account_lookup

import (
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/yohannesteshome/coreio/core"
)

type Account interface {
	LookupAccountByPhone(ctx context.Context, phone string) (bool, error)
	LookupAccountByAccountNumber(ctx context.Context, account model.AccountLookUpRequest) (*model.AccountDetail, error)
}

type CoreAccountLookupAdapter struct {
	coreAPI core.CBECoreAPIInterface
}

func NewCoreAccountLookupAdapter(coreAPI core.CBECoreAPIInterface) Account {
	return &CoreAccountLookupAdapter{
		coreAPI: coreAPI,
	}
}

func (a *CoreAccountLookupAdapter) LookupAccountByPhone(ctx context.Context, phone string) (bool, error) {
	return false, errors.New(localization.ErrorExternalServiceError.Code)
}

func (a *CoreAccountLookupAdapter) LookupAccountByAccountNumber(ctx context.Context, account model.AccountLookUpRequest) (*model.AccountDetail, error) {
	response, err := a.coreAPI.AccountLookup(core.AccountLookupParam{
		AccountNumber: account.AccountNumber,
	})
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, errors.New(localization.ErrorExternalServiceError.Code)
	}
	if response.Detail == nil {
		return nil, errors.New(localization.ErrorAccountNumberNotFound.Code)
	}

	detail := response.Detail
	return &model.AccountDetail{
		AccountNumber:  detail.AccountNumber,
		CustomerName:   detail.CustomerName,
		Restriction:    detail.Restriction,
		Currency:       detail.Currency,
		WorkingBalance: detail.WorkingBalance,
		CustomerID:     detail.CustomerID,
		AccountType:    detail.AccountType,
	}, nil
}
