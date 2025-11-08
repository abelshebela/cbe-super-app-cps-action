package external_call

import (
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/yohannesteshome/coreio/core"
)

type AccountLookUp struct {
	CoreApi core.CBECoreAPIInterface
}

func NewAccountLookUp(coreCfg core.CBECoreCredential) AccountLookUp {
	coreapi := core.NewCBECoreAPI(coreCfg)
	return AccountLookUp{
		CoreApi: coreapi,
	}
}

func (a *AccountLookUp) LookupAccountByAccountNumber(account string) (*types.AccountLookupData, error) {
	response, err := a.CoreApi.AccountLookup(core.AccountLookupParam{
		AccountNumber: account,
	})
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, err
	}
	return &types.AccountLookupData{
		AccountNumber:  response.Detail.AccountNumber,
		CustomerName:   response.Detail.CustomerName,
		Restriction:    response.Detail.Restriction,
		Currency:       response.Detail.Currency,
		AccountType:    response.Detail.AccountType,
		WorkingBalance: response.Detail.WorkingBalance,
		CustomerID:     response.Detail.CustomerID,
		// AccountStatus: response.Detail.AccountStatus,
		// AccountHolder: response.AccountHolder,
	}, nil
}
