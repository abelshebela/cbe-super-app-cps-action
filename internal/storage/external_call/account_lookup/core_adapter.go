package account_lookup

import (
	"cbe-super-app-cps-action/internal/constants"
	accountLookup "cbe-super-app-cps-action/internal/constants/dto/account_lookup"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants/localization"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/hugokessem/coreio/core"
)

type Account interface {
	LookupAccountByPhone(ctx context.Context, phone string) (bool, error)
	LookupAccountByAccountNumber(ctx context.Context, account model.AccountLookUpRequest) (*model.AccountDetail, error)
	LookupAccountByAccountNumberFromBps(ctx context.Context, accountNumber string) (accountLookup.AccountResponse, error)
	CreateAccountWithFayda(ctx context.Context, account accountLookup.AccountCreateParams) (types.Account, error)
	CifSearch(ctx context.Context, cif string) ([]imodel.AccountData, error)
}

type CoreAccountLookupAdapter struct {
	coreAPI        core.CBECoreAPIInterface
	Client         *http.Client
	BaseUrl        string
	PhoneUrlPath   string
	AccountUrlPath string
	FaydaUrlPath   string
	Logger         utils.Logger
}

func NewCoreAccountLookupAdapter(coreAPI core.CBECoreAPIInterface, baseUrl string, timeout time.Duration, logger utils.Logger) Account {
	return &CoreAccountLookupAdapter{
		coreAPI:        coreAPI,
		BaseUrl:        baseUrl,
		Client:         &http.Client{Timeout: timeout},
		FaydaUrlPath:   "bps_banking/core/with_fayda",
		PhoneUrlPath:   "bps_banking/core/phone_number",
		AccountUrlPath: "bps_banking/core/account_number",
		Logger:         logger,
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

func (a *CoreAccountLookupAdapter) LookupAccountByAccountNumber(ctx context.Context, account model.AccountLookUpRequest) (accountDetail *model.AccountDetail, err error) {
	a.Logger.Infof("Looking up account by number: %s", account.AccountNumber)
	defer func() {
		if recovered := recover(); recovered != nil {
			a.Logger.Errorf("Panic occurred while validating account number %s: %v", account.AccountNumber, recovered)
			accountDetail = nil
			err = fmt.Errorf("core name lookup panic: %v", recovered)
			err = fmt.Errorf("Unable to lookup account number %s at the moment, please update the number", account.AccountNumber)
		}

	}()

	response, err := a.coreAPI.NameLookup(ctx, core.NameLookupParam{
		AccountNumber: account.AccountNumber,
	})

	if err != nil {
		a.Logger.Errorf("Error occurred while validating account number: %v", err)
		return nil, err
	}
	if !response.Success || response.Detail == nil {
		a.Logger.Errorf("Account not found for number: %s response : %v", account.AccountNumber, err)
		return nil, errors.New(localization.ErrorAccountNumberNotFound.Code)
	}

	// if response.Detail == nil {
	// 	a.Logger.Errorf("Account not found for number: %s", account.AccountNumber)
	// 	return nil, errors.New(localization.ErrorAccountNotFound.Code)
	// }

	detail := response.Detail
	return &model.AccountDetail{
		AccountNumber:  detail.AccountNumber,
		CustomerName:   detail.AccountName,
		Restriction:    detail.RestrictionType,
		Currency:       detail.Currency,
		WorkingBalance: "",
		CustomerID:     detail.CustomerNumber,
		AccountType:    detail.RestrictionType,
	}, nil
}

func (a *CoreAccountLookupAdapter) LookupAccountByAccountNumberFromBps(ctx context.Context, accountNumber string) (accountLookup.AccountResponse, error) {
	res, accountInfo, err := BPSBankingClient(ctx, a.Client, constants.AccountNumber, accountNumber, a.BaseUrl+a.AccountUrlPath)
	if err != nil {
		return accountLookup.AccountResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		a.Logger.Errorf("Failed to lookup account by number: %s", accountNumber)
		return accountLookup.AccountResponse{}, nil
	}

	return accountInfo, nil
}

func (a *CoreAccountLookupAdapter) CreateAccountWithFayda(ctx context.Context, account accountLookup.AccountCreateParams) (types.Account, error) {
	res, _, err := BPSBankingClient(ctx, a.Client, constants.WithFayda, account, a.BaseUrl+a.FaydaUrlPath)
	if err != nil {
		a.Logger.Errorf("Failed to create account with Fayda: %v", err)
		return types.Account{}, err
	}
	defer res.Body.Close()

	// to avoid import error for customer_creation package, can be removed when the response is used

	if res.StatusCode != http.StatusOK {
		a.Logger.Errorf("Failed to create account with Fayda, status code: %d", res.StatusCode)
		return types.Account{}, nil
	}

	return types.Account{}, nil
}

func (s *CoreAccountLookupAdapter) CifSearch(ctx context.Context, cif string) ([]imodel.AccountData, error) {

	search, err := s.coreAPI.AccountList(ctx, core.AccountListParam{
		ColumnName:    "CUS.ID",
		CriteriaValue: cif,
	})
	if err != nil {
		s.Logger.Errorf("[AccountLookup][CifSearch] failed to search CIF: %v", err)
		return nil, err
	}

	// Log the details count before processing
	detailsCount := len(search.Details)

	// Debug: Log the full search structure to understand the response format
	s.Logger.Infof("CIF search response received")

	if detailsCount == 0 {
		s.Logger.Warnf("CIF search successful but search.Details is empty")
		return []imodel.AccountData{}, nil
	}

	var accounts []imodel.AccountData
	for _, account := range search.Details {
		s.Logger.Infof("Processing account from CIF search")

		phoneNumberTrimed := strings.ReplaceAll(account.PhoneNo, "+", "")
		accounts = append(accounts, imodel.AccountData{
			AccountNumber:   account.AccountNumber,
			CustomerID:      cif,
			CustomerName:    account.CustomerName,
			AccountType:     account.AccountType,
			PhoneNumber:     phoneNumberTrimed,
			BirthOfDate:     account.BirthOfDate,
			Gender:          account.Gender,
			Currency:        account.Currency,
			BranchCode:      account.BranchCode,
			Branch:          account.BranchName,
			CustomerSegment: account.CustomerSegment,
			Restriction:     account.Restriction,
			RestrictionType: account.RestrictionType,
			Email:           account.Email,
		})
	}

	s.Logger.Infof("CIF search completed successfully")
	return accounts, nil
}
