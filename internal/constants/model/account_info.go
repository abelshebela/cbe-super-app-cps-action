package model

type AccountLookUpRequest struct {
	AccountNumber string `json:"account_number,omitempty"`
	PhoneNumber   string `json:"phone_number,omitempty"`
}

// type AccountDetail struct {
// 	AccountNumber  string `json:"account_number"`
// 	CustomerName   string `json:"customer_name"`
// 	Restriction    string `json:"restriction"`
// 	Currency       string `json:"currency"`
// 	WorkingBalance string `json:"working_balance"`
// 	CustomerID     string `json:"customer_id"`
// 	AccountType    string `json:"account_type"`
// }

type AccountDetail struct {
	AccountNumber        string `json:"account_number"`
	CustomerName         string `json:"customer_name"`
	CustomerNumber       string `json:"customer_number"`
	PhoneNumber          string `json:"phone_number"`
	BranchCode           string `json:"branch_code"`
	BranchName           string `json:"branch_name"`
	Restriction          string `json:"restriction"`
	Currency             string `json:"currency"`
	WorkingBalance       string `json:"working_balance"`
	CustomerSegmentation string `json:"customer_segmentation"`
	CustomerID           string `json:"customer_id"`
	AccountType          string `json:"account_type"`
}
