package model

type AccountLookUpRequest struct {
	AccountNumber string `json:"account_number,omitempty"`
	PhoneNumber   string `json:"phone_number,omitempty"`
}

type AccountDetail struct {
	AccountNumber  string `json:"account_number"`
	CustomerName   string `json:"customer_name"`
	Restriction    string `json:"restriction"`
	Currency       string `json:"currency"`
	WorkingBalance string `json:"working_balance"`
	CustomerID     string `json:"customer_id"`
	AccountType    string `json:"account_type"`
}

type AccountInfo struct {
	ID                 string `json:"id"`
	AccountBranchType  string `json:"account_branchtype"`
	AccountBranchCode  string `json:"account_branchcode"`
	AccountNumber      string `json:"account_number"`
	CustomerNumber     string `json:"customer_number"`
	CustomerName       string `json:"customer_name"`
	AccountDescription string `json:"account_description"`
	PhoneNumber        string `json:"phone_number"`
	CustomerAddress    string `json:"customer_address"`
	DebitAllowed       bool   `json:"debit_allowed"`
	CreditAllowed      bool   `json:"credit_allowed"`
	AccountType        string `json:"account_type"`
	AccountFrozen      bool   `json:"account_frozen"`
	AccountDormant     bool   `json:"account_dormant"`
	ActiveAccount      bool   `json:"active_account"`
	AccountCurrency    string `json:"account_currency"`
}
