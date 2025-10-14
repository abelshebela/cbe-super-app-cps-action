package account_lookup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"time"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"

	mock "cbe-super-app-cps-action/internal/constants/mocks"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accountAPIClient struct {
	logger     utils.Logger
	cbeBaseURL string
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
	}
}

func (b *accountAPIClient) LookupAccountByAccountNumber(ctx context.Context, account model.AccountLookUpRequest) (*model.AccountInfo, error) {
	b.logger.Infof("Looking up account number: %s using mock data", account.AccountNumber)

	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Get the root folder (assuming executable is in a subfolder like ./bin)
	rootDir := filepath.Dir(execPath)

	// Build the full path to the JSON file in the root folder
	jsonPath := filepath.Join(rootDir, "corebanking.json")

	jsonFile, err := os.Open(jsonPath)

	// Read the mock JSON file
	// jsonFile, err := os.Open("./corebanking.json")
	if err != nil {
		b.logger.Errorf("failed to open corebanking.json: %v on path:%v", err, jsonPath)
		return nil, errors.New(localization.ErrorExternalServiceError.Code)
	}
	defer jsonFile.Close()

	// Read the file content
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		b.logger.Errorf("failed to read corebanking.json: %v", err)
		return nil, errors.New(localization.ErrorExternalServiceError.Code)
	}

	// Parse the JSON data
	var mockData mock.MockAccountData
	if err := json.Unmarshal(jsonData, &mockData); err != nil {
		b.logger.Errorf("failed to unmarshal mock account data: %v", err)
		return nil, errors.New(localization.ErrorExternalServiceError.Code)
	}

	// Look for the account number in the mock data
	for _, mockAccount := range mockData.Accounts {
		if mockAccount.AccountNumber == account.AccountNumber {
			b.logger.Infof("Found account: %s - %s", mockAccount.AccountNumber, mockAccount.AccountName)

			// Convert mock account to AccountInfo
			accountInfo := &model.AccountInfo{
				ID:                 mockAccount.AccountNumber, // Using account number as ID
				AccountNumber:      mockAccount.AccountNumber,
				CustomerName:       mockAccount.AccountName,
				AccountType:        mockAccount.AccountType,
				AccountDormant:     mockAccount.AccountDormant,
				AccountCurrency:    mockAccount.Currency,
				AccountDescription: fmt.Sprintf("%s - %s", mockAccount.AccountType, mockAccount.Status),
				AccountFrozen:      mockAccount.AccountFrozen,
				ActiveAccount:      mockAccount.ActiveAccount,
			}

			return accountInfo, nil
		}
	}

	// Account not found in mock data
	b.logger.Warnf("Account number %s not found in mock data", account.AccountNumber)
	return nil, errors.New(localization.ErrorAccountNumberNotFound.Code)
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
