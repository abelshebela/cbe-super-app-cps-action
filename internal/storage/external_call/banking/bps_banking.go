package mockbanking

import (
	"cbe-super-app-cps-action/internal/constants"
	accountLookup "cbe-super-app-cps-action/internal/constants/dto/account_lookup"
	"cbe-super-app-cps-action/internal/constants/types"
	"context"
	"net/http"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BpsBankingClient struct {
	BaseUrl        string
	Client         *http.Client
	Logger         utils.Logger
	FaydaUrlPath   string
	PhoneUrlPath   string
	AccountUrlPath string
}

func NewBpsBankingClient(baseUrl string, timeout time.Duration, logger utils.Logger) BpsBankingClient {
	return BpsBankingClient{
		BaseUrl:        baseUrl,
		Client:         &http.Client{Timeout: timeout},
		Logger:         logger,
		FaydaUrlPath:   "bps_banking/core/with_fayda",
		PhoneUrlPath:   "bps_banking/core/phone_number",
		AccountUrlPath: "bps_banking/core/account_number",
	}
}

func (b *BpsBankingClient) LookupAccountByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	res, _, err := BPSBankingClient(ctx, b.Client, constants.PhoneNumber, phoneNumber, b.BaseUrl+b.PhoneUrlPath)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func (b *BpsBankingClient) LookupAccountByAccountNumber(ctx context.Context, accountNumber string) (accountLookup.AccountResponse, error) {
	res, accountInfo, err := BPSBankingClient(ctx, b.Client, constants.AccountNumber, accountNumber, b.BaseUrl+b.AccountUrlPath)
	if err != nil {
		return accountLookup.AccountResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return accountLookup.AccountResponse{}, nil
	}

	return accountInfo, nil
}

func (b *BpsBankingClient) CreateAccountWithFayda(ctx context.Context, account accountLookup.CreateAccountRequest) (*types.Account, error) {
	res, _, err := BPSBankingClient(ctx, b.Client, constants.WithFayda, account, b.BaseUrl+b.FaydaUrlPath)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil
	}

	return nil, nil
}
