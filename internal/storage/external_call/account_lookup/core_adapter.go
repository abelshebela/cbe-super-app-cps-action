package account_lookup

import (
	"cbe-super-app-cps-action/internal/constants"
	accountLookup "cbe-super-app-cps-action/internal/constants/dto/account_lookup"
	"cbe-super-app-cps-action/internal/constants/types"
	"context"
	"errors"
	"net/http"
	"time"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"

	"github.com/hugokessem/coreio/core"
)

type Account interface {
	LookupAccountByPhone(ctx context.Context, phone string) (bool, error)
	LookupAccountByAccountNumber(ctx context.Context, account model.AccountLookUpRequest) (*model.AccountDetail, error)
	LookupAccountByAccountNumberFromBps(ctx context.Context, accountNumber string) (accountLookup.AccountResponse, error)
	CreateAccountWithFayda(ctx context.Context, account accountLookup.CreateAccountRequest) (types.Account, error)
}

type CoreAccountLookupAdapter struct {
	coreAPI        core.CBECoreAPIInterface
	Client         *http.Client
	BaseUrl        string
	PhoneUrlPath   string
	AccountUrlPath string
	FaydaUrlPath   string
}

func NewCoreAccountLookupAdapter(coreAPI core.CBECoreAPIInterface, baseUrl string, timeout time.Duration) Account {
	return &CoreAccountLookupAdapter{
		coreAPI:        coreAPI,
		BaseUrl:        baseUrl,
		Client:         &http.Client{Timeout: timeout},
		FaydaUrlPath:   "bps_banking/core/with_fayda",
		PhoneUrlPath:   "bps_banking/core/phone_number",
		AccountUrlPath: "bps_banking/core/account_number",
	}
}

func (a *CoreAccountLookupAdapter) LookupAccountByPhone(ctx context.Context, phone string) (bool, error) {
	res, _, err := BPSBankingClient(ctx, a.Client, constants.PhoneNumber, phone, a.BaseUrl+a.PhoneUrlPath)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
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

func (a *CoreAccountLookupAdapter) LookupAccountByAccountNumberFromBps(ctx context.Context, accountNumber string) (accountLookup.AccountResponse, error) {
	res, accountInfo, err := BPSBankingClient(ctx, a.Client, constants.AccountNumber, accountNumber, a.BaseUrl+a.AccountUrlPath)
	if err != nil {
		return accountLookup.AccountResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return accountLookup.AccountResponse{}, nil
	}

	return accountInfo, nil
}

func (a *CoreAccountLookupAdapter) CreateAccountWithFayda(ctx context.Context, account accountLookup.CreateAccountRequest) (types.Account, error) {
	res, _, err := BPSBankingClient(ctx, a.Client, constants.WithFayda, account, a.BaseUrl+a.FaydaUrlPath)
	if err != nil {
		return types.Account{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return types.Account{}, nil
	}

	return types.Account{}, nil
}
