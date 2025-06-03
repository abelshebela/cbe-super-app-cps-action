package users

type FullName struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

type LinkedAccountDetail struct {
	AccountNumber     string `json:"account_number"`
	AccountBranchCode string `json:"account_branch_code"`
	LinkedBranch      string `json:"linked_branch"`
	IsAccountActive   bool   `json:"is_account_active"`
	LinkedStatus      bool   `json:"linked_status"`
	CurrencyCode      string `json:"currency_code"`
}

type User struct {
	ID        string
	FullName  FullName
	IsDeleted bool
}

type LinkedAccountResponse struct {
	UserID         string                `json:"user_id"`
	FullName       FullName              `json:"full_name"`
	LinkedAccounts []LinkedAccountDetail `json:"linked_accounts"`
}