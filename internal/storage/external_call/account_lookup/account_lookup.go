package account_lookup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"time"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"

	mock "cbe-super-app-cps-action/internal/constants/mocks"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accountAPIClient struct {
	logger     utils.Logger
	cbeBaseURL string
	client     *http.Client
}

type Account interface {
	LookupAccountByPhone(ctx context.Context, phone string) (bool, error)
	LookupAccountByAccountNumber(ctx context.Context, account model.AccountLookUpRequest) (*model.AccountInfo, error)
}

// NewAccountAPIClient creates a new instance of AccountAPIClient
func InitAccountAPIClient(cbeBaseUrl string, timeout time.Duration, logger utils.Logger) Account {
	return &accountAPIClient{
		logger:     logger,
		cbeBaseURL: cbeBaseUrl,
		client:     &http.Client{Timeout: timeout},
	}
}


func (b *accountAPIClient) LookupAccountByPhone(ctx context.Context, phone string) (bool, error) {
	b.logger.Infof("Looking up phone number: %s using mock data", phone)

	// Read the mock JSON file
	jsonFile, err := os.Open("corebanking.json")
	if err != nil {
		b.logger.Errorf("failed to open corebanking.json: %v", err)
		return false, errors.New(localization.ErrorExternalServiceError.Code)
	}
	defer jsonFile.Close()

	// Read the file content
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		b.logger.Errorf("failed to read corebanking.json: %v", err)
		return false, errors.New(localization.ErrorExternalServiceError.Code)
	}

	// Parse the JSON data
	var mockData mock.MockAccountData
	if err := json.Unmarshal(jsonData, &mockData); err != nil {
		b.logger.Errorf("failed to unmarshal mock account data: %v", err)
		return false, errors.New(localization.ErrorExternalServiceError.Code)
	}
	if len(mockData.Accounts) > 0 {
		b.logger.Infof("Phone number %s found in mock data", phone)
		return true, nil
	}

	b.logger.Warnf("Phone number %s not found in mock data", phone)
	return false, nil
}

// func (b *accountAPIClient) isTimeoutError(err error, ctx context.Context) bool {
// 	return errors.Is(err, context.DeadlineExceeded) ||
// 		errors.Is(ctx.Err(), context.DeadlineExceeded) ||
// 		strings.Contains(err.Error(), "timeout") ||
// 		strings.Contains(err.Error(), "deadline exceeded")
// }

func (b *accountAPIClient) LookupAccountByAccountNumber(ctx context.Context, account model.AccountLookUpRequest) (*model.AccountInfo, error) {

	data, err := json.Marshal(account)
	if err != nil {
		b.logger.Errorf("failed to marshal account request: %v", err)
		return nil, errors.New(localization.ErrorExternalServiceError.Code)
	}
	fmt.Println(b.cbeBaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/%s", b.cbeBaseURL, "/bps_banking/core/account_number"), bytes.NewBuffer(data))
	if err != nil {
		b.logger.Errorf("failed to create HTTP request: %v", err)
		return nil, errors.New(localization.ErrorExternalServiceError.Code)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := b.client.Do(req)
	if err != nil {
		if b.isTimeoutError(err, ctx) {
			b.logger.Errorf("request timeout for account %s: %v", account.AccountNumber, err)
					return nil, errors.New(localization.ErrorExternalServiceError.Code)
		}

		b.logger.Errorf("HTTP request failed for account %s: %v", account.AccountNumber, err)
				return nil, errors.New(localization.ErrorExternalServiceError.Code)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		b.logger.Errorf("account not found with status %v", res.StatusCode)
				return nil, errors.New(localization.ErrorAccountNumberNotFound.Code)
	}

	if res.StatusCode != http.StatusOK {
		b.logger.Errorf("request failed with status code %v", res.StatusCode)
				return nil, errors.New(localization.ErrorExternalServiceError.Code)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		b.logger.Errorf("failed to read response body for account %s: %v", account.AccountNumber, err)
				return nil, errors.New(localization.ErrorExternalServiceError.Code)
	}

	var accountResp model.AccountLookupResponse
if err := json.Unmarshal(resBody, &accountResp); err != nil {
	b.logger.Errorf("failed to unmarshal account info for account %s: %v", account.AccountNumber, err)
	return nil, errors.New(localization.ErrorExternalServiceError.Code)
}

return &accountResp.Data, nil
}

func (b *accountAPIClient) isTimeoutError(err error, ctx context.Context) bool {
	return errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(ctx.Err(), context.DeadlineExceeded) ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "deadline exceeded")
}

