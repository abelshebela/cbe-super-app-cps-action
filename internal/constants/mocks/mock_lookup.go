package mocks

import (
	"cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)
type accountAPIClient struct {
	logger     utils.Logger
	cbeBaseURL string
}

// MockAccountData represents the structure of our mock JSON file
type MockAccountData struct {
	Accounts []MockAccount `json:"accounts"`
}

type MockAccount struct {
	AccountNumber string  `json:"account_number"`
	AccountName   string  `json:"account_name"`
	AccountType   string  `json:"account_type"`
	Status        string  `json:"status"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
ActiveAccount 	bool `json:"active_account"`
AccountDormant  bool `json:"account_dormant"`
AccountFrozen   bool `json:"account_frozen"`


}

type AccountResponse struct {
	Data model.AccountInfo `json:"data"`
}