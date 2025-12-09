package model

// type Transaction struct {
// 	ID                      string            `json:"id"`
// 	TransactionID           string            `json:"transaction_id"`
// 	FTNumber                string            `json:"ft_number"`
// 	DebitBranchCode         string            `json:"debit_branch_code"`
// 	DebitDistrictCode       string            `json:"debit_district_code"`
// 	DebitUserID             string            `json:"debit_user_id"`
// 	DebitAccountNumber      string            `json:"debit_account_number"`
// 	DebitAccountHolderName  string            `json:"debit_account_holder_name"`
// 	CreditUserID            string            `json:"credit_user_id"`
// 	CreditAccountNumber     string            `json:"credit_account_number"`
// 	CreditAccountHolderName string            `json:"credit_account_holder_name"`
// 	InstitutionCode         string            `json:"institution_code"`
// 	InstitutionName         string            `json:"institution_name"`
// 	Currency                Currency          `json:"currency"`
// 	ServiceFee              decimal.Decimal   `json:"service_fee"`
// 	TipAmount               decimal.Decimal   `json:"tip_amount"`
// 	PaidAmount              decimal.Decimal   `json:"paid_amount"`
// 	VAT                     decimal.Decimal   `json:"vat"`
// 	Amount                  decimal.Decimal   `json:"amount"`
// 	TotalAmount             decimal.Decimal   `json:"total_amount"`
// 	ExternalReference       string            `json:"external_reference,omitempty"`
// 	TransactionReason       string            `json:"transaction_reason,omitempty"`
// 	TransactionType         TransactionType   `json:"transaction_type"`
// 	TransactionStatus       TransactionStatus `json:"transaction_status"`
// 	IsIFB                   bool              `json:"is_ifb"`
// 	IsReversed              bool              `json:"is_reversed"`
// 	PaidAt                  time.Time         `json:"paid_at,omitzero"`
// 	ReversedAt              time.Time         `json:"reversed_at,omitzero"`
// 	Metadata                json.RawMessage   `json:"metadata,omitempty"`
// 	CreatedAt               time.Time         `json:"created_at"`
// 	LastModifiedAt          time.Time         `json:"last_modified_at,omitzero"`
// }
