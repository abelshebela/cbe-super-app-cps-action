package accountlookup

import "time"

type SearchQuery struct {
	AccountNumber string
	PhoneNumber   string
}

type UserSearchResult struct {
	ID                 string    `json:"id"`
	AccountBranchType  string    `json:"account_branchtype"`
	AccountBranchCode  string    `json:"account_branchcode"`
	AccountNumber      string    `json:"account_number"`
	CustomerNumber     string    `json:"customer_number"`
	CustomerName       string    `json:"customer_name"`
	AccountDescription string    `json:"account_description"`
	PhoneNumber        string    `json:"phone_number"`
	CustomerAddress    string    `json:"customer_address"`
	DebitAllowed       bool      `json:"debit_allowed"`
	CreditAllowed      bool      `json:"credit_allowed"`
	AccountType        string    `json:"account_type"`
	AccountFrozen      bool      `json:"account_frozen"`
	AccountDormant     bool      `json:"account_dormant"`
	ActiveAccount      bool      `json:"active_account"`
	AccountCurrency    string    `json:"account_currency"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
