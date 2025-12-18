package accountlookup

import (
	"cbe-super-app-cps-action/internal/constants"
)

type AccountLookUpRequest struct {
	AccountNumber string `json:"account_number,omitempty"`
	PhoneNumber   string `json:"phone_number"`
}

type CreateAccountRequest struct {
	CustomerName       string                `json:"customer_name"`
	Gender             constants.Gender      `json:"gender"`
	PhoneNumber        string                `json:"phone_number"`
	CustomerMotherName string                `json:"customer_mother_name"`
	AccountType        string                `json:"account_type"`
	AccountBranchType  constants.AccountType `json:"account_branchtype"`
	Picture            string                `json:"picture"`
	CustomerAddress    string                `json:"customer_address"`
}
